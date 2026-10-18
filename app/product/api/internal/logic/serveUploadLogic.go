package logic

import (
	"context"

	"go-zero-demo/product-api/internal/svc"
	"go-zero-demo/product-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServeUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServeUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServeUploadLogic {
	return &ServeUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServeUploadLogic) ServeUpload(req *types.UploadNameReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
