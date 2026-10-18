package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-demo/product-api/internal/logic"
	"go-zero-demo/product-api/internal/svc"
)

func UploadBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maxBytes := svcCtx.Config.Upload.MaxSizeMB * 1024 * 1024
		if maxBytes <= 0 {
			maxBytes = 5 * 1024 * 1024
		}
		const maxFiles = 8
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes*maxFiles+1024*1024)
		if err := r.ParseMultipartForm(maxBytes * maxFiles); err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("请求参数错误"))
			return
		}
		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			files = r.MultipartForm.File["file"]
		}
		if len(files) == 0 || len(files) > maxFiles {
			httpx.ErrorCtx(r.Context(), w, errors.New("请求参数错误"))
			return
		}

		resp, err := logic.NewUploadLogic(r.Context(), svcCtx).SaveBatch(files)
		if err != nil {
			log.Printf("failed to save upload batch: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
