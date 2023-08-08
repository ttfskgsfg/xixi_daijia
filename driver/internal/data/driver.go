package data

import (
	"context"
	"driver/internal/biz"
)

type DriverData struct {
	data *Data
}

func NewDriverInterface(data *Data) biz.DriverInterface {
	return &DriverData{data: data}
}

func (dt *DriverData) GetVerifyCode(ctx context.Context, tel string) (string, error) {

}
