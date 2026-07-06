package svc

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"go-zero-demo/user-rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *sql.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
