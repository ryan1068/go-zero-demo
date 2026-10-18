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

type CreateAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAddressLogic {
	return &CreateAddressLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateAddressLogic) CreateAddress(in *user.SaveAddressReq) (*user.AddressResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &user.AddressResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	p, valid := normalizeAddress(in)
	if !valid {
		return &user.AddressResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var resp *user.AddressResp
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		count, err := l.svcCtx.UserAddressModel.CountByUser(ctx, session, userID)
		if err != nil {
			return err
		}
		if count >= maxAddressCount {
			resp = &user.AddressResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}
			return model.ErrAbortTransaction
		}

		isDefault := p.IsDefault
		if count == 0 {
			isDefault = 1
		}
		if isDefault == 1 {
			if err := l.svcCtx.UserAddressModel.ClearDefault(ctx, session, userID); err != nil {
				return err
			}
		}

		result, err := l.svcCtx.UserAddressModel.InsertTx(ctx, session, &model.UserAddress{
			UserId:        userID,
			ReceiverName:  p.ReceiverName,
			ReceiverPhone: p.ReceiverPhone,
			Province:      p.Province,
			City:          p.City,
			District:      p.District,
			DetailAddress: p.DetailAddress,
			PostalCode:    p.PostalCode,
			Tag:           p.Tag,
			IsDefault:     int64(isDefault),
		})
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		stored, err := l.svcCtx.UserAddressModel.FindByIDAndUser(ctx, session, uint64(id), userID)
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
		l.Errorf("create address failed: %v", err)
		return nil, err
	}
	return resp, nil
}
