package logic

import (
	"context"

	"go-zero-demo/order-rpc/internal/model"
	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAfterSalesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAfterSalesLogic {
	return &ListAfterSalesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListAfterSalesLogic) ListAfterSales(in *order.ListAfterSalesReq) (*order.AfterSaleListResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.AfterSaleListResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	filter := model.AfterSaleFilter{
		Admin:    in.Admin == 1,
		UserID:   userID,
		Status:   int64(in.Status),
		Page:     int64(page),
		PageSize: int64(pageSize),
	}
	total, err := l.svcCtx.AfterSaleModel.CountAfterSales(l.ctx, nil, filter)
	if err != nil {
		l.Errorf("count after sales failed: %v", err)
		return nil, err
	}
	rows, err := l.svcCtx.AfterSaleModel.ListAfterSales(l.ctx, nil, filter)
	if err != nil {
		l.Errorf("query after sales failed: %v", err)
		return nil, err
	}

	list := make([]*order.AfterSaleInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, toAfterSaleInfo(row))
	}
	return &order.AfterSaleListResp{
		Code:     errcode.Success,
		Msg:      errcode.Msg(errcode.Success),
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
