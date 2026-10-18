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

type AddCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCartLogic {
	return &AddCartLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddCartLogic) AddCart(in *order.AddCartReq) (*order.CartItemResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.CartItemResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.ProductId <= 0 || in.Quantity < 1 || in.Quantity > maxQuantity {
		return &order.CartItemResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var resp *order.CartItemResp
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		product, err := l.svcCtx.ProductModel.FindForUpdate(ctx, session, in.ProductId)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &order.CartItemResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}
		if product.Status != 1 {
			resp = &order.CartItemResp{Code: errcode.BadRequest, Msg: "商品已下架"}
			return model.ErrAbortTransaction
		}

		var cartID int64
		existing, err := l.svcCtx.CartItemModel.FindForUpdate(ctx, session, userID, in.ProductId)
		switch {
		case errors.Is(err, sqlx.ErrNotFound):
			if int64(in.Quantity) > product.Stock {
				resp = &order.CartItemResp{Code: errcode.BadRequest, Msg: "库存不足"}
				return model.ErrAbortTransaction
			}
			cartID, err = l.svcCtx.CartItemModel.InsertCartItem(ctx, session, model.NewCartItem{
				UserID:    userID,
				ProductID: in.ProductId,
				Quantity:  int64(in.Quantity),
			})
			if err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			quantity := existing.Quantity + int64(in.Quantity)
			if quantity > maxQuantity {
				quantity = maxQuantity
			}
			if quantity > product.Stock {
				resp = &order.CartItemResp{Code: errcode.BadRequest, Msg: "库存不足"}
				return model.ErrAbortTransaction
			}
			cartID = existing.Id
			if _, err := l.svcCtx.CartItemModel.UpdateCartItem(ctx, session, cartID, userID, quantity, 1); err != nil {
				return err
			}
		}

		row, err := l.svcCtx.CartItemModel.FindCartItemRow(ctx, session, cartID, userID)
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
		l.Errorf("add cart item failed: %v", err)
		return nil, err
	}
	return resp, nil
}
