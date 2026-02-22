# go-template

GitHub template for Go services. Click **"Use this template"** to generate a
new repository with CI/CD, Docker, linting, and database migrations already
wired up via [goship](https://github.com/YOUR_ORG/goship) reusable workflows.

---

## What's included

```
.github/workflows/
  pr-checks.yml      # CI on every pull request
  main-push.yml      # CI + staging deploy on push to main
  release.yml        # CI + GitHub Release + production deploy on semver tag
  nightly.yml        # Full suite nightly + Slack alert on failure

Dockerfile           # Multi-stage scratch image (Go 1.23, linux/amd64)
docker-compose.test.yml  # Postgres + Redis for integration tests
Makefile             # build, test, lint, docker, migrate targets
.golangci.yml        # golangci-lint ruleset
.gitignore
```

---

## Setup

### 1. Replace the org/repo placeholder

Find every `YOUR_ORG/goship` in `.github/workflows/` and replace it with the
actual location of your goship instance:

```bash
grep -rl 'YOUR_ORG/goship' .github/workflows/ \
  | xargs sed -i 's|YOUR_ORG/goship|myorg/goship|g'
```

### 2. Adjust project settings

Edit the top of `Makefile`:

```makefile
BINARY_NAME  ?= myapp
MAIN_PACKAGE ?= ./cmd/server
DOCKER_IMAGE ?= ghcr.io/myorg/myapp
```

### 3. Add secrets & variables

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

### 4. Push

CI runs automatically on the first pull request.

---

## Local development

```bash
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

## Reusable workflows

All CI/CD logic lives in [goship](https://github.com/YOUR_ORG/goship).
See that repository for the full list of inputs, outputs, and secrets for
each workflow.
