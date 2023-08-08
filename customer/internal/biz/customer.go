package biz

import (
	"context"
	"customer/api/valuation"
	"database/sql"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

// 模型
type Customer struct {
	//业务逻辑
	CustomerWork
	//token部分
	CustomerToken
	//嵌入4个字段
	gorm.Model
}

const CustomerSecret = "yourSecretKey"
const CustomerDuration = 2 * 30 * 24 * 3600

// 业务逻辑部分
type CustomerWork struct {
	Telephone string         `gorm:"type:varchar(15);uniqueIndex" json:"telephone"`
	Name      sql.NullString `gorm:"type:varchar(255);uniqueIndex" json:"name"`
	Email     sql.NullString `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Wechat    sql.NullString `gorm:"type:varchar(255);uniqueIndex" json:"wechat"`
	CityID    uint           `gorm:"index;" json:"city_id"`
}

// token部分
type CustomerToken struct {
	Token          string       `gorm:"type:varchar(4095);" json:"token"`
	TokenCreatedAt sql.NullTime `gorm:"" json:"tokenCreatedAt"`
}

type CustomerBiz struct {
}

func NewCustomerBiz() *CustomerBiz {
	return &CustomerBiz{}
}

func (cb *CustomerBiz) GetEstimatePrice(origin, destination string) (int64, error) {
	//grpc获取
	//使用服务发现
	//1、获取consul客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "localhost:8500"
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return 0, err
	}
	//2、获取consul发现管理器
	dis := consul.New(consulClient)
	if err != nil {
		return 0, err
	}
	endpoint := "discovery:///Valuation"
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(endpoint), //目标服务名字
		grpc.WithDiscovery(dis),
	)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = conn.Close()
	}()

	//2. 2发送获取验证码请求
	client := valuation.NewValuationClient(conn)
	reply, err := client.GetEstimatePrice(context.Background(),
		&valuation.GetEstimatePriceReq{
			Origin:      origin,
			Destination: destination,
		})
	if err != nil {
		return 0, err
	}
	return reply.Price, nil
}
