package logic

import (
	"context"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAddressesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAddressesLogic {
	return &ListAddressesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListAddressesLogic) ListAddresses(_ *user.EmptyReq) (*user.ListAddressesResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &user.ListAddressesResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}

	rows, err := l.svcCtx.UserAddressModel.ListByUser(l.ctx, userID)
	if err != nil {
		l.Errorf("query addresses failed: %v", err)
		return nil, err
	}

	list := make([]*user.AddressInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, toAddressInfo(row))
	}

	return &user.ListAddressesResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		List: list,
	}, nil
}
