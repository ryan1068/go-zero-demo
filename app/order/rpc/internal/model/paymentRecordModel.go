package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentRecordModel = (*customPaymentRecordModel)(nil)

type (
	// PaymentRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentRecordModel.
	PaymentRecordModel interface {
		paymentRecordModel
		withSession(session sqlx.Session) PaymentRecordModel
		InsertPaymentRecord(ctx context.Context, session sqlx.Session, data NewPaymentRecord) error
	}

	customPaymentRecordModel struct {
		*defaultPaymentRecordModel
	}
)

// NewPaymentRecordModel returns a model for the database table.
func NewPaymentRecordModel(conn sqlx.SqlConn) PaymentRecordModel {
	return &customPaymentRecordModel{
		defaultPaymentRecordModel: newPaymentRecordModel(conn),
	}
}

func (m *customPaymentRecordModel) withSession(session sqlx.Session) PaymentRecordModel {
	return NewPaymentRecordModel(sqlx.NewSqlConnFromSession(session))
}
