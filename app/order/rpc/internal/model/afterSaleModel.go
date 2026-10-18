package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AfterSaleModel = (*customAfterSaleModel)(nil)

type (
	// AfterSaleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAfterSaleModel.
	AfterSaleModel interface {
		afterSaleModel
		withSession(session sqlx.Session) AfterSaleModel
		LockAfterSale(ctx context.Context, session sqlx.Session, id int64) (*AfterSaleLock, error)
		CountActiveAfterSales(ctx context.Context, session sqlx.Session, orderItemID, userID int64, excluded []int64) (int64, error)
		InsertAfterSale(ctx context.Context, session sqlx.Session, data NewAfterSale) (int64, error)
		MarkAfterSaleReviewed(ctx context.Context, session sqlx.Session, id, status int64, adminRemark string) (int64, error)
		MarkAfterSaleBuyerShipped(ctx context.Context, session sqlx.Session, id, toStatus int64, trackingNo string) (int64, error)
		MarkAfterSaleSellerReceived(ctx context.Context, session sqlx.Session, id, toStatus int64, adminRemark string) (int64, error)
		MarkAfterSaleRefunded(ctx context.Context, session sqlx.Session, id, toStatus int64, adminRemark string) (int64, error)
		FindAfterSale(ctx context.Context, session sqlx.Session, id, userID int64, admin bool) (*AfterSaleRow, error)
		CountAfterSales(ctx context.Context, session sqlx.Session, f AfterSaleFilter) (int64, error)
		ListAfterSales(ctx context.Context, session sqlx.Session, f AfterSaleFilter) ([]*AfterSaleRow, error)
	}

	customAfterSaleModel struct {
		*defaultAfterSaleModel
	}
)

// NewAfterSaleModel returns a model for the database table.
func NewAfterSaleModel(conn sqlx.SqlConn) AfterSaleModel {
	return &customAfterSaleModel{
		defaultAfterSaleModel: newAfterSaleModel(conn),
	}
}

func (m *customAfterSaleModel) withSession(session sqlx.Session) AfterSaleModel {
	return NewAfterSaleModel(sqlx.NewSqlConnFromSession(session))
}
