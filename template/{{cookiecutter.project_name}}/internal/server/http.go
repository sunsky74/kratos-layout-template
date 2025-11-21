package server

import (
	"{{cookiecutter.project_name}}/configs/conf"
	r "{{cookiecutter.project_name}}/internal/router"
	"{{cookiecutter.project_name}}/internal/service"
	"{{cookiecutter.project_name}}/pkg/middleware"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Config, sh *service.Holder, log log.Logger) *http.Server {
	srv := initServer(c, log)
	r.Route(c, srv, log, sh)
	return srv
}

func initServer(c *conf.Config, log log.Logger) *http.Server {

	s := c.GetServer()
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(log),
			middleware.Cors(c.Server.HttpCors),
		),
	}
	if s.Http.Network != "" {
		opts = append(opts, http.Network(s.Http.Network))
	}
	if s.Http.Addr != "" {
		opts = append(opts, http.Address(s.Http.Addr))
	}
	if s.Http.Timeout != "" {
		duration, _ := time.ParseDuration(s.Http.Timeout)
		opts = append(opts, http.Timeout(duration))
	}

	return http.NewServer(opts...)

}
