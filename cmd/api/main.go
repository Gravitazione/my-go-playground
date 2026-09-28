package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/Gravitazione/go-fiber-playground/internal/config"
	"github.com/Gravitazione/go-fiber-playground/internal/database"
	"github.com/Gravitazione/go-fiber-playground/internal/health"
	"github.com/Gravitazione/go-fiber-playground/internal/logger"
	"github.com/Gravitazione/go-fiber-playground/internal/redis"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	if err := run(cfg, log); err != nil {
		log.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run(cfg config.Config, log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.NewPostgres(ctx, cfg.PostgresDSN)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("close postgres", "error", err)
		}
	}()
	log.Info("postgres connected")

	rdb, err := redis.New(ctx, cfg.RedisURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Error("close redis", "error", err)
		}
	}()
	log.Info("redis connected")

	healthService := health.NewService(db, rdb)
	healthHandler := health.NewHandler(healthService)

	app := fiber.New()
	app.Get("/health", healthHandler.Check)

	go func() {
		<-ctx.Done()
		log.Info("shutting down...")
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Error("shutdown server", "error", err)
		}
	}()

	log.Debug("starting server", "port", cfg.AppPort)

	return app.Listen(":" + cfg.AppPort)
}
