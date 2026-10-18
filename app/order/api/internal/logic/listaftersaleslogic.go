package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAfterSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAfterSalesLogic {
	return &ListAfterSalesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListAfterSalesLogic) ListAfterSales(req *types.ListAfterSalesReq) (*types.AfterSaleListResp, error) {
	return l.list(req, false)
}

func (l *ListAfterSalesLogic) list(req *types.ListAfterSalesReq, admin bool) (*types.AfterSaleListResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AfterSaleListResp{Code: 401, Msg: "未授权"}, nil
	}
	adminFlag := int32(0)
	if admin {
		adminFlag = 1
	}
	s, err := l.svcCtx.OrderRpc.ListAfterSales(ctx, &order.ListAfterSalesReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		Admin:    adminFlag,
	})
	if err != nil {
		l.Errorf("list after sales rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AfterSaleListResp{Code: s.Code, Msg: s.Msg}
	if s.Code == 200 {
		list := make([]types.AfterSaleInfo, 0, len(s.List))
		for _, item := range s.List {
			list = append(list, toAfterSaleInfo(item))
		}
		resp.Data = types.AfterSaleListData{
			List:     list,
			Total:    s.Total,
			Page:     s.Page,
			PageSize: s.PageSize,
		}
	}
	return resp, nil
}
