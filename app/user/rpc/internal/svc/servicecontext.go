package svc

import (
	"errors"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-demo/user-rpc/internal/config"
	"go-zero-demo/user-rpc/internal/model"
)

type ServiceContext struct {
	Config config.Config
	DB     sqlx.SqlConn

	UserModel          model.UserModel
	UserAddressModel   model.UserAddressModel
	AdminAuditLogModel model.AdminAuditLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.MustNewConn(c.Mysql, sqlx.WithAcceptable(acceptable))
	return &ServiceContext{
		Config:             c,
		DB:                 conn,
		UserModel:          model.NewUserModel(conn),
		UserAddressModel:   model.NewUserAddressModel(conn),
		AdminAuditLogModel: model.NewAdminAuditLogModel(conn),
	}
}

// acceptable keeps normal not-found and duplicate-key outcomes from opening
// the go-zero breaker.
func acceptable(err error) bool {
	if err == nil || errors.Is(err, sqlx.ErrNotFound) || errors.Is(err, model.ErrAbortTransaction) {
		return true
	}
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
