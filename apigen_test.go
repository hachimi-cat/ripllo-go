package ripllo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// Client.API (api_generated.go) goes through apigenRequest: the same signing,
// headers and envelope as every hand-written call.

func assertSigned(t *testing.T, req recordedRequest, wantIdem bool) {
	t.Helper()
	idem := req.Headers.Get("Idempotency-Key")
	if wantIdem != (idem != "") {
		t.Fatalf("Idempotency-Key = %q, want one: %v", idem, wantIdem)
	}
	// the server signs req.originalUrl as sent (still escaped), without the query
	pathToSign := strings.SplitN(req.Path, "?", 2)[0]
	want := computeExpectedSig(req.Method, pathToSign, "1750000000", req.Body, idem, "secret_test_xyz")
	auth := req.Headers.Get("Authorization")
	if auth != fmt.Sprintf("Ripllo-HMAC-SHA256 keyId=AKIARPLO_test, scope=*, signature=%s", want) {
		t.Fatalf("Authorization = %q, want signature %s", auth, want)
	}
	if req.Headers.Get("X-Ripllo-Timestamp") != "1750000000" {
		t.Fatalf("X-Ripllo-Timestamp = %q", req.Headers.Get("X-Ripllo-Timestamp"))
	}
}

func TestAPI_ListSendsQueryAndSignsPathOnly(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, http.StatusOK, json.RawMessage(`{"discountCodes":[{"id":"dc_1"}]}`), nil, "req_1")
	}
	c := mustClient(t, ts.server.URL)
	data, err := c.API.DiscountCodesList(context.Background(), &DiscountCodesListArgs{Active: true, Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"discountCodes":[{"id":"dc_1"}]}` {
		t.Fatalf("data = %s", data)
	}
	req := ts.requests[0]
	if req.Method != "GET" || req.Path != "/api/v1/discount-codes?active=true&limit=5" || len(req.Body) != 0 {
		t.Fatalf("request = %s %s %q", req.Method, req.Path, req.Body)
	}
	assertSigned(t, req, false)
}

func TestAPI_CreateSendsBodyWithIdempotencyKey(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, err := c.API.DiscountCodesCreate(context.Background(), &DiscountCodesCreateArgs{
		Code: "SPRING10", Type: "percent", Value: 10, Currency: "IDR",
		Active: Ptr(false),
		Body:   map[string]any{"description": "from Body", "code": "overridden"},
	})
	if err != nil {
		t.Fatal(err)
	}
	req := ts.requests[0]
	if req.Method != "POST" || req.Path != "/api/v1/discount-codes" {
		t.Fatalf("request = %s %s", req.Method, req.Path)
	}
	var got map[string]any
	if err := json.Unmarshal(req.Body, &got); err != nil {
		t.Fatal(err)
	}
	want := `{"active":false,"code":"SPRING10","currency":"IDR","description":"from Body","type":"percent","value":10}`
	if b, _ := json.Marshal(got); string(b) != want {
		t.Fatalf("body = %s, want %s", b, want)
	}
	if req.Headers.Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type = %q", req.Headers.Get("Content-Type"))
	}
	assertSigned(t, req, true)
}

func TestAPI_EmptyBodyIsNotSent(t *testing.T) {
	// The server hashes an empty body as "": sending "{}" would break the signature.
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	if _, err := c.API.DiscountCodesDelete(context.Background(), "dc/1", nil); err != nil {
		t.Fatal(err)
	}
	req := ts.requests[0]
	if req.Method != "DELETE" || req.Path != "/api/v1/discount-codes/dc%2F1" || len(req.Body) != 0 {
		t.Fatalf("request = %s %s %q", req.Method, req.Path, req.Body)
	}
	assertSigned(t, req, true)
}

func TestAPI_RequiredFieldMissing(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL)
	_, err := c.API.DiscountCodesCreate(context.Background(), &DiscountCodesCreateArgs{Code: "X", Type: "percent"})
	if err == nil || !strings.Contains(err.Error(), "Currency") {
		t.Fatalf("err = %v, want a missing Currency", err)
	}
	if len(ts.requests) != 0 {
		t.Fatalf("sent %d requests", len(ts.requests))
	}
}

func TestAPI_ErrorEnvelope(t *testing.T) {
	ts := newTestServer(t)
	ts.handler = func(w http.ResponseWriter, r *http.Request, body []byte) {
		writeEnvelope(w, http.StatusNotFound, nil, &envelopeError{Code: "NOT_FOUND", Message: "no such code"}, "req_2")
	}
	c := mustClient(t, ts.server.URL)
	_, err := c.API.DiscountCodesGet(context.Background(), "dc_missing")
	e, ok := err.(*Error)
	if !ok || e.Status != 404 || e.Code != "NOT_FOUND" || e.RequestID != "req_2" {
		t.Fatalf("err = %#v", err)
	}
}

func TestAPI_ForMerchantKeepsScope(t *testing.T) {
	ts := newTestServer(t)
	c := mustClient(t, ts.server.URL).ForMerchant("acc_m")
	if _, err := c.API.DiscountCodesGet(context.Background(), "dc_1"); err != nil {
		t.Fatal(err)
	}
	if ts.requests[0].OnBehalfOf != "acc_m" {
		t.Fatalf("X-Ripllo-On-Behalf-Of = %q", ts.requests[0].OnBehalfOf)
	}
}
