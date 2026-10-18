package logic

import (
	"context"
	"database/sql"

	"go-zero-demo/user-rpc/pb/user"
)

const (
	addressSelectColumns = "id, user_id, receiver_name, receiver_phone, COALESCE(province, ''), COALESCE(city, ''), COALESCE(district, ''), detail_address, COALESCE(postal_code, ''), COALESCE(tag, ''), is_default, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at)"
	maxAddressCount      = 20
	maxReceiverName      = 64
	maxReceiverPhone     = 32
	maxRegionName        = 64
	maxDetailAddress     = 255
	maxPostalCode        = 16
	maxAddressTag        = 16
)

func scanAddressInfo(row rowScanner) (*user.AddressInfo, error) {
	info := new(user.AddressInfo)
	err := row.Scan(
		&info.Id,
		&info.UserId,
		&info.ReceiverName,
		&info.ReceiverPhone,
		&info.Province,
		&info.City,
		&info.District,
		&info.DetailAddress,
		&info.PostalCode,
		&info.Tag,
		&info.IsDefault,
		&info.CreatedAt,
		&info.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func findAddress(ctx context.Context, q queryer, id, userID int64) (*user.AddressInfo, error) {
	return scanAddressInfo(q.QueryRowContext(ctx,
		"SELECT "+addressSelectColumns+" FROM user_address WHERE id = ? AND user_id = ? AND deleted_at IS NULL",
		id, userID,
	))
}

type queryer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
