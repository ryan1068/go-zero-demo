package logic

import (
	"context"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCartLogic {
	return &AddCartLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AddCartLogic) AddCart(req *types.AddCartReq) (*types.CartItemResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.CartItemResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.AddCart(ctx, &order.AddCartReq{
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	})
	if err != nil {
		l.Errorf("add cart rpc failed: %v", err)
		return nil, err
	}
	resp := &types.CartItemResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toCartItem(s.Data)
	}
	return resp, nil
}

type UpdateCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartLogic {
	return &UpdateCartLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateCartLogic) UpdateCart(req *types.UpdateCartReq) (*types.CartItemResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.CartItemResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.UpdateCart(ctx, &order.UpdateCartReq{
		Id:       req.Id,
		Quantity: req.Quantity,
		Selected: req.Selected,
	})
	if err != nil {
		l.Errorf("update cart rpc failed: %v", err)
		return nil, err
	}
	resp := &types.CartItemResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toCartItem(s.Data)
	}
	return resp, nil
}

type DeleteCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCartLogic {
	return &DeleteCartLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteCartLogic) DeleteCart(id int64) (*types.BaseResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.BaseResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.DeleteCart(ctx, &order.CartItemReq{Id: id})
	if err != nil {
		l.Errorf("delete cart rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}

type SelectCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectCartLogic {
	return &SelectCartLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SelectCartLogic) SelectCart(req *types.SelectCartReq) (*types.BaseResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.BaseResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.SelectCart(ctx, &order.SelectCartReq{Selected: req.Selected})
	if err != nil {
		l.Errorf("select cart rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}

type ListCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCartLogic {
	return &ListCartLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListCartLogic) ListCart() (*types.CartListResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.CartListResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.OrderRpc.ListCart(ctx, &order.EmptyReq{})
	if err != nil {
		l.Errorf("list cart rpc failed: %v", err)
		return nil, err
	}
	resp := &types.CartListResp{Code: s.Code, Msg: s.Msg}
	resp.Data.Total = s.Total
	resp.Data.List = make([]types.CartItemInfo, 0, len(s.List))
	for _, item := range s.List {
		if item != nil {
			resp.Data.List = append(resp.Data.List, toCartItem(item))
		}
	}
	return resp, nil
}
