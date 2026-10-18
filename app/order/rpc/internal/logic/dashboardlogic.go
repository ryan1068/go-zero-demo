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

	var resp order.DashboardResp
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL AND created_at >= CURDATE()",
	).Scan(&resp.TodayOrderCount); err != nil {
		l.Errorf("count today orders failed: %v", err)
		return nil, err
	}
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COALESCE(SUM(pay_amount), 0) FROM orders WHERE deleted_at IS NULL AND status IN (?, ?, ?) AND paid_at >= CURDATE()",
		orderStatusPendingShipment, orderStatusShipped, orderStatusCompleted,
	).Scan(&resp.TodayPaidAmount); err != nil {
		l.Errorf("sum today paid amount failed: %v", err)
		return nil, err
	}
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL AND status = ?",
		orderStatusPendingShipment,
	).Scan(&resp.PendingShipmentCount); err != nil {
		l.Errorf("count pending shipment orders failed: %v", err)
		return nil, err
	}
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM looklook_product.product WHERE deleted_at IS NULL",
	).Scan(&resp.ProductCount); err != nil {
		l.Errorf("count products failed: %v", err)
		return nil, err
	}
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM looklook_product.product WHERE deleted_at IS NULL AND status = 1 AND stock <= ?",
		lowStockThreshold,
	).Scan(&resp.LowStockCount); err != nil {
		l.Errorf("count low stock products failed: %v", err)
		return nil, err
	}

	resp.Code = errcode.Success
	resp.Msg = errcode.Msg(errcode.Success)
	return &resp, nil
}
