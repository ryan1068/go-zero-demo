package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListOrdersLogic) ListOrders(req *types.ListOrdersReq) (*types.OrderListResp, error) {
	return l.listOrders(req, false)
}

func (l *ListOrdersLogic) listOrders(req *types.ListOrdersReq, admin bool) (*types.OrderListResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.OrderListResp{Code: 401, Msg: "未授权"}, nil
	}
	adminFlag := int32(0)
	if admin {
		adminFlag = 1
	}
	s, err := l.svcCtx.OrderRpc.ListOrders(ctx, &order.ListOrdersReq{
		Page:      req.Page,
		PageSize:  req.PageSize,
		OrderNo:   req.OrderNo,
		Status:    req.Status,
		UserId:    req.UserId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Admin:     adminFlag,
	})
	if err != nil {
		l.Errorf("list orders rpc failed: %v", err)
		return nil, err
	}
	resp := &types.OrderListResp{Code: s.Code, Msg: s.Msg}
	resp.Data.Total = s.Total
	resp.Data.Page = s.Page
	resp.Data.PageSize = s.PageSize
	resp.Data.List = make([]types.OrderInfo, 0, len(s.List))
	for _, item := range s.List {
		if item != nil {
			resp.Data.List = append(resp.Data.List, toOrderInfo(item))
		}
	}
	return resp, nil
}
