package logic

import (
	"context"
	"sort"
	"time"

	"go-zero-demo/order-api/internal/svc"
	"go-zero-demo/order-api/internal/types"
	"go-zero-demo/order-rpc/pb/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderTracksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderTracksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderTracksLogic {
	return &GetOrderTracksLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetOrderTracks returns a timeline derived from the order timestamps. It is a
// demo-friendly logistics view until a real carrier API is integrated.
func (l *GetOrderTracksLogic) GetOrderTracks(id int64, admin bool) (*types.OrderTracksResp, error) {
	ctx, ok := injectUser(l.ctx)
	if !ok {
		return &types.OrderTracksResp{Code: 401, Msg: "未授权"}, nil
	}
	adminFlag := int32(0)
	if admin {
		adminFlag = 1
	}
	s, err := l.svcCtx.OrderRpc.GetOrder(ctx, &order.GetOrderReq{Id: id, Admin: adminFlag})
	if err != nil {
		l.Errorf("get order rpc failed: %v", err)
		return nil, err
	}
	if s.Code != 200 || s.Data == nil {
		return &types.OrderTracksResp{Code: s.Code, Msg: s.Msg}, nil
	}

	info := s.Data
	tracks := make([]types.OrderTrack, 0, 6)
	add := func(t int64, status, description string) {
		if t > 0 {
			tracks = append(tracks, types.OrderTrack{Time: t, Status: status, Description: description})
		}
	}

	add(info.CreatedAt, "created", "订单已创建")
	add(info.PaidAt, "paid", "支付成功")
	if info.ShippedAt > 0 {
		desc := "商家已发货"
		if info.LogisticsCompany != "" || info.TrackingNo != "" {
			desc += " · " + info.LogisticsCompany + " " + info.TrackingNo
		}
		add(info.ShippedAt, "shipped", desc)
		now := time.Now().Unix()
		if now >= info.ShippedAt+12*3600 {
			add(info.ShippedAt+12*3600, "transit", "快件已到达分拨中心")
		}
		if now >= info.ShippedAt+24*3600 {
			add(info.ShippedAt+24*3600, "delivering", "快件正在派送")
		}
	}
	if info.CompletedAt > 0 {
		add(info.CompletedAt, "completed", "已确认收货")
	}
	if info.CanceledAt > 0 {
		desc := "订单已取消"
		if info.CancelReason != "" {
			desc += " · " + info.CancelReason
		}
		add(info.CanceledAt, "canceled", desc)
	}

	sort.Slice(tracks, func(i, j int) bool { return tracks[i].Time > tracks[j].Time })
	return &types.OrderTracksResp{
		Code: 200,
		Msg:  "ok",
		Data: types.OrderTracksData{
			OrderId:          info.Id,
			OrderNo:          info.OrderNo,
			LogisticsCompany: info.LogisticsCompany,
			TrackingNo:       info.TrackingNo,
			Tracks:           tracks,
		},
	}, nil
}
