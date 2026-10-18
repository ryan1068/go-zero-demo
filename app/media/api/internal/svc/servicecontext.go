package svc

import (
	"go-zero-demo/media-api/internal/config"
	"go-zero-demo/media-api/internal/model"
	"go-zero-demo/media-api/internal/storage"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	DB         sqlx.SqlConn
	MediaFiles model.MediaFileModel
	Storage    storage.Storage
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.MustNewConn(c.Mysql)
	return &ServiceContext{
		Config:     c,
		DB:         conn,
		MediaFiles: model.NewMediaFileModel(conn),
		Storage:    storage.NewLocal(c.Upload.Dir),
	}
}
