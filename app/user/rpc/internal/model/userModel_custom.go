package model

import (
	"context"
	"strings"
)

// UserFilter is the filter for the admin user list.
type UserFilter struct {
	Username string
	Nickname string
	Email    string
	Phone    string
	// Status: 0 = all, 1 = active, 2 = disabled.
	Status   int64
	Page     int64
	PageSize int64
}

// FindUsers returns a page of users and the total count.
func (m *customUserModel) FindUsers(ctx context.Context, f UserFilter) ([]*User, int64, error) {
	where := []string{"1 = 1"}
	args := make([]any, 0, 5)
	if f.Username != "" {
		where = append(where, "`username` LIKE ?")
		args = append(args, f.Username+"%")
	}
	if f.Nickname != "" {
		where = append(where, "`nickname` LIKE ?")
		args = append(args, f.Nickname+"%")
	}
	if f.Email != "" {
		where = append(where, "`email` LIKE ?")
		args = append(args, f.Email+"%")
	}
	if f.Phone != "" {
		where = append(where, "`phone` LIKE ?")
		args = append(args, f.Phone+"%")
	}
	switch f.Status {
	case 1:
		where = append(where, "`status` = ?")
		args = append(args, 1)
	case 2:
		where = append(where, "`status` = ?")
		args = append(args, 0)
	}
	condition := strings.Join(where, " AND ")

	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total,
		"select count(*) from "+m.table+" where "+condition, args...,
	); err != nil {
		return nil, 0, err
	}

	listArgs := append(append([]any{}, args...), f.PageSize, (f.Page-1)*f.PageSize)
	var list []*User
	if err := m.conn.QueryRowsCtx(ctx, &list,
		"select "+userRows+" from "+m.table+" where "+condition+" order by `id` desc limit ? offset ?",
		listArgs...,
	); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// UpdateStatus updates a user status and returns affected rows.
func (m *customUserModel) UpdateStatus(ctx context.Context, id, status int64) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `status` = ? where `id` = ?", status, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateRole updates a user role and returns affected rows.
func (m *customUserModel) UpdateRole(ctx context.Context, id, role int64) (int64, error) {
	result, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `role` = ? where `id` = ?", role, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdatePassword updates the password hash of a user.
func (m *customUserModel) UpdatePassword(ctx context.Context, id int64, password string) error {
	_, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `password` = ?, `updated_at` = now() where `id` = ?", password, id,
	)
	return err
}

// UpdateLastLogin records the last successful login time.
func (m *customUserModel) UpdateLastLogin(ctx context.Context, id int64) error {
	_, err := m.conn.ExecCtx(ctx,
		"update "+m.table+" set `last_login_at` = now() where `id` = ?", id,
	)
	return err
}
