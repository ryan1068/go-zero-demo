package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAfterSaleLogic {
	return &GetAfterSaleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAfterSaleLogic) GetAfterSale(req *types.AfterSaleReq) (*types.AfterSaleResp, error) {
	return l.get(req.Id, false)
}

func (l *GetAfterSaleLogic) get(id int64, admin bool) (*types.AfterSaleResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AfterSaleResp{Code: 401, Msg: "未授权"}, nil
	}
	adminFlag := int32(0)
	if admin {
		adminFlag = 1
	}
	s, err := l.svcCtx.OrderRpc.GetAfterSale(ctx, &order.AfterSaleReq{Id: id, Admin: adminFlag})
	if err != nil {
		l.Errorf("get after sale rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AfterSaleResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toAfterSaleInfo(s.Data)
	}
	return resp, nil
}
