package logic

import (
	"context"

	"go-zero-demo/user-rpc/pb/user"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderReq) (*order.CommonResp, error) {
	_, err := l.svcCtx.UserRpc.User(l.ctx, &user.UserReq{
		Name: in.GoodsName,
	})

	if err != nil {
		l.Errorf("user rpc failed: %v", err)
	}

	return &order.CommonResp{
		Msg:  errcode.Msg(errcode.Success),
		Code: errcode.Success,
		Data: in.GoodsName,
	}, nil
}
