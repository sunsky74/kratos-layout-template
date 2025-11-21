package cmd

import (
	"fmt"
	"{{cookiecutter.project_name}}/configs/conf"
	pkg "{{cookiecutter.project_name}}/pkg/log"
	"{{cookiecutter.project_name}}/pkg/nacos"
	"{{cookiecutter.project_name}}/pkg/profile"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, cc *conf.Config, nac *nacos.Client) *kratos.App {
	options := []kratos.Option{
		kratos.ID(cc.GetGlobal().Id),
		kratos.Name(cc.GetGlobal().AppName),
		kratos.Version(cc.GetGlobal().Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	}
	if cc.Nacos.Enable && nac != nil {
		options = append(options, kratos.Registrar(nac))
	}
	return kratos.New(options...)

}

func NewApp() (*kratos.App, func()) {

	cc, err2 := loadLocalConfig()
	if err2 != nil {
		panic(err2)
	}

	logger, err := pkg.New(cc.Log, cc.Global)
	if err != nil {
		panic(err)
	}

	if cc.Nacos != nil && cc.Nacos.Enable {
		cs, e := newNacosConfigSource(cc, logger)
		if e != nil {
			panic(e)
		}

		c := config.New(
			config.WithSource(cs),
		)
		if err := c.Load(); err != nil {
			_ = c.Close()
			log.NewHelper(logger).Error("nacos config load failed", "error", err)
		} else {
			if err := c.Scan(cc); err != nil {
				_ = c.Close()
				log.NewHelper(logger).Error("nacos config load failed", "error", err)
			}
		}

		app, f, err := wireApp(cc, cs.Client, logger)
		if err != nil {
			_ = c.Close()
			panic(err)
		}
		return app, func() {
			_ = c.Close()
			defer f()
		}
	}
	app, f, err := wireApp(cc, nil, logger)
	if err != nil {
		panic(err)
	}
	return app, f

}

func loadLocalConfig() (*conf.Config, error) {
	pro := profile.LoadProfile()
	cfgs := make([]config.Source, 0, len(pro.FilePaths))
	for _, path := range pro.FilePaths {
		cfgs = append(cfgs, file.NewSource(path))
	}
	c := config.New(
		config.WithSource(cfgs...),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		return nil, err
	}

	var cfg conf.Config
	if err := c.Scan(&cfg); err != nil {
		return nil, err
	}
	if pro.ENV != "" && pro.ENV != cfg.Global.Env {
		cfg.Global.Env = pro.ENV
	}

	return &cfg, nil
}

func newNacosConfigSource(cc *conf.Config, logger log.Logger) (*nacos.ConfigSource, error) {
	var dataId string
	if cc.Global.Env == "" {
		dataId = fmt.Sprintf("%s.yaml", cc.Global.AppName)
	} else {
		dataId = fmt.Sprintf("%s-%s.yaml", cc.Global.AppName, cc.Global.Env)
	}

	client, err := nacos.NewClient(
		nacos.WithHost(fmt.Sprintf("%s:%d", cc.Nacos.GetIp(), cc.Nacos.GetPort())),
		nacos.WithCacheDir("./configs/nacos/cache"),
		nacos.WithLogDir("./logs/nacos/log"),
		nacos.WithConfigGroup(cc.Nacos.Discovery.GetGroupName()),
		nacos.WithNamespaceId(cc.Nacos.Config.GetNamespace()),
		nacos.WithUsername(cc.Nacos.Config.GetUsername()),
		nacos.WithPassword(cc.Nacos.Config.GetPassword()),
		nacos.WithLogger(logger),
		nacos.WithLogLevel("info"),
		nacos.WithConfigDataID(dataId),
	)
	if err != nil {
		return nil, err
	}

	return nacos.NewConfigSource(client), nil
}
