package logic

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"go-zero-demo/media-api/internal/storage"
	"go-zero-demo/media-api/internal/types"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx     context.Context
	storage storage.Storage
}

func NewDeleteFileLogic(ctx context.Context, store storage.Storage) *DeleteFileLogic {
	return &DeleteFileLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		storage: store,
	}
}

func (l *DeleteFileLogic) Delete(name string) (*types.BaseResp, error) {
	if name == "" || filepath.Base(name) != name || strings.Contains(name, "..") {
		return &types.BaseResp{Code: errcode.BadRequest, Msg: "文件名不合法"}, nil
	}

	if _, ok := auth.UserID(l.ctx); !ok {
		return &types.BaseResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	role, _ := auth.Role(l.ctx)
	if !auth.IsAdmin(role) {
		return &types.BaseResp{Code: errcode.Forbidden, Msg: errcode.Msg(errcode.Forbidden)}, nil
	}

	if err := l.storage.Delete(name); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &types.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
		}
		l.Errorf("delete media file failed: %v", err)
		return nil, err
	}

	return &types.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
