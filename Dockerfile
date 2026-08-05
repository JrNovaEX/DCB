# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# Install git for go mod download (private deps) and ca-certs for TLS.
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Cache module downloads before copying source.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and build.
COPY . .

ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags "-s -w \
        -X github.com/JrNovaEX/DCB/pkg/version.Version=${VERSION} \
        -X github.com/JrNovaEX/DCB/pkg/version.Commit=${COMMIT} \
        -X github.com/JrNovaEX/DCB/pkg/version.Date=${DATE}" \
      -o /dcb ./cmd/dcb

# ── Stage 2: Minimal runtime ─────────────────────────────────────────────────
FROM scratch

# Copy CA certificates for HTTPS (needed if DCB ever calls external APIs).
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the statically compiled binary.
COPY --from=builder /dcb /usr/local/bin/dcb

# DCB reads dcb.yaml from the working directory.
WORKDIR /workspace

ENTRYPOINT ["/usr/local/bin/dcb"]
CMD ["--help"]
