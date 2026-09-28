package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/Gravitazione/go-fiber-playground/internal/config"
	"github.com/Gravitazione/go-fiber-playground/internal/database"
	"github.com/Gravitazione/go-fiber-playground/internal/health"
	"github.com/Gravitazione/go-fiber-playground/internal/redis"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	db, err := database.NewPostgres(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close postgres: %v", err)
		}
	}()

	rdb, err := redis.New(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("close redis: %v", err)
		}
	}()

	healthService := health.NewService(db, rdb)
	healthHandler := health.NewHandler(healthService)

	app := fiber.New()
	app.Get("/health", healthHandler.Check)

	go func() {
		<-ctx.Done()
		log.Println("shutting down...")
		_ = app.ShutdownWithTimeout(10 * time.Second)
	}()

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Printf("server: %v", err)
	}
}
