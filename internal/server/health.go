package server

import "github.com/gofiber/fiber/v2"

type HealthCheckHandler struct {
}

func NewHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) check(c *fiber.Ctx) error {
	return c.SendString("ok")
}
