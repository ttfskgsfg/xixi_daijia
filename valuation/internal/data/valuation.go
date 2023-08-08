package data

import (
	"fmt"
	"valuation/internal/biz"
)

type PriceRuleData struct {
	data *Data
}

func NewPriceRuleInterface(data *Data) biz.PriceRuleInterface {
	return &PriceRuleData{data: data}
}

// data里PriceRuleData  实现 biz里PriceRuleInterface接口
// biz中定义接口 data里实现
func (prd *PriceRuleData) GetRule(cityid uint, curr int) (*biz.PriceRule, error) {
	pr := &biz.PriceRule{}
	result := prd.data.Mdb.Where("city_id=? AND start_at <= ? AND end_at > ?", cityid, curr,
		curr).First(pr)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println(prd, prd.data)
	return pr, nil
}
