package logic

import (
	"context"
	"errors"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProductLogic) GetProduct(in *product.ProductReq) (*product.ProductResp, error) {
	if in.Id <= 0 {
		return &product.ProductResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	row, err := l.svcCtx.ProductModel.FindProductRow(l.ctx, uint64(in.Id), in.Admin == 1)
	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		return &product.ProductResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query product failed: %v", err)
		return nil, err
	}
	return &product.ProductResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toProductInfo(row)}, nil
}
