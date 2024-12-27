package server

import (
	"github.com/labstack/echo/v4"
	httpCode "poll.ant/internal/libs/http/http-code"
)

type HealthCheckHandler struct {
}

func NewHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) check(c echo.Context) error {
	return c.String(httpCode.Ok.Code, "ok")
}
