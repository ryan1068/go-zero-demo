package handler

import (
	"net"
	"net/http"
	"strings"
)

// clientIP extracts the caller IP for rate limiting.
//
// nginx always overwrites X-Real-IP with the real peer, so it is trusted first.
// X-Forwarded-For is checked from the right because the leftmost entries are
// supplied by the client and can be spoofed; the last valid hop is the one
// appended by the trusted proxy.
func clientIP(r *http.Request) string {
	if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); net.ParseIP(ip) != nil {
		return ip
	}

	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		hops := strings.Split(forwarded, ",")
		for i := len(hops) - 1; i >= 0; i-- {
			if ip := strings.TrimSpace(hops[i]); net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	addr := strings.TrimSpace(r.RemoteAddr)
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}
