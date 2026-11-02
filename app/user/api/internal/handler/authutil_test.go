package handler

import (
	"net/http"
	"testing"
)

func TestClientIP_PrefersRealIPOverForwardedFor(t *testing.T) {
	req := mustRequest(t, "192.0.2.10:1234")
	req.Header.Set("X-Real-IP", "203.0.113.9")
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 203.0.113.9")

	if got := clientIP(req); got != "203.0.113.9" {
		t.Fatalf("clientIP() = %q, want %q", got, "203.0.113.9")
	}
}

func TestClientIP_UsesLastForwardedForHop(t *testing.T) {
	req := mustRequest(t, "192.0.2.10:1234")
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 203.0.113.9")

	if got := clientIP(req); got != "203.0.113.9" {
		t.Fatalf("clientIP() = %q, want the hop appended by the trusted proxy", got)
	}
}

func TestClientIP_IgnoresInvalidRealIP(t *testing.T) {
	req := mustRequest(t, "192.0.2.10:1234")
	req.Header.Set("X-Real-IP", "not-an-ip")
	req.Header.Set("X-Forwarded-For", "203.0.113.9")

	if got := clientIP(req); got != "203.0.113.9" {
		t.Fatalf("clientIP() = %q, want %q", got, "203.0.113.9")
	}
}

func TestClientIP_SkipsInvalidForwardedForHops(t *testing.T) {
	req := mustRequest(t, "192.0.2.10:1234")
	req.Header.Set("X-Forwarded-For", "spoofed, 203.0.113.9, garbage")

	if got := clientIP(req); got != "203.0.113.9" {
		t.Fatalf("clientIP() = %q, want %q", got, "203.0.113.9")
	}
}

func TestClientIP_FallsBackToRemoteAddr(t *testing.T) {
	req := mustRequest(t, "192.0.2.10:1234")

	if got := clientIP(req); got != "192.0.2.10" {
		t.Fatalf("clientIP() = %q, want %q", got, "192.0.2.10")
	}
}

func TestClientIP_FallsBackToRawRemoteAddr(t *testing.T) {
	req := mustRequest(t, "192.0.2.10")

	if got := clientIP(req); got != "192.0.2.10" {
		t.Fatalf("clientIP() = %q, want %q", got, "192.0.2.10")
	}
}

func mustRequest(t *testing.T, remoteAddr string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/user/login", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.RemoteAddr = remoteAddr
	return req
}
