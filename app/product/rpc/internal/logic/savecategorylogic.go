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

type SaveCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveCategoryLogic {
	return &SaveCategoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveCategoryLogic) SaveCategory(req *product.SaveCategoryReq) (*product.CategoryResp, error) {
	name := strings.TrimSpace(req.Name)
	icon := strings.TrimSpace(req.Icon)
	if name == "" || !validText(name, maxCategoryName) || !validText(icon, maxIconLength) ||
		req.ParentId < 0 || (req.Status != 0 && req.Status != 1) {
		return &product.CategoryResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	if req.Id > 0 && req.ParentId == req.Id {
		return &product.CategoryResp{Code: errcode.BadRequest, Msg: "父分类不能是自己"}, nil
	}

	data := &model.Category{
		ParentId: req.ParentId,
		Name:     name,
		Icon:     icon,
		Sort:     int64(req.Sort),
		Status:   int64(req.Status),
	}

	var id uint64
	if req.Id > 0 {
		data.Id = uint64(req.Id)
		affected, err := l.svcCtx.CategoryModel.UpdateCategory(l.ctx, data)
		if err != nil {
			l.Errorf("update category failed: %v", err)
			return nil, err
		}
		if affected == 0 {
			if _, err := l.svcCtx.CategoryModel.FindActive(l.ctx, data.Id); errors.Is(err, sqlx.ErrNotFound) {
				return &product.CategoryResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
			} else if err != nil {
				l.Errorf("query category failed: %v", err)
				return nil, err
			}
		}
		id = data.Id
	} else {
		inserted, err := l.svcCtx.CategoryModel.InsertCategory(l.ctx, data)
		if err != nil {
			l.Errorf("insert category failed: %v", err)
			return nil, err
		}
		id = inserted
	}

	stored, err := l.svcCtx.CategoryModel.FindActive(l.ctx, id)
	if err != nil {
		l.Errorf("query saved category failed: %v", err)
		return nil, err
	}
	return &product.CategoryResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toCategoryInfo(stored)}, nil
}
