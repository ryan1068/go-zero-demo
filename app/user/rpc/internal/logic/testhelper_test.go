package logic

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-demo/user-rpc/internal/config"
	"go-zero-demo/user-rpc/internal/model"
	"go-zero-demo/user-rpc/internal/svc"
)

const testSecret = "0123456789abcdef0123456789abcdef"

var userRowColumns = []string{
	"id", "username", "password", "nickname", "created_at", "updated_at",
	"email", "phone", "avatar", "status", "remark", "role", "last_login_at",
}

func newUserRows() *sqlmock.Rows {
	return sqlmock.NewRows(userRowColumns)
}

func addUserRow(rows *sqlmock.Rows, id int64, username, password, nickname, email, phone, avatar string, status int64, remark string, role int64) *sqlmock.Rows {
	return rows.AddRow(
		id, username, password, nickname,
		time.Unix(1704067200, 0), time.Unix(1704067201, 0),
		email, phone, avatar, status, remark, role, nil,
	)
}

func newUserServiceContext(t *testing.T, db *sql.DB) *svc.ServiceContext {
	t.Helper()
	conn := sqlx.NewSqlConnFromDB(db)
	return &svc.ServiceContext{
		DB:                 conn,
		UserModel:          model.NewUserModel(conn),
		UserAddressModel:   model.NewUserAddressModel(conn),
		AdminAuditLogModel: model.NewAdminAuditLogModel(conn),
		Config: config.Config{
			Jwt: config.JwtConfig{Secret: testSecret, ExpireHours: 24},
		},
	}
}
