package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddressLogic {
	return &DeleteAddressLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteAddressLogic) DeleteAddress(req *types.PathIdReq) (*types.BaseResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.BaseResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.UserRpc.DeleteAddress(ctx, &user.AddressReq{Id: req.Id})
	if err != nil {
		l.Errorf("delete address rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
