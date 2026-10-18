package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"go-zero-demo/order-rpc/internal/model"
	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ApplyAfterSaleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAfterSaleLogic {
	return &ApplyAfterSaleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ApplyAfterSaleLogic) ApplyAfterSale(in *order.ApplyAfterSaleReq) (*order.AfterSaleResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.AfterSaleResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	reason := strings.TrimSpace(in.Reason)
	description := strings.TrimSpace(in.Description)
	if in.OrderId <= 0 || (in.Type != 1 && in.Type != 2) || reason == "" ||
		!validAfterSaleText(reason, maxAfterSaleReason) || !validAfterSaleText(description, maxAfterSaleDescription) ||
		len(in.Images) > maxAfterSaleImages {
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	images := make([]string, 0, len(in.Images))
	for _, image := range in.Images {
		image = strings.TrimSpace(image)
		if image == "" {
			continue
		}
		if !validAfterSaleText(image, maxAfterSaleImageLength) {
			return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
		}
		images = append(images, image)
	}
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		l.Errorf("encode after sale images failed: %v", err)
		return nil, err
	}

	var resp *order.AfterSaleResp
	var afterSaleID int64
	err = l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		orderLock, err := l.svcCtx.OrdersModel.LockOrder(ctx, session, in.OrderId, nil)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}
		if orderLock.UserId != userID {
			resp = &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if orderLock.Status != orderStatusPendingShipment &&
			orderLock.Status != orderStatusShipped &&
			orderLock.Status != orderStatusCompleted {
			resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前订单状态不支持售后"}
			return model.ErrAbortTransaction
		}

		item, err := l.svcCtx.OrderItemModel.FindItemForAfterSale(ctx, session, in.OrderId, in.OrderItemId)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &order.AfterSaleResp{Code: errcode.NotFound, Msg: "订单明细不存在"}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}

		active, err := l.svcCtx.AfterSaleModel.CountActiveAfterSales(ctx, session, item.Id, userID,
			[]int64{afterSaleStatusRejected, afterSaleStatusRefunded, afterSaleStatusClosed})
		if err != nil {
			return err
		}
		if active > 0 {
			resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "该商品已有进行中的售后"}
			return model.ErrAbortTransaction
		}

		afterSaleID, err = l.svcCtx.AfterSaleModel.InsertAfterSale(ctx, session, model.NewAfterSale{
			AfterSaleNo: generateNo("AS"),
			OrderID:     in.OrderId,
			OrderItemID: item.Id,
			UserID:      userID,
			Type:        int64(in.Type),
			Reason:      reason,
			Description: description,
			ImagesJSON:  string(imagesJSON),
			Amount:      item.Price * item.Quantity,
			Status:      afterSaleStatusPending,
		})
		return err
	})
	switch {
	case errors.Is(err, model.ErrAbortTransaction):
		return resp, nil
	case err != nil:
		l.Errorf("apply after sale failed: %v", err)
		return nil, err
	}

	row, err := l.svcCtx.AfterSaleModel.FindAfterSale(l.ctx, nil, afterSaleID, userID, false)
	if err != nil {
		l.Errorf("query saved after sale failed: %v", err)
		return nil, err
	}
	return &order.AfterSaleResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toAfterSaleInfo(row)}, nil
}
