package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// NewPaymentRecord is a payment record to insert.
type NewPaymentRecord struct {
	PaymentNo     string
	OrderID       int64
	OrderNo       string
	UserID        int64
	Amount        int64
	PayType       int64
	Status        int64
	TransactionID string
}

// InsertPaymentRecord inserts a successful payment record.
func (m *customPaymentRecordModel) InsertPaymentRecord(ctx context.Context, session sqlx.Session, data NewPaymentRecord) error {
	_, err := session.ExecCtx(ctx,
		"insert into "+m.table+" (`payment_no`, `order_id`, `order_no`, `user_id`, `amount`,"+
			" `pay_type`, `status`, `transaction_id`, `paid_at`)"+
			" values (?, ?, ?, ?, ?, ?, ?, ?, now())",
		data.PaymentNo, data.OrderID, data.OrderNo, data.UserID, data.Amount,
		data.PayType, data.Status, data.TransactionID,
	)
	return err
}
