package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CategoryModel = (*customCategoryModel)(nil)

type (
	// CategoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCategoryModel.
	CategoryModel interface {
		categoryModel
		withSession(session sqlx.Session) CategoryModel
		List(ctx context.Context, admin bool) ([]*Category, error)
		FindActive(ctx context.Context, id uint64) (*Category, error)
		InsertCategory(ctx context.Context, data *Category) (uint64, error)
		UpdateCategory(ctx context.Context, data *Category) (int64, error)
		SoftDeleteCategory(ctx context.Context, id uint64) (int64, error)
		CountProducts(ctx context.Context, id uint64) (int64, error)
		CountChildren(ctx context.Context, id uint64) (int64, error)
		Exists(ctx context.Context, id uint64) (bool, error)
	}

	customCategoryModel struct {
		*defaultCategoryModel
	}
)

// NewCategoryModel returns a model for the database table.
func NewCategoryModel(conn sqlx.SqlConn) CategoryModel {
	return &customCategoryModel{
		defaultCategoryModel: newCategoryModel(conn),
	}
}

func (m *customCategoryModel) withSession(session sqlx.Session) CategoryModel {
	return NewCategoryModel(sqlx.NewSqlConnFromSession(session))
}
