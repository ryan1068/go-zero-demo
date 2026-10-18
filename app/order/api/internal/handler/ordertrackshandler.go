package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-demo/order-api/internal/logic"
	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
)

func GetOrderTracksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return getTracks(svcCtx, false)
}

func AdminGetOrderTracksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return getTracks(svcCtx, true)
}

func getTracks(svcCtx *svc.ServiceContext, admin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		resp, err := logic.NewGetOrderTracksLogic(r.Context(), svcCtx).GetOrderTracks(req.Id, admin)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
