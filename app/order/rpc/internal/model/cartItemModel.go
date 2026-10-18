package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CartItemModel = (*customCartItemModel)(nil)

type (
	// CartItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCartItemModel.
	CartItemModel interface {
		cartItemModel
		withSession(session sqlx.Session) CartItemModel
		FindSelectedForUpdate(ctx context.Context, session sqlx.Session, userID int64) ([]*CartSelection, error)
		FindForUpdate(ctx context.Context, session sqlx.Session, userID, productID int64) (*CartItemLock, error)
		LockCartProduct(ctx context.Context, session sqlx.Session, id, userID int64) (*CartProductLock, error)
		InsertCartItem(ctx context.Context, session sqlx.Session, data NewCartItem) (int64, error)
		UpdateCartItem(ctx context.Context, session sqlx.Session, id, userID, quantity, selected int64) (int64, error)
		DeleteCartItem(ctx context.Context, session sqlx.Session, id, userID int64) (int64, error)
		UpdateCartSelectedByUser(ctx context.Context, session sqlx.Session, userID, selected int64) error
		DeleteSelectedCartItems(ctx context.Context, session sqlx.Session, userID int64) error
		FindCartItemRow(ctx context.Context, session sqlx.Session, id, userID int64) (*CartItemRow, error)
		ListCartItemRows(ctx context.Context, session sqlx.Session, userID int64) ([]*CartItemRow, error)
	}

	customCartItemModel struct {
		*defaultCartItemModel
	}
)

// NewCartItemModel returns a model for the database table.
func NewCartItemModel(conn sqlx.SqlConn) CartItemModel {
	return &customCartItemModel{
		defaultCartItemModel: newCartItemModel(conn),
	}
}

func (m *customCartItemModel) withSession(session sqlx.Session) CartItemModel {
	return NewCartItemModel(sqlx.NewSqlConnFromSession(session))
}
