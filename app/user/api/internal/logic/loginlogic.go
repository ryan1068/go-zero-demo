package logic

import (
	"context"
	"strconv"
	"strings"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	loginFailLimit         = 5
	loginFailWindowSeconds = 15 * 60
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
	key := loginFailKey(req.IP, req.Username)
	if l.svcCtx.Redis != nil && req.IP != "" {
		exists, checkErr := l.svcCtx.Redis.ExistsCtx(l.ctx, key)
		if checkErr != nil {
			l.Errorf("check login fail count failed: %v", checkErr)
		} else if exists {
			countStr, readErr := l.svcCtx.Redis.GetCtx(l.ctx, key)
			if readErr != nil {
				l.Errorf("read login fail count failed: %v", readErr)
			} else if count, _ := strconv.Atoi(countStr); count >= loginFailLimit {
				return &types.LoginResp{Code: 429, Msg: "登录失败次数过多，请 15 分钟后再试"}, nil
			}
		}
	}

	s, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		l.Errorf("login rpc failed: %v", err)
		return nil, err
	}

	if l.svcCtx.Redis != nil && req.IP != "" {
		if s.Code == 200 {
			if _, delErr := l.svcCtx.Redis.DelCtx(l.ctx, key); delErr != nil {
				l.Errorf("clear login fail count failed: %v", delErr)
			}
		} else {
			count, incrErr := l.svcCtx.Redis.IncrCtx(l.ctx, key)
			if incrErr != nil {
				l.Errorf("increase login fail count failed: %v", incrErr)
			} else if count == 1 {
				if expireErr := l.svcCtx.Redis.ExpireCtx(l.ctx, key, loginFailWindowSeconds); expireErr != nil {
					l.Errorf("expire login fail count failed: %v", expireErr)
				}
			}
		}
	}

	return &types.LoginResp{
		Code: s.Code,
		Msg:  s.Msg,
		Data: types.LoginData{
			Token:    s.Token,
			UserId:   s.UserId,
			Username: s.Username,
			Role:     s.Role,
		},
	}, nil
}

func loginFailKey(ip, username string) string {
	return "login:fail:" + ip + ":" + strings.ToLower(strings.TrimSpace(username))
}
