package logic

import (
	"context"

	"go-zero-demo/user-api/internal/svc"
	"go-zero-demo/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// 演示用精简行政区划数据。生产环境可替换为完整的省市区数据源。
var regionData = []types.RegionInfo{
	{Id: 110000, ParentId: 0, Name: "北京市", Level: 1},
	{Id: 110100, ParentId: 110000, Name: "北京市", Level: 2},
	{Id: 110101, ParentId: 110100, Name: "东城区", Level: 3},
	{Id: 110102, ParentId: 110100, Name: "西城区", Level: 3},
	{Id: 110105, ParentId: 110100, Name: "朝阳区", Level: 3},
	{Id: 110108, ParentId: 110100, Name: "海淀区", Level: 3},

	{Id: 310000, ParentId: 0, Name: "上海市", Level: 1},
	{Id: 310100, ParentId: 310000, Name: "上海市", Level: 2},
	{Id: 310101, ParentId: 310100, Name: "黄浦区", Level: 3},
	{Id: 310104, ParentId: 310100, Name: "徐汇区", Level: 3},
	{Id: 310106, ParentId: 310100, Name: "静安区", Level: 3},
	{Id: 310115, ParentId: 310100, Name: "浦东新区", Level: 3},

	{Id: 440000, ParentId: 0, Name: "广东省", Level: 1},
	{Id: 440100, ParentId: 440000, Name: "广州市", Level: 2},
	{Id: 440104, ParentId: 440100, Name: "越秀区", Level: 3},
	{Id: 440105, ParentId: 440100, Name: "海珠区", Level: 3},
	{Id: 440106, ParentId: 440100, Name: "天河区", Level: 3},
	{Id: 440300, ParentId: 440000, Name: "深圳市", Level: 2},
	{Id: 440303, ParentId: 440300, Name: "罗湖区", Level: 3},
	{Id: 440304, ParentId: 440300, Name: "福田区", Level: 3},
	{Id: 440305, ParentId: 440300, Name: "南山区", Level: 3},

	{Id: 330000, ParentId: 0, Name: "浙江省", Level: 1},
	{Id: 330100, ParentId: 330000, Name: "杭州市", Level: 2},
	{Id: 330102, ParentId: 330100, Name: "上城区", Level: 3},
	{Id: 330106, ParentId: 330100, Name: "西湖区", Level: 3},
	{Id: 330108, ParentId: 330100, Name: "滨江区", Level: 3},
	{Id: 330200, ParentId: 330000, Name: "宁波市", Level: 2},
	{Id: 330203, ParentId: 330200, Name: "海曙区", Level: 3},
	{Id: 330212, ParentId: 330200, Name: "鄞州区", Level: 3},

	{Id: 320000, ParentId: 0, Name: "江苏省", Level: 1},
	{Id: 320100, ParentId: 320000, Name: "南京市", Level: 2},
	{Id: 320102, ParentId: 320100, Name: "玄武区", Level: 3},
	{Id: 320104, ParentId: 320100, Name: "秦淮区", Level: 3},
	{Id: 320106, ParentId: 320100, Name: "鼓楼区", Level: 3},
	{Id: 320500, ParentId: 320000, Name: "苏州市", Level: 2},
	{Id: 320506, ParentId: 320500, Name: "吴中区", Level: 3},
	{Id: 320508, ParentId: 320500, Name: "姑苏区", Level: 3},

	{Id: 510000, ParentId: 0, Name: "四川省", Level: 1},
	{Id: 510100, ParentId: 510000, Name: "成都市", Level: 2},
	{Id: 510104, ParentId: 510100, Name: "锦江区", Level: 3},
	{Id: 510105, ParentId: 510100, Name: "青羊区", Level: 3},
	{Id: 510107, ParentId: 510100, Name: "武侯区", Level: 3},

	{Id: 420000, ParentId: 0, Name: "湖北省", Level: 1},
	{Id: 420100, ParentId: 420000, Name: "武汉市", Level: 2},
	{Id: 420102, ParentId: 420100, Name: "江岸区", Level: 3},
	{Id: 420106, ParentId: 420100, Name: "武昌区", Level: 3},
	{Id: 420111, ParentId: 420100, Name: "洪山区", Level: 3},

	{Id: 610000, ParentId: 0, Name: "陕西省", Level: 1},
	{Id: 610100, ParentId: 610000, Name: "西安市", Level: 2},
	{Id: 610103, ParentId: 610100, Name: "碑林区", Level: 3},
	{Id: 610104, ParentId: 610100, Name: "莲湖区", Level: 3},
	{Id: 610113, ParentId: 610100, Name: "雁塔区", Level: 3},
}

type ListRegionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionsLogic {
	return &ListRegionsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListRegionsLogic) ListRegions(req *types.ListRegionsReq) (*types.ListRegionsResp, error) {
	data := make([]types.RegionInfo, 0, len(regionData))
	for _, region := range regionData {
		if region.ParentId == req.ParentId {
			data = append(data, region)
		}
	}
	return &types.ListRegionsResp{Code: 200, Msg: "ok", Data: data}, nil
}
