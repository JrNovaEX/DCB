**Languages:** [English](../README.md) | [日本語](README.ja.md) | [فارسی](README.fa.md) | [Deutsch](README.de.md) | [简体中文](README.zh-CN.md) | [Русский](README.ru.md) | [Español](README.es.md) | [Português](README.pt-BR.md) | [Français](README.fr.md) | [한국어](README.ko.md)

# DCB — Docker Compose Builder

> **YAML를 덜 작성하세요. 더 나은 컨테이너를 실행하세요.**

DCB는 간소화된 `dcb.yaml`을 읽어 자동 헬스 체크, 명명된 네트워크, 재시작 정책, 의존성 순서를 포함한
프로덕션 준비된 `docker-compose.yml`을 생성합니다.

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)

---

## 왜 DCB인가?

`docker-compose.yml`은 강력하지만 장황합니다. 일반적인 프로덕션 스택에서는 헬스 체크, 재시작 정책,
네트워크 선언, 볼륨 이름 지정을 위한 동일한 보일러플레이트가 반복됩니다.

DCB를 사용하면 이렇게 작성할 수 있습니다:

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

그리고 완전히 구성된 `docker-compose.yml`을 얻습니다:

- ✅ postgres용 올바른 `pg_isready` 헬스 체크 (이미지 인식)
- ✅ `condition: service_healthy`가 포함된 `depends_on`
- ✅ 명명된 네트워크 `webapp_default`
- ✅ 최상위 `volumes:` 선언
- ✅ 모든 서비스에 `restart: unless-stopped`
- ✅ 파싱 시 순환 의존성 감지

---

## 설치

### Homebrew (macOS / Linux)
```bash
brew install JrNovaEX/tap/dcb
```

### Go install
```bash
go install github.com/JrNovaEX/DCB/cmd/dcb@latest
```

### 바이너리 다운로드

