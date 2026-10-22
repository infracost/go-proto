package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const jwtKey = "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0.sig"

func roundTrip(t *testing.T, at *AuthTransport) http.Header {
	t.Helper()
	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
	}))
	t.Cleanup(srv.Close)

	resp, err := (&http.Client{Transport: at}).Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	return got
}

// absent marks a header that must not be sent at all.
const absent = "<absent>"

func assertHeaders(t *testing.T, got http.Header, want map[string]string) {
	t.Helper()
	for k, v := range want {
		vals, ok := got[http.CanonicalHeaderKey(k)]
		switch {
		case v == absent && ok:
			t.Errorf("%s: want absent, got %q", k, vals)
		case v != absent && !ok:
			t.Errorf("%s: missing, want %q", k, v)
		case v != absent && got.Get(k) != v:
			t.Errorf("%s: got %q, want %q", k, got.Get(k), v)
		}
	}
}

func TestAuthTransportHeaders(t *testing.T) {
	tests := []struct {
		name string
		at   AuthTransport
		want map[string]string
	}{
		{
			name: "plain api key",
			at:   AuthTransport{APIKey: "key", OrgID: "org"},
			want: map[string]string{
				"X-Api-Key":                      "key",
				"Authorization":                  absent,
				"User-Agent":                     "infracost-runner",
				"X-Infracost-Org-Id":             absent,
				"X-Infracost-Trace-ID":           absent,
				"X-Infracost-Replay-Delivery-Id": absent,
			},
		},
		{
			name: "service key",
			at:   AuthTransport{APIKey: "ics_123", TraceID: "trace"},
			want: map[string]string{
				"Authorization":        "Bearer ics_123",
				"X-Api-Key":            absent,
				"User-Agent":           "infracost-runner",
				"X-Infracost-Trace-ID": "trace",
			},
		},
		{
			name: "jwt with org id",
			at:   AuthTransport{APIKey: jwtKey, OrgID: "org"},
			want: map[string]string{
				"Authorization":      "Bearer " + jwtKey,
				"X-Api-Key":          absent,
				"User-Agent":         "infracost-cli",
				"X-Infracost-Org-Id": "org",
			},
		},
		{
			name: "jwt without org id sends an empty org header",
			at:   AuthTransport{APIKey: jwtKey},
			want: map[string]string{
				"Authorization":      "Bearer " + jwtKey,
				"User-Agent":         "infracost-cli",
				"X-Infracost-Org-Id": "",
			},
		},
		{
			name: "replay id set",
			at:   AuthTransport{APIKey: "key", ReplayID: "delivery"},
			want: map[string]string{"X-Infracost-Replay-Delivery-Id": "delivery"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertHeaders(t, roundTrip(t, &tt.at), tt.want)
		})
	}
}

type recordingTransport struct {
	req *http.Request
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r.req = req
	return httptest.NewRecorder().Result(), nil
}

func TestAuthTransportDelegatesToBase(t *testing.T) {
	base := &recordingTransport{}
	at := &AuthTransport{APIKey: "key", Base: base}

	req := httptest.NewRequest(http.MethodGet, "http://example.invalid/", nil)
	resp, err := at.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	if base.req == nil {
		t.Fatal("base transport not called")
	}
	if got := base.req.Header.Get("X-Api-Key"); got != "key" {
		t.Errorf("base saw X-Api-Key %q, want %q", got, "key")
	}
	if got := req.Header.Get("X-Api-Key"); got != "" {
		t.Errorf("caller's request was modified: X-Api-Key %q", got)
	}
}
