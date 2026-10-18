package logic

import (
	"context"
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

type HandleAfterSaleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleAfterSaleLogic {
	return &HandleAfterSaleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *HandleAfterSaleLogic) HandleAfterSale(in *order.HandleAfterSaleReq) (*order.AfterSaleResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.AfterSaleResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	action := strings.ToLower(strings.TrimSpace(in.Action))
	adminRemark := strings.TrimSpace(in.AdminRemark)
	trackingNo := strings.TrimSpace(in.TrackingNo)
	if in.Id <= 0 || !validAfterSaleText(adminRemark, maxAfterSaleRemark) || !validAfterSaleText(trackingNo, maxAfterSaleTrackingNo) {
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var resp *order.AfterSaleResp
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		lock, err := l.svcCtx.AfterSaleModel.LockAfterSale(ctx, session, in.Id)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}
		if in.Admin != 1 && lock.UserId != userID {
			resp = &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if in.Admin != 1 && action != "ship" {
			resp = &order.AfterSaleResp{Code: errcode.Unauthorized, Msg: "需要管理员权限"}
			return model.ErrAbortTransaction
		}
		if in.Admin == 1 && action == "ship" {
			resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "寄回操作由买家发起"}
			return model.ErrAbortTransaction
		}

		switch action {
		case "approve":
			if lock.Status != afterSaleStatusPending {
				resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可审核"}
				return model.ErrAbortTransaction
			}
			if _, err := l.svcCtx.AfterSaleModel.MarkAfterSaleReviewed(ctx, session, in.Id, afterSaleStatusApproved, adminRemark); err != nil {
				return err
			}
		case "reject":
			if lock.Status != afterSaleStatusPending {
				resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可审核"}
				return model.ErrAbortTransaction
			}
			if _, err := l.svcCtx.AfterSaleModel.MarkAfterSaleReviewed(ctx, session, in.Id, afterSaleStatusRejected, adminRemark); err != nil {
				return err
			}
		case "ship":
			if lock.Status != afterSaleStatusApproved || lock.Type != 2 {
				resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可寄回"}
				return model.ErrAbortTransaction
			}
			if trackingNo == "" {
				resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "请填写退回物流单号"}
				return model.ErrAbortTransaction
			}
			if _, err := l.svcCtx.AfterSaleModel.MarkAfterSaleBuyerShipped(ctx, session, in.Id, afterSaleStatusBuyerShipped, trackingNo); err != nil {
				return err
			}
		case "receive":
			if lock.Status != afterSaleStatusBuyerShipped {
				resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可收货"}
				return model.ErrAbortTransaction
			}
			if _, err := l.svcCtx.AfterSaleModel.MarkAfterSaleSellerReceived(ctx, session, in.Id, afterSaleStatusSellerReceived, adminRemark); err != nil {
				return err
			}
		case "refund":
			if lock.Status != afterSaleStatusApproved && lock.Status != afterSaleStatusSellerReceived {
				resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "当前状态不可退款"}
				return model.ErrAbortTransaction
			}
			if _, err := l.svcCtx.AfterSaleModel.MarkAfterSaleRefunded(ctx, session, in.Id, afterSaleStatusRefunded, adminRemark); err != nil {
				return err
			}
			if err := l.svcCtx.RefundRecordModel.InsertRefundRecord(ctx, session, model.NewRefundRecord{
				RefundNo:    generateNo("R"),
				AfterSaleID: in.Id,
				OrderID:     lock.OrderId,
				UserID:      lock.UserId,
				Amount:      lock.Amount,
				Channel:     "mock",
				Status:      1,
			}); err != nil {
				return err
			}
		default:
			resp = &order.AfterSaleResp{Code: errcode.BadRequest, Msg: "不支持的售后操作"}
			return model.ErrAbortTransaction
		}
		return nil
	})
	switch {
	case errors.Is(err, model.ErrAbortTransaction):
		return resp, nil
	case err != nil:
		l.Errorf("handle after sale failed: %v", err)
		return nil, err
	}

	row, err := l.svcCtx.AfterSaleModel.FindAfterSale(l.ctx, nil, in.Id, userID, true)
	if err != nil {
		l.Errorf("query handled after sale failed: %v", err)
		return nil, err
	}
	return &order.AfterSaleResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toAfterSaleInfo(row)}, nil
}
