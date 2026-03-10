# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Cache dependency downloads separately from source
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and build
COPY . .

ARG VERSION=dev
ARG BUILD_TIME
ARG COMMIT_SHA

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w \
      -X main.Version=${VERSION} \
      -X main.BuildTime=${BUILD_TIME} \
      -X main.CommitSHA=${COMMIT_SHA}" \
    -o /build/bin/app \
    ./cmd/server

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
FROM alpine:3.21

# ca-certificates: TLS for outbound requests; wget: Docker healthcheck
RUN apk add --no-cache ca-certificates wget tzdata

# nobody (65534) already exists in Alpine — use it as the non-root user
USER nobody:nobody

WORKDIR /app

COPY --from=builder /build/bin/app /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]
