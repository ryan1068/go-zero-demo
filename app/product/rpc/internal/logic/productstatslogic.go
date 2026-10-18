package logic

import (
	"context"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductStatsLogic {
	return &ProductStatsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ProductStatsLogic) ProductStats(_ *product.EmptyReq) (*product.ProductStatsResp, error) {
	stats, err := l.svcCtx.ProductModel.Stats(l.ctx, lowStockValue)
	if err != nil {
		l.Errorf("query product stats failed: %v", err)
		return nil, err
	}
	return &product.ProductStatsResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		Data: &product.ProductStats{
			Total:    stats.Total,
			OnSale:   stats.OnSale,
			OffSale:  stats.OffSale,
			LowStock: stats.LowStock,
		},
	}, nil
}
