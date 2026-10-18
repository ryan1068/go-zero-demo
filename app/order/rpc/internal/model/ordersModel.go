package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrdersModel = (*customOrdersModel)(nil)

type (
	// OrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrdersModel.
	OrdersModel interface {
		ordersModel
		withSession(session sqlx.Session) OrdersModel
		FindOrder(ctx context.Context, session sqlx.Session, id, userID int64, admin bool) (*OrderRow, error)
		LockOrder(ctx context.Context, session sqlx.Session, id int64, userID *int64) (*OrderLock, error)
		CountOrders(ctx context.Context, session sqlx.Session, f OrderFilter) (int64, error)
		ListOrders(ctx context.Context, session sqlx.Session, f OrderFilter) ([]*OrderRow, error)
		FindExpiredPendingOrderIDs(ctx context.Context, session sqlx.Session, status, minutes, limit int64) ([]int64, error)
		FindExpiredShippedOrderIDs(ctx context.Context, session sqlx.Session, status, days, limit int64) ([]int64, error)
		CountTodayOrders(ctx context.Context, session sqlx.Session) (int64, error)
		SumTodayPaidAmount(ctx context.Context, session sqlx.Session, statuses []int64) (int64, error)
		CountPendingShipment(ctx context.Context, session sqlx.Session, status int64) (int64, error)
		InsertOrder(ctx context.Context, session sqlx.Session, data NewOrder) (int64, error)
		UpdateOrderPaid(ctx context.Context, session sqlx.Session, id, userID, fromStatus, toStatus, payType int64, transactionID string) (int64, error)
		UpdateOrderShipped(ctx context.Context, session sqlx.Session, id, toStatus int64, company, trackingNo string) (int64, error)
		UpdateOrderCompletedByUser(ctx context.Context, session sqlx.Session, id, userID, toStatus int64) (int64, error)
		UpdateOrderCompletedByStatus(ctx context.Context, session sqlx.Session, id, fromStatus, toStatus int64) (int64, error)
		UpdateOrderCanceled(ctx context.Context, session sqlx.Session, id, toStatus int64, reason string) (int64, error)
	}

	customOrdersModel struct {
		*defaultOrdersModel
	}
)

// NewOrdersModel returns a model for the database table.
func NewOrdersModel(conn sqlx.SqlConn) OrdersModel {
	return &customOrdersModel{
		defaultOrdersModel: newOrdersModel(conn),
	}
}

func (m *customOrdersModel) withSession(session sqlx.Session) OrdersModel {
	return NewOrdersModel(sqlx.NewSqlConnFromSession(session))
}
