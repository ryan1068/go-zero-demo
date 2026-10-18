package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAddressesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAddressesLogic {
	return &ListAddressesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListAddressesLogic) ListAddresses() (*types.ListAddressesResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.ListAddressesResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.UserRpc.ListAddresses(ctx, &user.EmptyReq{})
	if err != nil {
		l.Errorf("list addresses rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ListAddressesResp{Code: s.Code, Msg: s.Msg}
	resp.Data.List = make([]types.AddressInfo, 0, len(s.List))
	for _, item := range s.List {
		if item != nil {
			resp.Data.List = append(resp.Data.List, toAddressInfo(item))
		}
	}
	return resp, nil
}
