package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetDefaultAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDefaultAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultAddressLogic {
	return &SetDefaultAddressLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SetDefaultAddressLogic) SetDefaultAddress(req *types.PathIdReq) (*types.BaseResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.BaseResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.UserRpc.SetDefaultAddress(ctx, &user.AddressReq{Id: req.Id})
	if err != nil {
		l.Errorf("set default address rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
