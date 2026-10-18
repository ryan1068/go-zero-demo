package handler

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
)

// LogoutHandler revokes the current token until its natural expiry.
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := auth.BearerToken(r)
		if token == "" {
			writeUnauthorized(w, "未授权")
			return
		}
		claims, err := auth.Parse(svcCtx.Config.Jwt.Secret, token)
		if err != nil {
			writeUnauthorized(w, "未授权")
			return
		}
		if claims.ExpiresAt != nil {
			ttl := int(time.Until(claims.ExpiresAt.Time).Seconds())
			if ttl > 0 {
				if err := svcCtx.Redis.SetexCtx(r.Context(), auth.BlacklistKey(token), "1", ttl); err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}
		}
		httpx.OkJsonCtx(r.Context(), w, &types.LogoutResp{Code: http.StatusOK, Msg: "已退出登录"})
	}
}
