package model

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const afterSaleSelectColumns = "a.`id`, a.`after_sale_no`, a.`order_id`, a.`order_item_id`, a.`user_id`," +
	" a.`type`, a.`reason`, a.`description`, coalesce(cast(a.`images` as char), '[]') as images," +
	" a.`amount`, a.`status`, a.`admin_remark`, a.`tracking_no`," +
	" unix_timestamp(a.`apply_at`) as apply_at, coalesce(unix_timestamp(a.`review_at`), 0) as review_at," +
	" coalesce(unix_timestamp(a.`buyer_shipped_at`), 0) as buyer_shipped_at," +
	" coalesce(unix_timestamp(a.`seller_received_at`), 0) as seller_received_at," +
	" coalesce(unix_timestamp(a.`refunded_at`), 0) as refunded_at," +
	" unix_timestamp(a.`created_at`) as created_at, unix_timestamp(a.`updated_at`) as updated_at," +
	" coalesce(o.`order_no`, '') as order_no, coalesce(i.`product_name`, '') as product_name," +
	" coalesce(i.`product_cover`, '') as product_cover"

const afterSaleFromClause = " from `after_sale` a left join `orders` o on o.`id` = a.`order_id`" +
	" left join `order_item` i on i.`id` = a.`order_item_id`"

// AfterSaleRow is an after-sale joined with its order and item.
type AfterSaleRow struct {
	Id               int64  `db:"id"`
	AfterSaleNo      string `db:"after_sale_no"`
	OrderId          int64  `db:"order_id"`
	OrderItemId      int64  `db:"order_item_id"`
	UserId           int64  `db:"user_id"`
	Type             int64  `db:"type"`
	Reason           string `db:"reason"`
	Description      string `db:"description"`
	Images           string `db:"images"`
	Amount           int64  `db:"amount"`
	Status           int64  `db:"status"`
	AdminRemark      string `db:"admin_remark"`
	TrackingNo       string `db:"tracking_no"`
	ApplyAt          int64  `db:"apply_at"`
	ReviewAt         int64  `db:"review_at"`
	BuyerShippedAt   int64  `db:"buyer_shipped_at"`
	SellerReceivedAt int64  `db:"seller_received_at"`
	RefundedAt       int64  `db:"refunded_at"`
	CreatedAt        int64  `db:"created_at"`
	UpdatedAt        int64  `db:"updated_at"`
	OrderNo          string `db:"order_no"`
	ProductName      string `db:"product_name"`
	ProductCover     string `db:"product_cover"`
}

// AfterSaleLock is the minimal after-sale row used under FOR UPDATE.
type AfterSaleLock struct {
	UserId  int64 `db:"user_id"`
	Type    int64 `db:"type"`
	Status  int64 `db:"status"`
	Amount  int64 `db:"amount"`
	OrderId int64 `db:"order_id"`
}

// AfterSaleFilter is the filter for the after-sale list.
type AfterSaleFilter struct {
	Admin    bool
	UserID   int64
	Status   int64
	Page     int64
	PageSize int64
}

// NewAfterSale is an after-sale application to insert.
type NewAfterSale struct {
	AfterSaleNo string
	OrderID     int64
	OrderItemID int64
	UserID      int64
	Type        int64
	Reason      string
	Description string
	ImagesJSON  string
	Amount      int64
	Status      int64
}

// LockAfterSale locks one after-sale row for a state transition.
func (m *customAfterSaleModel) LockAfterSale(ctx context.Context, session sqlx.Session, id int64) (*AfterSaleLock, error) {
	var row AfterSaleLock
	err := session.QueryRowCtx(ctx, &row,
		"select `user_id`, `type`, `status`, `amount`, `order_id` from "+m.table+
			" where `id` = ? and `deleted_at` is null for update", id,
	)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// CountActiveAfterSales counts non-terminal after-sales of an order item.
func (m *customAfterSaleModel) CountActiveAfterSales(ctx context.Context, session sqlx.Session, orderItemID, userID int64, excluded []int64) (int64, error) {
	placeholders := make([]string, 0, len(excluded))
	args := make([]any, 0, len(excluded)+2)
	args = append(args, orderItemID, userID)
	for _, status := range excluded {
		placeholders = append(placeholders, "?")
		args = append(args, status)
	}
	var total int64
	query := "select count(*) from " + m.table +
		" where `order_item_id` = ? and `user_id` = ? and `status` not in (" +
		strings.Join(placeholders, ",") + ") and `deleted_at` is null"
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &total, query, args...)
	return total, err
}

