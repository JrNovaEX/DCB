**Languages:** [English](../README.md) | [日本語](README.ja.md) | [فارسی](README.fa.md) | [Deutsch](README.de.md) | [简体中文](README.zh-CN.md) | [Русский](README.ru.md) | [Español](README.es.md) | [Português](README.pt-BR.md) | [Français](README.fr.md) | [한국어](README.ko.md)

# DCB — Docker Compose Builder

> **YAML کمتر بنویس. کانتینرهای بهتری اجرا کن.**

DCB یک فایل ساده‌شده‌ی `dcb.yaml` را می‌خواند و یک `docker-compose.yml` آماده برای محیط production تولید می‌کند؛
با health check خودکار، networkهای نام‌گذاری‌شده، سیاست‌های restart و ترتیب وابستگی‌ها.

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)

---

## چرا DCB؟

`docker-compose.yml` قدرتمند است اما پرحرف. در یک استک معمولی production، boilerplateهای تکراری برای health check، restart policy، تعریف network و نام‌گذاری volume بارها تکرار می‌شود.

با DCB می‌توانی این را بنویسی:

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

و یک `docker-compose.yml` کاملاً پیکربندی‌شده دریافت کنی با:

- ✅ health check صحیح `pg_isready` برای postgres (آگاه از image)
- ✅ `depends_on` با `condition: service_healthy`
- ✅ network نام‌گذاری‌شده‌ی `webapp_default`
- ✅ تعریف top-level برای `volumes:`
- ✅ `restart: unless-stopped` روی تمام سرویس‌ها
- ✅ تشخیص وابستگی دایره‌ای در زمان parse

---

## نصب

### Homebrew (macOS / Linux)
```bash
brew install JrNovaEX/tap/dcb
```

### Go install
```bash
go install github.com/JrNovaEX/DCB/cmd/dcb@latest
```

### دانلود باینری

