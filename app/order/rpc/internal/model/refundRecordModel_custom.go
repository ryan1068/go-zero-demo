package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// NewRefundRecord is a refund record to insert.
type NewRefundRecord struct {
	RefundNo    string
	AfterSaleID int64
	OrderID     int64
	UserID      int64
	Amount      int64
	Channel     string
	Status      int64
}

// InsertRefundRecord inserts a mock refund record.
func (m *customRefundRecordModel) InsertRefundRecord(ctx context.Context, session sqlx.Session, data NewRefundRecord) error {
	_, err := session.ExecCtx(ctx,
		"insert into "+m.table+" (`refund_no`, `after_sale_id`, `order_id`, `user_id`,"+
			" `amount`, `channel`, `status`, `refunded_at`) values (?, ?, ?, ?, ?, ?, ?, now())",
		data.RefundNo, data.AfterSaleID, data.OrderID, data.UserID, data.Amount,
		data.Channel, data.Status,
	)
	return err
}
