package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CrossUserAddressModel reads the receiver snapshot from the user database.
type CrossUserAddressModel struct {
	conn sqlx.SqlConn
}

// NewCrossUserAddressModel returns a cross-schema user address model.
func NewCrossUserAddressModel(conn sqlx.SqlConn) *CrossUserAddressModel {
	return &CrossUserAddressModel{conn: conn}
}

// ReceiverSnapshot is the receiver data copied into an order.
type ReceiverSnapshot struct {
	Name     string `db:"name"`
	Phone    string `db:"phone"`
	Province string `db:"province"`
	City     string `db:"city"`
	District string `db:"district"`
	Address  string `db:"address"`
	Postal   string `db:"postal"`
}

// FindReceiver returns the receiver of an address owned by the user.
func (m *CrossUserAddressModel) FindReceiver(ctx context.Context, session sqlx.Session, addressID, userID int64) (*ReceiverSnapshot, error) {
	var row ReceiverSnapshot
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &row,
		"select `receiver_name` as name, `receiver_phone` as phone,"+
			" coalesce(`province`, '') as province, coalesce(`city`, '') as city,"+
			" coalesce(`district`, '') as district, `detail_address` as address,"+
			" coalesce(`postal_code`, '') as postal"+
			" from looklook_usercenter.user_address"+
			" where `id` = ? and `user_id` = ? and `deleted_at` is null",
		addressID, userID,
	)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
