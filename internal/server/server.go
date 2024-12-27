package server

import (
	"github.com/labstack/echo/v4"
	"poll.ant/internal/middlewares"
	user "poll.ant/internal/services/users/presentation"
)

type Server struct {
	engine *echo.Echo
}

func NewServer(
	healthCheckHandler *HealthCheckHandler,
	userController *user.UserController,
) *Server {
	engine := echo.New()
	engine.HTTPErrorHandler = middlewares.ErrorHandler

	engine.GET("/health", healthCheckHandler.check)
	userGroup := engine.Group("/users")

	userController.Route(userGroup)

	return &Server{engine: engine}
}

func (s *Server) Run(addr string) error {
	return s.engine.Start(addr)
}
