package health

import (
	"context"
	"sync"
	"time"
)

const checkTimeout = 2 * time.Second

type Status string

const (
	StatusOK    Status = "ok"
	StatusError Status = "error"
)

type Checker interface {
	Ping(ctx context.Context) error
}

type Services struct {
	Postgres Status `json:"postgres"`
	Redis    Status `json:"redis"`
}

type Response struct {
	Status   Status   `json:"status"`
	Services Services `json:"services"`
}

type Service struct {
	postgres Checker
	redis    Checker
}

func NewService(postgres, redis Checker) *Service {
	return &Service{
		postgres: postgres,
		redis:    redis,
	}
}

func (s *Service) Check(ctx context.Context) Response {
	var (
		wg       sync.WaitGroup
		services Services
	)

	wg.Go(func() { services.Postgres = check(ctx, s.postgres) })
	wg.Go(func() { services.Redis = check(ctx, s.redis) })
	wg.Wait()

	status := StatusOK
	if services.Postgres != StatusOK || services.Redis != StatusOK {
		status = StatusError
	}

	return Response{
		Status:   status,
		Services: services,
	}
}

func (r Response) Healthy() bool {
	return r.Status == StatusOK
}

func check(ctx context.Context, c Checker) Status {
	ctx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	if err := c.Ping(ctx); err != nil {
		return StatusError
	}

	return StatusOK
}
