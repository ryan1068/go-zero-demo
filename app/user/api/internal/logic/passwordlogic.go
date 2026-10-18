package logic

import (
	"context"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (*types.BaseResp, error) {
	userID, ok := auth.UserID(l.ctx)
	if !ok {
		return &types.BaseResp{Code: 401, Msg: "未授权"}, nil
	}
	ctx := auth.InjectUserID(l.ctx, userID)
	s, err := l.svcCtx.UserRpc.ChangePassword(ctx, &user.ChangePasswordReq{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		l.Errorf("change password rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
