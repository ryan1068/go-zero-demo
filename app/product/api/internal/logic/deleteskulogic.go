package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSkuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSkuLogic {
	return &DeleteSkuLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteSkuLogic) DeleteSku(req *types.SkuIdReq) (*types.BaseResp, error) {
	s, err := l.svcCtx.ProductRpc.DeleteSku(l.ctx, &product.SkuIdReq{Id: req.Id})
	if err != nil {
		l.Errorf("delete sku rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}
