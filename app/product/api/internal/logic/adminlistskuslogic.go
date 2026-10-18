package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSkusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListSkusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSkusLogic {
	return &AdminListSkusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminListSkusLogic) AdminListSkus(req *types.PathIdReq) (*types.SkuListResp, error) {
	return NewListSkusLogic(l.ctx, l.svcCtx).listSkus(req.Id)
}
