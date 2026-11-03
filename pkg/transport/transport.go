// Package transport provides the http.RoundTripper that authenticates calls
// to the Infracost API.
package transport

import (
	"net/http"
	"strings"
)

// AuthTransport adds Infracost API auth and tracing headers to each request
// before passing it to Base.
//
// The key is sent as a bearer token when it is a JWT or an "ics" service key,
// and as X-Api-Key otherwise. A JWT is a user's CLI login, so it identifies as
// infracost-cli and always carries X-Infracost-Org-Id, empty or not, leaving
// the API to pick a default org. Other keys identify as infracost-runner.
// Empty trace and replay IDs send no header.
type AuthTransport struct {
	// APIKey authenticates the request.
	APIKey string
	// OrgID selects the organization for a JWT key; it is ignored for other
	// keys.
	OrgID string
	// TraceID is sent as X-Infracost-Trace-ID.
	TraceID string
	// ReplayID is sent as X-Infracost-Replay-Delivery-Id.
	ReplayID string
	// Base performs the request; nil means http.DefaultTransport.
	Base http.RoundTripper
}

// RoundTrip implements http.RoundTripper. It clones req rather than modifying
// it.
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	jwt := isJWT(t.APIKey)
	if jwt || strings.HasPrefix(t.APIKey, "ics") {
		req.Header.Set("Authorization", "Bearer "+t.APIKey)
	} else {
		req.Header.Set("X-Api-Key", t.APIKey)
	}
	if jwt {
		req.Header.Set("User-Agent", "infracost-cli")
		req.Header.Set("X-Infracost-Org-Id", t.OrgID)
	} else {
		req.Header.Set("User-Agent", "infracost-runner")
	}
	setIfNotEmpty(req.Header, "X-Infracost-Trace-ID", t.TraceID)
	setIfNotEmpty(req.Header, "X-Infracost-Replay-Delivery-Id", t.ReplayID)

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

func isJWT(key string) bool {
	return strings.HasPrefix(key, "eyJ") && strings.Count(key, ".") == 2
}

func setIfNotEmpty(h http.Header, key, value string) {
	if value != "" {
		h.Set(key, value)
	}
}
