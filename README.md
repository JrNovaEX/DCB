**Languages:** [English](README.md) | [日本語](i18n/README.ja.md) | [فارسی](i18n/README.fa.md) | [Deutsch](i18n/README.de.md) | [简体中文](i18n/README.zh-CN.md) | [Русский](i18n/README.ru.md) | [Español](i18n/README.es.md) | [Português](i18n/README.pt-BR.md) | [Français](i18n/README.fr.md) | [한국어](i18n/README.ko.md)

# DCB — Docker Compose Builder

> **Write less YAML. Run better containers.**

DCB reads a simplified `dcb.yaml` and generates a production-ready `docker-compose.yml`
with automatic health checks, named networks, restart policies, and dependency ordering.

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## Why DCB?

`docker-compose.yml` is powerful but verbose. A typical production stack repeats the same
boilerplate for health checks, restart policies, network declarations, and volume naming.

DCB lets you write this:

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

And get a fully configured `docker-compose.yml` with:

- ✅ Correct `pg_isready` health check for postgres (image-aware)
- ✅ `depends_on` with `condition: service_healthy`
- ✅ Named network `webapp_default`
- ✅ Top-level `volumes:` declaration
- ✅ `restart: unless-stopped` on every service
- ✅ Circular dependency detection at parse time

---

## Installation

### Homebrew (macOS / Linux)
```bash
brew install JrNovaEX/tap/dcb
```

### Go install
```bash
go install github.com/JrNovaEX/DCB/cmd/dcb@latest
```

### Download binary

Grab the latest release for your platform from the
[Releases page](https://github.com/JrNovaEX/DCB/releases).

```bash
# Linux amd64 example
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### Build from source
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## Quick Start

```bash
# 1. Create a new dcb.yaml interactively
dcb init

# 2. Validate your config and check port availability
dcb check

# 3. Generate docker-compose.yml
dcb build

# 4. Start everything in the background
dcb up --detach

# 5. Follow logs from all services
dcb logs -f

# 6. Tear down
dcb down
```

---

## dcb.yaml Reference

```yaml
project: myapp          # Required. Used to namespace networks and volumes.
version: "3.8"          # Optional. Docker Compose schema version (default: 3.8).

services:
  <service-name>:
    # Source — pick ONE of image or build:
    image: nginx:alpine

    build: ./myservice          # Short form: path to context
    build:                      # Long form:
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # Exposed as host:container (same value both sides).
    env_file: .env       # Path to an env file.
    env:                 # Inline environment variables.
      KEY: value

    depends_on:          # Service name(s) to wait for (condition: service_healthy).
      - db

    healthcheck: true    # Auto-generates a health check based on the image name.
                         # Supported: postgres, mysql, mariadb, redis, nginx,
                         #            mongo, rabbitmq. Others get CMD-SHELL exit 0.

    volume: pgdata:/var/lib/postgresql/data   # Named or bind mount.
    restart: unless-stopped                   # Default. Override if needed.
```

---

## Commands

| Command | Description |
|---------|-------------|
| `dcb init` | Interactively scaffold a new `dcb.yaml` |
| `dcb check` | Validate `dcb.yaml` + check host port availability |
| `dcb build` | Generate `docker-compose.yml` from `dcb.yaml` |
| `dcb up` | Build then start services (`docker compose up`) |
| `dcb down` | Stop and remove services (`docker compose down`) |
| `dcb logs` | Stream service logs |
| `dcb version` | Print version, commit, and build date |

### Global flags
```
--verbose, -V    Enable debug logging
--json           Output logs in JSON (for CI/CD pipelines)
```

### dcb build
```
--file,     -f   Input file (default: dcb.yaml)
--output,   -o   Output file (default: docker-compose.yml)
--env,      -e   Target environment (default: dev)
--validate       Run docker compose config after generation (default: true)
```

### dcb up
```
--file,     -f   Input dcb.yaml (default: dcb.yaml)
--env,      -e   Target environment
--detach,   -d   Run in the background
--build          Force rebuild of images
```

### dcb down
```
--file,     -f   Compose file (default: docker-compose.yml)
--project,  -p   Project name override
--volumes,  -v   Remove named volumes
--remove-orphans Remove orphan containers
```

### dcb logs
```
--file           Compose file (default: docker-compose.yml)
--follow,   -f   Stream log output
--tail,     -n   Number of lines to show (-1 = all)
[services...]    Filter to specific service names
```

---

## Multi-environment Support

DCB merges a base `dcb.yaml` with an environment-specific overlay:

```
dcb build --env prod
```

Loads `dcb.yaml`, then merges `dcb.prod.yaml` on top. Keys in `dcb.prod.yaml`
override the base; missing keys fall back to the base.

---

## Architecture

DCB follows **Clean Architecture** (Hexagonal):

```
cmd/dcb/
  commands/         CLI layer — Cobra commands, flag parsing, DI wiring
internal/
  domain/           Core entities — ProjectConfig, Service, Volume, Port
  ports/            Interfaces — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         Application logic — Build, Check, Up, Down, Logs
  adapters/
    parser/         YAML parser (gopkg.in/yaml.v3, strict mode)
    generator/      docker-compose.yml templater (text/template)
    runner/         Docker CLI executor (os/exec)
    config/         Multi-env config loader (Viper)
  infrastructure/
    logger/         slog initialiser
    validator/      go-playground/validator wrapper
pkg/version/        Build-time version info
```

Dependencies flow inward: **cmd → usecases → domain**. Adapters depend on ports,
not on each other.

---

## Development

```bash
# Run all tests
make test

# Run tests with race detector and coverage
make cover

# Lint
make lint

# Format
make fmt

# Build local binary
make build
./dist/dcb version
```

### Running tests

```bash
go test ./...                          # all packages
go test ./internal/domain/...          # domain only
go test ./internal/adapters/parser/... # parser only
```

---

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Commit using [Conventional Commits](https://www.conventionalcommits.org/)
4. Open a pull request

Please make sure `make test lint` passes before opening a PR.

---

## License

MIT — see [LICENSE](LICENSE).
