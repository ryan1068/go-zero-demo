package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListOrdersLogic {
	return &AdminListOrdersLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminListOrdersLogic) ListOrders(req *types.ListOrdersReq) (*types.OrderListResp, error) {
	return NewListOrdersLogic(l.ctx, l.svcCtx).listOrders(req, true)
}

type AdminGetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetOrderLogic {
	return &AdminGetOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminGetOrderLogic) GetOrder(id int64) (*types.OrderResp, error) {
	return NewGetOrderLogic(l.ctx, l.svcCtx).get(id, true)
}

type ShipOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShipOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipOrderLogic {
	return &ShipOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ShipOrderLogic) ShipOrder(req *types.ShipOrderReq) (*types.OrderResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.OrderResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.ShipOrder(ctx, &order.ShipOrderReq{
		Id:               req.Id,
		LogisticsCompany: req.LogisticsCompany,
		TrackingNo:       req.TrackingNo,
	})
	if err != nil {
		l.Errorf("ship order rpc failed: %v", err)
		return nil, err
	}
	resp := &types.OrderResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toOrderInfo(s.Data)
	}
	return resp, nil
}

type AdminCancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCancelOrderLogic {
	return &AdminCancelOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminCancelOrderLogic) CancelOrder(req *types.CancelOrderReq) (*types.OrderResp, error) {
	return NewCancelOrderLogic(l.ctx, l.svcCtx).CancelOrder(req, true)
}

type DashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DashboardLogic {
	return &DashboardLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DashboardLogic) Dashboard() (*types.DashboardResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.DashboardResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.Dashboard(ctx, &order.EmptyReq{})
	if err != nil {
		l.Errorf("dashboard rpc failed: %v", err)
		return nil, err
	}
	return &types.DashboardResp{
		Code: s.Code,
		Msg:  s.Msg,
		Data: types.DashboardData{
			TodayOrderCount:      s.TodayOrderCount,
			TodayPaidAmount:      s.TodayPaidAmount,
			PendingShipmentCount: s.PendingShipmentCount,
			ProductCount:         s.ProductCount,
			LowStockCount:        s.LowStockCount,
		},
	}, nil
}
