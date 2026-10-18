package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

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
