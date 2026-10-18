package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/model"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SaveSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveSkuLogic {
	return &SaveSkuLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SaveSkuLogic) SaveSku(in *product.SaveSkuReq) (*product.SkuResp, error) {
	skuCode := strings.TrimSpace(in.SkuCode)
	specs := strings.TrimSpace(in.Specs)
	if specs == "" {
		specs = "{}"
	}
	if in.ProductId <= 0 || skuCode == "" || !validText(skuCode, maxSkuCode) ||
		!validText(specs, maxSkuSpecs) || in.Price <= 0 || in.OriginalPrice < 0 ||
		in.Stock < 0 || (in.Status != 0 && in.Status != 1) || !json.Valid([]byte(specs)) {
		return &product.SkuResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	exists, err := l.svcCtx.ProductModel.Exists(l.ctx, uint64(in.ProductId))
	if err != nil {
		l.Errorf("query product for sku failed: %v", err)
		return nil, err
	}
	if !exists {
		return &product.SkuResp{Code: errcode.BadRequest, Msg: "商品不存在"}, nil
	}

	data := &model.ProductSku{
		ProductId:     in.ProductId,
		SkuCode:       skuCode,
		Specs:         specs,
		Price:         in.Price,
		OriginalPrice: in.OriginalPrice,
		Stock:         int64(in.Stock),
		Status:        int64(in.Status),
	}

	var id uint64
	if in.Id > 0 {
		data.Id = uint64(in.Id)
		affected, err := l.svcCtx.ProductSkuModel.UpdateSku(l.ctx, data)
		if err != nil {
			l.Errorf("update sku failed: %v", err)
			return nil, err
		}
		if affected == 0 {
			if _, err := l.svcCtx.ProductSkuModel.FindActive(l.ctx, data.Id); errors.Is(err, sqlx.ErrNotFound) {
				return &product.SkuResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
			} else if err != nil {
				l.Errorf("query sku failed: %v", err)
				return nil, err
			}
		}
		id = data.Id
	} else {
		inserted, err := l.svcCtx.ProductSkuModel.InsertSku(l.ctx, data)
		if err != nil {
			l.Errorf("insert sku failed: %v", err)
			return nil, err
		}
		id = inserted
	}

	stored, err := l.svcCtx.ProductSkuModel.FindActive(l.ctx, id)
	if err != nil {
		l.Errorf("query saved sku failed: %v", err)
		return nil, err
	}
	return &product.SkuResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toSkuInfo(stored)}, nil
}
