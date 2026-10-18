package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderItemModel = (*customOrderItemModel)(nil)

type (
	// OrderItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderItemModel.
	OrderItemModel interface {
		orderItemModel
		withSession(session sqlx.Session) OrderItemModel
		FindOrderItems(ctx context.Context, session sqlx.Session, orderID int64) ([]*OrderItemRow, error)
		FindItemForAfterSale(ctx context.Context, session sqlx.Session, orderID, itemID int64) (*AfterSaleItem, error)
		InsertOrderItems(ctx context.Context, session sqlx.Session, orderID int64, items []NewOrderItem) error
		StockLines(ctx context.Context, session sqlx.Session, orderID int64) ([]*StockLine, error)
		FindOrderItemsByOrderIDs(ctx context.Context, session sqlx.Session, orderIDs []int64) (map[int64][]*OrderItemRow, error)
	}

	customOrderItemModel struct {
		*defaultOrderItemModel
	}
)

// NewOrderItemModel returns a model for the database table.
func NewOrderItemModel(conn sqlx.SqlConn) OrderItemModel {
	return &customOrderItemModel{
		defaultOrderItemModel: newOrderItemModel(conn),
	}
}

func (m *customOrderItemModel) withSession(session sqlx.Session) OrderItemModel {
	return NewOrderItemModel(sqlx.NewSqlConnFromSession(session))
}
