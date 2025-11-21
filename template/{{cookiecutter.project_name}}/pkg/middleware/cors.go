package middleware

import (
	"context"
	"{{cookiecutter.project_name}}/configs/conf"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func Cors(cors *conf.Server_Cors) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 获取 HTTP transport 信息
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tr.(*kratoshttp.Transport); ok {
					httpReq := ht.Request()
					httpResp := ht.Response()

					origin := httpReq.Header.Get("Origin")

					if cors.Mode == "allow-all" || origin == "" || origin == "null" {
						// 允许所有来源的配置
						httpResp.Header().Set("Access-Control-Allow-Origin", "*")
						httpResp.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
						httpResp.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
						httpResp.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
						httpResp.Header().Set("Access-Control-Allow-Credentials", "true")

						// 处理预检请求
						if httpReq.Method == http.MethodOptions {
							httpResp.WriteHeader(http.StatusNoContent)
							return nil, nil
						}
					} else {
						// 白名单模式
						whitelist := checkOrigin(origin, cors.Whitelist)
						if whitelist != nil {
							httpResp.Header().Set("Access-Control-Allow-Origin", whitelist.AllowOrigin)
							httpResp.Header().Set("Access-Control-Allow-Headers", whitelist.AllowHeaders)
							httpResp.Header().Set("Access-Control-Allow-Methods", whitelist.AllowMethods)
							httpResp.Header().Set("Access-Control-Expose-Headers", whitelist.ExposeHeaders)
							if whitelist.AllowCredentials {
								httpResp.Header().Set("Access-Control-Allow-Credentials", "true")
							}

							// 处理预检请求
							if httpReq.Method == http.MethodOptions {
								httpResp.WriteHeader(http.StatusNoContent)
								return nil, nil
							}
						} else {
							// 不在白名单中，返回 403
							httpResp.WriteHeader(http.StatusForbidden)
							return nil, nil
						}
					}
				}
			}

			// 继续执行后续处理器
			return handler(ctx, req)
		}
	}
}

func checkOrigin(origin string, whitelist []*conf.Server_Cors_Whitelist) *conf.Server_Cors_Whitelist {
	for _, v := range whitelist {
		if origin == v.AllowOrigin {
			return v
		}
	}
	return nil
}
