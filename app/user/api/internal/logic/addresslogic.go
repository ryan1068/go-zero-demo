package logic

import (
	"context"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

func injectUser(ctx context.Context) (context.Context, bool) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return ctx, false
	}
	return auth.InjectUserID(ctx, userID), true
}

type ListAddressesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAddressesLogic {
	return &ListAddressesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListAddressesLogic) ListAddresses() (*types.ListAddressesResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.ListAddressesResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.UserRpc.ListAddresses(ctx, &user.EmptyReq{})
	if err != nil {
		l.Errorf("list addresses rpc failed: %v", err)
		return nil, err
	}
	resp := &types.ListAddressesResp{Code: s.Code, Msg: s.Msg}
	resp.Data.List = make([]types.AddressInfo, 0, len(s.List))
	for _, item := range s.List {
		if item != nil {
			resp.Data.List = append(resp.Data.List, toAddressInfo(item))
		}
	}
	return resp, nil
}

type CreateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAddressLogic {
	return &CreateAddressLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateAddressLogic) CreateAddress(req *types.SaveAddressReq) (*types.AddressResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AddressResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.UserRpc.CreateAddress(ctx, toSaveAddressReq(req))
	if err != nil {
		l.Errorf("create address rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AddressResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toAddressInfo(s.Data)
	}
	return resp, nil
}

type UpdateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateAddressLogic) UpdateAddress(req *types.SaveAddressReq) (*types.AddressResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.AddressResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.UserRpc.UpdateAddress(ctx, toSaveAddressReq(req))
	if err != nil {
		l.Errorf("update address rpc failed: %v", err)
		return nil, err
	}
	resp := &types.AddressResp{Code: s.Code, Msg: s.Msg}
	if s.Data != nil {
		resp.Data = toAddressInfo(s.Data)
	}
	return resp, nil
}

type DeleteAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddressLogic {
	return &DeleteAddressLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteAddressLogic) DeleteAddress(id int64) (*types.BaseResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.BaseResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.UserRpc.DeleteAddress(ctx, &user.AddressReq{Id: id})
	if err != nil {
		l.Errorf("delete address rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}

type SetDefaultAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDefaultAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultAddressLogic {
	return &SetDefaultAddressLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SetDefaultAddressLogic) SetDefaultAddress(id int64) (*types.BaseResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.BaseResp{Code: 401, Msg: "未授权"}, nil
	}
	s, err := l.svcCtx.UserRpc.SetDefaultAddress(ctx, &user.AddressReq{Id: id})
	if err != nil {
		l.Errorf("set default address rpc failed: %v", err)
		return nil, err
	}
	return &types.BaseResp{Code: s.Code, Msg: s.Msg}, nil
}

func toSaveAddressReq(req *types.SaveAddressReq) *user.SaveAddressReq {
	return &user.SaveAddressReq{
		Id:            req.Id,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		DetailAddress: req.DetailAddress,
		PostalCode:    req.PostalCode,
		Tag:           req.Tag,
		IsDefault:     req.IsDefault,
	}
}
