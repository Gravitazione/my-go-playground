package health

import "github.com/gofiber/fiber/v3"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Check(c fiber.Ctx) error {
	res := h.service.Check(c.Context())

	if !res.Healthy() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(res)
	}

	return c.JSON(res)
}
