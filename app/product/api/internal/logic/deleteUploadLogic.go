package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUploadLogic {
	return &DeleteUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUploadLogic) DeleteUpload(req *types.UploadNameReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
