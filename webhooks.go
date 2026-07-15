package ripllo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// VerifyWebhookOptions tunes the webhook verifier. Pass nil to use the
// defaults (300 s tolerance, time.Now clock).
type VerifyWebhookOptions struct {
	// ToleranceSeconds rejects signatures with a timestamp older than this
	// many seconds. Zero or negative means "use the default" (300).
	ToleranceSeconds int64
	// Now is the clock used for the freshness check. Nil means time.Now.
	// Useful in tests.
	Now func() time.Time
}

// VerifyWebhook verifies an inbound Ripllo webhook signature and
// returns the parsed event envelope.
//
// Header: `Ripllo-Signature: t=<unix>,v1=<hex>`
// Signed payload: `${t}.${rawBody}` with HMAC-SHA256 keyed on the
// endpoint's shared secret.
//
// IMPORTANT: pass the raw bytes you received over the wire. If you let
// a framework parse JSON first and re-stringify, whitespace drift will
// break the signature.
//
// Mirrors Node `verifyWebhook` + Python `verify_webhook`.
func VerifyWebhook(rawBody []byte, signatureHeader, secret string, opts *VerifyWebhookOptions) (*WebhookEvent, error) {
	if signatureHeader == "" {
		return nil, newErr(0, "missing_signature", "missing Ripllo-Signature header")
	}
	if secret == "" {
		return nil, newErr(0, "missing_secret", "secret is required")
	}

	tolerance := int64(300)
	nowFn := time.Now
	if opts != nil {
		if opts.ToleranceSeconds > 0 {
			tolerance = opts.ToleranceSeconds
		}
		if opts.Now != nil {
			nowFn = opts.Now
		}
	}

	timestamp, v1, ok := parseRiplloSignature(signatureHeader)
	if !ok {
		return nil, newErr(0, "malformed_signature", "malformed signature header")
	}

	nowTs := nowFn().Unix()
	drift := nowTs - timestamp
	if drift < 0 {
		drift = -drift
	}
	if drift > tolerance {
		return nil, newErr(0, "signature_expired", "signature timestamp "+strconv.FormatInt(drift, 10)+"s out of tolerance")
	}

	payload := strconv.FormatInt(timestamp, 10) + "." + string(rawBody)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedHex := hex.EncodeToString(mac.Sum(nil))

	expected, err := hex.DecodeString(expectedHex)
	if err != nil {
		return nil, newErr(0, "bad_signature", err.Error())
	}
	actual, err := hex.DecodeString(v1)
	if err != nil {
		return nil, newErr(0, "bad_signature", "v1 is not hex")
	}
	if !hmac.Equal(expected, actual) {
		return nil, newErr(0, "bad_signature", "bad signature")
	}

	var ev WebhookEvent
	if err := json.Unmarshal(rawBody, &ev); err != nil {
		return nil, newErr(0, "invalid_body", "webhook body is not valid JSON: "+err.Error())
	}
	return &ev, nil
}

func parseRiplloSignature(header string) (timestamp int64, v1 string, ok bool) {
	parts := strings.Split(header, ",")
	var tSeen, v1Seen bool
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		k, v, found := strings.Cut(p, "=")
		if !found {
			continue
		}
		switch strings.TrimSpace(k) {
		case "t":
			if ts, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
				timestamp = ts
				tSeen = true
			}
		case "v1":
			v1 = strings.TrimSpace(v)
			v1Seen = true
		}
	}
	return timestamp, v1, tSeen && v1Seen
}
