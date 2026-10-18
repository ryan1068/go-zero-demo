package logic

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"go-zero-demo/media-api/internal/model"
	"go-zero-demo/media-api/internal/storage"
	"go-zero-demo/media-api/internal/types"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx     context.Context
	storage storage.Storage
	files   model.MediaFileModel
}

func NewDeleteFileLogic(ctx context.Context, store storage.Storage, files model.MediaFileModel) *DeleteFileLogic {
	return &DeleteFileLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		storage: store,
		files:   files,
	}
}

func (l *DeleteFileLogic) Delete(name string) (*types.BaseResp, error) {
	if name == "" || filepath.Base(name) != name || strings.Contains(name, "..") {
		return &types.BaseResp{Code: errcode.BadRequest, Msg: "文件名不合法"}, nil
	}

	userID, ok := auth.UserID(l.ctx)
	if !ok {
		return &types.BaseResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	role, _ := auth.Role(l.ctx)

	file, err := l.files.FindActiveByName(l.ctx, name)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return &types.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
		}
		l.Errorf("find media file failed: %v", err)
		return nil, err
	}

	if file.UploaderId != userID && !auth.IsAdmin(role) {
		return &types.BaseResp{Code: errcode.Forbidden, Msg: errcode.Msg(errcode.Forbidden)}, nil
	}

	if err := l.storage.Delete(name); err != nil {
		l.Errorf("delete media file failed: %v", err)
		return nil, err
	}

	affected, err := l.files.SoftDelete(l.ctx, name)
	if err != nil {
		l.Errorf("soft delete media metadata failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &types.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}

	return &types.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}
