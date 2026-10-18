package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"go-zero-demo/media-api/internal/config"
	"go-zero-demo/media-api/internal/model"
	"go-zero-demo/media-api/internal/storage"
	"go-zero-demo/media-api/internal/types"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx     context.Context
	cfg     config.UploadConfig
	storage storage.Storage
	files   model.MediaFileModel
}

func NewUploadLogic(ctx context.Context, cfg config.UploadConfig, store storage.Storage, files model.MediaFileModel) *UploadLogic {
	return &UploadLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		cfg:     cfg,
		storage: store,
		files:   files,
	}
}

func (l *UploadLogic) Save(file multipart.File, header *multipart.FileHeader) (*types.UploadResp, error) {
	limitMB := l.cfg.MaxSizeMB
	if limitMB <= 0 {
		limitMB = 10
	}
	maxBytes := limitMB * 1024 * 1024
	if header.Size > maxBytes {
		return badUploadResp(errcode.BadRequest, fmt.Sprintf("图片不能超过 %dMB", limitMB)), nil
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExt(ext) {
		return badUploadResp(errcode.BadRequest, "只支持 jpg/jpeg/png/gif/webp 图片"), nil
	}

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		l.Errorf("read upload header failed: %v", err)
		return nil, err
	}
	contentType := http.DetectContentType(buffer[:n])
	if !strings.HasPrefix(contentType, "image/") {
		return badUploadResp(errcode.BadRequest, "文件不是有效图片"), nil
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		l.Errorf("seek upload file failed: %v", err)
		return nil, err
	}

	uploaderID, ok := auth.UserID(l.ctx)
	if !ok {
		return badUploadResp(errcode.Unauthorized, errcode.Msg(errcode.Unauthorized)), nil
	}

	var random [8]byte
	if _, err := rand.Read(random[:]); err != nil {
		l.Errorf("generate upload name failed: %v", err)
		return nil, err
	}
	name := hex.EncodeToString(random[:]) + ext

	written, err := l.storage.Save(l.ctx, name, file, maxBytes)
	if err != nil {
		l.Errorf("save upload file failed: %v", err)
		return nil, err
	}

	if _, err := l.files.Insert(l.ctx, &model.MediaFile{
		FileName:     name,
		OriginalName: header.Filename,
		Mime:         contentType,
		Size:         written,
		UploaderId:   uploaderID,
	}); err != nil {
		if deleteErr := l.storage.Delete(name); deleteErr != nil {
			l.Errorf("remove upload after metadata failure: %v", deleteErr)
		}
		l.Errorf("write media metadata failed: %v", err)
		return nil, err
	}

	prefix := l.cfg.URLPrefix
	if prefix == "" {
		prefix = "/api/media/files/"
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	return &types.UploadResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		Url:  prefix + name,
	}, nil
}

func (l *UploadLogic) SaveBatch(headers []*multipart.FileHeader) (*types.UploadBatchResp, error) {
	if len(headers) == 0 {
		return &types.UploadBatchResp{Code: errcode.BadRequest, Msg: "请选择图片"}, nil
	}

	urls := make([]string, 0, len(headers))
	for _, header := range headers {
		file, err := header.Open()
		if err != nil {
			l.Errorf("open upload file failed: %v", err)
			return nil, err
		}
		resp, saveErr := l.Save(file, header)
		closeErr := file.Close()
		if saveErr != nil {
			return nil, saveErr
		}
		if closeErr != nil {
			l.Errorf("close upload file failed: %v", closeErr)
		}
		if resp.Code != errcode.Success {
			return &types.UploadBatchResp{Code: resp.Code, Msg: resp.Msg}, nil
		}
		urls = append(urls, resp.Url)
	}

	return &types.UploadBatchResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		Urls: urls,
	}, nil
}

func badUploadResp(code int32, msg string) *types.UploadResp {
	return &types.UploadResp{Code: code, Msg: msg}
}

func allowedImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}
