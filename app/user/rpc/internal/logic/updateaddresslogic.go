package logic

import (
	"context"
	"errors"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
	"go-zero-demo/user-rpc/internal/model"
	"go-zero-demo/user-rpc/internal/svc"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateAddressLogic) UpdateAddress(in *user.SaveAddressReq) (*user.AddressResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &user.AddressResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	p, valid := normalizeAddress(in)
	if !valid || in.Id <= 0 {
		return &user.AddressResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var resp *user.AddressResp
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := l.svcCtx.UserAddressModel.FindByIDAndUser(ctx, session, uint64(in.Id), userID); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				resp = &user.AddressResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
				return model.ErrAbortTransaction
			}
			return err
		}
		if p.IsDefault == 1 {
			if err := l.svcCtx.UserAddressModel.ClearDefault(ctx, session, userID); err != nil {
				return err
			}
		}

		if err := l.svcCtx.UserAddressModel.UpdateFields(ctx, session, &model.UserAddress{
			Id:            uint64(in.Id),
			UserId:        userID,
			ReceiverName:  p.ReceiverName,
			ReceiverPhone: p.ReceiverPhone,
			Province:      p.Province,
			City:          p.City,
			District:      p.District,
			DetailAddress: p.DetailAddress,
			PostalCode:    p.PostalCode,
			Tag:           p.Tag,
			IsDefault:     int64(p.IsDefault),
		}); err != nil {
			return err
		}

		stored, err := l.svcCtx.UserAddressModel.FindByIDAndUser(ctx, session, uint64(in.Id), userID)
		if err != nil {
			return err
		}
		resp = &user.AddressResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: toAddressInfo(stored)}
		return nil
	})
	switch {
	case errors.Is(err, model.ErrAbortTransaction):
		return resp, nil
	case err != nil:
		l.Errorf("update address failed: %v", err)
		return nil, err
	}
	return resp, nil
}
