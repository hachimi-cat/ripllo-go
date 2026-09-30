// Package ripllo is the official Go SDK for Ripllo.
//
// It mirrors the developer-facing surface of `@forjio/ripllo-node`
// 0.2.2 and `ripllo` (PyPI) 0.1.x: HMAC-SHA256 partner-billing auth
// (Pattern 2), resource namespaces for every route group, generic
// `Passthrough` for partners that need to relay arbitrary
// merchant-portal requests, and inbound webhook verification.
//
// Auth scheme
//
//	Authorization: Ripllo-HMAC-SHA256 keyId=<keyId>, scope=*, signature=<hex>
//	X-Ripllo-Timestamp: <unix_seconds>
//
// where
//
//	signature = HMAC-SHA256(secret, "{METHOD}\n{PATH}\n{TIMESTAMP}\n{BODY_SHA256}")
//
// with an optional trailing "\n<idempotencyKey>" segment when the
// request carries an Idempotency-Key header.
//
// IMPORTANT: the path is signed WITHOUT the query string. The backend
// verifier strips the query before reconstructing the string-to-sign;
// signing the full path causes a mismatch on every request that
// carries query params (`?limit=`, `?status=`, etc.). This is the
// same fix as Node SDK 0.2.x and Python SDK 0.1.x.
package ripllo

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// DefaultBaseURL is the production Ripllo API base.
const DefaultBaseURL = "https://ripllo.com"

// Client is the high-level SDK surface. Construct via NewClient.
//
// Concurrency-safe: every method is safe for use from multiple
// goroutines. ForMerchant returns a shallow clone that shares the
// underlying *http.Client.
type Client struct {
	keyId             string
	secret            string
	baseURL           string
	defaultOnBehalfOf string
	timeout           time.Duration
	http              *http.Client
	nowFn             func() time.Time
	idemFn            func() string

	// Resource namespaces — every API surface is grouped here, matching
	// the Node SDK's nested object shape.
	DiscountCodes      *DiscountCodesResource
	Pixels             *PixelsResource
	Feeds              *FeedsResource
	Blog               *BlogResource
	AbandonedCart      *AbandonedCartResource
	Referrals          *ReferralsResource
	APIKeys            *APIKeysResource
	Webhooks           *WebhooksResource
	AuditLog           *AuditLogResource
	Integrations       *IntegrationsResource
	Billing            *BillingResource
	Uploads            *UploadsResource
	Campaigns          *CampaignsResource
	Programs           *ProgramsResource
	Collaborations     *CollaborationsResource
	Insights           *InsightsResource
	Channels           *ChannelsResource
	Contacts           *ContactsResource
	ContactLists       *ContactListsResource
	Broadcasts *BroadcastsResource
	// MarketingCampaigns is the top-level marketing-campaign hub
	// (/api/v1/marketing-campaigns). It is now a SEPARATE resource
	// from Broadcasts — pre-v0.5 it was a deprecated alias of
	// Broadcasts (when /api/v1/marketing-campaigns served email
	// blasts). The hub took over that URL in PR #11; the alias was
	// removed in PR #13.
	MarketingCampaigns *MarketingCampaignsResource
	Funnels            *FunnelsResource
	Inbox              *InboxResource
	AudienceSegments   *AudienceSegmentsResource
	MerchantProfile    *MerchantProfileResource
	CreatorProfile    *CreatorProfileResource
	AffiliatorProfile *AffiliatorProfileResource
	Marketplace       *MarketplaceResource
	Affiliates        *AffiliatesResource
	KYC               *KYCResource
	CreatorStats      *CreatorStatsResource
	Admin             *AdminResource

	// API has every feature route, one method each (generated from the API
	// spec: api_generated.go), signed like every other call.
	API *GeneratedAPI
}

// ClientOptions matches the env-var defaults of the Node + Python
// SDKs. KeyID + Secret default to RIPLLO_KEY_ID + RIPLLO_SECRET when
// blank. BaseURL defaults to RIPLLO_BASE_URL, then DefaultBaseURL.
type ClientOptions struct {
	KeyID       string        // defaults to RIPLLO_KEY_ID
	Secret      string        // defaults to RIPLLO_SECRET
	BaseURL     string        // defaults to RIPLLO_BASE_URL, then DefaultBaseURL
	OnBehalfOf  string        // optional merchant accountId for X-Ripllo-On-Behalf-Of
	Timeout     time.Duration // per-request timeout, default 30s
	HTTP        *http.Client  // custom client; defaults to http.DefaultClient
}

