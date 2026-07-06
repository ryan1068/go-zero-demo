package middleware

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gateway/internal/types"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取 token
		token := r.Header.Get("Authorization")
		if token == "" {
			httpx.Error(w, types.NewBizError(http.StatusUnauthorized, "Unauthorized"))
			return
		}

		// 验证 token (示例逻辑，实际应调用认证服务)
		if token != "valid_token" {
			httpx.Error(w, types.NewBizError(http.StatusForbidden, "invalid token"))
			return
		}

		// 将用户信息存入 context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", 123)
		ctx = context.WithValue(ctx, "username", "doubao")
		r = r.WithContext(ctx)

		// 继续处理请求
		next(w, r)
	}
}