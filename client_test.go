package ripllo

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

// ─── Test helpers ──────────────────────────────────────────────────────

type recordedRequest struct {
	Method     string
	Path       string // includes querystring
	PathOnly   string
	Headers    http.Header
	Body       []byte
	OnBehalfOf string
}

type testServer struct {
	t        *testing.T
	server   *httptest.Server
	requests []recordedRequest
	// route handlers keyed by "METHOD PATH-without-querystring"
	handler func(w http.ResponseWriter, r *http.Request, body []byte)
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()
	ts := &testServer{t: t}
	ts.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()
		pathOnly := r.URL.Path
		ts.requests = append(ts.requests, recordedRequest{
			Method:     r.Method,
			Path:       r.URL.RequestURI(),
			PathOnly:   pathOnly,
			Headers:    r.Header.Clone(),
			Body:       body,
			OnBehalfOf: r.Header.Get("X-Ripllo-On-Behalf-Of"),
		})
		if ts.handler != nil {
			ts.handler(w, r, body)
			return
		}
		// Default: echo a success envelope with empty data.
		writeEnvelope(w, http.StatusOK, json.RawMessage(`{}`), nil, "")
	}))
	t.Cleanup(ts.server.Close)
	return ts
}

func writeEnvelope(w http.ResponseWriter, status int, data json.RawMessage, errEnv *envelopeError, requestID string) {
	body := map[string]any{
		"data":  data,
		"error": errEnv,
		"meta":  map[string]any{"requestId": requestID},
	}
	if errEnv == nil {
		body["error"] = nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func mustClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	c, err := NewClient(ClientOptions{
		KeyID:   "AKIARPLO_test",
		Secret:  "secret_test_xyz",
		BaseURL: baseURL,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	// Deterministic clock + idem keys so we can verify signatures
	// exactly.
	c.nowFn = func() time.Time { return time.Unix(1_750_000_000, 0) }
	c.idemFn = func() string { return "idem_test_constant" }
	return c
}

// computeExpectedSig replicates the SDK's signing exactly. Tests use
// it to compare the Authorization header against what the backend
// verifier would compute.
func computeExpectedSig(method, pathToSign, ts string, body []byte, idem, secret string) string {
	bh := sha256.Sum256(body)
	str := strings.ToUpper(method) + "\n" + pathToSign + "\n" + ts + "\n" + hex.EncodeToString(bh[:])
	if idem != "" {
		str += "\n" + idem
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(str))
	return hex.EncodeToString(mac.Sum(nil))
}

// ─── Tests ─────────────────────────────────────────────────────────────

func TestNewClient_RequiresCredentials(t *testing.T) {
	t.Setenv("RIPLLO_KEY_ID", "")
	t.Setenv("RIPLLO_SECRET", "")
	if _, err := NewClient(ClientOptions{}); err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestNewClient_EnvFallback(t *testing.T) {
	t.Setenv("RIPLLO_KEY_ID", "AKIARPLO_env")
	t.Setenv("RIPLLO_SECRET", "secret_env")
	t.Setenv("RIPLLO_BASE_URL", "https://example.test")
	c, err := NewClient(ClientOptions{})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if c.keyId != "AKIARPLO_env" || c.secret != "secret_env" {
		t.Fatalf("env fallback didn't apply: %+v", c)
	}
	if c.baseURL != "https://example.test" {
		t.Fatalf("baseURL: %q", c.baseURL)
	}
}

func TestNewClient_DefaultBaseURL(t *testing.T) {
	t.Setenv("RIPLLO_BASE_URL", "")
	c, err := NewClient(ClientOptions{KeyID: "k", Secret: "s"})
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != DefaultBaseURL {
		t.Fatalf("baseURL = %q, want %q", c.baseURL, DefaultBaseURL)
	}
}

func TestNewClient_TrimsTrailingSlash(t *testing.T) {
	c, err := NewClient(ClientOptions{KeyID: "k", Secret: "s", BaseURL: "https://x.test/"})
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != "https://x.test" {
		t.Fatalf("baseURL = %q", c.baseURL)
	}
}

func TestNewClient_MountsAllResources(t *testing.T) {
	c, err := NewClient(ClientOptions{KeyID: "k", Secret: "s"})
	if err != nil {
		t.Fatal(err)
	}
	// Spot-check a handful of resources; if any is nil we know
	// mountResources missed it.
	if c.DiscountCodes == nil || c.Pixels == nil || c.Blog == nil ||
		c.Referrals == nil || c.Admin == nil || c.Broadcasts == nil ||
		c.MarketingCampaigns == nil ||
		c.Funnels == nil || c.AudienceSegments == nil || c.KYC == nil ||
		c.Marketplace == nil {
		t.Fatal("not all resources mounted")
	}
	// MarketingCampaigns is now a SEPARATE resource (the hub) from
	// Broadcasts. They share the *Client but are distinct namespaces
	// pointing at different URLs (/api/v1/marketing-campaigns vs
	// /api/v1/broadcasts). The pre-v0.5 alias was removed when the hub
	// took over the legacy URL.
	if c.MarketingCampaigns == nil || c.Broadcasts == nil {
		t.Fatal("Broadcasts + MarketingCampaigns must both be mounted")
	}
}

func TestForMerchant_ClonesAndScopes(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	scoped := c.ForMerchant("acc_test_merchant")

	// Each is a distinct *Client.
	if c == scoped {
		t.Fatal("ForMerchant returned same pointer")
	}
	if c.defaultOnBehalfOf != "" {
		t.Fatalf("original client mutated: %q", c.defaultOnBehalfOf)
	}
	if scoped.defaultOnBehalfOf != "acc_test_merchant" {
		t.Fatalf("scoped clone wrong: %q", scoped.defaultOnBehalfOf)
	}
	// Resources on the scoped clone fire with the on-behalf header.
	_, _ = scoped.DiscountCodes.List(context.Background(), DiscountCodesListParams{})
	if len(ts.requests) != 1 {
		t.Fatalf("requests = %d", len(ts.requests))
	}
	if got := ts.requests[0].OnBehalfOf; got != "acc_test_merchant" {
		t.Fatalf("X-Ripllo-On-Behalf-Of header = %q", got)
	}
}

func TestSign_PathWithoutQuerystring(t *testing.T) {
	// This is the specific bug guarded by the Node SDK 0.2.x + Python
	// SDK 0.1.x fix: signing the full path causes a mismatch on every
	// request with query params. Verify the SDK signs ONLY the path
	// portion.
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	limit := 5
	_, _ = c.DiscountCodes.List(context.Background(), DiscountCodesListParams{Limit: &limit})
	if len(ts.requests) != 1 {
		t.Fatalf("requests = %d", len(ts.requests))
	}
	req := ts.requests[0]
	// The request URL must include the querystring.
	if !strings.Contains(req.Path, "?limit=5") {
		t.Fatalf("expected ?limit=5 in path, got %q", req.Path)
	}
	// But the signature must be computed against the path WITHOUT it.
	auth := req.Headers.Get("Authorization")
	tsHdr := req.Headers.Get("X-Ripllo-Timestamp")
	want := computeExpectedSig("GET", "/api/v1/discount-codes", tsHdr, nil, "", "secret_test_xyz")
	if !strings.Contains(auth, "signature="+want) {
		t.Fatalf("path-without-querystring signature mismatch.\n got %s\nwant signature=%s", auth, want)
	}
}

func TestSign_HeaderFormat(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.Pixels.Get(context.Background())
	if len(ts.requests) != 1 {
		t.Fatalf("requests = %d", len(ts.requests))
	}
	auth := ts.requests[0].Headers.Get("Authorization")
	// Strict format check — backend regexes on this shape.
	if !strings.HasPrefix(auth, "Ripllo-HMAC-SHA256 keyId=AKIARPLO_test, scope=*, signature=") {
		t.Fatalf("Authorization wrong shape: %q", auth)
	}
	// Timestamp header is the deterministic nowFn value.
	if got := ts.requests[0].Headers.Get("X-Ripllo-Timestamp"); got != "1750000000" {
		t.Fatalf("X-Ripllo-Timestamp = %q", got)
	}
}

func TestSign_BodyHashAndIdempotency(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	input := DiscountCodeCreateInput{Code: "SAVE10", Type: DiscountTypePercent, Value: 10, Currency: "USD"}
	_, _ = c.DiscountCodes.Create(context.Background(), input)
	if len(ts.requests) != 1 {
		t.Fatalf("requests = %d", len(ts.requests))
	}
	req := ts.requests[0]
	if req.Method != "POST" || req.PathOnly != "/api/v1/discount-codes" {
		t.Fatalf("method/path: %s %s", req.Method, req.PathOnly)
	}
	if req.Headers.Get("Idempotency-Key") != "idem_test_constant" {
		t.Fatalf("Idempotency-Key header = %q", req.Headers.Get("Idempotency-Key"))
	}
	if req.Headers.Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type = %q", req.Headers.Get("Content-Type"))
	}
	// Verify signature mirrors body-hash + idem.
	want := computeExpectedSig("POST", "/api/v1/discount-codes", req.Headers.Get("X-Ripllo-Timestamp"), req.Body, "idem_test_constant", "secret_test_xyz")
	if !strings.Contains(req.Headers.Get("Authorization"), "signature="+want) {
		t.Fatalf("signature mismatch for create-with-idem; auth=%s want=%s", req.Headers.Get("Authorization"), want)
	}
}

func TestSign_EmptyBodyMatchesEmptyHash(t *testing.T) {
	// `sha256("")` is well-known. The signing path must use the same
	// empty-string sentinel as the JS/Python SDKs.
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.Referrals.GetProgram(context.Background())
	req := ts.requests[0]
	emptySHA := sha256.Sum256(nil)
	expectedHashHex := hex.EncodeToString(emptySHA[:])
	// Sanity: the sha256 of empty body is e3b0...b855.
	if expectedHashHex != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Fatalf("sha256 of empty bytes wrong: %s", expectedHashHex)
	}
	// And the request signature decodes cleanly against that hash.
	want := computeExpectedSig("GET", "/api/v1/referrals/program", req.Headers.Get("X-Ripllo-Timestamp"), nil, "", "secret_test_xyz")
	if !strings.Contains(req.Headers.Get("Authorization"), "signature="+want) {
		t.Fatalf("empty-body signature mismatch")
	}
}

func TestOnBehalfOf_PerRequestOverride(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL).ForMerchant("acc_default")
	// Direct Do() call with explicit override beats the default.
	_ = c.Do(context.Background(), RequestOptions{Method: "GET", Path: "/api/v1/pixels", OnBehalfOf: "acc_override"}, nil)
	if got := ts.requests[0].OnBehalfOf; got != "acc_override" {
		t.Fatalf("override OnBehalfOf = %q", got)
	}
}

// ─── Resource round-trips ──────────────────────────────────────────────

func TestDiscountCodes_GetRoundTrip(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`{"id":"dc_1","accountId":"acc_x","code":"SAVE10","type":"percent","value":10,"currency":"USD","scope":"cart","productIds":[],"tagFilter":[],"active":true,"public":false,"redemptionCount":3,"createdAt":"t","updatedAt":"t"}`), nil, "req_123")
	}
	c := mustClient(t, ts.server.URL)
	dc, err := c.DiscountCodes.Get(context.Background(), "dc_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if dc.Code != "SAVE10" || dc.Type != DiscountTypePercent || dc.Value != 10 {
		t.Fatalf("decoded shape wrong: %+v", dc)
	}
	if ts.requests[0].PathOnly != "/api/v1/discount-codes/dc_1" {
		t.Fatalf("path = %q", ts.requests[0].PathOnly)
	}
}

