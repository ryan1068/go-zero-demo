package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

func toSkuInfo(info *product.SkuInfo) types.SkuInfo {
	return types.SkuInfo{
		Id:            info.Id,
		ProductId:     info.ProductId,
		SkuCode:       info.SkuCode,
		Specs:         info.Specs,
		Price:         info.Price,
		OriginalPrice: info.OriginalPrice,
		Stock:         info.Stock,
		Sales:         info.Sales,
		Status:        info.Status,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.UpdatedAt,
	}
}

type ListSkusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSkusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSkusLogic {
	return &ListSkusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListSkusLogic) ListSkus(productID int64) (*types.SkuListResp, error) {
	s, err := l.svcCtx.ProductRpc.ListSkus(l.ctx, &product.ProductIdReq{Id: productID})
	if err != nil {
		l.Errorf("list skus rpc failed: %v", err)
		return nil, err
	}
	list := make([]types.SkuInfo, 0, len(s.List))
	for _, item := range s.List {
		list = append(list, toSkuInfo(item))
	}
	return &types.SkuListResp{Code: s.Code, Msg: s.Msg, List: list}, nil
}

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

type DeleteSkuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSkuLogic {
	return &DeleteSkuLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteSkuLogic) DeleteSku(id int64) (*types.BaseResp, error) {
	s, err := l.svcCtx.ProductRpc.DeleteSku(l.ctx, &product.SkuIdReq{Id: id})
	if err != nil {
		l.Errorf("delete sku rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
