package svc

import (
	"errors"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-demo/product-rpc/internal/config"
	"go-zero-demo/product-rpc/internal/model"
)

type ServiceContext struct {
	Config config.Config
	DB     sqlx.SqlConn

	CategoryModel   model.CategoryModel
	ProductModel    model.ProductModel
	ProductSkuModel model.ProductSkuModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.MustNewConn(c.Mysql, sqlx.WithAcceptable(acceptable))
	return &ServiceContext{
		Config:          c,
		DB:              conn,
		CategoryModel:   model.NewCategoryModel(conn),
		ProductModel:    model.NewProductModel(conn),
		ProductSkuModel: model.NewProductSkuModel(conn),
	}
}

// acceptable keeps normal not-found and duplicate-key outcomes from opening
// the go-zero breaker.
func acceptable(err error) bool {
	if err == nil || errors.Is(err, sqlx.ErrNotFound) {
		return true
	}
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
