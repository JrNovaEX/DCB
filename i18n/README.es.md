**Languages:** [English](../README.md) | [日本語](README.ja.md) | [فارسی](README.fa.md) | [Deutsch](README.de.md) | [简体中文](README.zh-CN.md) | [Русский](README.ru.md) | [Español](README.es.md) | [Português](README.pt-BR.md) | [Français](README.fr.md) | [한국어](README.ko.md)

# DCB — Docker Compose Builder

> **Escribe menos YAML. Ejecuta mejores contenedores.**

DCB lee un `dcb.yaml` simplificado y genera un `docker-compose.yml` listo para producción
con health checks automáticos, redes con nombre, políticas de reinicio y ordenamiento de dependencias.

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)

---

## ¿Por qué DCB?

`docker-compose.yml` es potente pero verboso. Un stack de producción típico repite el mismo
boilerplate para health checks, políticas de reinicio, declaraciones de red y nombres de volúmenes.

Con DCB escribes esto:

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

Y obtienes un `docker-compose.yml` completamente configurado con:

- ✅ Health check correcto `pg_isready` para postgres (consciente de la imagen)
- ✅ `depends_on` con `condition: service_healthy`
- ✅ Red con nombre `webapp_default`
- ✅ Declaración top-level de `volumes:`
- ✅ `restart: unless-stopped` en todos los servicios
- ✅ Detección de dependencias circulares en tiempo de parseo

---

## Instalación

### Homebrew (macOS / Linux)
```bash
brew install JrNovaEX/tap/dcb
```

### Go install
```bash
go install github.com/JrNovaEX/DCB/cmd/dcb@latest
```

### Descargar binario

Obtén la última versión para tu plataforma desde la
[página de Releases](https://github.com/JrNovaEX/DCB/releases).

```bash
# Ejemplo Linux amd64
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### Compilar desde el código fuente
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## Inicio rápido

```bash
# 1. Crear un nuevo dcb.yaml de forma interactiva
dcb init

# 2. Validar la configuración y comprobar disponibilidad de puertos
dcb check

# 3. Generar docker-compose.yml
dcb build

# 4. Iniciar todo en segundo plano
dcb up --detach

# 5. Seguir los logs de todos los servicios
dcb logs -f

# 6. Detener y limpiar
dcb down
```

---

## Referencia de dcb.yaml

```yaml
project: myapp          # Obligatorio. Se usa para el namespace de redes y volúmenes.
version: "3.8"          # Opcional. Versión del schema de Docker Compose (por defecto: 3.8).

services:
  <service-name>:
    # Origen — elige UNO entre image o build:
    image: nginx:alpine

    build: ./myservice          # Forma corta: ruta al context
    build:                      # Forma larga:
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # Expuesto como host:container (mismo valor en ambos lados).
    env_file: .env       # Ruta a un archivo env.
    env:                 # Variables de entorno en línea.
      KEY: value

    depends_on:          # Nombres de servicios a esperar (condition: service_healthy).
      - db

    healthcheck: true    # Genera automáticamente un health check según el nombre de la imagen.
                         # Soportados: postgres, mysql, mariadb, redis, nginx,
                         #             mongo, rabbitmq. Otros reciben CMD-SHELL exit 0.

    volume: pgdata:/var/lib/postgresql/data   # Named o bind mount.
    restart: unless-stopped                   # Por defecto. Se puede sobrescribir.
```

---

## Comandos

| Comando | Descripción |
|---------|-------------|
| `dcb init` | Crear un nuevo `dcb.yaml` de forma interactiva |
| `dcb check` | Validar `dcb.yaml` + comprobar disponibilidad de puertos del host |
| `dcb build` | Generar `docker-compose.yml` a partir de `dcb.yaml` |
| `dcb up` | Construir e iniciar servicios (`docker compose up`) |
| `dcb down` | Detener y eliminar servicios (`docker compose down`) |
| `dcb logs` | Transmitir logs de los servicios |
| `dcb version` | Mostrar versión, commit y fecha de build |

### Flags globales
```
--verbose, -V    Activar logging de depuración
--json           Salida de logs en JSON (para pipelines CI/CD)
```

### dcb build
```
--file,     -f   Archivo de entrada (por defecto: dcb.yaml)
--output,   -o   Archivo de salida (por defecto: docker-compose.yml)
--env,      -e   Entorno objetivo (por defecto: dev)
--validate       Ejecutar docker compose config tras la generación (por defecto: true)
```

### dcb up
```
--file,     -f   dcb.yaml de entrada (por defecto: dcb.yaml)
--env,      -e   Entorno objetivo
--detach,   -d   Ejecutar en segundo plano
--build          Forzar reconstrucción de imágenes
```

### dcb down
```
--file,     -f   Archivo Compose (por defecto: docker-compose.yml)
--project,  -p   Sobrescribir el nombre del proyecto
--volumes,  -v   Eliminar volúmenes con nombre
--remove-orphans Eliminar contenedores huérfanos
```

### dcb logs
```
--file           Archivo Compose (por defecto: docker-compose.yml)
--follow,   -f   Transmitir la salida de logs
--tail,     -n   Número de líneas a mostrar (-1 = todas)
[services...]    Filtrar por nombres de servicios específicos
```

---

## Soporte multi-entorno

DCB combina un `dcb.yaml` base con un overlay específico del entorno:

```
dcb build --env prod
```

Carga `dcb.yaml` y luego fusiona `dcb.prod.yaml` encima. Las claves de `dcb.prod.yaml` sobrescriben la base; las claves ausentes se toman de la base.

---

## Arquitectura

DCB sigue **Clean Architecture** (Hexagonal):

```
cmd/dcb/
  commands/         Capa CLI — comandos Cobra, parsing de flags, cableado DI
internal/
  domain/           Entidades principales — ProjectConfig, Service, Volume, Port
  ports/            Interfaces — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         Lógica de aplicación — Build, Check, Up, Down, Logs
  adapters/
    parser/         Parser YAML (gopkg.in/yaml.v3, modo strict)
    generator/      Generador de plantillas docker-compose.yml (text/template)
    runner/         Ejecutor de Docker CLI (os/exec)
    config/         Cargador de configuración multi-entorno (Viper)
  infrastructure/
    logger/         Inicializador de slog
    validator/      Wrapper de go-playground/validator
pkg/version/        Información de versión en tiempo de build
```

Las dependencias fluyen hacia adentro: **cmd → usecases → domain**. Los adapters dependen de los ports, no entre sí.

---

## Desarrollo

```bash
# Ejecutar todas las pruebas
make test

# Pruebas con race detector y coverage
make cover

# Lint
make lint

# Format
make fmt

# Construir el binario local
make build
./dist/dcb version
```

### Ejecutar pruebas

```bash
go test ./...                          # todos los paquetes
go test ./internal/domain/...          # solo domain
go test ./internal/adapters/parser/... # solo parser
```

---

## Contribuir

1. Haz Fork del repositorio
2. Crea una rama de funcionalidad (`git checkout -b feat/my-feature`)
3. Haz commit usando [Conventional Commits](https://www.conventionalcommits.org/)
4. Abre un pull request

Asegúrate de que `make test lint` pase antes de abrir un PR.

---

## Licencia

MIT — consulta [LICENSE](../LICENSE).
