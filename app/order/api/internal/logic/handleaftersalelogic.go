package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleAfterSaleLogic {
	return &HandleAfterSaleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *HandleAfterSaleLogic) HandleAfterSale(req *types.HandleAfterSaleReq) (*types.AfterSaleResp, error) {
	return l.handle(req.Id, req.Action, req.AdminRemark, req.TrackingNo, true)
}

func (l *HandleAfterSaleLogic) handle(id int64, action, adminRemark, trackingNo string, admin bool) (*types.AfterSaleResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AfterSaleResp{Code: 401, Msg: "未授权"}, nil
	}
	adminFlag := int32(0)
	if admin {
		adminFlag = 1
	}
	s, err := l.svcCtx.OrderRpc.HandleAfterSale(ctx, &order.HandleAfterSaleReq{
		Id:          id,
		Action:      action,
		AdminRemark: adminRemark,
		TrackingNo:  trackingNo,
		Admin:       adminFlag,
	})
	if err != nil {
		l.Errorf("handle after sale rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AfterSaleResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toAfterSaleInfo(s.Data)
	}
	return resp, nil
}
