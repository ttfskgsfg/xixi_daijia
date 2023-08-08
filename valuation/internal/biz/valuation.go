package biz

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
	"log"
	"strconv"
	"valuation/api/mapService"
)

type PriceRule struct {
	gorm.Model
	PriceRuleWork
}

type PriceRuleWork struct {
	CityID      uint  `gorm:"" json:"city_id,omitempty"`
	StartFee    int64 `gorm:"" json:"start_fee,omitempty"`
	DistanceFee int64 `gorm:"" json:"distance_fee,omitempty"`
	DurationFee int64 `gorm:"" json:"duration_fee,omitempty"`
	StartAt     int   `gorm:"type:int" json:"start_at,omitempty"`
	EndAt       int   `gorm:"type:int" json:"end_at,omitempty"`
}

// 定义操作priceRule的接口
type PriceRuleInterface interface {
	//获取规则
	GetRule(cityid uint, curr int) (*PriceRule, error)
}

type ValuationBiz struct {
	pri PriceRuleInterface
}

func NewValuationBiz(pri PriceRuleInterface) *ValuationBiz {
	return &ValuationBiz{
		pri: pri,
	}
}

// 获取价格
func (vb *ValuationBiz) GetPrice(ctx context.Context, distance, duration string, cityId uint,
	curr int) (int64, error) {
	//一、获取规则
	rule, err := vb.pri.GetRule(cityId, curr)
	if err != nil {
		return 0, err
	}
	//将距离与时长转换
	distancInt64, _ := strconv.ParseInt(distance, 10, 64)
	durationInt64, _ := strconv.ParseInt(distance, 10, 64)
	//二、基于rule计算
	distancInt64 /= 1000
	durationInt64 /= 60
	var startDistance int64 = 5
	total := rule.StartFee +
		rule.DurationFee*(distancInt64-startDistance) +
		rule.DurationFee*durationInt64
	return total, nil
}

// 获取时长和距离
func (*ValuationBiz) GetDrivingInfo(ctx context.Context, origin, destination string) (distance string, duration string, err error) {
	//1、获取consul客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "localhost:8500"
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}
	//2、获取consul发现管理器
	dis := consul.New(consulClient)
	if err != nil {
		log.Fatal(err)
	}
	selector.SetGlobalSelector(random.NewBuilder())
	//selector.SetGlobalSelector(wrr.NewBuilder())
	//selector.SetGlobalSelector(p2c.NewBuilder())
	//连接目标grpc服务器
	endpoint := "discovery:///Map"
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(endpoint), //目标服务名字
		//grpc.WithEndpoint("localhost:9000"))  //verifyCode grpc service 地址
		//使用服务发现
		grpc.WithDiscovery(dis),
		//中间件
		grpc.WithMiddleware(
			//tracing中间件
			tracing.Client(),
		),
	)
	if err != nil {
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	//2.2发送获取驾驶距离和时长的请求，rpc调用
	client := mapService.NewMapServiceClient(conn)
	reply, err := client.GetDrivingInfo(context.Background(),
		&mapService.GetDrivingInfoReq{
			Origin:       origin,
			Destionation: destination,
		})
	if err != nil {
		return
	}
	distance, duration = reply.Distance, reply.Duration
	//返回正确信息
	return
}
