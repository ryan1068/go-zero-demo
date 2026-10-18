package handler

import (
	"net/http"

	"go-zero-demo/product-api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/categories", Handler: ListCategoriesHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/products", Handler: ListProductsHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/products/suggest", Handler: SuggestProductsHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/products/hot", Handler: HotWordsHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/product/:id", Handler: GetProductHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/product/:id/skus", Handler: ListSkusHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/uploads/:name", Handler: ServeUploadHandler(serverCtx)},

		{Method: http.MethodGet, Path: "/admin/categories", Handler: adminOnly(AdminListCategoriesHandler(serverCtx))},
		{Method: http.MethodPost, Path: "/admin/categories", Handler: adminOnly(SaveCategoryHandler(serverCtx))},
		{Method: http.MethodPut, Path: "/admin/categories/:id", Handler: adminOnly(SaveCategoryHandler(serverCtx))},
		{Method: http.MethodDelete, Path: "/admin/categories/:id", Handler: superAdminOnly(DeleteCategoryHandler(serverCtx))},

		{Method: http.MethodGet, Path: "/admin/products/stats", Handler: adminOnly(ProductStatsHandler(serverCtx))},
		{Method: http.MethodGet, Path: "/admin/products", Handler: adminOnly(AdminListProductsHandler(serverCtx))},
		{Method: http.MethodGet, Path: "/admin/products/:id", Handler: adminOnly(AdminGetProductHandler(serverCtx))},
		{Method: http.MethodGet, Path: "/admin/products/:id/skus", Handler: adminOnly(ListSkusHandler(serverCtx))},
		{Method: http.MethodPost, Path: "/admin/skus", Handler: adminOnly(SaveSkuHandler(serverCtx))},
		{Method: http.MethodDelete, Path: "/admin/skus/:id", Handler: superAdminOnly(DeleteSkuHandler(serverCtx))},
		{Method: http.MethodPost, Path: "/admin/products", Handler: adminOnly(SaveProductHandler(serverCtx))},
		{Method: http.MethodPut, Path: "/admin/products/:id", Handler: adminOnly(SaveProductHandler(serverCtx))},
		{Method: http.MethodDelete, Path: "/admin/products/:id", Handler: superAdminOnly(DeleteProductHandler(serverCtx))},
		{Method: http.MethodPut, Path: "/admin/products/:id/status", Handler: adminOnly(UpdateProductStatusHandler(serverCtx))},
		{Method: http.MethodPost, Path: "/admin/upload", Handler: adminOnly(UploadHandler(serverCtx))},
		{Method: http.MethodPost, Path: "/admin/uploads/batch", Handler: adminOnly(UploadBatchHandler(serverCtx))},
		{Method: http.MethodDelete, Path: "/admin/uploads/:name", Handler: adminOnly(DeleteUploadHandler(serverCtx))},
	})
}