func TestDiscountCodes_Redeem(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`{"id":"red_1","created":true}`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	res, err := c.DiscountCodes.Redeem(context.Background(), RedeemInput{
		AccountID:         "acc_x",
		DiscountCodeID:    "dc_1",
		CheckoutSessionID: "ses_1",
		AppliedAmount:     500,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Created || res.ID != "red_1" {
		t.Fatalf("redeem result %+v", res)
	}
}

func TestPixels_PublicReturnsNullSafely(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`null`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	pp, err := c.Pixels.Public(context.Background(), "acc_abc")
	if err != nil {
		t.Fatalf("Public: %v", err)
	}
	// Null envelope → zero-value struct (Enabled=false, all nil ptrs).
	if pp.Enabled {
		t.Fatalf("expected zero-value PublicPixels, got %+v", pp)
	}
}

func TestFeeds_GoogleFeedURL(t *testing.T) {
	c := mustClient(t, "https://ripllo.test")
	got := c.Feeds.GoogleFeedURL("acc_xyz")
	want := "https://ripllo.test/api/v1/feeds/google/acc_xyz.xml"
	if got != want {
		t.Fatalf("GoogleFeedURL = %q, want %q", got, want)
	}
}

func TestBlog_CreateSendsIdempotency(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`{"post":{"id":"bp_1","accountId":"a","slug":"hello","title":"Hi","body":"x","status":"draft","tags":[],"createdAt":"t","updatedAt":"t"}}`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	_, err := c.Blog.Create(context.Background(), BlogPostInput{Slug: "hello", Title: "Hi", Body: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if got := ts.requests[0].Headers.Get("Idempotency-Key"); got != "idem_test_constant" {
		t.Fatalf("Idempotency-Key = %q", got)
	}
}

func TestAbandonedCart_StatsQueryParams(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	d := 30
	_, _ = c.AbandonedCart.Stats(context.Background(), StatsParams{WindowDays: &d})
	if !strings.Contains(ts.requests[0].Path, "windowDays=30") {
		t.Fatalf("querystring missing windowDays: %q", ts.requests[0].Path)
	}
}

func TestReferrals_PutProgramRoundTrip(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		// Echo back as program.
		writeEnvelope(w, 200, json.RawMessage(`{"id":"rp_1","accountId":"acc_x","enabled":true,"rewardType":"percent","referrerValue":10,"refereeValue":5,"currency":"USD","rewardExpiryDays":30,"attributionWindowDays":7,"createdAt":"t","updatedAt":"t"}`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	prog, err := c.Referrals.PutProgram(context.Background(), ReferralProgramInput{
		Enabled:       BoolPtr(true),
		RewardType:    DiscountTypePercent,
		ReferrerValue: 10,
		RefereeValue:  5,
		Currency:      "USD",
	})
	if err != nil {
		t.Fatal(err)
	}
	if prog.ID != "rp_1" || !prog.Enabled {
		t.Fatalf("program: %+v", prog)
	}
	if ts.requests[0].Method != "PUT" {
		t.Fatalf("method = %s", ts.requests[0].Method)
	}
}

func TestCampaigns_NestedRoute(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.Campaigns.AcceptApplication(context.Background(), "cmp_1", "app_2")
	if ts.requests[0].PathOnly != "/api/v1/campaigns/cmp_1/applications/app_2/accept" {
		t.Fatalf("path = %q", ts.requests[0].PathOnly)
	}
	if ts.requests[0].Method != "POST" {
		t.Fatalf("method = %s", ts.requests[0].Method)
	}
}

func TestPrograms_ShortAliasSameAsLong(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.Programs.Approve(context.Background(), "prg_1", "enr_2")
	if ts.requests[0].PathOnly != "/api/v1/programs/prg_1/enrollments/enr_2/approve" {
		t.Fatalf("short-alias path = %q", ts.requests[0].PathOnly)
	}
}

func TestCollaborations_RejectDeliverable_ReasonBody(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.Collaborations.RejectDeliverable(context.Background(), "col_1", "del_2", "needs work")
	var body map[string]string
	_ = json.Unmarshal(ts.requests[0].Body, &body)
	if body["reason"] != "needs work" {
		t.Fatalf("reason body = %v", body)
	}
}

func TestChannels_DNSAliasMatchesDNSRecords(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.Channels.DNS(context.Background(), "ch_1")
	if ts.requests[0].PathOnly != "/api/v1/channels/ch_1/dns-records" {
		t.Fatalf("DNS alias path = %q", ts.requests[0].PathOnly)
	}
}

func TestContacts_ListWithQueryParams(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	lim := 10
	_, _ = c.Contacts.List(context.Background(), ContactsListParams{Limit: &lim, Search: StrPtr("john")})
	if !strings.Contains(ts.requests[0].Path, "limit=10") || !strings.Contains(ts.requests[0].Path, "search=john") {
		t.Fatalf("contacts query: %q", ts.requests[0].Path)
	}
}

func TestBroadcasts_SendUsesIdempotency(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.Broadcasts.Send(context.Background(), "mc_1", nil)
	if got := ts.requests[0].Headers.Get("Idempotency-Key"); got != "idem_test_constant" {
		t.Fatalf("Idempotency-Key = %q", got)
	}
	if want := "/api/v1/broadcasts/mc_1/send"; ts.requests[0].PathOnly != want {
		t.Fatalf("path = %q, want %q", ts.requests[0].PathOnly, want)
	}
}

func TestBroadcasts_ListHitsBroadcastsURL(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`[]`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	_, err := c.Broadcasts.List(context.Background())
	if err != nil {
		t.Fatalf("Broadcasts.List: %v", err)
	}
	if got, want := ts.requests[0].PathOnly, "/api/v1/broadcasts"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}

func TestBroadcasts_CompileTemplateHitsBroadcastsURL(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.Broadcasts.CompileTemplate(context.Background(), JSON{"template": "X"})
	if got, want := ts.requests[0].PathOnly, "/api/v1/broadcasts/templates/compile"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}

// ─── Marketing campaigns (hub) ─────────────────────────────────────────

func TestMarketingCampaigns_ListHitsHubURL(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`{"campaigns":[]}`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	_, err := c.MarketingCampaigns.List(context.Background(), MarketingCampaignsListParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if got, want := ts.requests[0].PathOnly, "/api/v1/marketing-campaigns"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
	if ts.requests[0].Method != "GET" {
		t.Fatalf("method = %q", ts.requests[0].Method)
	}
}

func TestMarketingCampaigns_ListSerializesStatusArrayAsCSV(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`{"campaigns":[]}`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	_, _ = c.MarketingCampaigns.List(context.Background(), MarketingCampaignsListParams{
		Status: []MarketingCampaignStatus{MarketingCampaignStatusDraft, MarketingCampaignStatusLive},
	})
	if !strings.Contains(ts.requests[0].Path, "status=draft%2Clive") &&
		!strings.Contains(ts.requests[0].Path, "status=draft,live") {
		t.Fatalf("expected comma-joined status, got %q", ts.requests[0].Path)
	}
}

func TestMarketingCampaigns_GetByID(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.MarketingCampaigns.Get(context.Background(), "mc_abc")
	if got, want := ts.requests[0].PathOnly, "/api/v1/marketing-campaigns/mc_abc"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}

func TestMarketingCampaigns_GetFull(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.MarketingCampaigns.GetFull(context.Background(), "mc_abc")
	if got, want := ts.requests[0].PathOnly, "/api/v1/marketing-campaigns/mc_abc/full"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}

func TestMarketingCampaigns_Selector(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.MarketingCampaigns.Selector(context.Background())
	if got, want := ts.requests[0].PathOnly, "/api/v1/marketing-campaigns/_/selector"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}

func TestMarketingCampaigns_CreatePostsBodyAndIdempotency(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.MarketingCampaigns.Create(context.Background(), MarketingCampaignCreateInput{
		Name: "Q3 Push",
		Goal: MarketingCampaignGoalConversion,
	})
	if got, want := ts.requests[0].PathOnly, "/api/v1/marketing-campaigns"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
	if ts.requests[0].Method != "POST" {
		t.Fatalf("method = %q", ts.requests[0].Method)
	}
	if got := ts.requests[0].Headers.Get("Idempotency-Key"); got != "idem_test_constant" {
		t.Fatalf("Idempotency-Key = %q", got)
	}
	var body map[string]any
	_ = json.Unmarshal(ts.requests[0].Body, &body)
	if body["name"] != "Q3 Push" || body["goal"] != "conversion" {
		t.Fatalf("body = %v", body)
	}
}

func TestMarketingCampaigns_Update(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.MarketingCampaigns.Update(context.Background(), "mc_abc", MarketingCampaignUpdateInput{
		Status: MarketingCampaignStatusLive,
	})
	if got, want := ts.requests[0].PathOnly, "/api/v1/marketing-campaigns/mc_abc"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
	if ts.requests[0].Method != "PATCH" {
		t.Fatalf("method = %q", ts.requests[0].Method)
	}
}

func TestMarketingCampaigns_Delete(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.MarketingCampaigns.Delete(context.Background(), "mc_abc")
	if got, want := ts.requests[0].PathOnly, "/api/v1/marketing-campaigns/mc_abc"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
	if ts.requests[0].Method != "DELETE" {
		t.Fatalf("method = %q", ts.requests[0].Method)
	}
}

func TestFunnels_SetStepsWrapsBareArray(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, _ = c.Funnels.SetSteps(context.Background(), "fn_1", []any{
		map[string]any{"kind": "email"},
	})
	var body map[string]any
	_ = json.Unmarshal(ts.requests[0].Body, &body)
	if _, ok := body["steps"].([]any); !ok {
		t.Fatalf("expected steps key wrapping the array: %v", body)
	}
}

func TestAdmin_ProvisionWorkspaceRoundTrip(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`{"accountId":"acc_partner","partner":"storlaunch","discountRate":0.3,"brandName":null,"businessEmail":null,"createdAt":"t"}`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	ws, err := c.Admin.ProvisionWorkspace(context.Background(), ProvisionWorkspaceInput{
		AccountID:    "acc_partner",
		Partner:      "storlaunch",
		DiscountRate: 0.3,
	})
	if err != nil {
		t.Fatalf("provision: %v", err)
	}
	if ws.AccountID != "acc_partner" || ws.DiscountRate != 0.3 {
		t.Fatalf("workspace: %+v", ws)
	}
}

func TestAdmin_PartnerUsage_QueryParams(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`{"partner":"storlaunch","period":{"from":"2026-05-01","to":"2026-05-12"},"totals":{"redemptions":0,"referralRewards":0,"reminders":0,"chargeableCents":0},"byMerchant":[]}`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	_, err := c.Admin.PartnerUsage(context.Background(), PartnerUsageQuery{
		Partner: "storlaunch", From: "2026-05-01", To: "2026-05-12",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := ts.requests[0]
	if !strings.Contains(req.Path, "from=2026-05-01") || !strings.Contains(req.Path, "to=2026-05-12") || !strings.Contains(req.Path, "partner=storlaunch") {
		t.Fatalf("partner usage query: %q", req.Path)
	}
}

func TestPassthrough_Generic(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 200, json.RawMessage(`{"ok":true,"echoed":"hi"}`), nil, "")
	}
	c := mustClient(t, ts.server.URL)
	var out map[string]any
	if err := c.Passthrough(context.Background(), "POST", "/api/v1/anything", map[string]string{"x": "y"}, &out); err != nil {
		t.Fatal(err)
	}
	if out["ok"] != true {
		t.Fatalf("passthrough decode: %v", out)
	}
}

// ─── Error envelope ────────────────────────────────────────────────────

func TestError_EnvelopeSurfaces(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, 422, json.RawMessage(`null`),
			&envelopeError{Code: "invalid_input", Message: "code is required"},
			"req_abc",
		)
	}
	c := mustClient(t, ts.server.URL)
	_, err := c.DiscountCodes.Get(context.Background(), "dc_x")
	if err == nil {
		t.Fatal("expected error")
	}
	re, ok := err.(*Error)
	if !ok {
		t.Fatalf("not *ripllo.Error: %T", err)
	}
	if re.Status != 422 || re.Code != "invalid_input" || re.RequestID != "req_abc" {
		t.Fatalf("decoded error: %+v", re)
	}
	if !strings.Contains(re.Error(), "invalid_input") || !strings.Contains(re.Error(), "code is required") {
		t.Fatalf("Error() string: %q", re.Error())
	}
}

func TestError_NonJSON(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		w.WriteHeader(502)
		_, _ = w.Write([]byte("<html>upstream broken</html>"))
	}
	c := mustClient(t, ts.server.URL)
	_, err := c.Pixels.Get(context.Background())
	re, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if re.Status != 502 || re.Code != "invalid_response" {
		t.Fatalf("bad error: %+v", re)
	}
}

func TestQS_SkipsNils(t *testing.T) {
	got := qs(map[string]any{"a": "1", "b": nil, "c": 2})
	if strings.Contains(got, "b=") {
		t.Fatalf("nil value should be skipped: %q", got)
	}
	if !strings.Contains(got, "a=1") || !strings.Contains(got, "c=2") {
		t.Fatalf("missing keys: %q", got)
	}
}

func TestQS_EmptyReturnsEmpty(t *testing.T) {
	if got := qs(nil); got != "" {
		t.Fatalf("nil map -> %q", got)
	}
	if got := qs(map[string]any{}); got != "" {
		t.Fatalf("empty map -> %q", got)
	}
}

// ─── Webhook verification ──────────────────────────────────────────────

func TestVerifyWebhook_RoundTrip(t *testing.T) {
	secret := "whsec_test_xyz"
	body := []byte(`{"id":"evt_1","type":"ripllo.discount.redeemed.v1","occurredAt":"2026-05-12T00:00:00Z","accountId":"acc_x","data":{"discountCodeId":"dc_1"},"metadata":{}}`)
	ts := int64(1_750_000_000)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d.%s", ts, body)))
	header := fmt.Sprintf("t=%d,v1=%s", ts, hex.EncodeToString(mac.Sum(nil)))

	ev, err := VerifyWebhook(body, header, secret, &VerifyWebhookOptions{
		Now: func() time.Time { return time.Unix(ts+5, 0) },
	})
	if err != nil {
		t.Fatalf("VerifyWebhook: %v", err)
	}
	if ev.ID != "evt_1" || ev.Type != "ripllo.discount.redeemed.v1" {
		t.Fatalf("event decode: %+v", ev)
	}
}

func TestVerifyWebhook_BadSignature(t *testing.T) {
	secret := "s"
	body := []byte(`{"id":"evt"}`)
	ts := int64(1_750_000_000)
	header := fmt.Sprintf("t=%d,v1=%s", ts, strings.Repeat("ab", 32))
	_, err := VerifyWebhook(body, header, secret, &VerifyWebhookOptions{
		Now: func() time.Time { return time.Unix(ts, 0) },
	})
	if err == nil {
		t.Fatal("expected bad_signature error")
	}
	re, ok := err.(*Error)
	if !ok || re.Code != "bad_signature" {
		t.Fatalf("got %v", err)
	}
}

func TestVerifyWebhook_Replay(t *testing.T) {
	secret := "s"
	body := []byte(`{}`)
	ts := int64(1_750_000_000)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d.%s", ts, body)))
	header := fmt.Sprintf("t=%d,v1=%s", ts, hex.EncodeToString(mac.Sum(nil)))
	// 10 minutes later — outside the default 5-minute tolerance.
	_, err := VerifyWebhook(body, header, secret, &VerifyWebhookOptions{
		Now: func() time.Time { return time.Unix(ts+600, 0) },
	})
	if err == nil {
		t.Fatal("expected replay error")
	}
	re, _ := err.(*Error)
	if re == nil || re.Code != "signature_expired" {
		t.Fatalf("got %v", err)
	}
}

func TestVerifyWebhook_Malformed(t *testing.T) {
	if _, err := VerifyWebhook([]byte("{}"), "", "s", nil); err == nil {
		t.Fatal("missing header should error")
	}
	if _, err := VerifyWebhook([]byte("{}"), "garbage", "s", nil); err == nil {
		t.Fatal("garbage header should error")
	}
	if _, err := VerifyWebhook([]byte("{}"), "t=abc,v1=deadbeef", "s", nil); err == nil {
		t.Fatal("non-numeric t should error")
	}
}

func TestVerifyWebhook_TamperedBody(t *testing.T) {
	secret := "s"
	ts := int64(1_750_000_000)
	good := []byte(`{"id":"a"}`)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(ts, 10) + "." + string(good)))
	header := fmt.Sprintf("t=%d,v1=%s", ts, hex.EncodeToString(mac.Sum(nil)))
	tampered := []byte(`{"id":"b"}`)
	_, err := VerifyWebhook(tampered, header, secret, &VerifyWebhookOptions{
		Now: func() time.Time { return time.Unix(ts, 0) },
	})
	if err == nil {
		t.Fatal("tampered body should fail")
	}
}

// ─── Misc ──────────────────────────────────────────────────────────────

func TestDefaultIdemFn_ProducesPrefixed(t *testing.T) {
	id := defaultIdemFn()
	if !strings.HasPrefix(id, "idem_") || len(id) < 10 {
		t.Fatalf("defaultIdemFn output: %q", id)
	}
}

func TestAsJSON_RoundTrip(t *testing.T) {
	type X struct {
		A int    `json:"a"`
		B string `json:"b"`
	}
	got, err := AsJSON(X{A: 1, B: "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if got["a"].(float64) != 1 || got["b"].(string) != "hi" {
		t.Fatalf("AsJSON: %v", got)
	}
}
