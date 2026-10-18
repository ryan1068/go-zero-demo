package model

import (
	"context"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const orderSelectColumns = "`id`, `order_no`, `user_id`, `total_amount`, `freight_amount`," +
	" `discount_amount`, `pay_amount`, `status`, `pay_type`, coalesce(`transaction_id`, '') as transaction_id," +
	" `receiver_name`, `receiver_phone`, coalesce(`receiver_province`, '') as receiver_province," +
	" coalesce(`receiver_city`, '') as receiver_city, coalesce(`receiver_district`, '') as receiver_district," +
	" `receiver_address`, coalesce(`receiver_postal_code`, '') as receiver_postal_code," +
	" coalesce(`remark`, '') as remark, coalesce(`admin_remark`, '') as admin_remark," +
	" coalesce(`logistics_company`, '') as logistics_company, coalesce(`tracking_no`, '') as tracking_no," +
	" coalesce(unix_timestamp(`paid_at`), 0) as paid_at, coalesce(unix_timestamp(`shipped_at`), 0) as shipped_at," +
	" coalesce(unix_timestamp(`completed_at`), 0) as completed_at," +
	" coalesce(unix_timestamp(`canceled_at`), 0) as canceled_at," +
	" coalesce(`cancel_reason`, '') as cancel_reason," +
	" unix_timestamp(`created_at`) as created_at, unix_timestamp(`updated_at`) as updated_at"

// OrderRow is an order header with unix timestamps.
type OrderRow struct {
	Id                 int64  `db:"id"`
	OrderNo            string `db:"order_no"`
	UserId             int64  `db:"user_id"`
	TotalAmount        int64  `db:"total_amount"`
	FreightAmount      int64  `db:"freight_amount"`
	DiscountAmount     int64  `db:"discount_amount"`
	PayAmount          int64  `db:"pay_amount"`
	Status             int64  `db:"status"`
	PayType            int64  `db:"pay_type"`
	TransactionId      string `db:"transaction_id"`
	ReceiverName       string `db:"receiver_name"`
	ReceiverPhone      string `db:"receiver_phone"`
	ReceiverProvince   string `db:"receiver_province"`
	ReceiverCity       string `db:"receiver_city"`
	ReceiverDistrict   string `db:"receiver_district"`
	ReceiverAddress    string `db:"receiver_address"`
	ReceiverPostalCode string `db:"receiver_postal_code"`
	Remark             string `db:"remark"`
	AdminRemark        string `db:"admin_remark"`
	LogisticsCompany   string `db:"logistics_company"`
	TrackingNo         string `db:"tracking_no"`
	PaidAt             int64  `db:"paid_at"`
	ShippedAt          int64  `db:"shipped_at"`
	CompletedAt        int64  `db:"completed_at"`
	CanceledAt         int64  `db:"canceled_at"`
	CancelReason       string `db:"cancel_reason"`
	CreatedAt          int64  `db:"created_at"`
	UpdatedAt          int64  `db:"updated_at"`
}

// OrderLock is the minimal order row used under FOR UPDATE.
type OrderLock struct {
	Id        int64  `db:"id"`
	UserId    int64  `db:"user_id"`
	OrderNo   string `db:"order_no"`
	PayAmount int64  `db:"pay_amount"`
	Status    int64  `db:"status"`
}

// IDRow is used for id-only queries.
type IDRow struct {
	Id int64 `db:"id"`
}

// OrderFilter is the filter for the order list.
type OrderFilter struct {
	Admin     bool
	UserID    int64
	OrderNo   string
	Status    int64
	StartTime int64
	EndTime   int64
	Page      int64
	PageSize  int64
	MaxRows   int64
}

// NewOrder is the insert-only order parameter set.
type NewOrder struct {
	OrderNo            string
	UserID             int64
	TotalAmount        int64
	PayAmount          int64
	Status             int64
	ReceiverName       string
	ReceiverPhone      string
	ReceiverProvince   string
	ReceiverCity       string
	ReceiverDistrict   string
	ReceiverAddress    string
	ReceiverPostalCode string
	Remark             string
}

func orderCondition(f OrderFilter) (string, []any) {
	where := []string{"`deleted_at` is null"}
	args := make([]any, 0, 6)
	if !f.Admin {
		where = append(where, "`user_id` = ?")
		args = append(args, f.UserID)
	} else if f.UserID > 0 {
		where = append(where, "`user_id` = ?")
		args = append(args, f.UserID)
	}
	if f.OrderNo != "" {
		where = append(where, "`order_no` = ?")
		args = append(args, f.OrderNo)
	}
	if f.Status > 0 {
		where = append(where, "`status` = ?")
		args = append(args, f.Status)
	}
	if f.StartTime > 0 {
		where = append(where, "`created_at` >= from_unixtime(?)")
		args = append(args, f.StartTime)
	}
	if f.EndTime > 0 {
		where = append(where, "`created_at` <= from_unixtime(?)")
		args = append(args, f.EndTime)
	}
	return strings.Join(where, " AND "), args
}

// FindOrder returns an order header, scoped to the user unless admin.
func (m *customOrdersModel) FindOrder(ctx context.Context, session sqlx.Session, id, userID int64, admin bool) (*OrderRow, error) {
	query := "select " + orderSelectColumns + " from " + m.table + " where `id` = ? and `deleted_at` is null"
	args := []any{id}
	if !admin {
		query += " and `user_id` = ?"
		args = append(args, userID)
	}
	query += " limit 1"
	var row OrderRow
	if err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &row, query, args...); err != nil {
		return nil, err
	}
	return &row, nil
}

