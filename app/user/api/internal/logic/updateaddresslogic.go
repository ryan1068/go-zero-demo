package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateAddressLogic) UpdateAddress(req *types.SaveAddressReq) (*types.AddressResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AddressResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.UserRpc.UpdateAddress(ctx, toSaveAddressReq(req))
	if err != nil {
		l.Errorf("update address rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AddressResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toAddressInfo(s.Data)
	}
	return resp, nil
}
