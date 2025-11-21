package server

import (
	v1 "{{cookiecutter.project_name}}/api/v1/helloworld"
	"{{cookiecutter.project_name}}/configs/conf"
	"{{cookiecutter.project_name}}/internal/service"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Config, greeter *service.GreeterService, logger log.Logger) *grpc.Server {
	s := c.GetServer()
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	}
	if s.Grpc.Network != "" {
		opts = append(opts, grpc.Network(s.Grpc.Network))
	}
	if s.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(s.Grpc.Addr))
	}
	if s.Grpc.Timeout != "" {
		duration, _ := time.ParseDuration(s.Http.Timeout)
		opts = append(opts, grpc.Timeout(duration))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServer(srv, greeter)
	return srv
}
