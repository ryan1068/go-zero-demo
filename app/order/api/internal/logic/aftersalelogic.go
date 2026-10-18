package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

func toAfterSaleInfo(info *order.AfterSaleInfo) types.AfterSaleInfo {
	return types.AfterSaleInfo{
		Id:               info.Id,
		AfterSaleNo:      info.AfterSaleNo,
		OrderId:          info.OrderId,
		OrderItemId:      info.OrderItemId,
		UserId:           info.UserId,
		Type:             info.Type,
		Reason:           info.Reason,
		Description:      info.Description,
		Images:           info.Images,
		Amount:           info.Amount,
		Status:           info.Status,
		AdminRemark:      info.AdminRemark,
		TrackingNo:       info.TrackingNo,
		ApplyAt:          info.ApplyAt,
		ReviewAt:         info.ReviewAt,
		BuyerShippedAt:   info.BuyerShippedAt,
		SellerReceivedAt: info.SellerReceivedAt,
		RefundedAt:       info.RefundedAt,
		CreatedAt:        info.CreatedAt,
		UpdatedAt:        info.UpdatedAt,
		OrderNo:          info.OrderNo,
		ProductName:      info.ProductName,
		ProductCover:     info.ProductCover,
	}
}

type ApplyAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAfterSaleLogic {
	return &ApplyAfterSaleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ApplyAfterSaleLogic) ApplyAfterSale(req *types.ApplyAfterSaleReq) (*types.AfterSaleResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AfterSaleResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.ApplyAfterSale(ctx, &order.ApplyAfterSaleReq{
		OrderId:     req.OrderId,
		OrderItemId: req.OrderItemId,
		Type:        req.Type,
		Reason:      req.Reason,
		Description: req.Description,
		Images:      req.Images,
	})
	if err != nil {
		l.Errorf("apply after sale rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AfterSaleResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toAfterSaleInfo(s.Data)
	}
	return resp, nil
}

type ListAfterSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAfterSalesLogic {
	return &ListAfterSalesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListAfterSalesLogic) ListAfterSales(req *types.ListAfterSalesReq, admin bool) (*types.AfterSaleListResp, error) {
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

type GetAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAfterSaleLogic {
	return &GetAfterSaleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAfterSaleLogic) GetAfterSale(id int64, admin bool) (*types.AfterSaleResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AfterSaleResp{Code: 401, Msg: "未授权"}, nil
	}
	adminFlag := int32(0)
	if admin {
		adminFlag = 1
	}
	s, err := l.svcCtx.OrderRpc.GetAfterSale(ctx, &order.AfterSaleReq{Id: id, Admin: adminFlag})
	if err != nil {
		l.Errorf("get after sale rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AfterSaleResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toAfterSaleInfo(s.Data)
	}
	return resp, nil
}

type HandleAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleAfterSaleLogic {
	return &HandleAfterSaleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *HandleAfterSaleLogic) HandleAfterSale(id int64, action, adminRemark, trackingNo string, admin bool) (*types.AfterSaleResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AfterSaleResp{Code: 401, Msg: "未授权"}, nil
	}
	adminFlag := int32(0)
	if admin {
		adminFlag = 1
	}
	s, err := l.svcCtx.OrderRpc.HandleAfterSale(ctx, &order.HandleAfterSaleReq{
		Id:          id,
		Action:      action,
		AdminRemark: adminRemark,
		TrackingNo:  trackingNo,
		Admin:       adminFlag,
	})
	if err != nil {
		l.Errorf("handle after sale rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AfterSaleResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toAfterSaleInfo(s.Data)
	}
	return resp, nil
}
