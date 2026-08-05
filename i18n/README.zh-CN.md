**Languages:** [English](../README.md) | [日本語](README.ja.md) | [فارسی](README.fa.md) | [Deutsch](README.de.md) | [简体中文](README.zh-CN.md) | [Русский](README.ru.md) | [Español](README.es.md) | [Português](README.pt-BR.md) | [Français](README.fr.md) | [한국어](README.ko.md)

# DCB — Docker Compose Builder

> **少写 YAML。运行更好的容器。**

DCB 读取简化的 `dcb.yaml`，并生成适用于生产环境的 `docker-compose.yml`，
自动包含健康检查、命名网络、重启策略和依赖排序。

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)

---

## 为什么选择 DCB？

`docker-compose.yml` 功能强大但较为冗长。典型的生产环境技术栈会重复相同的样板代码，用于健康检查、重启策略、网络声明和卷命名。

使用 DCB，你可以这样写：

```yaml
# dcb.yaml
project: webapp

services:
  api:
    build: ./api
    port: 8080
    depends_on: [db]
    healthcheck: true
    env_file: .env

  db:
    image: postgres:15-alpine
    port: 5432
    env:
      POSTGRES_DB: webapp
    volume: pgdata:/var/lib/postgresql/data
    healthcheck: true
```

即可获得完整配置的 `docker-compose.yml`，包含：

- ✅ 针对 postgres 的正确 `pg_isready` 健康检查（感知镜像）
- ✅ 带有 `condition: service_healthy` 的 `depends_on`
- ✅ 命名网络 `webapp_default`
- ✅ 顶层 `volumes:` 声明
- ✅ 每个服务默认 `restart: unless-stopped`
- ✅ 解析时检测循环依赖

---

## 安装

### Homebrew (macOS / Linux)
```bash
brew install JrNovaEX/tap/dcb
```

### Go install
```bash
go install github.com/JrNovaEX/DCB/cmd/dcb@latest
```

### 下载二进制文件

从 [Releases 页面](https://github.com/JrNovaEX/DCB/releases) 获取适合你平台的最新版本。

```bash
# Linux amd64 示例
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### 从源码构建
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## 快速开始

```bash
# 1. 交互式创建新的 dcb.yaml
dcb init

# 2. 验证配置并检查主机端口可用性
dcb check

# 3. 生成 docker-compose.yml
dcb build

# 4. 在后台启动所有服务
dcb up --detach

# 5. 跟踪所有服务的日志
dcb logs -f

# 6. 停止并清理
dcb down
```

---

## dcb.yaml 参考

```yaml
project: myapp          # 必填。用于网络和卷的命名空间。
version: "3.8"          # 可选。Docker Compose schema 版本（默认：3.8）。

services:
  <service-name>:
    # 来源 — 选择 image 或 build 其中之一：
    image: nginx:alpine

    build: ./myservice          # 简写形式：context 路径
    build:                      # 完整形式：
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # 以 host:container 方式暴露（两侧使用相同值）。
    env_file: .env       # env 文件路径。
    env:                 # 内联环境变量。
      KEY: value

    depends_on:          # 需要等待的服务名称（condition: service_healthy）。
      - db

    healthcheck: true    # 根据镜像名称自动生成健康检查。
                         # 支持：postgres、mysql、mariadb、redis、nginx、
                         #       mongo、rabbitmq。其他镜像使用 CMD-SHELL exit 0。

    volume: pgdata:/var/lib/postgresql/data   # 命名卷或绑定挂载。
    restart: unless-stopped                   # 默认值。可按需覆盖。
```

---

## 命令

| 命令 | 说明 |
|---------|-------------|
| `dcb init` | 交互式生成新的 `dcb.yaml` |
| `dcb check` | 验证 `dcb.yaml` + 检查主机端口可用性 |
| `dcb build` | 从 `dcb.yaml` 生成 `docker-compose.yml` |
| `dcb up` | 构建并启动服务（`docker compose up`） |
| `dcb down` | 停止并移除服务（`docker compose down`） |
| `dcb logs` | 流式输出服务日志 |
| `dcb version` | 打印版本、commit 和构建日期 |

### 全局标志
```
--verbose, -V    启用调试日志
--json           以 JSON 格式输出日志（适用于 CI/CD 流水线）
```

### dcb build
```
--file,     -f   输入文件（默认：dcb.yaml）
--output,   -o   输出文件（默认：docker-compose.yml）
--env,      -e   目标环境（默认：dev）
--validate       生成后运行 docker compose config（默认：true）
```

### dcb up
```
--file,     -f   输入 dcb.yaml（默认：dcb.yaml）
--env,      -e   目标环境
--detach,   -d   在后台运行
--build          强制重新构建镜像
```

### dcb down
```
--file,     -f   Compose 文件（默认：docker-compose.yml）
--project,  -p   覆盖项目名称
--volumes,  -v   移除命名卷
--remove-orphans 移除孤立容器
```

### dcb logs
```
--file           Compose 文件（默认：docker-compose.yml）
--follow,   -f   流式输出日志
--tail,     -n   显示的行数（-1 = 全部）
[services...]    按特定服务名称过滤
```

---

## 多环境支持

DCB 会将基础 `dcb.yaml` 与特定环境的覆盖文件合并：

```
dcb build --env prod
```

先加载 `dcb.yaml`，再将 `dcb.prod.yaml` 合并到其上。`dcb.prod.yaml` 中的键会覆盖基础配置；缺失的键则回退到基础配置。

---

## 架构

DCB 遵循 **Clean Architecture**（六边形架构）：

```
cmd/dcb/
  commands/         CLI 层 — Cobra 命令、标志解析、依赖注入
internal/
  domain/           核心实体 — ProjectConfig, Service, Volume, Port
  ports/            接口 — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         应用逻辑 — Build, Check, Up, Down, Logs
  adapters/
    parser/         YAML 解析器（gopkg.in/yaml.v3，严格模式）
    generator/      docker-compose.yml 模板生成器（text/template）
    runner/         Docker CLI 执行器（os/exec）
    config/         多环境配置加载器（Viper）
  infrastructure/
    logger/         slog 初始化
    validator/      go-playground/validator 封装
pkg/version/        构建时版本信息
```

依赖关系向内流动：**cmd → usecases → domain**。适配器依赖端口，彼此之间不相互依赖。

---

## 开发

```bash
# 运行所有测试
make test

# 带竞态检测和覆盖率的测试
make cover

# Lint
make lint

# Format
make fmt

# 构建本地二进制
make build
./dist/dcb version
```

### 运行测试

```bash
go test ./...                          # 所有包
go test ./internal/domain/...          # 仅 domain
go test ./internal/adapters/parser/... # 仅 parser
```

---

## 贡献

1. Fork 本仓库
2. 创建功能分支（`git checkout -b feat/my-feature`）
3. 使用 [Conventional Commits](https://www.conventionalcommits.org/) 提交
4. 发起 Pull Request

请在打开 PR 前确保 `make test lint` 通过。

---

## 许可证

MIT — 详见 [LICENSE](../LICENSE)。