// LockOrder locks an order row for a state transition.
func (m *customOrdersModel) LockOrder(ctx context.Context, session sqlx.Session, id int64, userID *int64) (*OrderLock, error) {
	query := "select `id`, `user_id`, `order_no`, `pay_amount`, `status` from " + m.table +
		" where `id` = ? and `deleted_at` is null"
	args := []any{id}
	if userID != nil {
		query += " and `user_id` = ?"
		args = append(args, *userID)
	}
	query += " for update"
	var row OrderLock
	if err := session.QueryRowCtx(ctx, &row, query, args...); err != nil {
		return nil, err
	}
	return &row, nil
}

// CountOrders counts orders matching the filter.
func (m *customOrdersModel) CountOrders(ctx context.Context, session sqlx.Session, f OrderFilter) (int64, error) {
	condition, args := orderCondition(f)
	var total int64
	query := "select count(*) from (select 1 from " + m.table + " where " + condition +
		" limit " + strconv.FormatInt(f.MaxRows, 10) + ") as filtered"
	if err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &total, query, args...); err != nil {
		return 0, err
	}
	return total, nil
}

// ListOrders returns a page of order headers.
func (m *customOrdersModel) ListOrders(ctx context.Context, session sqlx.Session, f OrderFilter) ([]*OrderRow, error) {
	condition, args := orderCondition(f)
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)
	var list []*OrderRow
	query := "select " + orderSelectColumns + " from " + m.table + " where " + condition +
		" order by `id` desc limit ? offset ?"
	if err := sessionOrConn(m.conn, session).QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

// FindExpiredPendingOrderIDs returns unpaid orders older than the timeout.
func (m *customOrdersModel) FindExpiredPendingOrderIDs(ctx context.Context, session sqlx.Session, status, minutes, limit int64) ([]int64, error) {
	var rows []*IDRow
	query := "select `id` from " + m.table +
		" where `status` = ? and `deleted_at` is null and `created_at` < date_sub(now(), interval ? minute)" +
		" order by `id` limit ?"
	if err := sessionOrConn(m.conn, session).QueryRowsCtx(ctx, &rows, query, status, minutes, limit); err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.Id)
	}
	return ids, nil
}

// FindExpiredShippedOrderIDs returns shipped orders older than the confirm window.
func (m *customOrdersModel) FindExpiredShippedOrderIDs(ctx context.Context, session sqlx.Session, status, days, limit int64) ([]int64, error) {
	var rows []*IDRow
	query := "select `id` from " + m.table +
		" where `status` = ? and `deleted_at` is null and `shipped_at` is not null" +
		" and `shipped_at` < date_sub(now(), interval ? day) order by `id` limit ?"
	if err := sessionOrConn(m.conn, session).QueryRowsCtx(ctx, &rows, query, status, days, limit); err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.Id)
	}
	return ids, nil
}

