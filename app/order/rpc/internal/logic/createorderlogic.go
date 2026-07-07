package logic

import (
	"context"
	"fmt"
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
	// todo: add your logic here and delete this line

	user, err := l.svcCtx.UserRpc.User(l.ctx, &user.UserReq{
		Name: in.GoodsName,
	})
	fmt.Println(user, 111)

	if err != nil {
		fmt.Println(err)
	}

	return &order.CommonResp{
		Msg:  errcode.Msg(errcode.Success),
		Code: errcode.Success,
		Data: in.GoodsName,
	}, nil
}
