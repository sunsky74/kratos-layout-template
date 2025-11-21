//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package cmd

import (
	"{{cookiecutter.project_name}}/configs/conf"
	"{{cookiecutter.project_name}}/internal/biz"
	"{{cookiecutter.project_name}}/internal/data"
	"{{cookiecutter.project_name}}/internal/server"
	"{{cookiecutter.project_name}}/internal/service"
	"{{cookiecutter.project_name}}/pkg/nacos"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Config, *nacos.Client, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
