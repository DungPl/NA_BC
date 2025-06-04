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
	auth.Get("/exportExcel", handler.ExportExcelHandler)
	account := v1.Group("/account", logger.New())
	account.Get("/", middleware.Protected(), handler.Me)
	account.Post("/changePassword/:staffId", validate.AdminChangePassword(&fiber.Ctx{}), handler.AdminChangePassword)

	staffAdmin := v1.Group("/staff", middleware.Protected(), middleware.CheckAdminHcns)

	staffAdmin.Post("/addStaff", validate.CreateStaff(&fiber.Ctx{}), handler.CreateStaff)
	staffAdmin.Get("/getAllStaff", handler.GetAllStaff)
	staffAdmin.Get("/getStaffById/:staffId", handler.GetStaffById)
	staffAdmin.Put("/updateStaff/:staffId", validate.EditStaff("staffId"), handler.EditStaff)
	staffAdmin.Delete("/deleteStaff/:staffId", validate.DeleteStaff("staffId"), handler.DeleteStaff)
	staffAdmin.Patch("/active/:staffId", validate.ActiveStaff(&fiber.Ctx{}), handler.ActiveStaff)

	staff := v1.Group("/profile", middleware.Protected())

	staff.Post("/changePassword", validate.StaffChangePassword(&fiber.Ctx{}), handler.StaffChangePassword)

	customer := v1.Group("/customer", middleware.Protected(), middleware.AuthSale)
	customer.Get("/getCustomer", handler.GetAllCustomer)
	customer.Post("/addCustomer", validate.CreateCustomer(&fiber.Ctx{}), handler.CreateCustomer)
	customer.Put("/updateCustomer/:customerId", validate.EditCustomer("customerId"), handler.EditCustomer)
	customer.Delete("deleteCustomer/:customerId", validate.DeleteCustomer("customerId"), handler.DeleteCustomer)
	customer.Patch("/tranfer/:customerId", validate.TranferManager("customerId"), handler.TranferManager)

	order := v1.Group("/order", middleware.Protected(), middleware.AuthSale)

	order.Post("/createOrderDraft", validate.CreateOrder(&fiber.Ctx{}), handler.CreateOrder)
	order.Patch("/sendOrder/:orderId", validate.SendOrder("orderId"), handler.SendOrder)
	order.Put("/editOrder/:orderId", validate.EditDraftOrder("orderId"), handler.EditDraftOrder)
	order.Delete("/deleteDraftOrder/:orderId", validate.DeleteDraftOrder("orderId"), handler.DeleteDraftOrder)
	order.Get("previewDraftOrder/:orderId", validate.DeleteDraftOrder("orderId"), handler.PreviewDraftOrder)
	order.Get("downloadDraftOrder/:orderId", validate.DeleteDraftOrder("orderId"), handler.DownloadDraftOrder)
	order.Patch("cancelOrder/:orderId", validate.DeleteDraftOrder("orderId"), handler.CancelDraftOrder)
	order.Patch("updateStatusOrder/:orderId", validate.UpdateStatusOrder("orderId"), handler.UpdateStatusOrder)
	order.Post("createInvoice/:orderId", validate.CreateEditInvoice("orderId"), handler.CreateEditInvoice)

	orderAdmin := v1.Group("/Admin", middleware.Protected(), middleware.CheckAdmin)
	orderAdmin.Patch("/insertOrder/:orderId", validate.AdminEditOrder("orderId"), handler.AdminEditOrder)
	orderAdmin.Patch("/updateRevision/:revisionInvoiceId", validate.UpdateRevisionStatus("revisionInvoiceId"), handler.UpdateRevisionStatus)
	orderAdmin.Get("/listOrder", validate.ListOrder(&fiber.Ctx{}), handler.ListOrder)
	orderAdmin.Get("/listRevisionInvoice", handler.ListInvoice)

}
