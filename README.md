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

Dockerfile                       # multi-stage scratch image (Go 1.26)
docker-compose.yml               # production (uploaded to server on deploy)
docker-compose.dev.yml           # local dev (postgres with volume)
docker-compose.test.yml          # integration test services (postgres + redis)
Makefile                         # build, test, lint, docker, migrate targets
.golangci.yml                    # golangci-lint ruleset
.env.example                     # env var reference
```

---

## Get Started

### 1. Create a repository from this template

Click **"Use this template"** → **"Create a new repository"** on GitHub,
then clone your new repo:

```bash
git clone git@github.com:your-org/your-service.git
cd your-service
```

### 2. Rename the module

Replace the `YOUR_ORG/myapp` placeholder with your actual module path
everywhere — Go files, `go.mod`, and the lint config:

```bash
# macOS
find . -type f \( -name '*.go' -o -name 'go.mod' -o -name '*.yml' \) \
  | xargs sed -i '' 's|github.com/YOUR_ORG/myapp|github.com/your-org/your-service|g'

# Linux
find . -type f \( -name '*.go' -o -name 'go.mod' -o -name '*.yml' \) \
  | xargs sed -i 's|github.com/YOUR_ORG/myapp|github.com/your-org/your-service|g'
```

Then regenerate the lockfile:

```bash
go mod tidy
```

### 3. Adjust project settings

Edit the variables at the top of `Makefile`:

```makefile
BINARY_NAME  ?= your-service
DOCKER_IMAGE ?= ghcr.io/your-org/your-service
```

### 4. Add GitHub secrets & variables

| Name | Type | Used by |
|------|------|---------|
| `CODECOV_TOKEN` | Secret | pr-checks, nightly |
| `INTEGRATION_ENV` | Secret | pr-checks, nightly — [what is this?](https://github.com/CheeryProgrammer/goship#integration-testyml--integration-tests) |
| `STAGING_SSH_PRIVATE_KEY` | Secret | main-push |
| `STAGING_DATABASE_URL` | Secret | main-push |
| `PROD_SSH_PRIVATE_KEY` | Secret | release |
| `PROD_DATABASE_URL` | Secret | release |
| `SLACK_WEBHOOK_URL` | Secret | nightly |
| `STAGING_SSH_KNOWN_HOSTS` | Variable | main-push |
| `PROD_SSH_KNOWN_HOSTS` | Variable | release |

### 5. Open a pull request

Push a branch and open a PR — `pr-checks` runs CI automatically.

```bash
git checkout -b init
git add -A
git commit -m "init: rename module and configure project"
git push -u origin init
```

---

## Local development

```bash
# Start local postgres
make dev-up

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
without a database — `GET /health` will return `"database": "disabled"`.

To add queries, define methods on the `Store` interface in
`internal/store/store.go` and implement them in `internal/store/postgres.go`.

---

## Reusable workflows

All CI/CD logic lives in [goship](https://github.com/CheeryProgrammer/goship).
See that repository for the full list of inputs, outputs, and secrets.
