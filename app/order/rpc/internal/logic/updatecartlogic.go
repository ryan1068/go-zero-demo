package logic

import (
	"context"
	"errors"

	"go-zero-demo/order-rpc/internal/model"
	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartLogic {
	return &UpdateCartLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateCartLogic) UpdateCart(in *order.UpdateCartReq) (*order.CartItemResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.CartItemResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 || in.Quantity < 1 || in.Quantity > maxQuantity || (in.Selected != 0 && in.Selected != 1) {
		return &order.CartItemResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var resp *order.CartItemResp
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		lock, err := l.svcCtx.CartItemModel.LockCartProduct(ctx, session, in.Id, userID)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &order.CartItemResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}
		if lock.Status != 1 {
			resp = &order.CartItemResp{Code: errcode.BadRequest, Msg: "商品已下架"}
			return model.ErrAbortTransaction
		}
		if int64(in.Quantity) > lock.Stock {
			resp = &order.CartItemResp{Code: errcode.BadRequest, Msg: "库存不足"}
			return model.ErrAbortTransaction
		}
		if _, err := l.svcCtx.CartItemModel.UpdateCartItem(ctx, session, in.Id, userID, int64(in.Quantity), int64(in.Selected)); err != nil {
			return err
		}

		row, err := l.svcCtx.CartItemModel.FindCartItemRow(ctx, session, in.Id, userID)
		if err != nil {
			return err
		}
		resp = &order.CartItemResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toCartItemInfo(row)}
		return nil
	})
	switch {
	case errors.Is(err, model.ErrAbortTransaction):
		return resp, nil
	case err != nil:
		l.Errorf("update cart item failed: %v", err)
		return nil, err
	}
	return resp, nil
}
