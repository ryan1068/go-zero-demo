package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go-zero-demo/pkg/errcode"
	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UploadLogic) Save(file multipart.File, header *multipart.FileHeader) (*types.UploadResp, error) {
	maxBytes := l.svcCtx.Config.Upload.MaxSizeMB * 1024 * 1024
	if maxBytes <= 0 {
		maxBytes = 5 * 1024 * 1024
	}
	if header.Size > maxBytes {
		return &types.UploadResp{Code: errcode.BadRequest, Msg: "图片不能超过 5MB"}, nil
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExt(ext) {
		return &types.UploadResp{Code: errcode.BadRequest, Msg: "只支持 jpg/jpeg/png/gif/webp 图片"}, nil
	}

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		l.Errorf("read upload header failed: %v", err)
		return nil, err
	}
	contentType := http.DetectContentType(buffer[:n])
	if !strings.HasPrefix(contentType, "image/") {
		return &types.UploadResp{Code: errcode.BadRequest, Msg: "文件不是有效图片"}, nil
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		l.Errorf("seek upload file failed: %v", err)
		return nil, err
	}

	var random [8]byte
	if _, err := rand.Read(random[:]); err != nil {
		l.Errorf("generate upload name failed: %v", err)
		return nil, err
	}
	name := hex.EncodeToString(random[:]) + ext

	dir := l.svcCtx.Config.Upload.Dir
	if dir == "" {
		dir = "./data/uploads"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		l.Errorf("create upload dir failed: %v", err)
		return nil, err
	}
	path := filepath.Join(dir, name)
	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		l.Errorf("create upload file failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := dst.Close(); closeErr != nil {
			l.Errorf("close upload file failed: %v", closeErr)
		}
	}()

	written, err := io.Copy(dst, io.LimitReader(file, maxBytes+1))
	if err != nil {
		l.Errorf("save upload file failed: %v", err)
		if removeErr := os.Remove(path); removeErr != nil {
			l.Errorf("remove partial upload file failed: %v", removeErr)
		}
		return nil, err
	}
	if written > maxBytes {
		if removeErr := os.Remove(path); removeErr != nil {
			l.Errorf("remove oversized upload file failed: %v", removeErr)
		}
		return &types.UploadResp{Code: errcode.BadRequest, Msg: "图片不能超过 5MB"}, nil
	}

	prefix := l.svcCtx.Config.Upload.URLPrefix
	if prefix == "" {
		prefix = "/uploads/"
	}
	return &types.UploadResp{
		Code: errcode.Success,
		Msg:  errcode.Msg(errcode.Success),
		Url:  prefix + name,
	}, nil
}

func allowedImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}
