package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"poll.ant/internal/config"
	"poll.ant/internal/middlewares"
	auth "poll.ant/internal/services/auth/presentation"
	user "poll.ant/internal/services/users/presentation"
)

type Server struct {
	app *fiber.App
}

func NewServer(
	healthCheckHandler *HealthCheckHandler,
	userController *user.UserController,
	authController *auth.AuthController,
) *Server {

	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
	})

	app.Use(requestid.New(), logger.New(logger.Config{
		Format:     "${time} | ${pid} | ${locals:requestid} | ${status} - ${method} ${path}\u200b\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Origin,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, authorization, X-Requested-With",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))

	app.Get("/health", healthCheckHandler.check)
	userGroup := app.Group("/users")
	userController.Route(userGroup)

	authGroup := userGroup.Group("/auth")
	authController.Route(authGroup)
	return &Server{app: app}
}

func (s *Server) Run(addr string) error {
	return s.app.Listen(addr)
}
