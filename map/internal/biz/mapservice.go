package biz

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"net/http"
)

type MapServiceBiz struct {
	log *log.Helper
}

// 参考greeter代码
func NewMapServiceBiz(logger log.Logger) *MapServiceBiz {
	return &MapServiceBiz{log: log.NewHelper(logger)}
}

func (msbiz *MapServiceBiz) GetDriverInfo(origin, destination string) (string, string, error) {
	//一、请求获取
	key := "30d14843e60c30dfc447e85c367a6ef9"
	api := "https://restapi.amap.com/v3/direction/driving"
	parameter := fmt.Sprintf("origin=%s&destination=%s&extensions=base&output=JSON&key=%s", origin, destination, key)
	url := api + "?" + parameter
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body) // io.Reader
	if err != nil {
		return "", "", err
	}
	//二、解析
	ddResp := DirectionDrivingResp{}
	if err := json.Unmarshal(body, &ddResp); err != nil {
		return "", "", err
	}
	//三、判断lsb请求结果
	if ddResp.Status == "0" {
		return "", "", errors.New(1, "", ddResp.Info)
	}
	//正确返回
	path := ddResp.Route.Paths[0]
	return path.Distance, path.Duration, nil
}

type DirectionDrivingResp struct {
	Status   string `json:"status,omitempty"`
	Info     string `json:"info,omitempty"`
	Infocode string `json:"infocode,omitempty"`
	Count    string `json:"count,omitempty"`
	Route    struct {
		Origin      string `json:"origin,omitempty"`
		Destination string `json:"destination,omitempty"`
		Paths       []Path `json:"paths,omitempty"`
	} `json:"route"`
}

type Path struct {
	Distance string `json:"distance,omitempty"`
	Duration string `json:"duration,omitempty"`
	Strategy string `json:"strategy,omitempty"`
}
