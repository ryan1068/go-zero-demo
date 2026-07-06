package logic

import (
	"context"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderReq) (*order.CommonResp, error) {
	// todo: add your logic here and delete this line

	return &order.CommonResp{
		Msg:  "ok",
		Code: 200,
		Data: "1111",
	}, nil
}
