package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// sessionOrConn returns the transaction session when one is provided,
// otherwise the model's raw connection. Both satisfy sqlx.Session.
func sessionOrConn(conn sqlx.SqlConn, session sqlx.Session) sqlx.Session {
	if session != nil {
		return session
	}
	return conn
}

// ListByUser returns all active addresses of a user.
func (m *customUserAddressModel) ListByUser(ctx context.Context, userID int64) ([]*UserAddress, error) {
	var list []*UserAddress
	err := m.conn.QueryRowsCtx(ctx, &list,
		"select "+userAddressRows+" from "+m.table+
			" where `user_id` = ? and `deleted_at` is null"+
			" order by `is_default` desc, `updated_at` desc",
		userID,
	)
	return list, err
}

// CountByUser counts active addresses of a user inside or outside a transaction.
func (m *customUserAddressModel) CountByUser(ctx context.Context, session sqlx.Session, userID int64) (int64, error) {
	var total int64
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &total,
		"select count(*) from "+m.table+" where `user_id` = ? and `deleted_at` is null", userID,
	)
	return total, err
}

// FindByIDAndUser returns an active address owned by the user.
func (m *customUserAddressModel) FindByIDAndUser(ctx context.Context, session sqlx.Session, id uint64, userID int64) (*UserAddress, error) {
	var resp UserAddress
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &resp,
		"select "+userAddressRows+" from "+m.table+
			" where `id` = ? and `user_id` = ? and `deleted_at` is null limit 1",
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// FindForUpdate locks an active address owned by the user.
// It must be called with a real transaction session.
func (m *customUserAddressModel) FindForUpdate(ctx context.Context, session sqlx.Session, id uint64, userID int64) (*UserAddress, error) {
	var resp UserAddress
	err := session.QueryRowCtx(ctx, &resp,
		"select "+userAddressRows+" from "+m.table+
			" where `id` = ? and `user_id` = ? and `deleted_at` is null for update",
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// InsertTx inserts an address inside a transaction.
func (m *customUserAddressModel) InsertTx(ctx context.Context, session sqlx.Session, data *UserAddress) (sql.Result, error) {
	return m.withSession(session).Insert(ctx, data)
}

// UpdateFields updates the editable fields of an address.
func (m *customUserAddressModel) UpdateFields(ctx context.Context, session sqlx.Session, data *UserAddress) error {
	_, err := sessionOrConn(m.conn, session).ExecCtx(ctx,
		"update "+m.table+" set `receiver_name` = ?, `receiver_phone` = ?, `province` = ?,"+
			" `city` = ?, `district` = ?, `detail_address` = ?, `postal_code` = ?,"+
			" `tag` = ?, `is_default` = ? where `id` = ? and `user_id` = ? and `deleted_at` is null",
		data.ReceiverName, data.ReceiverPhone, data.Province, data.City, data.District,
		data.DetailAddress, data.PostalCode, data.Tag, data.IsDefault, data.Id, data.UserId,
	)
	return err
}

// SoftDelete marks an address as deleted.
func (m *customUserAddressModel) SoftDelete(ctx context.Context, session sqlx.Session, id uint64, userID int64) error {
	_, err := sessionOrConn(m.conn, session).ExecCtx(ctx,
		"update "+m.table+" set `deleted_at` = now() where `id` = ? and `user_id` = ?", id, userID,
	)
	return err
}

// ClearDefault clears the default flag of all active addresses of a user.
func (m *customUserAddressModel) ClearDefault(ctx context.Context, session sqlx.Session, userID int64) error {
	_, err := sessionOrConn(m.conn, session).ExecCtx(ctx,
		"update "+m.table+" set `is_default` = 0 where `user_id` = ? and `deleted_at` is null", userID,
	)
	return err
}

// SetDefault marks one address as default.
func (m *customUserAddressModel) SetDefault(ctx context.Context, session sqlx.Session, id uint64, userID int64) error {
	_, err := sessionOrConn(m.conn, session).ExecCtx(ctx,
		"update "+m.table+" set `is_default` = 1 where `id` = ? and `user_id` = ?", id, userID,
	)
	return err
}

// FindNextDefault returns the most recently updated active address of a user.
func (m *customUserAddressModel) FindNextDefault(ctx context.Context, session sqlx.Session, userID int64) (uint64, error) {
	var id uint64
	err := sessionOrConn(m.conn, session).QueryRowCtx(ctx, &id,
		"select `id` from "+m.table+
			" where `user_id` = ? and `deleted_at` is null order by `updated_at` desc limit 1",
		userID,
	)
	return id, err
}