آخرین نسخه را برای پلتفرم خود از
[صفحه Releases](https://github.com/JrNovaEX/DCB/releases) دانلود کنید.

```bash
# مثال Linux amd64
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### ساخت از سورس
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## شروع سریع

```bash
# ۱. ساخت تعاملی یک dcb.yaml جدید
dcb init

# ۲. اعتبارسنجی کانفیگ و بررسی در دسترس بودن پورت‌ها
dcb check

# ۳. تولید docker-compose.yml
dcb build

# ۴. اجرای همه چیز در پس‌زمینه
dcb up --detach

# ۵. دنبال کردن لاگ تمام سرویس‌ها
dcb logs -f

# ۶. خاموش کردن و پاک‌سازی
dcb down
```

---

## مرجع dcb.yaml

```yaml
project: myapp          # الزامی. برای namespace کردن network و volume استفاده می‌شود.
version: "3.8"          # اختیاری. نسخه schemaی Docker Compose (پیش‌فرض: 3.8).

services:
  <service-name>:
    # منبع — یکی از image یا build را انتخاب کنید:
    image: nginx:alpine

    build: ./myservice          # فرم کوتاه: مسیر context
    build:                      # فرم بلند:
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # به صورت host:container با همان مقدار در هر دو طرف expose می‌شود.
    env_file: .env       # مسیر فایل env.
    env:                 # متغیرهای محیطی درون‌خطی.
      KEY: value

    depends_on:          # نام سرویس‌هایی که باید منتظرشان بماند (با condition: service_healthy).
      - db

    healthcheck: true    # به صورت خودکار health check بر اساس نام image تولید می‌کند.
                         # پشتیبانی‌شده: postgres, mysql, mariadb, redis, nginx,
                         #            mongo, rabbitmq. بقیه CMD-SHELL exit 0 می‌گیرند.

    volume: pgdata:/var/lib/postgresql/data   # named یا bind mount.
    restart: unless-stopped                   # پیش‌فرض. در صورت نیاز override کنید.
```

---

## دستورات

| دستور | توضیح |
|---------|-------------|
| `dcb init` | ساخت تعاملی یک `dcb.yaml` جدید |
| `dcb check` | اعتبارسنجی `dcb.yaml` + بررسی در دسترس بودن پورت‌های میزبان |
| `dcb build` | تولید `docker-compose.yml` از روی `dcb.yaml` |
| `dcb up` | ساخت و سپس اجرای سرویس‌ها (`docker compose up`) |
| `dcb down` | توقف و حذف سرویس‌ها (`docker compose down`) |
| `dcb logs` | نمایش استریم لاگ سرویس‌ها |
| `dcb version` | چاپ نسخه، commit و تاریخ build |

### فلگ‌های سراسری
```
--verbose, -V    فعال‌سازی لاگ دیباگ
--json           خروجی لاگ به صورت JSON (برای پایپ‌لاین‌های CI/CD)
```

### dcb build
```
--file,     -f   فایل ورودی (پیش‌فرض: dcb.yaml)
--output,   -o   فایل خروجی (پیش‌فرض: docker-compose.yml)
--env,      -e   محیط هدف (پیش‌فرض: dev)
--validate       اجرای docker compose config بعد از تولید (پیش‌فرض: true)
```

### dcb up
```
--file,     -f   فایل dcb.yaml ورودی (پیش‌فرض: dcb.yaml)
--env,      -e   محیط هدف
--detach,   -d   اجرا در پس‌زمینه
--build          اجبار به rebuild کردن imageها
```

### dcb down
```
--file,     -f   فایل Compose (پیش‌فرض: docker-compose.yml)
--project,  -p   override نام پروژه
--volumes,  -v   حذف named volumeها
--remove-orphans حذف کانتینرهای یتیم
```

### dcb logs
```
--file           فایل Compose (پیش‌فرض: docker-compose.yml)
--follow,   -f   استریم خروجی لاگ
--tail,     -n   تعداد خطوط برای نمایش (-1 = همه)
[services...]    فیلتر بر اساس نام سرویس‌های خاص
```

---

## پشتیبانی از چند محیط (Multi-environment)

DCB فایل پایه‌ی `dcb.yaml` را با یک overlay مخصوص محیط ادغام می‌کند:

```
dcb build --env prod
```

ابتدا `dcb.yaml` را بارگذاری می‌کند، سپس `dcb.prod.yaml` را روی آن merge می‌کند. کلیدهای موجود در `dcb.prod.yaml` جایگزین می‌شوند؛ کلیدهای غایب از فایل پایه گرفته می‌شوند.

---

## معماری

DCB از **Clean Architecture** (Hexagonal) پیروی می‌کند:

```
cmd/dcb/
  commands/         لایه CLI — دستورات Cobra، پارس فلگ، سیم‌کشی DI
internal/
  domain/           موجودیت‌های اصلی — ProjectConfig, Service, Volume, Port
  ports/            اینترفیس‌ها — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         منطق برنامه — Build, Check, Up, Down, Logs
  adapters/
    parser/         پارسر YAML (gopkg.in/yaml.v3، حالت strict)
    generator/      قالب‌ساز docker-compose.yml (text/template)
    runner/         اجراکننده Docker CLI (os/exec)
    config/         بارگذار کانفیگ چندمحیطی (Viper)
  infrastructure/
    logger/         راه‌انداز slog
    validator/      پوشش go-playground/validator
pkg/version/        اطلاعات نسخه در زمان build
```

وابستگی‌ها به سمت داخل جریان دارند: **cmd → usecases → domain**. Adapterها به ports وابسته هستند، نه به یکدیگر.

---

## توسعه

```bash
# اجرای تمام تست‌ها
make test

# تست با race detector و coverage
make cover

# Lint
make lint

# Format
make fmt

# ساخت باینری محلی
make build
./dist/dcb version
```

### اجرای تست‌ها

```bash
go test ./...                          # همه پکیج‌ها
go test ./internal/domain/...          # فقط domain
go test ./internal/adapters/parser/... # فقط parser
```

---

## مشارکت

1. ریپو را Fork کنید
2. یک feature branch بسازید (`git checkout -b feat/my-feature`)
3. با [Conventional Commits](https://www.conventionalcommits.org/) کامیت بزنید
4. یک pull request باز کنید

لطفاً قبل از باز کردن PR مطمئن شوید که `make test lint` پاس می‌شود.

---

## لایسنس

MIT — به [LICENSE](../LICENSE) مراجعه کنید.
