package ripllo

import "fmt"

// Error is the single error type returned by every operation in this
// package. Mirrors the Node `RiplloError` and Python `RiplloError`
// shape so callers can branch on .Code without matching strings.
//
// Status is the HTTP status code, or 0 for transport-level errors
// (network, timeout, non-JSON response received before the request
// reached the API). Code is the machine-readable error code — for
// server-side errors this is `envelope.error.code`; for client-side
// failures one of `timeout` / `network_error` / `invalid_response` /
// `unknown`. RequestID is the `meta.requestId` field of the
// envelope, when present.
type Error struct {
	Status    int
	Code      string
	Message   string
	RequestID string
}

func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("ripllo: %s: %s (status=%d, requestId=%s)", e.Code, e.Message, e.Status, e.RequestID)
	}
	return fmt.Sprintf("ripllo: %s: %s (status=%d)", e.Code, e.Message, e.Status)
}

func newErr(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}
