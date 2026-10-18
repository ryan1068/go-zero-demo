package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProductLogic) GetProduct(req *types.PathIdReq) (*types.ProductResp, error) {
	return l.getProduct(req.Id, false)
}

func (l *GetProductLogic) getProduct(id int64, admin bool) (*types.ProductResp, error) {
	s, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.ProductReq{Id: id, Admin: adminFlag(admin)})
	if err != nil {
		l.Errorf("get product rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ProductResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toProductInfo(s.Data)
	}
	return resp, nil
}
