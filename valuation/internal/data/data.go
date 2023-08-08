package data

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"valuation/internal/biz"
	"valuation/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewPriceRuleInterface)

// Data .
type Data struct {
	// TODO wrapped database client
	Mdb *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{}
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
		log.NewHelper(logger).Info("closing the data resources")
	}
	return data, cleanup, nil
}

func migrateTable(db *gorm.DB) {
	if err := db.AutoMigrate(&biz.PriceRule{}); err != nil {
		log.Info("price_rule table migrate error:", err)
	}
	//插入一些priceRule的测试数据
	rules := []biz.PriceRule{
		{
			Model: gorm.Model{ID: 1},
			PriceRuleWork: biz.PriceRuleWork{
				CityID:      1,
				StartFee:    300,
				DistanceFee: 35,
				DurationFee: 10,
				StartAt:     7,
				EndAt:       23,
			},
		},
		{
			Model: gorm.Model{ID: 2},
			PriceRuleWork: biz.PriceRuleWork{
				CityID:      1,
				StartFee:    400,
				DistanceFee: 35,
				DurationFee: 10,
				StartAt:     23,
				EndAt:       24,
			},
		},
		{
			Model: gorm.Model{ID: 3},
			PriceRuleWork: biz.PriceRuleWork{
				CityID:      1,
				StartFee:    500,
				DistanceFee: 35,
				DurationFee: 10,
				StartAt:     0,
				EndAt:       7,
			},
		},
	}
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(rules)
}
