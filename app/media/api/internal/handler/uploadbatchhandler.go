package handler

import (
	"net/http"

	"go-zero-demo/media-api/internal/logic"
	"go-zero-demo/media-api/internal/svc"
	"go-zero-demo/media-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
			httpx.OkJsonCtx(r.Context(), w, &types.UploadBatchResp{Code: http.StatusBadRequest, Msg: "请选择要上传的图片"})
			return
		}
		if r.MultipartForm == nil || len(r.MultipartForm.File["files"]) == 0 {
			httpx.OkJsonCtx(r.Context(), w, &types.UploadBatchResp{Code: http.StatusBadRequest, Msg: "请选择要上传的图片"})
			return
		}

		resp, err := logic.NewUploadLogic(r.Context(), svcCtx.Config.Upload, svcCtx.Storage, svcCtx.MediaFiles).
			SaveBatch(r.MultipartForm.File["files"])
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
