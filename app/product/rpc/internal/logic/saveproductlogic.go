package logic

import (
	"context"
	"errors"
	"strings"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/model"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SaveProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveProductLogic {
	return &SaveProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveProductLogic) SaveProduct(req *product.SaveProductReq) (*product.ProductResp, error) {
	name := strings.TrimSpace(req.Name)
	subtitle := strings.TrimSpace(req.Subtitle)
	cover := strings.TrimSpace(req.CoverImage)
	unit := strings.TrimSpace(req.Unit)
	if unit == "" {
		unit = "件"
	}

	imagesList := make([]string, 0, len(req.Images))
	for _, img := range req.Images {
		img = strings.TrimSpace(img)
		if img == "" {
			continue
		}
		if !validText(img, maxCoverImage) {
			return &product.ProductResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
		}
		imagesList = append(imagesList, img)
	}
	if len(imagesList) > maxProductImages {
		return &product.ProductResp{Code: errcode.BadRequest, Msg: "商品图片不能超过 8 张"}, nil
	}
	if cover == "" && len(imagesList) > 0 {
		cover = imagesList[0]
	}

	if req.CategoryId <= 0 || name == "" || !validText(name, maxProductName) ||
		!validText(subtitle, maxSubtitle) || !validText(cover, maxCoverImage) ||
		!validText(unit, maxUnit) || !validText(req.Detail, maxDetailLength) ||
		req.Price <= 0 || req.OriginalPrice < 0 || req.Stock < 0 ||
		(req.Status != 0 && req.Status != 1) || (req.IsRecommend != 0 && req.IsRecommend != 1) {
		return &product.ProductResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	exists, err := l.svcCtx.CategoryModel.Exists(l.ctx, uint64(req.CategoryId))
	if err != nil {
		l.Errorf("query category failed: %v", err)
		return nil, err
	}
	if !exists {
		return &product.ProductResp{Code: errcode.BadRequest, Msg: "分类不存在"}, nil
	}

	images, err := encodeImages(imagesList)
	if err != nil {
		l.Errorf("encode product images failed: %v", err)
		return nil, err
	}

	data := &model.Product{
		CategoryId:    req.CategoryId,
		Name:          name,
		Subtitle:      subtitle,
		CoverImage:    cover,
		Images:        model.NullString(images),
		Detail:        model.NullString(req.Detail),
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         int64(req.Stock),
		Unit:          unit,
		Status:        int64(req.Status),
		IsRecommend:   int64(req.IsRecommend),
		Sort:          int64(req.Sort),
	}

	var id uint64
	if req.Id > 0 {
		data.Id = uint64(req.Id)
		affected, err := l.svcCtx.ProductModel.UpdateProduct(l.ctx, data)
		if err != nil {
			l.Errorf("update product failed: %v", err)
			return nil, err
		}
		if affected == 0 {
			if _, err := l.svcCtx.ProductModel.FindProductRow(l.ctx, data.Id, true); errors.Is(err, sqlx.ErrNotFound) {
				return &product.ProductResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
			} else if err != nil {
				l.Errorf("query product failed: %v", err)
				return nil, err
			}
		}
		id = data.Id
	} else {
		inserted, err := l.svcCtx.ProductModel.InsertProduct(l.ctx, data)
		if err != nil {
			l.Errorf("insert product failed: %v", err)
			return nil, err
		}
		id = inserted
	}

	row, err := l.svcCtx.ProductModel.FindProductRow(l.ctx, id, true)
	if err != nil {
		l.Errorf("query saved product failed: %v", err)
		return nil, err
	}
	return &product.ProductResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toProductInfo(row)}, nil
}
