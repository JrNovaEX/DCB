## DCB — Docker Compose Builder
## Usage: make <target>

BINARY      := dcb
MODULE      := github.com/JrNovaEX/DCB
CMD         := ./cmd/dcb
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE        := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS     := -ldflags "-s -w \
  -X $(MODULE)/pkg/version.Version=$(VERSION) \
  -X $(MODULE)/pkg/version.Commit=$(COMMIT) \
  -X $(MODULE)/pkg/version.Date=$(DATE)"

BUILD_DIR   := dist
GO          := go
GOTEST      := $(GO) test
COVER_OUT   := coverage.out
COVER_HTML  := coverage.html

.PHONY: all build clean test test-verbose cover lint fmt vet tidy install snapshot help

## all: build the binary (default target)
all: build

## build: compile DCB binary into ./dist/
build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) $(CMD)
	@echo "✓ Built $(BUILD_DIR)/$(BINARY) ($(VERSION))"

## install: install DCB to GOPATH/bin
install:
	$(GO) install $(LDFLAGS) $(CMD)
	@echo "✓ Installed $(BINARY)"

## test: run all tests
test:
	$(GOTEST) ./... -count=1 -timeout 60s

## test-verbose: run all tests with -v
test-verbose:
	$(GOTEST) ./... -v -count=1 -timeout 60s

## cover: run tests with coverage and open HTML report
cover:
	$(GOTEST) ./... -coverprofile=$(COVER_OUT) -covermode=atomic -timeout 60s
	$(GO) tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	@echo "✓ Coverage report: $(COVER_HTML)"
	@$(GO) tool cover -func=$(COVER_OUT) | tail -1

## lint: run golangci-lint (requires golangci-lint to be installed)
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/"; exit 1; }
	golangci-lint run ./...

## fmt: run gofmt on all source files
fmt:
	gofmt -w -s .

## vet: run go vet
vet:
	$(GO) vet ./...

## tidy: tidy and verify go.mod / go.sum
tidy:
	$(GO) mod tidy
	$(GO) mod verify

## snapshot: build all release targets locally via goreleaser (no publish)
snapshot:
	@command -v goreleaser >/dev/null 2>&1 || { \
		echo "goreleaser not found. Install: https://goreleaser.com/install/"; exit 1; }
	goreleaser release --snapshot --clean

## clean: remove build artefacts
clean:
	rm -rf $(BUILD_DIR) $(COVER_OUT) $(COVER_HTML)
	@echo "✓ Cleaned"

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
