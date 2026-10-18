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

type DeleteAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddressLogic {
	return &DeleteAddressLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteAddressLogic) DeleteAddress(in *user.AddressReq) (*user.BaseResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &user.BaseResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &user.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	var resp *user.BaseResp
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		locked, err := l.svcCtx.UserAddressModel.FindForUpdate(ctx, session, uint64(in.Id), userID)
		if errors.Is(err, sqlx.ErrNotFound) {
			resp = &user.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}
			return model.ErrAbortTransaction
		}
		if err != nil {
			return err
		}

		if err := l.svcCtx.UserAddressModel.SoftDelete(ctx, session, uint64(in.Id), userID); err != nil {
			return err
		}

		if locked.IsDefault == 1 {
			nextID, err := l.svcCtx.UserAddressModel.FindNextDefault(ctx, session, userID)
			switch {
			case errors.Is(err, sqlx.ErrNotFound):
				// No address left; nothing to promote.
			case err != nil:
				return err
			default:
				if err := l.svcCtx.UserAddressModel.SetDefault(ctx, session, nextID, userID); err != nil {
					return err
				}
			}
		}

		resp = &user.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}
		return nil
	})
	switch {
	case errors.Is(err, model.ErrAbortTransaction):
		return resp, nil
	case err != nil:
		l.Errorf("delete address failed: %v", err)
		return nil, err
	}
	return resp, nil
}
