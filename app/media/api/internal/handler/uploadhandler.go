package handler

import (
	"net/http"

	"go-zero-demo/media-api/internal/logic"
	"go-zero-demo/media-api/internal/svc"
	"go-zero-demo/media-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

const maxUploadMemory = 32 << 20

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
			httpx.OkJsonCtx(r.Context(), w, &types.UploadResp{Code: http.StatusBadRequest, Msg: "请选择要上传的图片"})
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, &types.UploadResp{Code: http.StatusBadRequest, Msg: "请选择要上传的图片"})
			return
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				httpx.ErrorCtx(r.Context(), w, closeErr)
			}
		}()

		resp, err := logic.NewUploadLogic(r.Context(), svcCtx.Config.Upload, svcCtx.Storage, svcCtx.MediaFiles).Save(file, header)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
