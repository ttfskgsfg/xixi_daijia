package biz

import "context"

type DriverBiz struct {
	di DriverInterface
}

type DriverInterface interface {
	GetVerifycode(ctx context.Context, string) (string, error)
}

func NewDriverBiz(di DriverInterface) *DriverBiz {
	return &DriverBiz{
		di: di,
	}
}
