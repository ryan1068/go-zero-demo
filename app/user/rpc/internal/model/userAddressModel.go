package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserAddressModel = (*customUserAddressModel)(nil)

type (
	// UserAddressModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserAddressModel.
	UserAddressModel interface {
		userAddressModel
		withSession(session sqlx.Session) UserAddressModel
		ListByUser(ctx context.Context, userID int64) ([]*UserAddress, error)
		CountByUser(ctx context.Context, session sqlx.Session, userID int64) (int64, error)
		FindByIDAndUser(ctx context.Context, session sqlx.Session, id uint64, userID int64) (*UserAddress, error)
		FindForUpdate(ctx context.Context, session sqlx.Session, id uint64, userID int64) (*UserAddress, error)
		InsertTx(ctx context.Context, session sqlx.Session, data *UserAddress) (sql.Result, error)
		UpdateFields(ctx context.Context, session sqlx.Session, data *UserAddress) error
		SoftDelete(ctx context.Context, session sqlx.Session, id uint64, userID int64) error
		ClearDefault(ctx context.Context, session sqlx.Session, userID int64) error
		SetDefault(ctx context.Context, session sqlx.Session, id uint64, userID int64) error
		FindNextDefault(ctx context.Context, session sqlx.Session, userID int64) (uint64, error)
	}

	customUserAddressModel struct {
		*defaultUserAddressModel
	}
)

// NewUserAddressModel returns a model for the database table.
func NewUserAddressModel(conn sqlx.SqlConn) UserAddressModel {
	return &customUserAddressModel{
		defaultUserAddressModel: newUserAddressModel(conn),
	}
}

func (m *customUserAddressModel) withSession(session sqlx.Session) UserAddressModel {
	return NewUserAddressModel(sqlx.NewSqlConnFromSession(session))
}
