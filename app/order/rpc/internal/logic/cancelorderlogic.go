package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"unicode/utf8"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CancelOrderLogic) CancelOrder(in *order.CancelOrderReq) (*order.OrderResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	reason := strings.TrimSpace(in.Reason)
	if utf8.RuneCountInString(reason) > maxRemarkLength {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	admin := in.Admin == 1
	if reason == "" {
		if admin {
			reason = "管理员取消"
		} else {
			reason = "用户取消"
		}
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin cancel transaction failed: %v", err)
		return nil, err
	}
	defer rollbackTx(l.Logger, tx)

	query := "SELECT status FROM orders WHERE id = ? AND deleted_at IS NULL FOR UPDATE"
	args := []any{in.Id}
	if !admin {
		query = "SELECT status FROM orders WHERE id = ? AND user_id = ? AND deleted_at IS NULL FOR UPDATE"
		args = append(args, userID)
	}
	var status int32
	err = tx.QueryRowContext(l.ctx, query, args...).Scan(&status)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.OrderResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query order failed: %v", err)
		return nil, err
	}

	if status == orderStatusCanceled {
		info, err := loadOrder(l.ctx, tx, in.Id, userID, admin)
		if err != nil {
			l.Errorf("query canceled order failed: %v", err)
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			l.Errorf("commit cancel transaction failed: %v", err)
			return nil, err
		}
		return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
	}
	if status != orderStatusPendingPayment && status != orderStatusPendingShipment {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: "订单已发货，不能取消"}, nil
	}

	if _, err := tx.ExecContext(l.ctx,
		"UPDATE orders SET status = ?, canceled_at = NOW(), cancel_reason = ? WHERE id = ?",
		orderStatusCanceled, reason, in.Id,
	); err != nil {
		l.Errorf("update order cancel status failed: %v", err)
		return nil, err
	}
	if err := restoreStock(l.ctx, tx, in.Id); err != nil {
		l.Errorf("restore product stock failed: %v", err)
		return nil, err
	}
	operatorType := int32(operatorUser)
	if admin {
		operatorType = operatorAdmin
	}
	if err := writeOrderLog(l.ctx, tx, in.Id, status, orderStatusCanceled, operatorType, userID, reason); err != nil {
		l.Errorf("write order log failed: %v", err)
		return nil, err
	}

	info, err := loadOrder(l.ctx, tx, in.Id, userID, admin)
	if err != nil {
		l.Errorf("query canceled order failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit cancel transaction failed: %v", err)
		return nil, err
	}
	return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}
