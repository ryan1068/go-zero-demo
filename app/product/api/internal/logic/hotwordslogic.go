package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HotWordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHotWordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HotWordsLogic {
	return &HotWordsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *HotWordsLogic) HotWords() (*types.HotWordsResp, error) {
	return &types.HotWordsResp{
		Code: 200,
		Msg:  "ok",
		Data: []string{"Go 语言", "蓝牙耳机", "手冲咖啡", "明前龙井"},
	}, nil
}
