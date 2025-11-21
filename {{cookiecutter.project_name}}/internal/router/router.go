package router

import (
	"{{cookiecutter.project_name}}/configs/conf"
	"{{cookiecutter.project_name}}/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func Route(c *conf.Config, srv *http.Server, logger log.Logger, h *service.Holder) {
	RegisterKnife4g(c, srv, logger)
	RegisterGreeterRouter(srv, h.GreeterService)
}
