package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		withSession(session sqlx.Session) UserModel
		FindUsers(ctx context.Context, f UserFilter) ([]*User, int64, error)
		UpdateStatus(ctx context.Context, id, status int64) (int64, error)
		UpdateRole(ctx context.Context, id, role int64) (int64, error)
		UpdatePassword(ctx context.Context, id int64, password string) error
		UpdateLastLogin(ctx context.Context, id int64) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) withSession(session sqlx.Session) UserModel {
	return NewUserModel(sqlx.NewSqlConnFromSession(session))
}
