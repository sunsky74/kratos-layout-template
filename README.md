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
     cookiecutter ./template --no-input project_name="example_project"
      - project_name：服务名（用于模块名、二进制名、Nacos serviceName 等）


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
