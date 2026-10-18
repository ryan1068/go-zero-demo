package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveSkuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveSkuLogic {
	return &SaveSkuLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveSkuLogic) SaveSku(req *types.SaveSkuReq) (*types.SkuResp, error) {
	s, err := l.svcCtx.ProductRpc.SaveSku(l.ctx, &product.SaveSkuReq{
		Id:            req.Id,
		ProductId:     req.ProductId,
		SkuCode:       req.SkuCode,
		Specs:         req.Specs,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		Status:        req.Status,
	})
	if err != nil {
		l.Errorf("save sku rpc failed: %v", err)
		return nil, err
	}
	resp := &types.SkuResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toSkuInfo(s.Data)
	}
	return resp, nil
}
