package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
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

func (l *ProductStatsLogic) ProductStats() (*types.ProductStatsResp, error) {
	s, err := l.svcCtx.ProductRpc.ProductStats(l.ctx, &product.EmptyReq{})
	if err != nil {
		l.Errorf("product stats rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ProductStatsResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = types.ProductStats{
			Total:    s.Data.Total,
			OnSale:   s.Data.OnSale,
			OffSale:  s.Data.OffSale,
			LowStock: s.Data.LowStock,
		}
	}
	return resp, nil
}
