package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRoleLogic {
	return &UpdateUserRoleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserRoleLogic) UpdateUserRole(req *types.UpdateUserRoleReq) (*types.BaseResp, error) {
	ctx, _ := injectAdmin(l.ctx)
	s, err := l.svcCtx.UserRpc.UpdateUserRole(ctx, &user.UpdateUserRoleReq{Id: req.Id, Role: req.Role})
	if err != nil {
		l.Errorf("update user role rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
