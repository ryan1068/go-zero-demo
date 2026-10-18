package logic

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"go-zero-demo/order-rpc/internal/model"
	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

	var resp *order.OrderResp
	unlocks := make([]func(), 0, 4)
	defer func() {
		for i := len(unlocks) - 1; i >= 0; i-- {
			unlocks[i]()
		}
	}()

	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		receiver, err := l.svcCtx.UserAddressModel.FindReceiver(ctx, session, in.AddressId, userID)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &order.OrderResp{Code: errcode.BadRequest, Msg: "收货地址不存在"}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}

		inputs, err := l.orderInputs(ctx, session, userID, in)
		if err != nil {
			return err
		}
		if len(inputs) == 0 {
			resp = &order.OrderResp{Code: errcode.BadRequest, Msg: "请选择要购买的商品"}
			return model.ErrAbortTransaction
		}

		items := make([]model.NewOrderItem, 0, len(inputs))
		var totalAmount int64
		for _, input := range inputs {
			unlocks = append(unlocks, l.lockProductStock(input.ProductId))

			snapshot, err := l.svcCtx.ProductModel.FindForUpdate(ctx, session, input.ProductId)
			if errors.Is(err, sqlx.ErrNotFound) {
				resp = &order.OrderResp{Code: errcode.NotFound, Msg: "商品不存在"}
				return model.ErrAbortTransaction
			}
			if err != nil {
				return err
			}
			if snapshot.Status != 1 {
				resp = &order.OrderResp{Code: errcode.BadRequest, Msg: "商品已下架"}
				return model.ErrAbortTransaction
			}
			if snapshot.Stock < input.Quantity {
				resp = &order.OrderResp{Code: errcode.BadRequest, Msg: "库存不足"}
				return model.ErrAbortTransaction
			}

			amount := snapshot.Price * input.Quantity
			totalAmount += amount
			items = append(items, model.NewOrderItem{
				ProductID:    input.ProductId,
				ProductName:  snapshot.Name,
				ProductCover: snapshot.Cover,
				Price:        snapshot.Price,
				Quantity:     input.Quantity,
				TotalAmount:  amount,
			})

			affected, err := l.svcCtx.ProductModel.DeductStock(ctx, session, input.ProductId, input.Quantity)
			if err != nil {
				return err
			}
			if affected == 0 {
				resp = &order.OrderResp{Code: errcode.BadRequest, Msg: "库存不足"}
				return model.ErrAbortTransaction
			}
		}

		orderID, err := l.svcCtx.OrdersModel.InsertOrder(ctx, session, model.NewOrder{
			OrderNo:            generateNo("O"),
			UserID:             userID,
			TotalAmount:        totalAmount,
			PayAmount:          totalAmount,
			Status:             orderStatusPendingPayment,
			ReceiverName:       receiver.Name,
			ReceiverPhone:      receiver.Phone,
			ReceiverProvince:   receiver.Province,
			ReceiverCity:       receiver.City,
			ReceiverDistrict:   receiver.District,
			ReceiverAddress:    receiver.Address,
			ReceiverPostalCode: receiver.Postal,
			Remark:             remark,
		})
		if err != nil {
			return err
		}
		if err := l.svcCtx.OrderItemModel.InsertOrderItems(ctx, session, orderID, items); err != nil {
			return err
		}
		if in.BuyNow == 0 {
			if err := l.svcCtx.CartItemModel.DeleteSelectedCartItems(ctx, session, userID); err != nil {
				return err
			}
		}
		if err := l.svcCtx.OrderLogModel.InsertOrderLog(ctx, session, model.NewOrderLog{
			OrderID:      orderID,
			FromStatus:   0,
			ToStatus:     orderStatusPendingPayment,
			OperatorType: operatorUser,
			OperatorID:   userID,
			Remark:       "创建订单",
		}); err != nil {
			return err
		}

		row, err := l.svcCtx.OrdersModel.FindOrder(ctx, session, orderID, userID, false)
		if err != nil {
			return err
		}
		orderItems, err := l.svcCtx.OrderItemModel.FindOrderItems(ctx, session, orderID)
		if err != nil {
			return err
		}
		resp = &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toOrderInfo(row, orderItems)}
		return nil
	})
	switch {
	case errors.Is(err, model.ErrAbortTransaction):
		return resp, nil
	case err != nil:
		l.Errorf("create order failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (l *CreateOrderLogic) orderInputs(ctx context.Context, session sqlx.Session, userID int64, in *order.CreateOrderReq) ([]*model.CartSelection, error) {
	if in.BuyNow == 1 {
		return []*model.CartSelection{{ProductId: in.ProductId, Quantity: int64(in.Quantity)}}, nil
	}
	inputs, err := l.svcCtx.CartItemModel.FindSelectedForUpdate(ctx, session, userID)
	if err != nil {
		l.Errorf("query selected cart items failed: %v", err)
		return nil, err
	}
	return inputs, nil
}
