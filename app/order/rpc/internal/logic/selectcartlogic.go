package logic

import (
	"context"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSelectCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectCartLogic {
	return &SelectCartLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SelectCartLogic) SelectCart(in *order.SelectCartReq) (*order.BaseResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.BaseResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Selected != 0 && in.Selected != 1 {
		return &order.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	if err := l.svcCtx.CartItemModel.UpdateCartSelectedByUser(l.ctx, nil, userID, int64(in.Selected)); err != nil {
		l.Errorf("select cart items failed: %v", err)
		return nil, err
	}
	return &order.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
