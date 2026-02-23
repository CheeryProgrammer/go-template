# go-template

[![Go](https://img.shields.io/github/go-mod/go-version/CheeryProgrammer/go-template)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-blue)](#)

GitHub template for Go services. Click **"Use this template"** to generate a
new repository with CI/CD, Docker, linting, and database migrations already
wired up via [goship](https://github.com/CheeryProgrammer/goship) reusable workflows.

---

## What's included

```
cmd/server/main.go               # entry point

internal/
  config/config.go               # configuration from env vars
  server/server.go               # HTTP server lifecycle (start / graceful shutdown)
  handler/
    handler.go                   # router setup (chi) + middleware
    health.go                    # GET /health
  store/
    store.go                     # Store interface — add your query methods here
    postgres.go                  # database/sql + pgx implementation (optional)

migrations/                      # SQL migration files (golang-migrate)

.github/workflows/
  pr-checks.yml                  # CI on every pull request
  main-push.yml                  # CI + staging deploy on push to main
  release.yml                    # CI + GitHub Release + production deploy on semver tag
  nightly.yml                    # full suite nightly + Slack alert on failure

Dockerfile                       # multi-stage scratch image (Go 1.24)
docker-compose.yml               # local dev (postgres with volume)
docker-compose.test.yml          # integration test services (postgres + redis)
Makefile                         # build, test, lint, docker, migrate targets
.golangci.yml                    # golangci-lint ruleset
.env.example                     # env var reference
```

---

## Setup

### 1. Rename the module

Replace `YOUR_ORG/myapp` with your actual module path everywhere:

```bash
find . -type f \( -name '*.go' -o -name 'go.mod' \) \
  | xargs sed -i 's|github.com/YOUR_ORG/myapp|github.com/myorg/myservice|g'
```

### 2. Replace the goship placeholder

Replace `YOUR_ORG/goship` in `.github/workflows/` with the actual location
of your goship instance:

```bash
grep -rl 'YOUR_ORG/goship' .github/workflows/ \
  | xargs sed -i 's|YOUR_ORG/goship|myorg/goship|g'
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Adjust project settings

Edit the top of `Makefile`:

```makefile
BINARY_NAME  ?= myservice
MAIN_PACKAGE ?= ./cmd/server
DOCKER_IMAGE ?= ghcr.io/myorg/myservice
```

### 5. Add secrets & variables

| Name | Type | Used by |
|------|------|---------|
| `CODECOV_TOKEN` | Secret | pr-checks, main-push, nightly |
| `INTEGRATION_ENV` | Secret | pr-checks, main-push, nightly |
| `STAGING_SSH_PRIVATE_KEY` | Secret | main-push |
| `STAGING_DATABASE_URL` | Secret | main-push |
| `PROD_SSH_PRIVATE_KEY` | Secret | release |
| `PROD_DATABASE_URL` | Secret | release |
| `SLACK_WEBHOOK_URL` | Secret | nightly |
| `STAGING_SSH_KNOWN_HOSTS` | Variable | main-push |
| `PROD_SSH_KNOWN_HOSTS` | Variable | release |

### 6. Push

CI runs automatically on the first pull request.

---

## Local development

```bash
# Start postgres
docker compose up -d

# Copy and edit env
cp .env.example .env

make help              # list all targets
make check             # fmt + vet + lint before committing
make test              # unit tests with race detector
make test-coverage     # tests + HTML coverage report
make test-integration  # integration tests (needs running services)
make build             # build binary → ./bin/app
make docker-build      # build Docker image
make migrate-up        # apply pending migrations
make migrate-create NAME=add_users_table
```

---

## Database (optional)

Set `DATABASE_URL` to connect to postgres. Leave it empty to start the app
without a database — only `GET /health` returns `"database": "disabled"`.

To add queries, define methods on the `Store` interface in
`internal/store/store.go` and implement them in `internal/store/postgres.go`.

---

## Reusable workflows

All CI/CD logic lives in [goship](https://github.com/CheeryProgrammer/goship).
See that repository for the full list of inputs, outputs, and secrets for
each workflow.
