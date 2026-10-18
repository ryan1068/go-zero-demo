package logic

import (
	"context"
	"strings"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"
	"go-zero-demo/product-rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type SuggestProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSuggestProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SuggestProductsLogic {
	return &SuggestProductsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SuggestProductsLogic) SuggestProducts(req *types.SuggestReq) (*types.SuggestResp, error) {
	keyword := strings.TrimSpace(req.Keyword)
	limit := req.Limit
	if limit < 1 {
		limit = 5
	}
	if limit > 10 {
		limit = 10
	}
	if keyword == "" {
		return &types.SuggestResp{Code: 200, Msg: "ok", Data: []string{}}, nil
	}
	s, err := l.svcCtx.ProductRpc.ListProducts(l.ctx, &product.ListProductsReq{
		Page:     1,
		PageSize: limit,
		Keyword:  keyword,
	})
	if err != nil {
		l.Errorf("suggest products rpc failed: %v", err)
		return nil, err
	}
	names := make([]string, 0, len(s.List))
	for _, item := range s.List {
		names = append(names, item.Name)
	}
	return &types.SuggestResp{Code: s.Code, Msg: s.Msg, Data: names}, nil
}
