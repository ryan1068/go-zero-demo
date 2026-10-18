package logic

import (
	"context"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSkusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSkusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSkusLogic {
	return &ListSkusLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListSkusLogic) ListSkus(in *product.ProductIdReq) (*product.SkuListResp, error) {
	if in.Id <= 0 {
		return &product.SkuListResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	rows, err := l.svcCtx.ProductSkuModel.ListByProduct(l.ctx, in.Id)
	if err != nil {
		l.Errorf("query skus failed: %v", err)
		return nil, err
	}
	list := make([]*product.SkuInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, toSkuInfo(row))
	}
	return &product.SkuListResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), List: list}, nil
}
