package svc

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-demo/order-rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *sql.DB
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		logx.Must(err)
	}
	if err := db.Ping(); err != nil {
		logx.Must(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  redis.MustNewRedis(c.CacheRedis),
	}
}
