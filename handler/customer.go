package handler

import (
	"fmt"
	"order-manager/helper"
	"order-manager/utils"

	"github.com/gofiber/fiber/v2"
)

func AddCustomer(c *fiber.Ctx) error {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	tokenModel, _, _, _ := helper.GetInfoAccountFromToken(c)
	fmt.Println("Account ID:", tokenModel.AccountId)
	fmt.Println("Username:", tokenModel.Username)

	return utils.SuccessResponse(c, fiber.StatusOK, tokenModel)

	// type CreateCusRequest struct {
	// 	Name     string `gorm:"not null" validate:"required" json:"name"`
	// 	Email    string `gorm:"uniqueIndex;not null" validate:"required,email" json:"email"`
	// 	WhatsApp string `gorm:"uniqueIndex;not null" validate:"required" json:"phone"`
	// 	Gender   string `json:"gender"`
	// }
	// var reqC CreateCusRequest
	// if err := c.BodyParser(&reqC); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	// }
	// var existingPhone model.Customer
	// if err := database.DB.Where("what_app=?", reqC.WhatsApp).First(&existingPhone).Error; err == nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Phone already exists",
	// 	})
	// }
	// var existingEmail model.Customer
	// if err := database.DB.Where("email = ?", reqC.Email).First(&existingEmail).Error; err == nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Email already exists",
	// 	})
	//}
}
