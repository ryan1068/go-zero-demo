package logic

import (
	"context"
	"time"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// Logout revokes the current token until its natural expiry.
func (l *LogoutLogic) Logout(token string) (*types.LogoutResp, error) {
	if token == "" {
		return &types.LogoutResp{Code: 401, Msg: "未授权"}, nil
	}
	claims, err := auth.Parse(l.svcCtx.Config.Jwt.Secret, token)
	if err != nil {
		return &types.LogoutResp{Code: 401, Msg: "未授权"}, nil
	}
	if claims.ExpiresAt != nil {
		ttl := int(time.Until(claims.ExpiresAt.Time).Seconds())
		if ttl > 0 {
			if err := l.svcCtx.Redis.SetexCtx(l.ctx, auth.BlacklistKey(token), "1", ttl); err != nil {
				l.Errorf("write token blacklist failed: %v", err)
				return nil, err
			}
		}
	}
	return &types.LogoutResp{Code: 200, Msg: "已退出登录"}, nil
}