// InsertAfterSale inserts an application and returns its id.
func (m *customAfterSaleModel) InsertAfterSale(ctx context.Context, session sqlx.Session, data NewAfterSale) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"insert into "+m.table+" (`after_sale_no`, `order_id`, `order_item_id`, `user_id`, `type`,"+
			" `reason`, `description`, `images`, `amount`, `status`, `apply_at`)"+
			" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())",
		data.AfterSaleNo, data.OrderID, data.OrderItemID, data.UserID, data.Type,
		data.Reason, data.Description, data.ImagesJSON, data.Amount, data.Status,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// MarkAfterSaleReviewed approves or rejects an application.
func (m *customAfterSaleModel) MarkAfterSaleReviewed(ctx context.Context, session sqlx.Session, id, status int64, adminRemark string) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `status` = ?, `review_at` = now(), `admin_remark` = ? where `id` = ?",
		status, adminRemark, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// MarkAfterSaleBuyerShipped records the buyer's return shipment.
func (m *customAfterSaleModel) MarkAfterSaleBuyerShipped(ctx context.Context, session sqlx.Session, id, toStatus int64, trackingNo string) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `status` = ?, `tracking_no` = ?, `buyer_shipped_at` = now() where `id` = ?",
		toStatus, trackingNo, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// MarkAfterSaleSellerReceived records the seller receiving the return.
func (m *customAfterSaleModel) MarkAfterSaleSellerReceived(ctx context.Context, session sqlx.Session, id, toStatus int64, adminRemark string) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `status` = ?, `seller_received_at` = now(), `admin_remark` = ? where `id` = ?",
		toStatus, adminRemark, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// MarkAfterSaleRefunded marks an application as refunded.
func (m *customAfterSaleModel) MarkAfterSaleRefunded(ctx context.Context, session sqlx.Session, id, toStatus int64, adminRemark string) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `status` = ?, `refunded_at` = now(), `admin_remark` = ? where `id` = ?",
		toStatus, adminRemark, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// FindAfterSale returns one after-sale, scoped to the user unless admin.
func (m *customAfterSaleModel) FindAfterSale(ctx context.Context, session sqlx.Session, id, userID int64, admin bool) (*AfterSaleRow, error) {
	query := "select " + afterSaleSelectColumns + afterSaleFromClause +
		" where a.`id` = ? and a.`deleted_at` is null"
	args := []any{id}
	if !admin {
		query += " and a.`user_id` = ?"
		args = append(args, userID)
	}
	var row AfterSaleRow
	if err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &row, query, args...); err != nil {
		return nil, err
	}
	return &row, nil
}

func afterSaleCondition(f AfterSaleFilter) (string, []any) {
	where := []string{"a.`deleted_at` is null"}
	args := make([]any, 0, 2)
	if !f.Admin {
		where = append(where, "a.`user_id` = ?")
		args = append(args, f.UserID)
	}
	if f.Status > 0 {
		where = append(where, "a.`status` = ?")
		args = append(args, f.Status)
	}
	return strings.Join(where, " AND "), args
}

// CountAfterSales counts after-sales matching the filter.
func (m *customAfterSaleModel) CountAfterSales(ctx context.Context, session sqlx.Session, f AfterSaleFilter) (int64, error) {
	condition, args := afterSaleCondition(f)
	var total int64
	query := "select count(*) from " + m.table + " a where " + condition
	if err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &total, query, args...); err != nil {
		return 0, err
	}
	return total, nil
}

// ListAfterSales returns a page of after-sales.
func (m *customAfterSaleModel) ListAfterSales(ctx context.Context, session sqlx.Session, f AfterSaleFilter) ([]*AfterSaleRow, error) {
	condition, args := afterSaleCondition(f)
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)
	var list []*AfterSaleRow
	query := "select " + afterSaleSelectColumns + afterSaleFromClause +
		" where " + condition + " order by a.`id` desc limit ? offset ?"
	if err := sessionOrConn(m.conn, session).QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}
