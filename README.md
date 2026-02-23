# go-template

[![CI](https://github.com/CheeryProgrammer/go-template/actions/workflows/main-push.yml/badge.svg)](https://github.com/CheeryProgrammer/go-template/actions/workflows/main-push.yml)
[![Go](https://img.shields.io/github/go-mod/go-version/CheeryProgrammer/go-template)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-blue)](#)

A good starting point for pet projects and prototypes. Click **"Use this
template"** and skip the boring setup — CI/CD pipelines, Docker, linting rules,
and deployment scripts are already configured.

Pull requests run lint and tests automatically. Merging to main builds a Docker
image and deploys to staging. Releases go to production. All CI/CD logic lives
in [goship](https://github.com/CheeryProgrammer/goship) reusable workflows —
configure once, reuse across all your repos.

**You write the code. The template handles the rest.**

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
docker-compose.yml               # production compose file (uploaded to server on deploy)
docker-compose.dev.yml           # local dev (postgres with volume)
docker-compose.test.yml          # integration test services (postgres + redis)
Makefile                         # build, test, lint, docker, migrate targets
.golangci.yml                    # golangci-lint v2 ruleset
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

CI runs immediately with no configuration. Secrets are only needed to enable
optional features.

**CI secrets** (all optional):

| Name | Used by | Description |
|------|---------|-------------|
| `CODECOV_TOKEN` | pr-checks, nightly | Codecov upload token |
| `INTEGRATION_ENV` | pr-checks, nightly | Newline-separated `KEY=VALUE` env vars for integration tests — [details](https://github.com/CheeryProgrammer/goship#integration-testyml--integration-tests) |
| `SLACK_WEBHOOK_URL` | nightly | Slack incoming webhook for failure alerts |

**CD secrets & variables** (required to enable deploys):

| Name | Type | Used by | Description |
|------|------|---------|-------------|
| `DEPLOY_ENABLED` | Variable | main-push, release | Set to `true` to enable the deploy job |
| `STAGING_SSH_KNOWN_HOSTS` | Variable | main-push | `known_hosts` entry for staging server |
| `PROD_SSH_KNOWN_HOSTS` | Variable | release | `known_hosts` entry for production server |
| `STAGING_SSH_PRIVATE_KEY` | Secret | main-push | SSH deploy key for staging |
| `STAGING_DATABASE_URL` | Secret | main-push | Full postgres URL for staging migrations |
| `PROD_SSH_PRIVATE_KEY` | Secret | release | SSH deploy key for production |
| `PROD_DATABASE_URL` | Secret | release | Full postgres URL for production migrations |

> **Deploy is skipped by default.** The deploy job in `main-push.yml` and
> `release.yml` only runs when `DEPLOY_ENABLED` is set to `true`. This means
> CI works out of the box after creating the repo from the template, without
> needing any secrets configured upfront.

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

---

## Migrations

The template uses [golang-migrate](https://github.com/golang-migrate/migrate).
Each migration is two SQL files: `up` (apply) and `down` (rollback).

### Adding a migration

```bash
make migrate-create NAME=create_users_table
```

Creates two files:

```
migrations/
  000001_create_users_table.up.sql
  000001_create_users_table.down.sql
```

Fill them in:

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

```sql
-- 000001_create_users_table.down.sql
DROP TABLE users;
```

### Running migrations locally

```bash
make dev-up        # start postgres (first time)
make migrate-up    # apply all pending migrations
make migrate-status  # check current version
make migrate-down  # roll back the last migration
```

### In CI/CD

Migrations run automatically before every deploy when `run-migrations: true`
is set in the workflow. The database URL is taken from the
`STAGING_DATABASE_URL` / `PROD_DATABASE_URL` secret.

### Rules

- **Never edit a migration file after it has been applied** — create a new one instead
- **Always write `down.sql`** — even if it's just `DROP TABLE`
- **One change per migration** — easier to review and roll back
- **Commit together with code** — the migration and the code that uses it go in the same PR

---

## Queries with sqlc (optional)

[sqlc](https://sqlc.dev) generates type-safe Go code from SQL queries.
Instead of writing `rows.Scan(...)` by hand, you write SQL and get Go for free.

### Setup

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Create `sqlc.yaml` in the repo root:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/store/queries/"
    schema:  "migrations/"
    gen:
      go:
        package:       "store"
        out:           "internal/store"
        sql_package:   "pgx/v5"
```

### Workflow

**1. Write a query in `internal/store/queries/users.sql`:**

```sql
-- name: GetUser :one
SELECT id, email, created_at FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, created_at FROM users ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (email) VALUES ($1) RETURNING *;
```

**2. Generate Go code:**

```bash
make sqlc   # runs: sqlc generate
```

sqlc creates `internal/store/users.sql.go` with fully typed functions — no
`Scan`, no string casting, compiler catches wrong argument types.

**3. Use the generated code:**

```go
user, err := store.GetUser(ctx, userID)
users, err := store.ListUsers(ctx)
user, err := store.CreateUser(ctx, "alice@example.com")
```

**4. Enable CI check** so generated files are always up to date.

In `pr-checks.yml` and `main-push.yml`, add `sqlc-enabled: true`:

```yaml
uses: CheeryProgrammer/goship/.github/workflows/ci-pipeline.yml@main
with:
  sqlc-enabled: true
  # ...
```

CI will fail if someone edits a query without regenerating the Go code.

---

## Reusable workflows

All CI/CD logic lives in [goship](https://github.com/CheeryProgrammer/goship).
See that repository for the full list of inputs, outputs, and secrets.
