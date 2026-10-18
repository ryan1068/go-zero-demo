package logic

import (
	"context"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteProductLogic) DeleteProduct(in *product.ProductIdReq) (*product.BaseResp, error) {
	if in.Id <= 0 {
		return &product.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	affected, err := l.svcCtx.ProductModel.SoftDeleteProduct(l.ctx, uint64(in.Id))
	if err != nil {
		l.Errorf("delete product failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &product.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	return &product.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
