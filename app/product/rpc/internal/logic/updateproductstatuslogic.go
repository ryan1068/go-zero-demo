package logic

import (
	"context"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductStatusLogic {
	return &UpdateProductStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateProductStatusLogic) UpdateProductStatus(req *product.UpdateProductStatusReq) (*product.BaseResp, error) {
	if req.Id <= 0 || (req.Status != 0 && req.Status != 1) {
		return &product.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	affected, err := l.svcCtx.ProductModel.UpdateProductStatus(l.ctx, uint64(req.Id), int64(req.Status))
	if err != nil {
		l.Errorf("update product status failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &product.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	return &product.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
