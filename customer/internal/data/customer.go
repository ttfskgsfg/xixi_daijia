package data

import (
	"context"
	"customer/internal/biz"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"log"
	"time"
)

// customer中与数据操作相关的代码
type CustomerData struct {
	data *Data
}

// new方法
func NewCustomerData(data *Data) *CustomerData {
	return &CustomerData{data: data}
}

// 设置验证码方法
func (cd CustomerData) SetVerifyCode(telephone, code string, ex int) error {
	//设置key
	status := cd.data.Rdb.Set(context.Background(), "CVC:"+telephone, code, time.Duration(ex)*time.Second)
	if _, err := status.Result(); err != nil {
		return err
	}
	return nil
}

// 获取号码对应的验证码
func (cd CustomerData) GetVerifyCode(telephone string) string {
	status := cd.data.Rdb.Get(context.Background(), "CVC:"+telephone)
	return status.String()
}

// 根据电话 获取顾客信息
func (cd CustomerData) GetCustomerByTelephone(telephone string) (*biz.Customer, error) {
	//查询给予电话
	customer := &biz.Customer{}
	result := cd.data.Mdb.Where("telephone=?", telephone).First(customer)
	//query执行成功，同时查询到记录 直接返回顾客信息
	if result.Error == nil && customer.ID > 0 {
		return customer, nil
	}
	//发生记录不存在的错误
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//创建customer并返回
		customer.Telephone = telephone
		customer.Name = sql.NullString{Valid: false}
		customer.Email = sql.NullString{Valid: false}
		customer.Wechat = sql.NullString{Valid: false}
		reusltCreate := cd.data.Mdb.Create(customer)
		//插入成功
		if reusltCreate.Error != nil {
			return customer, nil
		} else {
			return nil, reusltCreate.Error
		}
	}
	//有错误 但也不是记录不存在的错误，不做业务逻辑处理
	return nil, result.Error
}

func (cd CustomerData) GenerateTokenAndSave(c *biz.Customer, duration time.Duration, secret string) (string, error) {
	//一、生成token
	//处理token生成的数据
	//标准的jwt的payload
	claims := jwt.RegisteredClaims{
		//签发方
		Issuer: "XiXiDJ",
		//说明
		Subject: "customer-authentication",
		//签发给谁使用
		Audience: []string{"customer", "other"},
		//有效期至
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		//何时启用
		NotBefore: nil,
		//签发时间
		IssuedAt: jwt.NewNumericDate(time.Now()),
		//id 用户的id
		ID: fmt.Sprintf("%d", c.ID),
	}
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//签名token
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	//二、存储
	c.Token = signedToken
	c.TokenCreatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if result := cd.data.Mdb.Save(c); result.Error != nil {
		return "", result.Error
	}
	log.Println(signedToken)
	//操作成功
	return signedToken, nil
}

// 利用顾客id获取token
func (cd CustomerData) GetToken(id interface{}) (string, error) {
	c := &biz.Customer{}
	if result := cd.data.Mdb.First(c, id); result.Error != nil {
		return "", result.Error
	}
	return c.Token, nil
}

// 利用顾客id获取token
func (cd CustomerData) DelToken(id interface{}) error {
	c := &biz.Customer{}
	//找到customer
	if result := cd.data.Mdb.First(c, id); result.Error != nil {
		return result.Error
	}
	//删除customer的token
	c.Token = ""
	c.TokenCreatedAt = sql.NullTime{Valid: false}
	cd.data.Mdb.Save(c)

	return nil
}
