package router

import (
	"{{cookiecutter.project_name}}/configs/conf"
	"{{cookiecutter.project_name}}/pkg/knife4g"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func RegisterKnife4g(cfg *conf.Config, srv *http.Server, logger log.Logger) *http.Server {
	//if cfg.Server.Http.EnableDoc {
	registerKnife4gDoc(cfg, srv, logger)
	//}
	return srv
}

// registerKnife4gDoc 注册 knife4g 文档服务
func registerKnife4gDoc(c *conf.Config, srv *http.Server, logger log.Logger) {
	// 加载 OpenAPI 配置（默认使用 docs/api/openapi.yaml）
	config, err := knife4g.NewDefaultConfig(c.Global.AppName)
	if err != nil {
		log.NewHelper(logger).Warnf("Failed to load OpenAPI config: %v", err)
		return
	}
	docPath := "/doc.html"
	if config.RelativePath != "" {
		docPath = config.RelativePath + "/doc.html"
	}
	log.NewHelper(logger).Infof("Knife4g documentation available at: %s", docPath)

	handler := knife4g.Handler(config)

	// 使用 HandlePrefix 注册根路径，让 Knife4g handler 处理所有请求
	srv.HandlePrefix("/", handler)
}
