package router

import (
	"order-manager/handler"
	"order-manager/middleware"
	"order-manager/validate"

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
	account := v1.Group("/account", logger.New())
	account.Get("/", middleware.Protected(), handler.Me)

	auth.Post("/changePss", handler.ChangePassword)

	staff := v1.Group("/staff", middleware.Protected(), middleware.CheckAdmin)
	staff.Post("/addStaff", validate.CreateStaff(&fiber.Ctx{}), handler.CreateStaff)
	staff.Get("/getAllStaff", handler.GetAllStaff)
	staff.Get("/getStaffById/:staffId", handler.GetStaffById)
	staff.Put("/updateStaff/:staffId", validate.EditStaff("staffId"), handler.EditStaff)
	staff.Delete("/deleteStaff/:staffId", validate.DeleteStaff("staffId"), handler.DeleteStaff)
	staff.Patch("/active/:staffId", validate.ActiveStaff(&fiber.Ctx{}), handler.ActiveStaff)
}
