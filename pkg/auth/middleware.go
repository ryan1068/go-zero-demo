package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Middleware returns an HTTP middleware that requires a valid JWT on every path
// except the given skipPaths. A skip path ending in "/" is treated as a prefix.
// Authenticated identity is stored in the request context and forwarded as gRPC
// metadata headers for the gateway.
func Middleware(secret string, skipPaths ...string) func(http.HandlerFunc) http.HandlerFunc {
	skipExact := make(map[string]struct{}, len(skipPaths))
	skipPrefixes := make([]string, 0, len(skipPaths))
	for _, path := range skipPaths {
		if strings.HasSuffix(path, "/") {
			skipPrefixes = append(skipPrefixes, path)
			continue
		}
		skipExact[path] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if shouldSkip(r.URL.Path, skipExact, skipPrefixes) {
				// Never trust identity headers supplied by the client.
				r.Header.Del(metadataHeaderUserID)
				r.Header.Del(metadataHeaderUsername)
				r.Header.Del(metadataHeaderUserRole)
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

			ctx := WithUserRole(r.Context(), claims.UserID, claims.Username, claims.Role)
			r = r.WithContext(ctx)

			r.Header.Set(metadataHeaderUserID, strconv.FormatInt(claims.UserID, 10))
			if claims.Username != "" {
				r.Header.Set(metadataHeaderUsername, claims.Username)
			}
			r.Header.Set(metadataHeaderUserRole, strconv.FormatInt(int64(claims.Role), 10))

			next(w, r)
		}
	}
}

func shouldSkip(path string, exact map[string]struct{}, prefixes []string) bool {
	if _, ok := exact[path]; ok {
		return true
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
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
	log.Printf("unauthorized request: %v", err)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if encodeErr := json.NewEncoder(w).Encode(map[string]any{
		"code": http.StatusUnauthorized,
		"msg":  "未授权",
	}); encodeErr != nil {
		log.Printf("failed to write unauthorized response: %v", encodeErr)
	}
}
