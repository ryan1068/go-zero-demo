package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// NewOrderLog is an order status log to insert.
type NewOrderLog struct {
	OrderID      int64
	FromStatus   int64
	ToStatus     int64
	OperatorType int64
	OperatorID   int64
	Remark       string
}

// InsertOrderLog appends one order status log.
func (m *customOrderLogModel) InsertOrderLog(ctx context.Context, session sqlx.Session, data NewOrderLog) error {
	_, err := session.ExecCtx(ctx,
		"insert into "+m.table+" (`order_id`, `from_status`, `to_status`, `operator_type`,"+
			" `operator_id`, `remark`) values (?, ?, ?, ?, ?, ?)",
		data.OrderID, data.FromStatus, data.ToStatus, data.OperatorType,
		data.OperatorID, data.Remark,
	)
	return err
}
