package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/product-api/internal/logic"
	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
)

func adminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := auth.Role(r.Context())
		if !ok || !auth.IsAdmin(role) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(map[string]any{"code": http.StatusForbidden, "msg": "无权限"}); err != nil {
				log.Printf("failed to write forbidden response: %v", err)
			}
			return
		}
		next(w, r)
	}
}

func superAdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := auth.Role(r.Context())
		if !ok || !auth.IsSuperAdmin(role) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(map[string]any{"code": http.StatusForbidden, "msg": "需要超级管理员权限"}); err != nil {
				log.Printf("failed to write forbidden response: %v", err)
			}
			return
		}
		next(w, r)
	}
}

func listCategories(svcCtx *svc.ServiceContext, admin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := logic.NewListCategoriesLogic(r.Context(), svcCtx).ListCategories(admin)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func ListCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return listCategories(svcCtx, false)
}

func AdminListCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return listCategories(svcCtx, true)
}

func SaveCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := logic.NewSaveCategoryLogic(r.Context(), svcCtx).SaveCategory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func DeleteCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PathIdReq
		if err := httpx.Parse(r, &req); err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := logic.NewDeleteCategoryLogic(r.Context(), svcCtx).DeleteCategory(req.Id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func listProducts(svcCtx *svc.ServiceContext, admin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListProductsReq
		if err := httpx.Parse(r, &req); err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := logic.NewListProductsLogic(r.Context(), svcCtx).ListProducts(&req, admin)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func ListProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return listProducts(svcCtx, false)
}

func AdminListProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return listProducts(svcCtx, true)
}

func getProduct(svcCtx *svc.ServiceContext, admin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PathIdReq
		if err := httpx.Parse(r, &req); err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := logic.NewGetProductLogic(r.Context(), svcCtx).GetProduct(req.Id, admin)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func GetProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return getProduct(svcCtx, false)
}

func AdminGetProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return getProduct(svcCtx, true)
}

func SaveProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveProductReq
		if err := httpx.Parse(r, &req); err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := logic.NewSaveProductLogic(r.Context(), svcCtx).SaveProduct(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func DeleteProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PathIdReq
		if err := httpx.Parse(r, &req); err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := logic.NewDeleteProductLogic(r.Context(), svcCtx).DeleteProduct(req.Id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func UpdateProductStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateProductStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			writeBadRequest(w)
			return
		}
		resp, err := logic.NewUpdateProductStatusLogic(r.Context(), svcCtx).UpdateProductStatus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func ProductStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := logic.NewProductStatsLogic(r.Context(), svcCtx).ProductStats()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maxBytes := svcCtx.Config.Upload.MaxSizeMB * 1024 * 1024
		if maxBytes <= 0 {
			maxBytes = 5 * 1024 * 1024
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes+1024*1024)
		file, header, err := r.FormFile("file")
		if err != nil {
			writeBadRequest(w)
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

func UploadBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maxBytes := svcCtx.Config.Upload.MaxSizeMB * 1024 * 1024
		if maxBytes <= 0 {
			maxBytes = 5 * 1024 * 1024
		}
		const maxFiles = 8
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes*maxFiles+1024*1024)
		if err := r.ParseMultipartForm(maxBytes * maxFiles); err != nil {
			writeBadRequest(w)
			return
		}
		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			files = r.MultipartForm.File["file"]
		}
		if len(files) == 0 || len(files) > maxFiles {
			writeBadRequest(w)
			return
		}

		urls := make([]string, 0, len(files))
		for _, header := range files {
			file, err := header.Open()
			if err != nil {
				writeBadRequest(w)
				return
			}
			resp, saveErr := logic.NewUploadLogic(r.Context(), svcCtx).Save(file, header)
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("failed to close upload file: %v", closeErr)
			}
			if saveErr != nil {
				httpx.ErrorCtx(r.Context(), w, saveErr)
				return
			}
			if resp.Code != http.StatusOK {
				httpx.OkJsonCtx(r.Context(), w, resp)
				return
			}
			urls = append(urls, resp.Url)
		}
		httpx.OkJsonCtx(r.Context(), w, &types.UploadBatchResp{Code: http.StatusOK, Msg: "ok", Urls: urls})
	}
}

func DeleteUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadNameReq
		if err := httpx.Parse(r, &req); err != nil {
			writeBadRequest(w)
			return
		}
		name := filepath.Base(req.Name)
		if name != req.Name || !validUploadName(name) {
			http.NotFound(w, r)
			return
		}
		dir := svcCtx.Config.Upload.Dir
		if dir == "" {
			dir = "./data/uploads"
		}
		if err := os.Remove(filepath.Join(dir, name)); err != nil && !os.IsNotExist(err) {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, &types.UploadResp{Code: http.StatusOK, Msg: "ok"})
	}
}

func validUploadName(name string) bool {
	ok, err := regexp.MatchString(`^[a-f0-9]{16}\.(jpg|jpeg|png|gif|webp)$`, name)
	return err == nil && ok
}

func ServeUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadNameReq
		if err := httpx.Parse(r, &req); err != nil {
			writeBadRequest(w)
			return
		}
		name := filepath.Base(req.Name)
		if name != req.Name || strings.Contains(name, "..") || name == "." {
			http.NotFound(w, r)
			return
		}
		dir := svcCtx.Config.Upload.Dir
		if dir == "" {
			dir = "./data/uploads"
		}
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(w, r, filepath.Join(dir, name))
	}
}

func writeBadRequest(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(map[string]any{"code": http.StatusBadRequest, "msg": "请求参数错误"}); err != nil {
		log.Printf("failed to write bad request response: %v", err)
	}
}
