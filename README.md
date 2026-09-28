# go-fiber-playground

A personal playground for experimenting with Go and the [Fiber v3](https://github.com/gofiber/fiber) web framework.

## Requirements

- Go 1.27.1 (pinned in `mise.toml`)
- [mise](https://mise.jdx.dev/) (optional, for managing the Go version)
- Docker + Docker Compose (for Postgres and Redis)

## Getting started

```sh
cp .env.example .env
```

### Run everything with Docker Compose

```sh
docker compose up -d --build
```

### Run the API locally

Start only the dependencies, then run the server on the host:

```sh
# install the pinned Go version (if using mise)
mise install

docker compose up -d postgres redis

# the app reads config from environment variables (it does not load .env itself)
set -a && source .env && set +a

go run ./cmd/api
```

The server listens on `http://localhost:3000`.

## Configuration

Config is loaded from environment variables in `internal/config`.

| Variable       | Default                                                                 | Description                                   |
| -------------- | ----------------------------------------------------------------------- | --------------------------------------------- |
| `APP_PORT`     | `3000`                                                                  | HTTP port                                     |
| `LOG_LEVEL`    | `info`                                                                  | `debug`, `info`, `warn` or `error`            |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/playground?sslmode=disable` | Postgres DSN                                  |
| `REDIS_URL`    | `redis://localhost:6379/0`                                              | Redis URL                                     |

See `.env.example` for the full list, including the Docker Compose variables.

## Logging

`internal/logger` provides a shared [`log/slog`](https://pkg.go.dev/log/slog) logger that writes JSON to stdout. `logger.New(level)` also sets it as the slog default, so any package can call `slog.Info(...)` directly.

```go
log := logger.New(cfg.LogLevel)
log.Info("redis connected")
log.Error("close postgres", "error", err)
```

```json
{"time":"2026-09-28T10:00:00Z","level":"INFO","msg":"redis connected"}
```

An unknown `LOG_LEVEL` falls back to `info`.

## Endpoints

| Method | Path      | Description                          |
| ------ | --------- | ------------------------------------ |
| GET    | `/health` | Checks Postgres and Redis connectivity |

```sh
curl http://localhost:3000/health
```

```json
{"status":"ok","services":{"postgres":"ok","redis":"ok"}}
```

Returns `503` with `"status":"error"` if any service is unreachable.

## Development

```sh
golangci-lint run ./...
```

## Project structure

```
.
├── cmd/
│   └── api/
│       └── main.go        # Fiber app entrypoint
├── internal/
│   ├── config/            # env-based configuration
│   ├── database/          # Postgres (GORM) connection
│   ├── error/             # app error codes and wrapping
│   ├── health/            # /health handler and service
│   ├── logger/            # shared slog logger
│   └── redis/             # Redis client
├── .env.example
├── .golangci.yml
├── docker-compose.yml
├── Dockerfile
└── mise.toml              # tool versions (Go)
```
