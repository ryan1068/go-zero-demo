package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserStatusLogic) UpdateUserStatus(req *types.UpdateUserStatusReq) (*types.BaseResp, error) {
	ctx, _ := injectAdmin(l.ctx)
	s, err := l.svcCtx.UserRpc.UpdateUserStatus(ctx, &user.UpdateUserStatusReq{
		Id:     req.Id,
		Status: req.Status,
	})
	if err != nil {
		l.Errorf("update user status rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
