# syntax=docker/dockerfile:1

ARG GO_VERSION=1.27.1

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata wget \
    && adduser -D -u 10001 app
USER app
WORKDIR /app

COPY --from=build /out/api /app/api

EXPOSE 3000
ENTRYPOINT ["/app/api"]
