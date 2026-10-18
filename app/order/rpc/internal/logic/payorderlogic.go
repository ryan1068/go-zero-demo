package logic

import (
	"context"
	"database/sql"
	"errors"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *PayOrderLogic) PayOrder(in *order.PayOrderReq) (*order.OrderResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin payment transaction failed: %v", err)
		return nil, err
	}
	defer rollbackTx(l.Logger, tx)

	var (
		status    int32
		orderNo   string
		payAmount int64
	)
	err = tx.QueryRowContext(l.ctx,
		"SELECT status, order_no, pay_amount FROM orders WHERE id = ? AND user_id = ? AND deleted_at IS NULL FOR UPDATE",
		in.Id, userID,
	).Scan(&status, &orderNo, &payAmount)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.OrderResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query order failed: %v", err)
		return nil, err
	}

	switch status {
	case orderStatusPendingShipment, orderStatusShipped, orderStatusCompleted:
		info, err := loadOrder(l.ctx, tx, in.Id, userID, false)
		if err != nil {
			l.Errorf("query paid order failed: %v", err)
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			l.Errorf("commit payment transaction failed: %v", err)
			return nil, err
		}
		return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
	case orderStatusCanceled:
		return &order.OrderResp{Code: errcode.BadRequest, Msg: "订单已取消"}, nil
	case orderStatusPendingPayment:
		// Continue below.
	default:
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	paymentNo := generateNo("P")
	transactionID := "MOCK-" + generateNo("")
	if _, err := tx.ExecContext(l.ctx,
		"UPDATE orders SET status = ?, pay_type = ?, transaction_id = ?, paid_at = NOW() WHERE id = ? AND user_id = ? AND status = ?",
		orderStatusPendingShipment, payTypeMock, transactionID, in.Id, userID, orderStatusPendingPayment,
	); err != nil {
		l.Errorf("update order payment status failed: %v", err)
		return nil, err
	}
	if _, err := tx.ExecContext(l.ctx,
		"INSERT INTO payment_record (payment_no, order_id, order_no, user_id, amount, pay_type, status, transaction_id, paid_at) VALUES (?, ?, ?, ?, ?, ?, 1, ?, NOW())",
		paymentNo, in.Id, orderNo, userID, payAmount, payTypeMock, transactionID,
	); err != nil {
		l.Errorf("insert payment record failed: %v", err)
		return nil, err
	}
	if err := writeOrderLog(l.ctx, tx, in.Id, orderStatusPendingPayment, orderStatusPendingShipment, operatorUser, userID, "模拟支付成功"); err != nil {
		l.Errorf("write order log failed: %v", err)
		return nil, err
	}

	info, err := loadOrder(l.ctx, tx, in.Id, userID, false)
	if err != nil {
		l.Errorf("query paid order failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit payment transaction failed: %v", err)
		return nil, err
	}
	return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}
