**Languages:** [English](../README.md) | [日本語](README.ja.md) | [فارسی](README.fa.md) | [Deutsch](README.de.md) | [简体中文](README.zh-CN.md) | [Русский](README.ru.md) | [Español](README.es.md) | [Português](README.pt-BR.md) | [Français](README.fr.md) | [한국어](README.ko.md)

# DCB — Docker Compose Builder

> **YAMLを少なく書く。より良いコンテナを動かす。**

DCBは簡略化された `dcb.yaml` を読み取り、自動ヘルスチェック、名前付きネットワーク、再起動ポリシー、依存関係の順序付けを備えた本番対応の `docker-compose.yml` を生成します。

[![CI](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml/badge.svg)](https://github.com/JrNovaEX/DCB/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/JrNovaEX/DCB)](https://goreportcard.com/report/github.com/JrNovaEX/DCB)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)

---

## なぜDCBなのか？

`docker-compose.yml` は強力ですが冗長です。典型的な本番スタックでは、ヘルスチェック、再起動ポリシー、ネットワーク宣言、ボリューム命名のボイラープレートが繰り返されます。

DCBでは次のように書けます：

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

そして完全に構成された `docker-compose.yml` が得られます：

- ✅ postgres用の正しい `pg_isready` ヘルスチェック（イメージ認識）
- ✅ `condition: service_healthy` 付きの `depends_on`
- ✅ 名前付きネットワーク `webapp_default`
- ✅ トップレベルの `volumes:` 宣言
- ✅ すべてのサービスに `restart: unless-stopped`
- ✅ パース時の循環依存検出

---

## インストール

### Homebrew (macOS / Linux)
```bash
brew install JrNovaEX/tap/dcb
```

### Go install
```bash
go install github.com/JrNovaEX/DCB/cmd/dcb@latest
```

### バイナリのダウンロード

お使いのプラットフォーム用の最新リリースを
[Releasesページ](https://github.com/JrNovaEX/DCB/releases) から取得してください。

```bash
# Linux amd64 の例
curl -Lo dcb https://github.com/JrNovaEX/DCB/releases/latest/download/dcb_linux_amd64
chmod +x dcb && sudo mv dcb /usr/local/bin/
```

### ソースからビルド
```bash
git clone https://github.com/JrNovaEX/DCB
cd DCB
make install
```

---

## クイックスタート

```bash
# 1. 対話形式で新しい dcb.yaml を作成
dcb init

# 2. 設定を検証し、ホストのポート空き状況を確認
dcb check

# 3. docker-compose.yml を生成
dcb build

# 4. バックグラウンドですべてを起動
dcb up --detach

# 5. すべてのサービスのログを追跡
dcb logs -f

# 6. 終了・削除
dcb down
```

---

## dcb.yaml リファレンス

```yaml
project: myapp          # 必須。ネットワークとボリュームの名前空間に使用。
version: "3.8"          # 任意。Docker Composeスキーマバージョン（デフォルト: 3.8）。

services:
  <service-name>:
    # ソース — image または build のどちらか一方を選択：
    image: nginx:alpine

    build: ./myservice          # 短縮形：コンテキストへのパス
    build:                      # 詳細形：
      context: ./myservice
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production

    port: 8080           # host:container として公開（両側同じ値）。
    env_file: .env       # envファイルへのパス。
    env:                 # インラインの環境変数。
      KEY: value

    depends_on:          # 待機するサービス名（condition: service_healthy）。
      - db

    healthcheck: true    # イメージ名に基づいてヘルスチェックを自動生成。
                         # 対応: postgres, mysql, mariadb, redis, nginx,
                         #       mongo, rabbitmq。その他は CMD-SHELL exit 0。

    volume: pgdata:/var/lib/postgresql/data   # 名前付きまたはバインドマウント。
    restart: unless-stopped                   # デフォルト。必要に応じて上書き。
```

---

## コマンド

| コマンド | 説明 |
|---------|-------------|
| `dcb init` | 対話形式で新しい `dcb.yaml` を作成 |
| `dcb check` | `dcb.yaml` の検証 + ホストポートの空き確認 |
| `dcb build` | `dcb.yaml` から `docker-compose.yml` を生成 |
| `dcb up` | ビルドしてサービスを起動（`docker compose up`） |
| `dcb down` | サービスを停止・削除（`docker compose down`） |
| `dcb logs` | サービスログをストリーム表示 |
| `dcb version` | バージョン、コミット、ビルド日を表示 |

### グローバルフラグ
```
--verbose, -V    デバッグログを有効化
--json           ログをJSON形式で出力（CI/CDパイプライン向け）
```

### dcb build
```
--file,     -f   入力ファイル（デフォルト: dcb.yaml）
--output,   -o   出力ファイル（デフォルト: docker-compose.yml）
--env,      -e   対象環境（デフォルト: dev）
--validate       生成後に docker compose config を実行（デフォルト: true）
```

### dcb up
```
--file,     -f   入力 dcb.yaml（デフォルト: dcb.yaml）
--env,      -e   対象環境
--detach,   -d   バックグラウンドで実行
--build          イメージの強制再ビルド
```

### dcb down
```
--file,     -f   Composeファイル（デフォルト: docker-compose.yml）
--project,  -p   プロジェクト名の上書き
--volumes,  -v   名前付きボリュームを削除
--remove-orphans 孤立コンテナを削除
```

### dcb logs
```
--file           Composeファイル（デフォルト: docker-compose.yml）
--follow,   -f   ログ出力をストリーム
--tail,     -n   表示する行数（-1 = すべて）
[services...]    特定のサービス名でフィルタ
```

---

## マルチ環境サポート

DCBはベースの `dcb.yaml` と環境固有のオーバーレイをマージします：

```
dcb build --env prod
```

`dcb.yaml` を読み込み、その上に `dcb.prod.yaml` をマージします。`dcb.prod.yaml` のキーは上書きされ、欠けているキーはベースから引き継がれます。

---

## アーキテクチャ

DCBは **Clean Architecture**（ヘキサゴナル）に従います：

```
cmd/dcb/
  commands/         CLI層 — Cobraコマンド、フラグ解析、DI配線
internal/
  domain/           コアエンティティ — ProjectConfig, Service, Volume, Port
  ports/            インターフェース — ComposeGenerator, ConfigParser, DockerRunner
  usecases/         アプリケーションロジック — Build, Check, Up, Down, Logs
  adapters/
    parser/         YAMLパーサー（gopkg.in/yaml.v3、strictモード）
    generator/      docker-compose.yml テンプレーター（text/template）
    runner/         Docker CLI実行器（os/exec）
    config/         マルチ環境設定ローダー（Viper）
  infrastructure/
    logger/         slog初期化
    validator/      go-playground/validator ラッパー
pkg/version/        ビルド時バージョン情報
```

依存関係は内側に向かって流れます：**cmd → usecases → domain**。アダプターはポートに依存し、互いに依存しません。

---

## 開発

```bash
# すべてのテストを実行
make test

# レースディテクタとカバレッジ付きでテスト
make cover

# Lint
make lint

# Format
make fmt

# ローカルバイナリをビルド
make build
./dist/dcb version
```

### テストの実行

```bash
go test ./...                          # すべてのパッケージ
go test ./internal/domain/...          # domainのみ
go test ./internal/adapters/parser/... # parserのみ
```

---

## コントリビューション

1. リポジトリをフォーク
2. フィーチャーブランチを作成（`git checkout -b feat/my-feature`）
3. [Conventional Commits](https://www.conventionalcommits.org/) でコミット
4. プルリクエストを開く

PRを開く前に `make test lint` が通ることを確認してください。

---

## ライセンス

MIT — [LICENSE](../LICENSE) を参照。
