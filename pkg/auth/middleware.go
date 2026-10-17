package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

// Middleware returns an HTTP middleware that requires a valid JWT on every path
// except the given skipPaths. Authenticated identity is stored in the request
// context and forwarded as gRPC metadata headers for the gateway.
func Middleware(secret string, skipPaths ...string) func(http.HandlerFunc) http.HandlerFunc {
	skip := make(map[string]struct{}, len(skipPaths))
	for _, path := range skipPaths {
		skip[path] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if _, ok := skip[r.URL.Path]; ok {
				// Never trust identity headers supplied by the client.
				r.Header.Del(metadataHeaderUserID)
				r.Header.Del(metadataHeaderUsername)
				next(w, r)
				return
			}

			token, err := tokenFromRequest(r)
			if err != nil {
				writeUnauthorized(w, err)
				return
			}
			claims, err := Parse(secret, token)
			if err != nil {
				writeUnauthorized(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, usernameKey, claims.Username)
			r = r.WithContext(ctx)

			r.Header.Set(metadataHeaderUserID, strconv.FormatInt(claims.UserID, 10))
			if claims.Username != "" {
				r.Header.Set(metadataHeaderUsername, claims.Username)
			}

			next(w, r)
		}
	}
}

func tokenFromRequest(r *http.Request) (string, error) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return "", ErrInvalidToken
	}
	if len(header) > 7 && strings.EqualFold(header[:7], "Bearer ") {
		return strings.TrimSpace(header[7:]), nil
	}
	return header, nil
}

func writeUnauthorized(w http.ResponseWriter, err error) {
	logx.Errorf("unauthorized request: %v", err)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if encodeErr := json.NewEncoder(w).Encode(map[string]any{
		"code": http.StatusUnauthorized,
		"msg":  "未授权",
	}); encodeErr != nil {
		logx.Errorf("failed to write unauthorized response: %v", encodeErr)
	}
}
