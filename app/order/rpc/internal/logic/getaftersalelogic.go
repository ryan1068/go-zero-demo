package logic

import (
	"context"
	"errors"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetAfterSaleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAfterSaleLogic {
	return &GetAfterSaleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetAfterSaleLogic) GetAfterSale(in *order.AfterSaleReq) (*order.AfterSaleResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.AfterSaleResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &order.AfterSaleResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	row, err := l.svcCtx.AfterSaleModel.FindAfterSale(l.ctx, nil, in.Id, userID, in.Admin == 1)
	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		return &order.AfterSaleResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query after sale failed: %v", err)
		return nil, err
	}
	return &order.AfterSaleResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toAfterSaleInfo(row)}, nil
}