// NewClient constructs a Client. KeyID + Secret are required (env
// fallback is checked first).
func NewClient(opts ClientOptions) (*Client, error) {
	if opts.KeyID == "" {
		opts.KeyID = os.Getenv("RIPLLO_KEY_ID")
	}
	if opts.Secret == "" {
		opts.Secret = os.Getenv("RIPLLO_SECRET")
	}
	if opts.BaseURL == "" {
		opts.BaseURL = os.Getenv("RIPLLO_BASE_URL")
	}
	if opts.BaseURL == "" {
		opts.BaseURL = DefaultBaseURL
	}
	if opts.KeyID == "" || opts.Secret == "" {
		return nil, newErr(0, "MISSING_CREDENTIALS", "RiplloClient: KeyID and Secret are required")
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.HTTP == nil {
		opts.HTTP = &http.Client{Timeout: opts.Timeout}
	}

	c := &Client{
		keyId:             opts.KeyID,
		secret:            opts.Secret,
		baseURL:           strings.TrimRight(opts.BaseURL, "/"),
		defaultOnBehalfOf: opts.OnBehalfOf,
		timeout:           opts.Timeout,
		http:              opts.HTTP,
		nowFn:             time.Now,
		idemFn:            defaultIdemFn,
	}
	c.mountResources()
	return c, nil
}

// ForMerchant returns a clone scoped to a specific merchant accountId,
// added as `X-Ripllo-On-Behalf-Of` on every request. Use with
// platform-admin keys.
func (c *Client) ForMerchant(accountID string) *Client {
	clone := &Client{
		keyId:             c.keyId,
		secret:            c.secret,
		baseURL:           c.baseURL,
		defaultOnBehalfOf: accountID,
		timeout:           c.timeout,
		http:              c.http,
		nowFn:             c.nowFn,
		idemFn:            c.idemFn,
	}
	clone.mountResources()
	return clone
}

// BaseURL returns the configured base URL. Useful for building
// public-facing URLs (e.g. `feeds.GoogleFeedURL`).
func (c *Client) BaseURL() string {
	return c.baseURL
}

// ─── HMAC signing ──────────────────────────────────────────────────────

type signed struct {
	signature string
	timestamp string
}

func (c *Client) sign(method, path string, body []byte, idempotencyKey string) signed {
	ts := strconv.FormatInt(c.nowFn().Unix(), 10)

	bh := sha256.Sum256(body) // sum is zero on nil body, matches sha256("") for empty string too
	bodyHash := hex.EncodeToString(bh[:])

	// Sign the path WITHOUT the query string — see package doc.
	pathToSign := path
	if i := strings.IndexByte(path, '?'); i >= 0 {
		pathToSign = path[:i]
	}

	stringToSign := strings.ToUpper(method) + "\n" + pathToSign + "\n" + ts + "\n" + bodyHash
	if idempotencyKey != "" {
		stringToSign += "\n" + idempotencyKey
	}

	mac := hmac.New(sha256.New, []byte(c.secret))
	mac.Write([]byte(stringToSign))
	return signed{signature: hex.EncodeToString(mac.Sum(nil)), timestamp: ts}
}

// ─── Low-level request ─────────────────────────────────────────────────

// RequestOptions configures a single Do call.
type RequestOptions struct {
	Method         string // GET / POST / PATCH / PUT / DELETE
	Path           string // full URL path including any querystring (e.g. "/api/v1/x?limit=5")
	Body           any    // optional; JSON-encoded if non-nil
	IdempotencyKey string // optional
	OnBehalfOf     string // overrides Client.defaultOnBehalfOf
}

// Do dispatches a signed request and decodes the `data` envelope into
// out. Pass out=nil to discard the response. Errors surface as *Error.
//
// out may also be a *json.RawMessage to capture the raw `data` field
// without decoding.
func (c *Client) Do(ctx context.Context, opts RequestOptions, out any) error {
	if opts.Method == "" {
		return newErr(0, "INVALID_REQUEST", "method is required")
	}
	if opts.Path == "" || !strings.HasPrefix(opts.Path, "/") {
		return newErr(0, "INVALID_REQUEST", "path must start with /")
	}

	var bodyBytes []byte
	if opts.Body != nil {
		b, err := json.Marshal(opts.Body)
		if err != nil {
			return newErr(0, "SERIALIZE_FAILED", err.Error())
		}
		// Treat marshal-of-nil-interface ("null") as no-body to match the
		// JS `body === undefined` semantics on the server side.
		if string(b) != "null" {
			bodyBytes = b
		}
	}

	sig := c.sign(opts.Method, opts.Path, bodyBytes, opts.IdempotencyKey)

	var reqBody io.Reader
	if bodyBytes != nil {
		reqBody = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(opts.Method), c.baseURL+opts.Path, reqBody)
	if err != nil {
		return newErr(0, "INVALID_REQUEST", err.Error())
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Ripllo-HMAC-SHA256 keyId=%s, scope=*, signature=%s", c.keyId, sig.signature))
	req.Header.Set("X-Ripllo-Timestamp", sig.timestamp)
	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if opts.IdempotencyKey != "" {
		req.Header.Set("Idempotency-Key", opts.IdempotencyKey)
	}
	obo := opts.OnBehalfOf
	if obo == "" {
		obo = c.defaultOnBehalfOf
	}
	if obo != "" {
		req.Header.Set("X-Ripllo-On-Behalf-Of", obo)
	}

	res, err := c.http.Do(req)
	if err != nil {
		if errStr := err.Error(); strings.Contains(errStr, "context deadline") || strings.Contains(errStr, "Client.Timeout") {
			return newErr(0, "timeout", fmt.Sprintf("Ripllo request timed out after %s", c.timeout))
		}
		return newErr(0, "network_error", err.Error())
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return newErr(res.StatusCode, "network_error", err.Error())
	}

	// Empty 204-style body — succeed with out untouched.
	if len(raw) == 0 {
		if res.StatusCode >= 400 {
			return newErr(res.StatusCode, "unknown", fmt.Sprintf("HTTP %d", res.StatusCode))
		}
		return nil
	}

	var envelope apiEnvelope
	if jerr := json.Unmarshal(raw, &envelope); jerr != nil {
		preview := string(raw)
		if len(preview) > 200 {
			preview = preview[:200]
		}
		return newErr(res.StatusCode, "invalid_response", "Non-JSON response: "+preview)
	}

	var requestID string
	if envelope.Meta != nil {
		requestID = envelope.Meta.RequestID
	}

	if res.StatusCode >= 400 || envelope.Error != nil {
		code := "unknown"
		msg := fmt.Sprintf("HTTP %d", res.StatusCode)
		if envelope.Error != nil {
			if envelope.Error.Code != "" {
				code = envelope.Error.Code
			}
			if envelope.Error.Message != "" {
				msg = envelope.Error.Message
			}
		}
		return &Error{Status: res.StatusCode, Code: code, Message: msg, RequestID: requestID}
	}

	if out == nil || len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(envelope.Data, out); err != nil {
		return newErr(res.StatusCode, "invalid_response", "failed to decode data: "+err.Error())
	}
	return nil
}

// Passthrough forwards an arbitrary merchant-portal request to Ripllo
// and decodes the resulting envelope `data` into out. For partners
// (storlaunch / fulkruma) that need to relay arbitrary requests
// without hand-writing a typed method per resource.
func (c *Client) Passthrough(ctx context.Context, method, path string, body, out any) error {
	return c.Do(ctx, RequestOptions{Method: method, Path: path, Body: body}, out)
}

// apigenRequest is the call behind Client.API (api_generated.go): signed
// like every other request, with an idempotency key on writes. An empty
// body is not sent: the server hashes an empty JSON body as ""
// (backend middleware/hmac-auth.ts) while Do would hash "{}", and the
// signatures would not agree.
func (c *Client) apigenRequest(ctx context.Context, method, path string, query url.Values, body map[string]any) (json.RawMessage, error) {
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	opts := RequestOptions{Method: method, Path: path}
	if len(body) > 0 {
		opts.Body = body
	}
	if strings.ToUpper(method) != http.MethodGet {
		opts.IdempotencyKey = c.idemFn()
	}
	var out json.RawMessage
	if err := c.Do(ctx, opts, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Helpers ───────────────────────────────────────────────────────────

func defaultIdemFn() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// extremely unlikely; fall back to time-based suffix
		return "idem_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	// RFC 4122 v4-ish formatting; we don't strictly need a canonical UUID.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("idem_%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// qs renders a querystring from a map, skipping nil values. Mirrors
// the Node SDK's `qs()` helper. Keys are sorted to keep tests stable.
func qs(params map[string]any) string {
	if len(params) == 0 {
		return ""
	}
	v := url.Values{}
	for k, val := range params {
		if val == nil {
			continue
		}
		v.Set(k, fmt.Sprintf("%v", val))
	}
	enc := v.Encode()
	if enc == "" {
		return ""
	}
	return "?" + enc
}

// Pretty pointer helpers for callers building input structs.
func StrPtr(s string) *string         { return &s }
func IntPtr(i int) *int               { return &i }
func BoolPtr(b bool) *bool            { return &b }
func Float64Ptr(f float64) *float64   { return &f }
func StatusPtr(s BlogPostStatus) *BlogPostStatus { return &s }
