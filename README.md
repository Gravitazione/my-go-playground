# go-fiber-playground

A personal playground for experimenting with Go and the [Fiber v3](https://github.com/gofiber/fiber) web framework.

## Requirements

- Go 1.27.1 (pinned in `mise.toml`)
- [mise](https://mise.jdx.dev/) (optional, for managing the Go version)

## Getting started

```sh
# install the pinned Go version (if using mise)
mise install

# download dependencies
go mod download

# run the API server
go run ./cmd/api
```

The server listens on `http://localhost:3000`.

## Endpoints

| Method | Path      | Description  | Response           |
| ------ | --------- | ------------ | ------------------ |
| GET    | `/health` | Health check | `{"status":"ok"}`  |

```sh
curl http://localhost:3000/health
```

## Project structure

```
.
├── cmd/
│   └── api/
│       └── main.go   # Fiber app entrypoint
├── go.mod
├── go.sum
└── mise.toml         # tool versions (Go)
```
