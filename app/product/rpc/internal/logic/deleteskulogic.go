package logic

import (
	"context"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSkuLogic {
	return &DeleteSkuLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteSkuLogic) DeleteSku(in *product.SkuIdReq) (*product.BaseResp, error) {
	if in.Id <= 0 {
		return &product.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	affected, err := l.svcCtx.ProductSkuModel.SoftDeleteSku(l.ctx, uint64(in.Id))
	if err != nil {
		l.Errorf("delete sku failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &product.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	return &product.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
