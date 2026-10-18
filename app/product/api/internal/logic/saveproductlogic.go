package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveProductLogic {
	return &SaveProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveProductLogic) SaveProduct(req *types.SaveProductReq) (*types.ProductResp, error) {
	s, err := l.svcCtx.ProductRpc.SaveProduct(l.ctx, &product.SaveProductReq{
		Id:            req.Id,
		CategoryId:    req.CategoryId,
		Name:          req.Name,
		Subtitle:      req.Subtitle,
		CoverImage:    req.CoverImage,
		Images:        req.Images,
		Detail:        req.Detail,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		Unit:          req.Unit,
		Status:        req.Status,
		IsRecommend:   req.IsRecommend,
		Sort:          req.Sort,
	})
	if err != nil {
		l.Errorf("save product rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ProductResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toProductInfo(s.Data)
	}
	return resp, nil
}
