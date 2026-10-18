package logic

import (
	"context"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type DashboardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DashboardLogic {
	return &DashboardLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DashboardLogic) Dashboard(_ *order.EmptyReq) (*order.DashboardResp, error) {
	if _, ok := auth.UserIDFromIncoming(l.ctx); !ok {
		return &order.DashboardResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}

	todayOrders, err := l.svcCtx.OrdersModel.CountTodayOrders(l.ctx, nil)
	if err != nil {
		l.Errorf("count today orders failed: %v", err)
		return nil, err
	}
	todayPaid, err := l.svcCtx.OrdersModel.SumTodayPaidAmount(l.ctx, nil,
		[]int64{orderStatusPendingShipment, orderStatusShipped, orderStatusCompleted})
	if err != nil {
		l.Errorf("sum today paid amount failed: %v", err)
		return nil, err
	}
	pendingShipment, err := l.svcCtx.OrdersModel.CountPendingShipment(l.ctx, nil, orderStatusPendingShipment)
	if err != nil {
		l.Errorf("count pending shipment orders failed: %v", err)
		return nil, err
	}
	productCount, err := l.svcCtx.ProductModel.CountProducts(l.ctx, nil)
	if err != nil {
		l.Errorf("count products failed: %v", err)
		return nil, err
	}
	lowStockCount, err := l.svcCtx.ProductModel.CountLowStock(l.ctx, nil, lowStockThreshold)
	if err != nil {
		l.Errorf("count low stock products failed: %v", err)
		return nil, err
	}

	return &order.DashboardResp{
		Code:                 errcode.Success,
		Msg:                  errcode.Msg(errcode.Success),
		TodayOrderCount:      todayOrders,
		TodayPaidAmount:      todayPaid,
		PendingShipmentCount: pendingShipment,
		ProductCount:         productCount,
		LowStockCount:        lowStockCount,
	}, nil
}
