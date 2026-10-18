package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-demo/pkg/auth"
)

// BlacklistMiddleware rejects tokens that were revoked by logout.
type BlacklistMiddleware struct {
	redis *redis.Redis
}

func NewBlacklistMiddleware(client *redis.Redis) *BlacklistMiddleware {
	return &BlacklistMiddleware{redis: client}
}

func (m *BlacklistMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := auth.BearerToken(r)
		if token != "" && m.redis != nil {
			exists, err := m.redis.ExistsCtx(r.Context(), auth.BlacklistKey(token))
			if err != nil {
				log.Printf("check token blacklist failed: %v", err)
			} else if exists {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				if encodeErr := json.NewEncoder(w).Encode(map[string]any{
					"code": http.StatusUnauthorized,
					"msg":  "登录已失效",
				}); encodeErr != nil {
					log.Printf("failed to write blacklist response: %v", encodeErr)
				}
				return
			}
		}
		next(w, r)
	}
}