[Releases 페이지](https://github.com/JrNovaEX/DCB/releases)에서
플랫폼에 맞는 최신 릴리스를 받으세요.

```bash
# Linux amd64 예시
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### 소스에서 빌드
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## 빠른 시작

```bash
# 1. 대화형으로 새 dcb.yaml 생성
dcb init

# 2. 설정 검증 및 호스트 포트 사용 가능 여부 확인
dcb check

# 3. docker-compose.yml 생성
dcb build

# 4. 백그라운드에서 모두 시작
dcb up --detach

# 5. 모든 서비스 로그 추적
dcb logs -f

# 6. 중지 및 정리
dcb down
```

---

## dcb.yaml 참조

```yaml
project: myapp          # 필수. 네트워크와 볼륨의 네임스페이스에 사용됩니다.
version: "3.8"          # 선택. Docker Compose 스키마 버전 (기본값: 3.8).

services:
  <service-name>:
    # 소스 — image 또는 build 중 하나를 선택:
    image: nginx:alpine

    build: ./myservice          # 짧은 형식: context 경로
    build:                      # 긴 형식:
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # host:container로 노출 (양쪽 동일한 값).
    env_file: .env       # env 파일 경로.
    env:                 # 인라인 환경 변수.
      KEY: value

    depends_on:          # 대기할 서비스 이름 (condition: service_healthy).
      - db

    healthcheck: true    # 이미지 이름을 기반으로 헬스 체크를 자동 생성합니다.
                         # 지원: postgres, mysql, mariadb, redis, nginx,
                         #       mongo, rabbitmq. 그 외는 CMD-SHELL exit 0.

    volume: pgdata:/var/lib/postgresql/data   # 명명된 볼륨 또는 바인드 마운트.
    restart: unless-stopped                   # 기본값. 필요 시 재정의 가능.
```

---

## 명령어

| 명령어 | 설명 |
|---------|-------------|
| `dcb init` | 대화형으로 새 `dcb.yaml` 생성 |
| `dcb check` | `dcb.yaml` 검증 + 호스트 포트 사용 가능 여부 확인 |
| `dcb build` | `dcb.yaml`에서 `docker-compose.yml` 생성 |
| `dcb up` | 빌드 후 서비스 시작 (`docker compose up`) |
| `dcb down` | 서비스 중지 및 제거 (`docker compose down`) |
| `dcb logs` | 서비스 로그 스트림 출력 |
| `dcb version` | 버전, 커밋, 빌드 날짜 출력 |

### 전역 플래그
```
--verbose, -V    디버그 로깅 활성화
--json           JSON 형식 로그 출력 (CI/CD 파이프라인용)
```

### dcb build
```
--file,     -f   입력 파일 (기본값: dcb.yaml)
--output,   -o   출력 파일 (기본값: docker-compose.yml)
--env,      -e   대상 환경 (기본값: dev)
--validate       생성 후 docker compose config 실행 (기본값: true)
```

### dcb up
```
--file,     -f   입력 dcb.yaml (기본값: dcb.yaml)
--env,      -e   대상 환경
--detach,   -d   백그라운드에서 실행
--build          이미지 강제 재빌드
```

### dcb down
```
--file,     -f   Compose 파일 (기본값: docker-compose.yml)
--project,  -p   프로젝트 이름 재정의
--volumes,  -v   명명된 볼륨 제거
--remove-orphans 고아 컨테이너 제거
```

### dcb logs
```
--file           Compose 파일 (기본값: docker-compose.yml)
--follow,   -f   로그 출력 스트림
--tail,     -n   표시할 줄 수 (-1 = 전체)
[services...]    특정 서비스 이름으로 필터
```

---

## 다중 환경 지원

DCB는 기본 `dcb.yaml`을 환경별 오버레이와 병합합니다:

```
dcb build --env prod
```

`dcb.yaml`을 로드한 다음 `dcb.prod.yaml`을 위에 병합합니다. `dcb.prod.yaml`의 키는 기본값을 덮어쓰고, 없는 키는 기본값으로 돌아갑니다.

---

## 아키텍처

DCB는 **Clean Architecture**(헥사고날)를 따릅니다:

```
cmd/dcb/
  commands/         CLI 계층 — Cobra 명령, 플래그 파싱, DI 연결
internal/
  domain/           핵심 엔티티 — ProjectConfig, Service, Volume, Port
  ports/            인터페이스 — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         애플리케이션 로직 — Build, Check, Up, Down, Logs
  adapters/
    parser/         YAML 파서 (gopkg.in/yaml.v3, strict 모드)
    generator/      docker-compose.yml 템플릿 생성기 (text/template)
    runner/         Docker CLI 실행기 (os/exec)
    config/         다중 환경 설정 로더 (Viper)
  infrastructure/
    logger/         slog 초기화
    validator/      go-playground/validator 래퍼
pkg/version/        빌드 시 버전 정보
```

의존성은 안쪽으로 흐릅니다: **cmd → usecases → domain**. 어댑터는 포트에 의존하며, 서로에게는 의존하지 않습니다.

---

## 개발

```bash
# 모든 테스트 실행
make test

# 레이스 디텍터 및 커버리지와 함께 테스트
make cover

# Lint
make lint

# Format
make fmt

# 로컬 바이너리 빌드
make build
./dist/dcb version
```

### 테스트 실행

```bash
go test ./...                          # 모든 패키지
go test ./internal/domain/...          # domain만
go test ./internal/adapters/parser/... # parser만
```

---

## 기여하기

1. 저장소를 Fork하세요
2. 기능 브랜치를 만드세요 (`git checkout -b feat/my-feature`)
3. [Conventional Commits](https://www.conventionalcommits.org/)를 사용해 커밋하세요
4. Pull Request를 열어주세요

PR을 열기 전에 `make test lint`가 통과하는지 확인해 주세요.

---

## 라이선스

MIT — [LICENSE](../LICENSE)를 참조하세요.
