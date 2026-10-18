package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RefundRecordModel = (*customRefundRecordModel)(nil)

type (
	// RefundRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefundRecordModel.
	RefundRecordModel interface {
		refundRecordModel
		withSession(session sqlx.Session) RefundRecordModel
		InsertRefundRecord(ctx context.Context, session sqlx.Session, data NewRefundRecord) error
	}

	customRefundRecordModel struct {
		*defaultRefundRecordModel
	}
)

// NewRefundRecordModel returns a model for the database table.
func NewRefundRecordModel(conn sqlx.SqlConn) RefundRecordModel {
	return &customRefundRecordModel{
		defaultRefundRecordModel: newRefundRecordModel(conn),
	}
}

func (m *customRefundRecordModel) withSession(session sqlx.Session) RefundRecordModel {
	return NewRefundRecordModel(sqlx.NewSqlConnFromSession(session))
}
