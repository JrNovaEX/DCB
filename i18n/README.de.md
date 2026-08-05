**Languages:** [English](../README.md) | [日本語](README.ja.md) | [فارسی](README.fa.md) | [Deutsch](README.de.md) | [简体中文](README.zh-CN.md) | [Русский](README.ru.md) | [Español](README.es.md) | [Português](README.pt-BR.md) | [Français](README.fr.md) | [한국어](README.ko.md)

# DCB — Docker Compose Builder

> **Weniger YAML schreiben. Bessere Container betreiben.**

DCB liest eine vereinfachte `dcb.yaml` und erzeugt eine produktionsreife `docker-compose.yml`
mit automatischen Health Checks, benannten Netzwerken, Restart-Policies und Abhängigkeitsreihenfolge.

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)

---

## Warum DCB?

`docker-compose.yml` ist mächtig, aber ausführlich. In einem typischen Produktions-Stack wiederholen sich dieselben Boilerplates für Health Checks, Restart-Policies, Netzwerkdeklarationen und Volume-Benennung.

Mit DCB schreibst du das:

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

Und erhältst eine vollständig konfigurierte `docker-compose.yml` mit:

- ✅ Korrektem `pg_isready`-Health-Check für Postgres (image-aware)
- ✅ `depends_on` mit `condition: service_healthy`
- ✅ Benanntem Netzwerk `webapp_default`
- ✅ Top-Level-`volumes:`-Deklaration
- ✅ `restart: unless-stopped` bei jedem Service
- ✅ Erkennung zirkulärer Abhängigkeiten zur Parse-Zeit

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

### Binary herunterladen

Hole dir das neueste Release für deine Plattform von der
[Releases-Seite](https://github.com/JrNovaEX/DCB/releases).

```bash
# Linux-amd64-Beispiel
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### Aus dem Quellcode bauen
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## Schnellstart

```bash
# 1. Interaktiv eine neue dcb.yaml anlegen
dcb init

# 2. Konfiguration validieren und Port-Verfügbarkeit prüfen
dcb check

# 3. docker-compose.yml generieren
dcb build

# 4. Alles im Hintergrund starten
dcb up --detach

# 5. Logs aller Services verfolgen
dcb logs -f

# 6. Herunterfahren und aufräumen
dcb down
```

---

## dcb.yaml-Referenz

```yaml
project: myapp          # Pflicht. Wird zur Namespace-Bildung von Netzwerken und Volumes genutzt.
version: "3.8"          # Optional. Docker-Compose-Schema-Version (Standard: 3.8).

services:
  <service-name>:
    # Quelle — wähle ENTWEDER image ODER build:
    image: nginx:alpine

    build: ./myservice          # Kurzform: Pfad zum Context
    build:                      # Langform:
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # Als host:container exponiert (gleicher Wert auf beiden Seiten).
    env_file: .env       # Pfad zu einer Env-Datei.
    env:                 # Inline-Umgebungsvariablen.
      KEY: value

    depends_on:          # Service-Namen, auf die gewartet wird (condition: service_healthy).
      - db

    healthcheck: true    # Generiert automatisch einen Health Check basierend auf dem Image-Namen.
                         # Unterstützt: postgres, mysql, mariadb, redis, nginx,
                         #            mongo, rabbitmq. Andere erhalten CMD-SHELL exit 0.

    volume: pgdata:/var/lib/postgresql/data   # Named oder Bind-Mount.
    restart: unless-stopped                   # Standard. Bei Bedarf überschreiben.
```

---

## Befehle

| Befehl | Beschreibung |
|---------|-------------|
| `dcb init` | Interaktiv eine neue `dcb.yaml` anlegen |
| `dcb check` | `dcb.yaml` validieren + Host-Port-Verfügbarkeit prüfen |
| `dcb build` | `docker-compose.yml` aus `dcb.yaml` generieren |
| `dcb up` | Bauen und Services starten (`docker compose up`) |
| `dcb down` | Services stoppen und entfernen (`docker compose down`) |
| `dcb logs` | Service-Logs streamen |
| `dcb version` | Version, Commit und Build-Datum ausgeben |

### Globale Flags
```
--verbose, -V    Debug-Logging aktivieren
--json           Logs im JSON-Format ausgeben (für CI/CD-Pipelines)
```

### dcb build
```
--file,     -f   Eingabedatei (Standard: dcb.yaml)
--output,   -o   Ausgabedatei (Standard: docker-compose.yml)
--env,      -e   Zielumgebung (Standard: dev)
--validate       Nach Generierung docker compose config ausführen (Standard: true)
```

### dcb up
```
--file,     -f   Eingabe-dcb.yaml (Standard: dcb.yaml)
--env,      -e   Zielumgebung
--detach,   -d   Im Hintergrund ausführen
--build          Images erzwingen neu bauen
```

### dcb down
```
--file,     -f   Compose-Datei (Standard: docker-compose.yml)
--project,  -p   Projektname überschreiben
--volumes,  -v   Benannte Volumes entfernen
--remove-orphans Verwaiste Container entfernen
```

### dcb logs
```
--file           Compose-Datei (Standard: docker-compose.yml)
--follow,   -f   Log-Ausgabe streamen
--tail,     -n   Anzahl der anzuzeigenden Zeilen (-1 = alle)
[services...]    Nach bestimmten Service-Namen filtern
```

---

## Multi-Environment-Unterstützung

DCB merged eine Basis-`dcb.yaml` mit einem umgebungsspezifischen Overlay:

```
dcb build --env prod
```

Lädt `dcb.yaml` und merged anschließend `dcb.prod.yaml` darüber. Schlüssel in `dcb.prod.yaml` überschreiben die Basis; fehlende Schlüssel fallen auf die Basis zurück.

---

## Architektur

DCB folgt der **Clean Architecture** (Hexagonal):

```
cmd/dcb/
  commands/         CLI-Schicht — Cobra-Befehle, Flag-Parsing, DI-Verdrahtung
internal/
  domain/           Kern-Entitäten — ProjectConfig, Service, Volume, Port
  ports/            Interfaces — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         Anwendungslogik — Build, Check, Up, Down, Logs
  adapters/
    parser/         YAML-Parser (gopkg.in/yaml.v3, strict mode)
    generator/      docker-compose.yml-Templater (text/template)
    runner/         Docker-CLI-Ausführer (os/exec)
    config/         Multi-Env-Config-Loader (Viper)
  infrastructure/
    logger/         slog-Initialisierer
    validator/      go-playground/validator-Wrapper
pkg/version/        Build-Zeit-Versionsinfo
```

Abhängigkeiten fließen nach innen: **cmd → usecases → domain**. Adapter hängen von Ports ab, nicht voneinander.

---

## Entwicklung

```bash
# Alle Tests ausführen
make test

# Tests mit Race-Detector und Coverage
make cover

# Lint
make lint

# Format
make fmt

# Lokales Binary bauen
make build
./dist/dcb version
```

### Tests ausführen

```bash
go test ./...                          # alle Packages
go test ./internal/domain/...          # nur domain
go test ./internal/adapters/parser/... # nur parser
```

---

## Mitwirken

1. Repo forken
2. Feature-Branch erstellen (`git checkout -b feat/my-feature`)
3. Mit [Conventional Commits](https://www.conventionalcommits.org/) committen
4. Pull Request öffnen

Bitte stelle sicher, dass `make test lint` vor dem Öffnen eines PR erfolgreich ist.

---

## Lizenz

MIT — siehe [LICENSE](../LICENSE).
