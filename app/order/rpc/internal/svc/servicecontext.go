package svc

import (
	"errors"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-demo/order-rpc/internal/config"
	"go-zero-demo/order-rpc/internal/model"
)

type ServiceContext struct {
	Config config.Config
	DB     sqlx.SqlConn
	Redis  *redis.Redis

	OrdersModel        model.OrdersModel
	OrderItemModel     model.OrderItemModel
	CartItemModel      model.CartItemModel
	PaymentRecordModel model.PaymentRecordModel
	OrderLogModel      model.OrderLogModel
	AfterSaleModel     model.AfterSaleModel
	RefundRecordModel  model.RefundRecordModel

	ProductModel     *model.CrossProductModel
	UserAddressModel *model.CrossUserAddressModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.MustNewConn(c.Mysql, sqlx.WithAcceptable(acceptable))
	return &ServiceContext{
		Config:             c,
		DB:                 conn,
		Redis:              redis.MustNewRedis(c.CacheRedis),
		OrdersModel:        model.NewOrdersModel(conn),
		OrderItemModel:     model.NewOrderItemModel(conn),
		CartItemModel:      model.NewCartItemModel(conn),
		PaymentRecordModel: model.NewPaymentRecordModel(conn),
		OrderLogModel:      model.NewOrderLogModel(conn),
		AfterSaleModel:     model.NewAfterSaleModel(conn),
		RefundRecordModel:  model.NewRefundRecordModel(conn),
		ProductModel:       model.NewCrossProductModel(conn),
		UserAddressModel:   model.NewCrossUserAddressModel(conn),
	}
}

// acceptable keeps normal not-found, duplicate-key and business rollback
// outcomes from opening the go-zero breaker.
func acceptable(err error) bool {
	if err == nil || errors.Is(err, sqlx.ErrNotFound) || errors.Is(err, model.ErrAbortTransaction) {
		return true
	}
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
