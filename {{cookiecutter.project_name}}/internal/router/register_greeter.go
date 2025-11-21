package router

import (
	v1 "{{cookiecutter.project_name}}/api/v1/helloworld"
	"{{cookiecutter.project_name}}/internal/service"

	"github.com/go-kratos/kratos/v2/transport/http"
)

func RegisterGreeterRouter(srv *http.Server, greeter *service.GreeterService) *http.Server {

	// 注册业务路由
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv

}
