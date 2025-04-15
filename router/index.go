package router

import (
	"order-manager/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	// Middleware
	app.Static("/resources", "./resources")
	app.Static("/uploads", "./uploads")
	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1", logger.New())

	// Auth
	// auth := api.Group("/auth")
	// auth.Post("/login", controller.Login)

	// auth
	// authen := v1.Group("/authen", logger.New())
	// authen.Post("/register")

	// User
	auth := v1.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/refresh-token", handler.RefreshToken)
}
