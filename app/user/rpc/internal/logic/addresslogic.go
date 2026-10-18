package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"unicode/utf8"

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

	rows, err := l.svcCtx.DB.QueryContext(l.ctx,
		"SELECT "+addressSelectColumns+" FROM user_address WHERE user_id = ? AND deleted_at IS NULL ORDER BY is_default DESC, updated_at DESC",
		userID,
	)
	if err != nil {
		l.Errorf("query addresses failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			l.Errorf("close address rows failed: %v", closeErr)
		}
	}()

	list := make([]*user.AddressInfo, 0, 8)
	for rows.Next() {
		info, err := scanAddressInfo(rows)
		if err != nil {
			l.Errorf("scan address failed: %v", err)
			return nil, err
		}
		list = append(list, info)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate addresses failed: %v", err)
		return nil, err
	}

	return &user.ListAddressesResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		List: list,
	}, nil
}

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

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin address transaction failed: %v", err)
		return nil, err
	}
	defer rollbackAddressTx(l.Logger, tx)

	var count int32
	if err := tx.QueryRowContext(l.ctx,
		"SELECT COUNT(*) FROM user_address WHERE user_id = ? AND deleted_at IS NULL", userID,
	).Scan(&count); err != nil {
		l.Errorf("count addresses failed: %v", err)
		return nil, err
	}
	if count >= maxAddressCount {
		return &user.AddressResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	isDefault := p.IsDefault
	if count == 0 {
		isDefault = 1
	}
	if isDefault == 1 {
		if _, err := tx.ExecContext(l.ctx,
			"UPDATE user_address SET is_default = 0 WHERE user_id = ? AND deleted_at IS NULL", userID,
		); err != nil {
			l.Errorf("clear default address failed: %v", err)
			return nil, err
		}
	}

	result, err := tx.ExecContext(l.ctx,
		"INSERT INTO user_address (user_id, receiver_name, receiver_phone, province, city, district, detail_address, postal_code, tag, is_default) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		userID, p.ReceiverName, p.ReceiverPhone, p.Province, p.City, p.District, p.DetailAddress, p.PostalCode, p.Tag, isDefault,
	)
	if err != nil {
		l.Errorf("insert address failed: %v", err)
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		l.Errorf("read inserted address id failed: %v", err)
		return nil, err
	}

	info, err := findAddress(l.ctx, tx, id, userID)
	if err != nil {
		l.Errorf("query created address failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit address transaction failed: %v", err)
		return nil, err
	}

	return &user.AddressResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

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

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin address transaction failed: %v", err)
		return nil, err
	}
	defer rollbackAddressTx(l.Logger, tx)

	if _, err := findAddress(l.ctx, tx, in.Id, userID); err != nil {
		return addressLookupError(l, err)
	}
	if p.IsDefault == 1 {
		if _, err := tx.ExecContext(l.ctx,
			"UPDATE user_address SET is_default = 0 WHERE user_id = ? AND deleted_at IS NULL", userID,
		); err != nil {
			l.Errorf("clear default address failed: %v", err)
			return nil, err
		}
	}

	if _, err := tx.ExecContext(l.ctx,
		"UPDATE user_address SET receiver_name = ?, receiver_phone = ?, province = ?, city = ?, district = ?, detail_address = ?, postal_code = ?, tag = ?, is_default = ? WHERE id = ? AND user_id = ? AND deleted_at IS NULL",
		p.ReceiverName, p.ReceiverPhone, p.Province, p.City, p.District, p.DetailAddress, p.PostalCode, p.Tag, p.IsDefault, in.Id, userID,
	); err != nil {
		l.Errorf("update address failed: %v", err)
		return nil, err
	}

	info, err := findAddress(l.ctx, tx, in.Id, userID)
	if err != nil {
		l.Errorf("query updated address failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit address transaction failed: %v", err)
		return nil, err
	}

	return &user.AddressResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

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

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin address transaction failed: %v", err)
		return nil, err
	}
	defer rollbackAddressTx(l.Logger, tx)

	var isDefault int32
	err = tx.QueryRowContext(l.ctx,
		"SELECT is_default FROM user_address WHERE id = ? AND user_id = ? AND deleted_at IS NULL FOR UPDATE",
		in.Id, userID,
	).Scan(&isDefault)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &user.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query address failed: %v", err)
		return nil, err
	}

	if _, err := tx.ExecContext(l.ctx,
		"UPDATE user_address SET deleted_at = NOW() WHERE id = ? AND user_id = ?", in.Id, userID,
	); err != nil {
		l.Errorf("delete address failed: %v", err)
		return nil, err
	}

	if isDefault == 1 {
		var nextID int64
		err := tx.QueryRowContext(l.ctx,
			"SELECT id FROM user_address WHERE user_id = ? AND deleted_at IS NULL ORDER BY updated_at DESC LIMIT 1",
			userID,
		).Scan(&nextID)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// No address left; nothing to promote.
		case err != nil:
			l.Errorf("query next default address failed: %v", err)
			return nil, err
		default:
			if _, err := tx.ExecContext(l.ctx,
				"UPDATE user_address SET is_default = 1 WHERE id = ? AND user_id = ?", nextID, userID,
			); err != nil {
				l.Errorf("promote default address failed: %v", err)
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		l.Errorf("commit address transaction failed: %v", err)
		return nil, err
	}
	return &user.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}

type SetDefaultAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetDefaultAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultAddressLogic {
	return &SetDefaultAddressLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SetDefaultAddressLogic) SetDefaultAddress(in *user.AddressReq) (*user.BaseResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &user.BaseResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &user.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin address transaction failed: %v", err)
		return nil, err
	}
	defer rollbackAddressTx(l.Logger, tx)

	var id int64
	err = tx.QueryRowContext(l.ctx,
		"SELECT id FROM user_address WHERE id = ? AND user_id = ? AND deleted_at IS NULL FOR UPDATE",
		in.Id, userID,
	).Scan(&id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &user.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query address failed: %v", err)
		return nil, err
	}

	if _, err := tx.ExecContext(l.ctx,
		"UPDATE user_address SET is_default = 0 WHERE user_id = ? AND deleted_at IS NULL", userID,
	); err != nil {
		l.Errorf("clear default address failed: %v", err)
		return nil, err
	}
	if _, err := tx.ExecContext(l.ctx,
		"UPDATE user_address SET is_default = 1 WHERE id = ? AND user_id = ?", in.Id, userID,
	); err != nil {
		l.Errorf("set default address failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit address transaction failed: %v", err)
		return nil, err
	}
	return &user.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}

type addressParams struct {
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	DetailAddress string
	PostalCode    string
	Tag           string
	IsDefault     int32
}

func normalizeAddress(in *user.SaveAddressReq) (addressParams, bool) {
	p := addressParams{
		ReceiverName:  strings.TrimSpace(in.ReceiverName),
		ReceiverPhone: strings.TrimSpace(in.ReceiverPhone),
		Province:      strings.TrimSpace(in.Province),
		City:          strings.TrimSpace(in.City),
		District:      strings.TrimSpace(in.District),
		DetailAddress: strings.TrimSpace(in.DetailAddress),
		PostalCode:    strings.TrimSpace(in.PostalCode),
		Tag:           strings.TrimSpace(in.Tag),
		IsDefault:     in.IsDefault,
	}
	if p.ReceiverName == "" || utf8.RuneCountInString(p.ReceiverName) > maxReceiverName ||
		p.ReceiverPhone == "" || utf8.RuneCountInString(p.ReceiverPhone) > maxReceiverPhone ||
		utf8.RuneCountInString(p.Province) > maxRegionName ||
		utf8.RuneCountInString(p.City) > maxRegionName ||
		utf8.RuneCountInString(p.District) > maxRegionName ||
		p.DetailAddress == "" || utf8.RuneCountInString(p.DetailAddress) > maxDetailAddress ||
		utf8.RuneCountInString(p.PostalCode) > maxPostalCode ||
		utf8.RuneCountInString(p.Tag) > maxAddressTag ||
		(p.IsDefault != 0 && p.IsDefault != 1) {
		return p, false
	}
	return p, true
}

func addressLookupError(l logx.Logger, err error) (*user.AddressResp, error) {
	if errors.Is(err, sql.ErrNoRows) {
		return &user.AddressResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	l.Errorf("query address failed: %v", err)
	return nil, err
}

func rollbackAddressTx(logger logx.Logger, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		logger.Errorf("rollback address transaction failed: %v", err)
	}
}
