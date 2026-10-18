package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetProductLogic {
	return &AdminGetProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminGetProductLogic) AdminGetProduct(req *types.PathIdReq) (*types.ProductResp, error) {
	return NewGetProductLogic(l.ctx, l.svcCtx).getProduct(req.Id, true)
}
