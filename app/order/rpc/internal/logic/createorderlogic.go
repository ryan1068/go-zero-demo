package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

const stockLockTTL = 10

const releaseStockLockScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
  return redis.call("del", KEYS[1])
else
  return 0
end
`

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// lockProductStock serializes stock deduction across order-rpc instances. The
// database row lock remains the final guard if Redis is unavailable or busy.
func (l *CreateOrderLogic) lockProductStock(productID int64) func() {
	if l.svcCtx.Redis == nil {
		return func() {}
	}
	key := fmt.Sprintf("lock:product:stock:%d", productID)
	token := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
	acquired := false
	for attempt := 0; attempt < 3; attempt++ {
		ok, err := l.svcCtx.Redis.SetnxExCtx(l.ctx, key, token, stockLockTTL)
		if err != nil {
			l.Errorf("acquire product stock lock %d failed: %v", productID, err)
			return func() {}
		}
		if ok {
			acquired = true
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if !acquired {
		l.Errorf("product stock lock %d busy, using database row lock", productID)
		return func() {}
	}
	return func() {
		if _, err := l.svcCtx.Redis.EvalCtx(l.ctx, releaseStockLockScript, []string{key}, token); err != nil {
			l.Errorf("release product stock lock %d failed: %v", productID, err)
		}
	}
}

type orderItemInput struct {
	productID int64
	quantity  int32
}

type orderItemSnapshot struct {
	productID int64
	name      string
	cover     string
	price     int64
	quantity  int32
	amount    int64
}

type receiverSnapshot struct {
	name     string
	phone    string
	province string
	city     string
	district string
	address  string
	postal   string
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderReq) (*order.OrderResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	remark := strings.TrimSpace(in.Remark)
	if in.AddressId <= 0 || utf8.RuneCountInString(remark) > maxRemarkLength {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	if in.BuyNow != 0 && in.BuyNow != 1 {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	if in.BuyNow == 1 {
		if in.ProductId <= 0 || in.Quantity < 1 || in.Quantity > maxQuantity {
			return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
		}
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin order transaction failed: %v", err)
		return nil, err
	}
	defer rollbackTx(l.Logger, tx)

	receiver, err := l.loadReceiver(l.ctx, tx, in.AddressId, userID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.OrderResp{Code: errcode.BadRequest, Msg: "收货地址不存在"}, nil
	case err != nil:
		l.Errorf("query receiver address failed: %v", err)
		return nil, err
	}

	inputs, err := l.orderInputs(l.ctx, tx, userID, in)
	if err != nil {
		return nil, err
	}
	if len(inputs) == 0 {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: "请选择要购买的商品"}, nil
	}

	items := make([]orderItemSnapshot, 0, len(inputs))
	var totalAmount int64
	for _, input := range inputs {
		var snap orderItemSnapshot
		var status int32
		var stock int32
		unlock := l.lockProductStock(input.productID)
		defer unlock()
		err := tx.QueryRowContext(l.ctx,
			"SELECT name, COALESCE(cover_image, ''), price, stock, status FROM looklook_product.product WHERE id = ? AND deleted_at IS NULL FOR UPDATE",
			input.productID,
		).Scan(&snap.name, &snap.cover, &snap.price, &stock, &status)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &order.OrderResp{Code: errcode.NotFound, Msg: "商品不存在"}, nil
		case err != nil:
			l.Errorf("query product failed: %v", err)
			return nil, err
		}
		if status != 1 {
			return &order.OrderResp{Code: errcode.BadRequest, Msg: "商品已下架"}, nil
		}
		if stock < input.quantity {
			return &order.OrderResp{Code: errcode.BadRequest, Msg: "库存不足"}, nil
		}

		snap.productID = input.productID
		snap.quantity = input.quantity
		snap.amount = snap.price * int64(input.quantity)
		totalAmount += snap.amount
		items = append(items, snap)

		result, err := tx.ExecContext(l.ctx,
			"UPDATE looklook_product.product SET stock = stock - ?, sales = sales + ? WHERE id = ? AND stock >= ?",
			input.quantity, input.quantity, input.productID, input.quantity,
		)
		if err != nil {
			l.Errorf("deduct product stock failed: %v", err)
			return nil, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			l.Errorf("read stock affected rows failed: %v", err)
			return nil, err
		}
		if affected == 0 {
			return &order.OrderResp{Code: errcode.BadRequest, Msg: "库存不足"}, nil
		}
	}

	orderNo := generateNo("O")
	result, err := tx.ExecContext(l.ctx,
		"INSERT INTO orders (order_no, user_id, total_amount, freight_amount, discount_amount, pay_amount, status, pay_type, receiver_name, receiver_phone, receiver_province, receiver_city, receiver_district, receiver_address, receiver_postal_code, remark) VALUES (?, ?, ?, 0, 0, ?, ?, 0, ?, ?, ?, ?, ?, ?, ?, ?)",
		orderNo, userID, totalAmount, totalAmount, orderStatusPendingPayment,
		receiver.name, receiver.phone, receiver.province, receiver.city, receiver.district, receiver.address, receiver.postal, remark,
	)
	if err != nil {
		l.Errorf("insert order failed: %v", err)
		return nil, err
	}
	orderID, err := result.LastInsertId()
	if err != nil {
		l.Errorf("read inserted order id failed: %v", err)
		return nil, err
	}

	for _, item := range items {
		if _, err := tx.ExecContext(l.ctx,
			"INSERT INTO order_item (order_id, product_id, product_name, product_cover, price, quantity, total_amount) VALUES (?, ?, ?, ?, ?, ?, ?)",
			orderID, item.productID, item.name, item.cover, item.price, item.quantity, item.amount,
		); err != nil {
			l.Errorf("insert order item failed: %v", err)
			return nil, err
		}
	}

	if in.BuyNow == 0 {
		if _, err := tx.ExecContext(l.ctx,
			"DELETE FROM cart_item WHERE user_id = ? AND selected = 1", userID,
		); err != nil {
			l.Errorf("delete selected cart items failed: %v", err)
			return nil, err
		}
	}

	if err := writeOrderLog(l.ctx, tx, orderID, 0, orderStatusPendingPayment, operatorUser, userID, "创建订单"); err != nil {
		l.Errorf("write order log failed: %v", err)
		return nil, err
	}

	info, err := loadOrder(l.ctx, tx, orderID, userID, false)
	if err != nil {
		l.Errorf("query created order failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit order transaction failed: %v", err)
		return nil, err
	}

	return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

func (l *CreateOrderLogic) loadReceiver(ctx context.Context, tx *sql.Tx, addressID, userID int64) (*receiverSnapshot, error) {
	receiver := new(receiverSnapshot)
	err := tx.QueryRowContext(ctx,
		"SELECT receiver_name, receiver_phone, COALESCE(province, ''), COALESCE(city, ''), COALESCE(district, ''), detail_address, COALESCE(postal_code, '') FROM looklook_usercenter.user_address WHERE id = ? AND user_id = ? AND deleted_at IS NULL",
		addressID, userID,
	).Scan(&receiver.name, &receiver.phone, &receiver.province, &receiver.city, &receiver.district, &receiver.address, &receiver.postal)
	if err != nil {
		return nil, err
	}
	return receiver, nil
}

func (l *CreateOrderLogic) orderInputs(ctx context.Context, tx *sql.Tx, userID int64, in *order.CreateOrderReq) ([]orderItemInput, error) {
	if in.BuyNow == 1 {
		return []orderItemInput{{productID: in.ProductId, quantity: in.Quantity}}, nil
	}

	rows, err := tx.QueryContext(ctx,
		"SELECT product_id, quantity FROM cart_item WHERE user_id = ? AND selected = 1 FOR UPDATE", userID,
	)
	if err != nil {
		l.Errorf("query selected cart items failed: %v", err)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	inputs := make([]orderItemInput, 0, 4)
	for rows.Next() {
		var input orderItemInput
		if err := rows.Scan(&input.productID, &input.quantity); err != nil {
			l.Errorf("scan selected cart item failed: %v", err)
			return nil, err
		}
		inputs = append(inputs, input)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate selected cart items failed: %v", err)
		return nil, err
	}
	return inputs, nil
}
