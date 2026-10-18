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

type ShipOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShipOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipOrderLogic {
	return &ShipOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ShipOrderLogic) ShipOrder(in *order.ShipOrderReq) (*order.OrderResp, error) {
	adminID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.OrderResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	company := strings.TrimSpace(in.LogisticsCompany)
	trackingNo := strings.TrimSpace(in.TrackingNo)
	if in.Id <= 0 || company == "" || !trackNoValid(company, trackingNo) {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin ship transaction failed: %v", err)
		return nil, err
	}
	defer rollbackTx(l.Logger, tx)

	var status int32
	err = tx.QueryRowContext(l.ctx,
		"SELECT status FROM orders WHERE id = ? AND deleted_at IS NULL FOR UPDATE", in.Id,
	).Scan(&status)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.OrderResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query order failed: %v", err)
		return nil, err
	}

	if status == orderStatusShipped {
		info, err := loadOrder(l.ctx, tx, in.Id, 0, true)
		if err != nil {
			l.Errorf("query shipped order failed: %v", err)
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			l.Errorf("commit ship transaction failed: %v", err)
			return nil, err
		}
		return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
	}
	if status != orderStatusPendingShipment {
		return &order.OrderResp{Code: errcode.BadRequest, Msg: "当前订单状态不能发货"}, nil
	}

	if _, err := tx.ExecContext(l.ctx,
		"UPDATE orders SET status = ?, logistics_company = ?, tracking_no = ?, shipped_at = NOW() WHERE id = ?",
		orderStatusShipped, company, trackingNo, in.Id,
	); err != nil {
		l.Errorf("update order ship status failed: %v", err)
		return nil, err
	}
	if err := writeOrderLog(l.ctx, tx, in.Id, orderStatusPendingShipment, orderStatusShipped, operatorAdmin, adminID, "管理员发货"); err != nil {
		l.Errorf("write order log failed: %v", err)
		return nil, err
	}

	info, err := loadOrder(l.ctx, tx, in.Id, 0, true)
	if err != nil {
		l.Errorf("query shipped order failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit ship transaction failed: %v", err)
		return nil, err
	}
	return &order.OrderResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

func trackNoValid(company, trackingNo string) bool {
	return trackingNo != "" &&
		utf8.RuneCountInString(company) <= maxLogisticsCompany &&
		utf8.RuneCountInString(trackingNo) <= maxTrackingNo
}
