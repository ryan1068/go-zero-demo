package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-demo/product-api/internal/logic"
	"go-zero-demo/product-api/internal/svc"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maxBytes := svcCtx.Config.Upload.MaxSizeMB * 1024 * 1024
		if maxBytes <= 0 {
			maxBytes = 5 * 1024 * 1024
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes+1024*1024)
		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("请求参数错误"))
			return
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("failed to close upload file: %v", closeErr)
			}
		}()

		resp, err := logic.NewUploadLogic(r.Context(), svcCtx).Save(file, header)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
