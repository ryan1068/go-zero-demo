package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-demo/order-api/internal/logic"
	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
)

func ApplyAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplyAfterSaleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		resp, err := logic.NewApplyAfterSaleLogic(r.Context(), svcCtx).ApplyAfterSale(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func listAfterSales(svcCtx *svc.ServiceContext, admin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListAfterSalesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		resp, err := logic.NewListAfterSalesLogic(r.Context(), svcCtx).ListAfterSales(&req, admin)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func ListAfterSalesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return listAfterSales(svcCtx, false)
}

func AdminListAfterSalesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return listAfterSales(svcCtx, true)
}

func getAfterSale(svcCtx *svc.ServiceContext, admin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AfterSaleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		resp, err := logic.NewGetAfterSaleLogic(r.Context(), svcCtx).GetAfterSale(req.Id, admin)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func GetAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return getAfterSale(svcCtx, false)
}

func AdminGetAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return getAfterSale(svcCtx, true)
}

func HandleAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HandleAfterSaleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		resp, err := logic.NewHandleAfterSaleLogic(r.Context(), svcCtx).HandleAfterSale(
			req.Id, req.Action, req.AdminRemark, req.TrackingNo, true,
		)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func ShipAfterSaleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShipAfterSaleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		resp, err := logic.NewHandleAfterSaleLogic(r.Context(), svcCtx).HandleAfterSale(
			req.Id, "ship", "", req.TrackingNo, false,
		)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
