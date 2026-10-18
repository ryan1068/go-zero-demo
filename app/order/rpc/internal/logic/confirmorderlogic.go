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

type ConfirmOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmOrderLogic {
	return &ConfirmOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConfirmOrderLogic) ConfirmOrder(in *order.ConfirmOrderReq) (*order.OrderResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin confirm transaction failed: %v", err)
		return nil, err
	}
	defer rollbackTx(l.Logger, tx)

	var status int32
	err = tx.QueryRowContext(l.ctx,
		"SELECT status FROM orders WHERE id = ? AND user_id = ? AND deleted_at IS NULL FOR UPDATE",
		in.Id, userID,
	).Scan(&status)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.OrderResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query order failed: %v", err)
		return nil, err
	}

	if status == orderStatusCompleted {
		info, err := loadOrder(l.ctx, tx, in.Id, userID, false)
		if err != nil {
			l.Errorf("query completed order failed: %v", err)
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			l.Errorf("commit confirm transaction failed: %v", err)
			return nil, err
		}
		return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
	}
	if status != orderStatusShipped {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: "订单未发货，不能确认收货"}, nil
	}

	if _, err := tx.ExecContext(l.ctx,
		"UPDATE orders SET status = ?, completed_at = NOW() WHERE id = ? AND user_id = ?",
		orderStatusCompleted, in.Id, userID,
	); err != nil {
		l.Errorf("update order completed status failed: %v", err)
		return nil, err
	}
	if err := writeOrderLog(l.ctx, tx, in.Id, orderStatusShipped, orderStatusCompleted, operatorUser, userID, "确认收货"); err != nil {
		l.Errorf("write order log failed: %v", err)
		return nil, err
	}

	info, err := loadOrder(l.ctx, tx, in.Id, userID, false)
	if err != nil {
		l.Errorf("query completed order failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit confirm transaction failed: %v", err)
		return nil, err
	}
	return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}
