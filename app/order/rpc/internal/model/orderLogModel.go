package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderLogModel = (*customOrderLogModel)(nil)

type (
	// OrderLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderLogModel.
	OrderLogModel interface {
		orderLogModel
		withSession(session sqlx.Session) OrderLogModel
		InsertOrderLog(ctx context.Context, session sqlx.Session, data NewOrderLog) error
	}

	customOrderLogModel struct {
		*defaultOrderLogModel
	}
)

// NewOrderLogModel returns a model for the database table.
func NewOrderLogModel(conn sqlx.SqlConn) OrderLogModel {
	return &customOrderLogModel{
		defaultOrderLogModel: newOrderLogModel(conn),
	}
}

func (m *customOrderLogModel) withSession(session sqlx.Session) OrderLogModel {
	return NewOrderLogModel(sqlx.NewSqlConnFromSession(session))
}
