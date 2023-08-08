package server

import (
	"context"
	"customer/internal/service"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"strings"
)

// 用来生成中间件
func customerJWT(customerService *service.CustomerService) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			//一、获取存储在jwt中的用户(顾客)id
			claims, ok := jwt.FromContext(ctx)
			if !ok {
				//没有获取到claims
				return nil, errors.Unauthorized("UNAUTHORIZED",
					"claims not found")
			}
			//1.2 断言使用
			claimsMap := claims.(jwtv4.MapClaims)
			id := claimsMap["jti"]
			//二、获取id对应顾客token
			token, err := customerService.CD.GetToken(id)
			if err != nil {
				return nil, errors.Unauthorized("UNAUTHORIZED",
					"customer not found")
			}
			//三、比对数据表中token与请求token是否一致
			//获取请求头
			header, _ := transport.FromServerContext(ctx)
			//从header获取token
			auths := strings.SplitN(header.RequestHeader().Get("Authorization"), " ", 2)
			jwtToken := auths[1]
			if jwtToken != token {
				return nil, errors.Unauthorized("UNAUTHORIZED",
					"token was updated")
			}

			//四、校验通过放行
			//交由下个中间件handler处理
			return handler(ctx, req)
		}
	}
}
