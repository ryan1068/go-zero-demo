package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-rpc/internal/svc"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoriesLogic {
	return &ListCategoriesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListCategoriesLogic) ListCategories(in *product.ListCategoriesReq) (*product.ListCategoriesResp, error) {
	admin := in.Admin == 1
	query := "SELECT " + categorySelectColumns + " FROM category WHERE deleted_at IS NULL"
	if !admin {
		query += " AND status = 1"
	}
	query += " ORDER BY sort DESC, id ASC"

	rows, err := l.svcCtx.DB.QueryContext(l.ctx, query)
	if err != nil {
		l.Errorf("query categories failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			l.Errorf("close category rows failed: %v", closeErr)
		}
	}()

	list := make([]*product.CategoryInfo, 0, 8)
	for rows.Next() {
		item, err := scanCategory(rows)
		if err != nil {
			l.Errorf("scan category failed: %v", err)
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate categories failed: %v", err)
		return nil, err
	}
	return &product.ListCategoriesResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), List: list}, nil
}

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

	var id int64
	var err error
	if req.Id > 0 {
		result, err := l.svcCtx.DB.ExecContext(l.ctx,
			"UPDATE category SET parent_id = ?, name = ?, icon = ?, sort = ?, status = ? WHERE id = ? AND deleted_at IS NULL",
			req.ParentId, name, icon, req.Sort, req.Status, req.Id,
		)
		if err != nil {
			l.Errorf("update category failed: %v", err)
			return nil, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			l.Errorf("read affected rows failed: %v", err)
			return nil, err
		}
		if affected == 0 {
			if _, err := l.findCategory(req.Id); errors.Is(err, sql.ErrNoRows) {
				return &product.CategoryResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
			} else if err != nil {
				l.Errorf("query category failed: %v", err)
				return nil, err
			}
		}
		id = req.Id
	} else {
		result, err := l.svcCtx.DB.ExecContext(l.ctx,
			"INSERT INTO category (parent_id, name, icon, sort, status) VALUES (?, ?, ?, ?, ?)",
			req.ParentId, name, icon, req.Sort, req.Status,
		)
		if err != nil {
			l.Errorf("insert category failed: %v", err)
			return nil, err
		}
		id, err = result.LastInsertId()
		if err != nil {
			l.Errorf("read inserted category id failed: %v", err)
			return nil, err
		}
	}

	info, err := l.findCategory(id)
	if err != nil {
		l.Errorf("query saved category failed: %v", err)
		return nil, err
	}
	return &product.CategoryResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

func (l *SaveCategoryLogic) findCategory(id int64) (*product.CategoryInfo, error) {
	return scanCategory(l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT "+categorySelectColumns+" FROM category WHERE id = ? AND deleted_at IS NULL", id,
	))
}

type DeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategoryLogic {
	return &DeleteCategoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteCategoryLogic) DeleteCategory(in *product.CategoryReq) (*product.BaseResp, error) {
	id := in.Id
	if id <= 0 {
		return &product.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}
	var count int64
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM product WHERE category_id = ? AND deleted_at IS NULL", id,
	).Scan(&count); err != nil {
		l.Errorf("count category products failed: %v", err)
		return nil, err
	}
	if count > 0 {
		return &product.BaseResp{Code: errcode.Conflict, Msg: "分类下还有商品，不能删除"}, nil
	}
	var childCount int64
	if err := l.svcCtx.DB.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM category WHERE parent_id = ? AND deleted_at IS NULL", id,
	).Scan(&childCount); err != nil {
		l.Errorf("count child categories failed: %v", err)
		return nil, err
	}
	if childCount > 0 {
		return &product.BaseResp{Code: errcode.Conflict, Msg: "分类下还有子分类，不能删除"}, nil
	}

	result, err := l.svcCtx.DB.ExecContext(l.ctx,
		"UPDATE category SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id,
	)
	if err != nil {
		l.Errorf("delete category failed: %v", err)
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		l.Errorf("read affected rows failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &product.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	return &product.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
