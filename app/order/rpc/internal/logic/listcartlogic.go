package logic

import (
	"context"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCartLogic {
	return &ListCartLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListCartLogic) ListCart(_ *order.EmptyReq) (*order.CartListResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.CartListResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}

	rows, err := l.svcCtx.CartItemModel.ListCartItemRows(l.ctx, nil, userID)
	if err != nil {
		l.Errorf("query cart items failed: %v", err)
		return nil, err
	}

	list := make([]*order.CartItemInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCartItemInfo(row))
	}
	return &order.CartListResp{
		Code:  errcode.Success,
		Msg:   errcode.Msg(errcode.Success),
		List:  list,
		Total: int64(len(list)),
	}, nil
}
