  # Kratos Layout Template

Kratos Layout Template 是一个通过 Cookiecutter 生成的 Kratos 微服务模板，内置 HTTP/gRPC 双
  协议服务、Nacos 配置中心与服务注册、Knife4g API 文档、Zap + Lumberjack 日志、可扩展的中间件及指标
  组件。本 README 覆盖了模板使用方式、项目目录说明、配置方法与常见开发流程。

  > 如果你仍在模板仓库中，请先执行 `cookiecutter template/` 生成真正的服务项目，再按下文说明进行
  开发。

  ## 目录

  1. [特性概览](#特性概览)
  2. [环境要求](#环境要求)
  3. [模板使用方式](#模板使用方式)
  4. [项目目录说明](#项目目录说明)
  5. [配置说明](#配置说明)
  6. [常用命令（Makefile）](#常用命令makefile)
  7. [API、Proto 与文档](#apiproto-与文档)
  8. [运行、调试与测试](#运行调试与测试)
  9. [Docker 部署](#docker-部署)
  10. [扩展指南](#扩展指南)
  11. [故障排查](#故障排查)

  ## 特性概览

  - ✅ Cookiecutter 模板化生成，支持自定义 `project_name` 与 `nacos_addr`
  - ✅ Kratos 推荐分层（cmd/internal/pkg）+ Wire 依赖注入
  - ✅ HTTP/gRPC 双协议服务，内置 CORS 中间件与统一路由注册
  - ✅ 可选的 Nacos 配置中心与服务注册，实现动态配置热加载
  - ✅ Zap 日志 + Lumberjack 切割，默认输出到 stdout 与 `./logs`
  - ✅ Makefile 驱动的代码生成：Proto、HTTP/GRPC 桩、OpenAPI、配置 PB
  - ✅ Knife4g (Knife4j) API 文档，监听 `docs/api/openapi.yaml` 自动服务
  - ✅ Prometheus 指标封装、统一 Response 结构、可扩展中间件占位
  - ✅ 多环境配置 (`APP_ENV`) & Docker 部署示例

  ## 环境要求

| 组件 | 说明 |
| --- | --- |
| Go | 1.21+（模板当前 `go.mod` 写为 1.24.10，建议使用最新版 Go 发行版） |
| Python | 3.8+，用于运行 Cookiecutter |
| Cookiecutter | `pip install cookiecutter` or `brew install cookiecutter` |
| Protocol Buffers | `protoc` 3.21+ |
| Kratos CLI | `go install github.com/go-kratos/kratos/cmd/kratos/v2@latest` |
| Proto 插件 | `protoc-gen-go`、`protoc-gen-go-grpc`、`protoc-gen-go-http`、`protoc-gen-openapi`|
  |
  | Wire | `go install github.com/google/wire/cmd/wire@latest` |
  | 可选工具 | Docker、Nacos 服务、Prometheus/Grafana 等 |

  ### 安装示例

  ```bash
  # Cookiecutter
  pip install --upgrade cookiecutter

  # Go 生态工具
  go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
  go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
  go install github.com/google/wire/cmd/wire@latest

  ## 模板使用方式

  1. 获取模板代码并进入仓库根目录。
  2. 执行 Cookiecutter（本地模板用 template/ 目录）：

     cookiecutter template/
     输入提示项：
      - project_name：服务名（用于模块名、二进制名、Nacos serviceName 等）
      - nacos_addr：默认 Nacos 地址（如 127.0.0.1）
  3. 生成完成后进入新项目目录：

     cd {{cookiecutter.project_name}}
  4. 首次执行 make init 安装所有 proto/wire 工具。
  5. hooks/post_gen_project.py 会尝试运行：
      - make api
      - make config
      - wire ./cmd/{{cookiecutter.project_name}}
        若缺少依赖可手动重新执行。

  ## 项目目录说明

  | 路径 | 说明 |
  | --- | --- |
  | api/ | Proto 描述（示例 v1/helloworld） |
  | assets/knife4g/ | 内嵌 Knife4g 前端静态资源 |
  | cmd/{{cookiecutter.project_name}}/ | 入口、配置加载、Wire 组装 |
  | configs/ | app.yaml、多环境配置与 conf.proto |
  | docs/api/ | make api/make apiDoc 生成的 OpenAPI 文档 |
  | internal/biz | 业务用例层（Usecase + 接口定义） |
  | internal/data | 数据访问层（Repo + 数据源初始化） |
  | internal/service | gRPC/HTTP Handler，持有 biz 用例 |
  | internal/router | HTTP 路由、Knife4g 文档注册 |
  | internal/server | HTTP/GRPC Server 构造与中间件 |
  | logs/ | 默认日志目录（运行时生成） |
  | pkg/log | Zap 日志封装 |
  | pkg/nacos | Nacos config/registry 实现 |
  | pkg/profile | 配置文件选择（APP_ENV） |
  | pkg/middleware | 自定义中间件（如 CORS） |
  | pkg/knife4g | Knife4g 文档 server |
  | pkg/metric | Prometheus 指标封装 |
  | third_party/ | Proto 依赖（errors/openapi/validate 等） |
  | Makefile | 常用命令集合 |
  | Dockerfile | 多阶段容器构建示例 |
  | main.go | 启动 cmd/{{cookiecutter.project_name}}.NewApp() |

  ## 配置说明

  ### configs/

  - configs/app.yaml：默认配置，建议复制为 app-dev.yaml、app-prod.yaml 等。
  - pkg/profile 根据 APP_ENV 环境变量，依次读取：
      1. configs/app.yaml
      2. configs/app-<env>.yaml（若 APP_ENV 设置）
  - 配置结构定义在 configs/conf/conf.proto，变更后运行 make config 生成 conf.pb.go。

  示例 app.yaml 片段：

  global:
    appName: {{cookiecutter.project_name}}
    env: test
    version: v1
    id: 127.0.0.1

  nacos:
    enable: true
    ip: {{cookiecutter.nacos_addr}}
    port: 8848
    discovery:
      serviceName: {{cookiecutter.project_name}}
      groupName: DEFAULT_GROUP
      clusterName: {{cookiecutter.project_name}}

  log:
    level: info
    format: json
    filename: ./logs
    maxSize: 10
    maxBackups: 5
    maxAge: 30
    compress: true

  server:
    http:
      network: tcp
      addr: 0.0.0.0:8000
      timeout: 2s
      enableDoc: true
    grpc:
      network: tcp
      addr: 0.0.0.0:9000
      timeout: 2s
    httpCors:
      mode: allow-all

  ### Nacos 集成

  - cmd/{{cookiecutter.project_name}}/cmd.go 会在 nacos.enable=true 时：
      1. 使用 pkg/nacos 创建自定义客户端 (config + naming)
      2. 挂载为 Kratos config.Source，自动热加载配置
      3. 将 Nacos 注册器注入 kratos.App，提供服务注册/发现
  - DataID: {{cookiecutter.project_name}}.yaml 或 {{cookiecutter.project_name}}-<env>.yaml（根
    据 global.env）
  - 缓存与日志目录：./configs/nacos/cache、./logs/nacos/log
  - 启动本地 Nacos（示例）：

    docker run -d --name nacos -p 8848:8848 nacos/nacos-server:v2.3.2
  - 将相同内容的配置写入 Nacos，模板会自动覆盖本地 app.yaml 内容。

  ### 日志与观测

  - pkg/log 使用 Zap + Lumberjack，日志输出字段附带 env/service_id/service_name/trace_id。
  - global、log 节点控制日志级别与文件大小。
  - pkg/metric 提供 counter/gauge/histogram 封装，可在 Handler 中调用记录 HTTP/GRPC 指标。
  - Prometheus 集成可通过 Kratos 官方中间件或自定义 /metrics 路由扩展。

  ### 服务与中间件

  - HTTP/GRPC Server 均位于 internal/server，默认挂载 Recovery + Logging。
  - pkg/middleware/cors 读取 server.httpCors 配置，支持 allow-all 或白名单模式。
  - internal/router/router.go 负责注册业务路由及 Knife4g 静态服务；需要新增路由时在此扩展。

  ### 数据层

  - internal/data 中的 Data 结构是数据库/缓存等客户端的封装，使用 Wire 注入。
  - 默认 internal/data/greeter.go 提供示例 repo，实现 biz.GreeterRepo 接口。
  - 可在 data 包内初始化 sql.DB、gorm.DB、redis.Client 等，并通过 Data 传入 Usecase。

  ## 常用命令（Makefile）

  | 命令 | 作用 |
  | --- | --- |
  | make help | 查看全部命令说明 |
  | make init | 安装 proto/kratos/wire 等依赖 |
  | make api | 根据 api/**.proto 生成 pb、http、grpc、openapi 文件 |
  | make httpApi / make grpcApi | 仅生成 HTTP 或 gRPC 相关代码 |
  | make apiDoc | 只生成 OpenAPI 文档到 docs/api/ |
  | make config | 根据 configs/conf/conf.proto 生成配置结构 |
  | make generate | 执行 go generate ./... 与 go mod tidy |
  | make build | 编译二进制至 ./bin/{{cookiecutter.project_name}} |
  | make all | 依次执行 make api && make config && make generate |

  ## API、Proto 与文档

  1. 新增/修改 api/**.proto：
      - 定义 HTTP rule、验证、错误码（见 api/v1/helloworld 示例）
      - 执行 make api 生成 pb.go、pb.http.go、pb.grpc.go、docs/api/openapi.yaml
  2. third_party/ 目录已包含常用 proto 依赖，可按需添加。
  3. Knife4g 文档：
      - pkg/knife4g 自动加载 docs/api/openapi.yaml
      - 启动服务后访问 http://localhost:8000/doc.html（默认 HTTP 端口 8000）即可查看文档
      - 若需要自定义访问前缀，可修改 pkg/knife4g/config.go 或在 RegisterKnife4g 中设定

  ## 运行、调试与测试

  # 本地运行（默认读取 configs/app.yaml）
  go run ./cmd/{{cookiecutter.project_name}}

  # 指定环境（先在 configs/ 下创建 app-prod.yaml 并配置 APP_ENV）
  APP_ENV=prod go run ./cmd/{{cookiecutter.project_name}}

  # 构建 & 运行
  make build
  ./bin/{{cookiecutter.project_name}}

  # 单元测试
  go test ./...

  - 当 nacos.enable=true 且 Nacos 可用时，服务启动会自动向 Nacos 注册（HTTP/GRPC 分别注册为
    {{cookiecutter.project_name}}.http、.grpc）。
  - 默认 HTTP 监听 0.0.0.0:8000，GRPC 监听 0.0.0.0:9000，可在 server.http/addr、server.grpc/addr 中
    修改。
  - 日志输出目录为 ./logs，Nacos 缓存位于 ./configs/nacos/cache。

  ## Docker 部署

  模板自带多阶段 Dockerfile：

  # 构建镜像
  docker build -t {{cookiecutter.project_name}}:latest .

  # 运行容器，并挂载配置
  docker run --rm \
    -p 8000:8000 -p 9000:9000 \
    -e APP_ENV=prod \
    -v $(pwd)/configs:/data/conf \
    {{cookiecutter.project_name}}:latest

  > 注意：容器默认命令为 ./server -conf /configs/app.yaml。若需改为多环境文件，可在 Dockerfile 或启
  > 动命令中调整，并确保容器内存在对应配置。

  ## 扩展指南

  1. 新增业务模块
      - 定义 api proto → make api
      - 在 internal/biz 增加 Usecase 接口与实现
      - 在 internal/data 实现 repo
      - 在 internal/service 编写 Handler 并通过 service.Holder 暴露
      - 更新 internal/router 注册 HTTP/GRPC
      - 运行 wire ./cmd/{{cookiecutter.project_name}} 重新生成依赖注入代码
  2. 接入数据库/缓存
      - 修改 configs/app.yaml 中的 data.database/data.redis
      - 在 internal/data/data.go 初始化对应客户端（sql.DB、gorm.DB、redis.Client 等）
      - 在 repo 中复用 Data 获取到的连接
  3. 新增中间件或指标
      - 在 pkg/middleware 新增实现（可复用 Kratos middleware 包）
      - 在 internal/server/http.go/grpc.go 中注入
      - 使用 pkg/metric 提供的 ReqCount、RespDurationHistogram 等指标
      - 部署时设置 APP_ENV=<env>，模板会自动先加载 app.yaml 再 overlay app-<env>.yaml
      - 若需要支持命令行参数，可在 pkg/profile 中扩展 Flag 逻辑

  ## 故障排查

  - wire 生成失败：确认已执行 go install github.com/google/wire/cmd/wire@latest 并在 cmd/
    {{cookiecutter.project_name}} 目录下运行 wire。
  - Nacos 连接失败：检查 nacos.ip、nacos.port、config.namespace、discovery.groupName 等配置；确保
    DataID 命名符合 appName[-env].yaml。
  - Knife4g 文档无法访问：确认 make api 已生成 docs/api/openapi.yaml，并检查
    server.http.enableDoc（默认 true）。
  - CORS 403：若启用白名单模式，请在 server.httpCors.whitelist 中添加正确的 allowOrigin、
    allowHeaders、allowMethods。
  - proto 代码未更新：清理 api/**.pb.* 后重新执行 make api，或直接运行 make all。
  - Docker 中读取配置失败：挂载目录时请确保容器内路径（/configs）与 pkg/profile 预期一致，或在镜像
    中调整 configs 目录。
