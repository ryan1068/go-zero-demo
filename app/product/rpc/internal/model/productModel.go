package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductModel = (*customProductModel)(nil)

type (
	// ProductModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductModel.
	ProductModel interface {
		productModel
		withSession(session sqlx.Session) ProductModel
		ListProducts(ctx context.Context, f ProductFilter) ([]*ProductRow, int64, error)
		FindProductRow(ctx context.Context, id uint64, admin bool) (*ProductRow, error)
		InsertProduct(ctx context.Context, data *Product) (uint64, error)
		UpdateProduct(ctx context.Context, data *Product) (int64, error)
		SoftDeleteProduct(ctx context.Context, id uint64) (int64, error)
		UpdateProductStatus(ctx context.Context, id uint64, status int64) (int64, error)
		Exists(ctx context.Context, id uint64) (bool, error)
		Stats(ctx context.Context, lowStock int64) (ProductStats, error)
	}

	customProductModel struct {
		*defaultProductModel
	}
)

// NewProductModel returns a model for the database table.
func NewProductModel(conn sqlx.SqlConn) ProductModel {
	return &customProductModel{
		defaultProductModel: newProductModel(conn),
	}
}

func (m *customProductModel) withSession(session sqlx.Session) ProductModel {
	return NewProductModel(sqlx.NewSqlConnFromSession(session))
}
