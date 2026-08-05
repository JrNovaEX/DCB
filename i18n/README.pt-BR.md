**Languages:** [English](../README.md) | [日本語](README.ja.md) | [فارسی](README.fa.md) | [Deutsch](README.de.md) | [简体中文](README.zh-CN.md) | [Русский](README.ru.md) | [Español](README.es.md) | [Português](README.pt-BR.md) | [Français](README.fr.md) | [한국어](README.ko.md)

# DCB — Docker Compose Builder

> **Escreva menos YAML. Execute containers melhores.**

O DCB lê um `dcb.yaml` simplificado e gera um `docker-compose.yml` pronto para produção
com health checks automáticos, redes nomeadas, políticas de restart e ordenação de dependências.

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)

---

## Por que DCB?

O `docker-compose.yml` é poderoso, mas verboso. Um stack de produção típico repete o mesmo
boilerplate para health checks, políticas de restart, declarações de rede e nomes de volumes.

Com o DCB você escreve isto:

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

E recebe um `docker-compose.yml` totalmente configurado com:

- ✅ Health check correto `pg_isready` para postgres (consciente da imagem)
- ✅ `depends_on` com `condition: service_healthy`
- ✅ Rede nomeada `webapp_default`
- ✅ Declaração top-level de `volumes:`
- ✅ `restart: unless-stopped` em todos os serviços
- ✅ Detecção de dependências circulares no momento do parse

---

## Instalação

### Homebrew (macOS / Linux)
```bash
brew install JrNovaEX/tap/dcb
```

### Go install
```bash
go install github.com/JrNovaEX/DCB/cmd/dcb@latest
```

### Baixar o binário

Pegue a versão mais recente para a sua plataforma na
[página de Releases](https://github.com/JrNovaEX/DCB/releases).

```bash
# Exemplo Linux amd64
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### Compilar a partir do código-fonte
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## Início rápido

```bash
# 1. Criar um novo dcb.yaml de forma interativa
dcb init

# 2. Validar a configuração e verificar disponibilidade de portas
dcb check

# 3. Gerar o docker-compose.yml
dcb build

# 4. Iniciar tudo em segundo plano
dcb up --detach

# 5. Acompanhar os logs de todos os serviços
dcb logs -f

# 6. Parar e limpar
dcb down
```

---

## Referência do dcb.yaml

```yaml
project: myapp          # Obrigatório. Usado para o namespace de redes e volumes.
version: "3.8"          # Opcional. Versão do schema do Docker Compose (padrão: 3.8).

services:
  <service-name>:
    # Fonte — escolha UM entre image ou build:
    image: nginx:alpine

    build: ./myservice          # Forma curta: caminho do context
    build:                      # Forma longa:
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # Exposto como host:container (mesmo valor nos dois lados).
    env_file: .env       # Caminho para um arquivo env.
    env:                 # Variáveis de ambiente inline.
      KEY: value

    depends_on:          # Nomes dos serviços a aguardar (condition: service_healthy).
      - db

    healthcheck: true    # Gera automaticamente um health check com base no nome da imagem.
                         # Suportados: postgres, mysql, mariadb, redis, nginx,
                         #             mongo, rabbitmq. Outros recebem CMD-SHELL exit 0.

    volume: pgdata:/var/lib/postgresql/data   # Named ou bind mount.
    restart: unless-stopped                   # Padrão. Pode ser sobrescrito.
```

---

## Comandos

| Comando | Descrição |
|---------|-------------|
| `dcb init` | Criar um novo `dcb.yaml` de forma interativa |
| `dcb check` | Validar `dcb.yaml` + verificar disponibilidade de portas do host |
| `dcb build` | Gerar `docker-compose.yml` a partir do `dcb.yaml` |
| `dcb up` | Construir e iniciar serviços (`docker compose up`) |
| `dcb down` | Parar e remover serviços (`docker compose down`) |
| `dcb logs` | Transmitir logs dos serviços |
| `dcb version` | Exibir versão, commit e data de build |

### Flags globais
```
--verbose, -V    Ativar logging de debug
--json           Saída de logs em JSON (para pipelines CI/CD)
```

### dcb build
```
--file,     -f   Arquivo de entrada (padrão: dcb.yaml)
--output,   -o   Arquivo de saída (padrão: docker-compose.yml)
--env,      -e   Ambiente alvo (padrão: dev)
--validate       Executar docker compose config após a geração (padrão: true)
```

### dcb up
```
--file,     -f   dcb.yaml de entrada (padrão: dcb.yaml)
--env,      -e   Ambiente alvo
--detach,   -d   Executar em segundo plano
--build          Forçar rebuild das imagens
```

### dcb down
```
--file,     -f   Arquivo Compose (padrão: docker-compose.yml)
--project,  -p   Sobrescrever o nome do projeto
--volumes,  -v   Remover volumes nomeados
--remove-orphans Remover containers órfãos
```

### dcb logs
```
--file           Arquivo Compose (padrão: docker-compose.yml)
--follow,   -f   Transmitir a saída de logs
--tail,     -n   Número de linhas a exibir (-1 = todas)
[services...]    Filtrar por nomes de serviços específicos
```

---

## Suporte a múltiplos ambientes

O DCB mescla um `dcb.yaml` base com um overlay específico do ambiente:

```
dcb build --env prod
```

Carrega o `dcb.yaml` e depois mescla o `dcb.prod.yaml` por cima. As chaves do `dcb.prod.yaml` sobrescrevem a base; as chaves ausentes caem de volta para a base.

---

## Arquitetura

O DCB segue a **Clean Architecture** (Hexagonal):

```
cmd/dcb/
  commands/         Camada CLI — comandos Cobra, parsing de flags, DI
internal/
  domain/           Entidades principais — ProjectConfig, Service, Volume, Port
  ports/            Interfaces — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         Lógica de aplicação — Build, Check, Up, Down, Logs
  adapters/
    parser/         Parser YAML (gopkg.in/yaml.v3, modo strict)
    generator/      Gerador de templates docker-compose.yml (text/template)
    runner/         Executor do Docker CLI (os/exec)
    config/         Carregador de config multi-ambiente (Viper)
  infrastructure/
    logger/         Inicializador do slog
    validator/      Wrapper do go-playground/validator
pkg/version/        Informações de versão em tempo de build
```

As dependências fluem para dentro: **cmd → usecases → domain**. Os adapters dependem dos ports, não uns dos outros.

---

## Desenvolvimento

```bash
# Executar todos os testes
make test

# Testes com race detector e coverage
make cover

# Lint
make lint

# Format
make fmt

# Compilar o binário local
make build
./dist/dcb version
```

### Executar testes

```bash
go test ./...                          # todos os pacotes
go test ./internal/domain/...          # apenas domain
go test ./internal/adapters/parser/... # apenas parser
```

---

## Contribuindo

1. Faça Fork do repositório
2. Crie uma branch de feature (`git checkout -b feat/my-feature`)
3. Faça commit usando [Conventional Commits](https://www.conventionalcommits.org/)
4. Abra um pull request

Certifique-se de que `make test lint` passe antes de abrir um PR.

---

## Licença

MIT — veja [LICENSE](../LICENSE).
