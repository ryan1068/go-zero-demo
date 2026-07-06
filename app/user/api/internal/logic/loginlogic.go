package logic

import (
	"context"
	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	//s, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
	//	Username: req.Username,
	//	Password: req.Password,
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, err
	//}

	resp = new(types.LoginResp)
	resp.Msg = "s.Msg"
	resp.Code = 200

	return
}
