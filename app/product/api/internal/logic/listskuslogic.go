package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSkusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSkusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSkusLogic {
	return &ListSkusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListSkusLogic) ListSkus(req *types.PathIdReq) (*types.SkuListResp, error) {
	return l.listSkus(req.Id)
}

func (l *ListSkusLogic) listSkus(productID int64) (*types.SkuListResp, error) {
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
