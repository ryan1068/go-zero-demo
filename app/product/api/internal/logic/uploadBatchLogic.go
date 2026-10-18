package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadBatchLogic {
	return &UploadBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadBatchLogic) UploadBatch() (resp *types.UploadBatchResp, err error) {
	// todo: add your logic here and delete this line

	return
}
