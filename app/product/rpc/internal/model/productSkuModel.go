package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductSkuModel = (*customProductSkuModel)(nil)

type (
	// ProductSkuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductSkuModel.
	ProductSkuModel interface {
		productSkuModel
		withSession(session sqlx.Session) ProductSkuModel
		ListByProduct(ctx context.Context, productID int64) ([]*ProductSku, error)
		FindActive(ctx context.Context, id uint64) (*ProductSku, error)
		InsertSku(ctx context.Context, data *ProductSku) (uint64, error)
		UpdateSku(ctx context.Context, data *ProductSku) (int64, error)
		SoftDeleteSku(ctx context.Context, id uint64) (int64, error)
	}

	customProductSkuModel struct {
		*defaultProductSkuModel
	}
)

// NewProductSkuModel returns a model for the database table.
func NewProductSkuModel(conn sqlx.SqlConn) ProductSkuModel {
	return &customProductSkuModel{
		defaultProductSkuModel: newProductSkuModel(conn),
	}
}

func (m *customProductSkuModel) withSession(session sqlx.Session) ProductSkuModel {
	return NewProductSkuModel(sqlx.NewSqlConnFromSession(session))
}
