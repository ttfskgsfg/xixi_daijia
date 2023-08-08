package data

import (
	"customer/internal/biz"
	"customer/internal/conf"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewCustomerData)

// Data .
type Data struct {
	// TODO wrapped database client
	Rdb *redis.Client
	//操作mysql客户端
	Mdb *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{}
	//连接redis 使用服务配置 c就是解析后的配置信息
	redisURL := fmt.Sprintf("redis://%s/1?dial_timeout=%d", c.Redis.Addr, 1)
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		data.Rdb = nil
		log.Fatal(err)
	}
	data.Rdb = redis.NewClient(options)

	//二、初始mdb
	//连接mysql使用配置
	dsn := c.Database.Source
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		data.Mdb = nil
		log.Fatal(err)
	}
	data.Mdb = db
	//三、开发阶段 自动迁移表  发布阶段表结构稳定
	migrateTable(db)

	cleanup := func() {
		//清理redis连接
		_ = data.Rdb.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return data, cleanup, nil
}

func migrateTable(db *gorm.DB) {
	if err := db.AutoMigrate(&biz.Customer{}); err != nil {
		log.Info("customer table migrate error:", err)
	}
}
