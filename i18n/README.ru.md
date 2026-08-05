**Languages:** [English](../README.md) | [日本語](README.ja.md) | [فارسی](README.fa.md) | [Deutsch](README.de.md) | [简体中文](README.zh-CN.md) | [Русский](README.ru.md) | [Español](README.es.md) | [Português](README.pt-BR.md) | [Français](README.fr.md) | [한국어](README.ko.md)

# DCB — Docker Compose Builder

> **Пишите меньше YAML. Запускайте лучшие контейнеры.**

DCB читает упрощённый `dcb.yaml` и генерирует production-ready `docker-compose.yml`
с автоматическими health check, именованными сетями, политиками перезапуска и упорядочиванием зависимостей.

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)

---

## Почему DCB?

`docker-compose.yml` мощный, но многословный. В типичном production-стеке постоянно повторяется один и тот же boilerplate для health check, restart policy, объявления сетей и именования volumes.

С DCB вы пишете так:

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

И получаете полностью настроенный `docker-compose.yml` с:

- ✅ Правильным health check `pg_isready` для postgres (с учётом образа)
- ✅ `depends_on` с `condition: service_healthy`
- ✅ Именованной сетью `webapp_default`
- ✅ Объявлением top-level `volumes:`
- ✅ `restart: unless-stopped` для каждого сервиса
- ✅ Обнаружением циклических зависимостей на этапе парсинга

---

## Установка

### Homebrew (macOS / Linux)
```bash
brew install JrNovaEX/tap/dcb
```

### Go install
```bash
go install github.com/JrNovaEX/DCB/cmd/dcb@latest
```

### Скачать бинарник

Возьмите последний релиз для вашей платформы со
[страницы Releases](https://github.com/JrNovaEX/DCB/releases).

```bash
# Пример для Linux amd64
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### Сборка из исходников
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## Быстрый старт

```bash
# 1. Интерактивно создать новый dcb.yaml
dcb init

# 2. Проверить конфигурацию и доступность портов
dcb check

# 3. Сгенерировать docker-compose.yml
dcb build

# 4. Запустить всё в фоне
dcb up --detach

# 5. Следить за логами всех сервисов
dcb logs -f

# 6. Остановить и удалить
dcb down
```

---

## Справка по dcb.yaml

```yaml
project: myapp          # Обязательно. Используется для пространства имён сетей и volumes.
version: "3.8"          # Необязательно. Версия схемы Docker Compose (по умолчанию: 3.8).

services:
  <service-name>:
    # Источник — выберите ОДИН из image или build:
    image: nginx:alpine

    build: ./myservice          # Короткая форма: путь к context
    build:                      # Полная форма:
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # Публикуется как host:container (одинаковое значение с обеих сторон).
    env_file: .env       # Путь к env-файлу.
    env:                 # Встроенные переменные окружения.
      KEY: value

    depends_on:          # Имена сервисов, которых нужно ждать (condition: service_healthy).
      - db

    healthcheck: true    # Автоматически генерирует health check на основе имени образа.
                         # Поддерживаются: postgres, mysql, mariadb, redis, nginx,
                         #                 mongo, rabbitmq. Остальные получают CMD-SHELL exit 0.

    volume: pgdata:/var/lib/postgresql/data   # Именованный или bind-mount.
    restart: unless-stopped                   # По умолчанию. Можно переопределить.
```

---

## Команды

| Команда | Описание |
|---------|-------------|
| `dcb init` | Интерактивно создать новый `dcb.yaml` |
| `dcb check` | Проверить `dcb.yaml` + доступность портов на хосте |
| `dcb build` | Сгенерировать `docker-compose.yml` из `dcb.yaml` |
| `dcb up` | Собрать и запустить сервисы (`docker compose up`) |
| `dcb down` | Остановить и удалить сервисы (`docker compose down`) |
| `dcb logs` | Потоковый вывод логов сервисов |
| `dcb version` | Вывести версию, commit и дату сборки |

### Глобальные флаги
```
--verbose, -V    Включить отладочные логи
--json           Выводить логи в формате JSON (для CI/CD)
```

### dcb build
```
--file,     -f   Входной файл (по умолчанию: dcb.yaml)
--output,   -o   Выходной файл (по умолчанию: docker-compose.yml)
--env,      -e   Целевое окружение (по умолчанию: dev)
--validate       Запустить docker compose config после генерации (по умолчанию: true)
```

### dcb up
```
--file,     -f   Входной dcb.yaml (по умолчанию: dcb.yaml)
--env,      -e   Целевое окружение
--detach,   -d   Запуск в фоне
--build          Принудительная пересборка образов
```

### dcb down
```
--file,     -f   Compose-файл (по умолчанию: docker-compose.yml)
--project,  -p   Переопределить имя проекта
--volumes,  -v   Удалить именованные volumes
--remove-orphans Удалить orphan-контейнеры
```

### dcb logs
```
--file           Compose-файл (по умолчанию: docker-compose.yml)
--follow,   -f   Потоковый вывод логов
--tail,     -n   Количество строк (-1 = все)
[services...]    Фильтр по именам сервисов
```

---

## Поддержка нескольких окружений

DCB объединяет базовый `dcb.yaml` с окруженческим overlay:

```
dcb build --env prod
```

Загружает `dcb.yaml`, затем накладывает `dcb.prod.yaml`. Ключи из `dcb.prod.yaml` перекрывают базовые; отсутствующие ключи берутся из базы.

---

## Архитектура

DCB следует **Clean Architecture** (Hexagonal):

```
cmd/dcb/
  commands/         CLI-слой — команды Cobra, разбор флагов, DI
internal/
  domain/           Основные сущности — ProjectConfig, Service, Volume, Port
  ports/            Интерфейсы — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         Прикладная логика — Build, Check, Up, Down, Logs
  adapters/
    parser/         YAML-парсер (gopkg.in/yaml.v3, strict mode)
    generator/      Шаблонизатор docker-compose.yml (text/template)
    runner/         Исполнитель Docker CLI (os/exec)
    config/         Загрузчик мульти-окруженческой конфигурации (Viper)
  infrastructure/
    logger/         Инициализация slog
    validator/      Обёртка go-playground/validator
pkg/version/        Информация о версии на этапе сборки
```

Зависимости направлены внутрь: **cmd → usecases → domain**. Адаптеры зависят от портов, а не друг от друга.

---

## Разработка

```bash
# Запустить все тесты
make test

# Тесты с race detector и coverage
make cover

# Lint
make lint

# Format
make fmt

# Собрать локальный бинарник
make build
./dist/dcb version
```

### Запуск тестов

```bash
go test ./...                          # все пакеты
go test ./internal/domain/...          # только domain
go test ./internal/adapters/parser/... # только parser
```

---

## Участие в разработке

1. Сделайте Fork репозитория
2. Создайте feature-ветку (`git checkout -b feat/my-feature`)
3. Коммитьте с использованием [Conventional Commits](https://www.conventionalcommits.org/)
4. Откройте pull request

Перед открытием PR убедитесь, что `make test lint` проходит успешно.

---

## Лицензия

MIT — см. [LICENSE](../LICENSE).
