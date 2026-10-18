package logic

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"sync/atomic"
	"time"

	"go-zero-demo/order-rpc/pb/order"
)

const (
	orderStatusPendingPayment  = 10
	orderStatusPendingShipment = 20
	orderStatusShipped         = 30
	orderStatusCompleted       = 40
	orderStatusCanceled        = 50

	payTypeMock = 1

	operatorUser   = 1
	operatorAdmin  = 2
	operatorSystem = 3

	defaultPageSize = 10
	maxPageSize     = 100
	maxOrderCount   = 10000

	maxQuantity          = 99
	maxRemarkLength      = 255
	maxReceiverName      = 64
	maxReceiverPhone     = 32
	maxReceiverRegion    = 64
	maxReceiverAddress   = 255
	maxReceiverPostal    = 16
	maxLogisticsCompany  = 64
	maxTrackingNo        = 64
	lowStockThreshold    = 10
	timeoutCancelMinutes = 30
)

const orderSelectColumns = "id, order_no, user_id, total_amount, freight_amount, discount_amount, pay_amount, status, pay_type, COALESCE(transaction_id, ''), receiver_name, receiver_phone, COALESCE(receiver_province, ''), COALESCE(receiver_city, ''), COALESCE(receiver_district, ''), receiver_address, COALESCE(receiver_postal_code, ''), COALESCE(remark, ''), COALESCE(admin_remark, ''), COALESCE(logistics_company, ''), COALESCE(tracking_no, ''), COALESCE(UNIX_TIMESTAMP(paid_at), 0), COALESCE(UNIX_TIMESTAMP(shipped_at), 0), COALESCE(UNIX_TIMESTAMP(completed_at), 0), COALESCE(UNIX_TIMESTAMP(canceled_at), 0), COALESCE(cancel_reason, ''), UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at)"

const cartSelectColumns = "ci.id, ci.product_id, COALESCE(p.name, ''), COALESCE(p.cover_image, ''), COALESCE(p.price, 0), COALESCE(p.stock, 0), ci.quantity, ci.selected, COALESCE(p.status, 0), UNIX_TIMESTAMP(ci.created_at), UNIX_TIMESTAMP(ci.updated_at)"

const cartFromClause = " FROM cart_item ci LEFT JOIN looklook_product.product p ON p.id = ci.product_id AND p.deleted_at IS NULL "

type rowScanner interface {
	Scan(dest ...any) error
}

type queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func scanOrderHeader(row rowScanner) (*order.OrderInfo, error) {
	info := new(order.OrderInfo)
	err := row.Scan(
		&info.Id,
		&info.OrderNo,
		&info.UserId,
		&info.TotalAmount,
		&info.FreightAmount,
		&info.DiscountAmount,
		&info.PayAmount,
		&info.Status,
		&info.PayType,
		&info.TransactionId,
		&info.ReceiverName,
		&info.ReceiverPhone,
		&info.ReceiverProvince,
		&info.ReceiverCity,
		&info.ReceiverDistrict,
		&info.ReceiverAddress,
		&info.ReceiverPostalCode,
		&info.Remark,
		&info.AdminRemark,
		&info.LogisticsCompany,
		&info.TrackingNo,
		&info.PaidAt,
		&info.ShippedAt,
		&info.CompletedAt,
		&info.CanceledAt,
		&info.CancelReason,
		&info.CreatedAt,
		&info.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func loadOrder(ctx context.Context, q queryer, id, userID int64, admin bool) (*order.OrderInfo, error) {
	query := "SELECT " + orderSelectColumns + " FROM orders WHERE id = ? AND deleted_at IS NULL"
	args := []any{id}
	if !admin {
		query += " AND user_id = ?"
		args = append(args, userID)
	}

	info, err := scanOrderHeader(q.QueryRowContext(ctx, query, args...))
	if err != nil {
		return nil, err
	}

	items, err := loadOrderItems(ctx, q, id)
	if err != nil {
		return nil, err
	}
	info.Items = items
	return info, nil
}

func loadOrderItems(ctx context.Context, q queryer, orderID int64) ([]*order.OrderItemInfo, error) {
	rows, err := q.QueryContext(ctx,
		"SELECT id, product_id, product_name, product_cover, price, quantity, total_amount FROM order_item WHERE order_id = ? ORDER BY id",
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	items := make([]*order.OrderItemInfo, 0, 4)
	for rows.Next() {
		item := new(order.OrderItemInfo)
		if err := rows.Scan(
			&item.Id,
			&item.ProductId,
			&item.ProductName,
			&item.ProductCover,
			&item.Price,
			&item.Quantity,
			&item.TotalAmount,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func scanCartItem(row rowScanner) (*order.CartItemInfo, error) {
	info := new(order.CartItemInfo)
	err := row.Scan(
		&info.Id,
		&info.ProductId,
		&info.ProductName,
		&info.ProductCover,
		&info.Price,
		&info.Stock,
		&info.Quantity,
		&info.Selected,
		&info.ProductStatus,
		&info.CreatedAt,
		&info.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func findCartItem(ctx context.Context, q queryer, id, userID int64) (*order.CartItemInfo, error) {
	return scanCartItem(q.QueryRowContext(ctx,
		"SELECT "+cartSelectColumns+cartFromClause+" WHERE ci.id = ? AND ci.user_id = ?",
		id, userID,
	))
}

func writeOrderLog(ctx context.Context, tx *sql.Tx, orderID int64, from, to, operatorType int32, operatorID int64, remark string) error {
	_, err := tx.ExecContext(ctx,
		"INSERT INTO order_log (order_id, from_status, to_status, operator_type, operator_id, remark) VALUES (?, ?, ?, ?, ?, ?)",
		orderID, from, to, operatorType, operatorID, remark,
	)
	return err
}

func restoreStock(ctx context.Context, tx *sql.Tx, orderID int64) error {
	rows, err := tx.QueryContext(ctx,
		"SELECT product_id, quantity FROM order_item WHERE order_id = ?", orderID,
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = rows.Close()
	}()

	type item struct {
		productID int64
		quantity  int32
	}
	items := make([]item, 0, 4)
	for rows.Next() {
		var it item
		if err := rows.Scan(&it.productID, &it.quantity); err != nil {
			return err
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, it := range items {
		if _, err := tx.ExecContext(ctx,
			"UPDATE looklook_product.product SET stock = stock + ?, sales = GREATEST(sales - ?, 0) WHERE id = ?",
			it.quantity, it.quantity, it.productID,
		); err != nil {
			return err
		}
	}
	return nil
}

var noSequence uint64

func generateNo(prefix string) string {
	seq := atomic.AddUint64(&noSequence, 1)
	return fmt.Sprintf("%s%d%05d%03d", prefix, time.Now().UnixNano(), seq%100000, rand.IntN(1000))
}