// CountTodayOrders counts orders created today.
func (m *customOrdersModel) CountTodayOrders(ctx context.Context, session sqlx.Session) (int64, error) {
	var total int64
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &total,
		"select count(*) from "+m.table+" where `deleted_at` is null and `created_at` >= curdate()",
	)
	return total, err
}

// SumTodayPaidAmount sums paid amounts of today's orders in the given statuses.
func (m *customOrdersModel) SumTodayPaidAmount(ctx context.Context, session sqlx.Session, statuses []int64) (int64, error) {
	placeholders := make([]string, 0, len(statuses))
	args := make([]any, 0, len(statuses))
	for _, status := range statuses {
		placeholders = append(placeholders, "?")
		args = append(args, status)
	}
	var total int64
	query := "select coalesce(sum(`pay_amount`), 0) from " + m.table +
		" where `deleted_at` is null and `status` in (" + strings.Join(placeholders, ",") + ")" +
		" and `paid_at` >= curdate()"
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &total, query, args...)
	return total, err
}

// CountPendingShipment counts orders waiting to be shipped.
func (m *customOrdersModel) CountPendingShipment(ctx context.Context, session sqlx.Session, status int64) (int64, error) {
	var total int64
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &total,
		"select count(*) from "+m.table+" where `deleted_at` is null and `status` = ?", status,
	)
	return total, err
}

// InsertOrder inserts an order and returns its id.
func (m *customOrdersModel) InsertOrder(ctx context.Context, session sqlx.Session, data NewOrder) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"insert into "+m.table+" (`order_no`, `user_id`, `total_amount`, `freight_amount`,"+
			" `discount_amount`, `pay_amount`, `status`, `pay_type`, `receiver_name`, `receiver_phone`,"+
			" `receiver_province`, `receiver_city`, `receiver_district`, `receiver_address`,"+
			" `receiver_postal_code`, `remark`) values (?, ?, ?, 0, 0, ?, ?, 0, ?, ?, ?, ?, ?, ?, ?, ?)",
		data.OrderNo, data.UserID, data.TotalAmount, data.PayAmount, data.Status,
		data.ReceiverName, data.ReceiverPhone, data.ReceiverProvince, data.ReceiverCity,
		data.ReceiverDistrict, data.ReceiverAddress, data.ReceiverPostalCode, data.Remark,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateOrderPaid marks an order as paid.
func (m *customOrdersModel) UpdateOrderPaid(ctx context.Context, session sqlx.Session, id, userID, fromStatus, toStatus, payType int64, transactionID string) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `status` = ?, `pay_type` = ?, `transaction_id` = ?, `paid_at` = now()"+
			" where `id` = ? and `user_id` = ? and `status` = ?",
		toStatus, payType, transactionID, id, userID, fromStatus,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateOrderShipped marks an order as shipped.
func (m *customOrdersModel) UpdateOrderShipped(ctx context.Context, session sqlx.Session, id, toStatus int64, company, trackingNo string) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `status` = ?, `logistics_company` = ?, `tracking_no` = ?, `shipped_at` = now()"+
			" where `id` = ?",
		toStatus, company, trackingNo, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateOrderCompletedByUser marks an order as completed by the buyer.
func (m *customOrdersModel) UpdateOrderCompletedByUser(ctx context.Context, session sqlx.Session, id, userID, toStatus int64) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `status` = ?, `completed_at` = now() where `id` = ? and `user_id` = ?",
		toStatus, id, userID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateOrderCompletedByStatus marks an order as completed from a known status.
func (m *customOrdersModel) UpdateOrderCompletedByStatus(ctx context.Context, session sqlx.Session, id, fromStatus, toStatus int64) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `status` = ?, `completed_at` = now() where `id` = ? and `status` = ?",
		toStatus, id, fromStatus,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateOrderCanceled marks an order as canceled.
func (m *customOrdersModel) UpdateOrderCanceled(ctx context.Context, session sqlx.Session, id, toStatus int64, reason string) (int64, error) {
	result, err := session.ExecCtx(ctx,
		"update "+m.table+" set `status` = ?, `canceled_at` = now(), `cancel_reason` = ? where `id` = ?",
		toStatus, reason, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
