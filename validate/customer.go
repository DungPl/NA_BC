package validate

import (
	"errors"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/model"
	"order-manager/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetAllCustomer(c *fiber.Ctx) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var customer []model.Customer

		// Truy vấn tất cả nhân viên từ cơ sở dữ liệu
		if err := database.DB.Find(&customer).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to fetch staff data",
			})
		}

		// Trả về danh sách nhân viên dưới dạng JSON
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Customer list fetched successfully",
			"data":    customer,
		})
	}
}
func CreateCustomer(c *fiber.Ctx) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var customer model.Customer
		if err := c.BodyParser(&customer); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		var validate = validator.New()
		if err := validate.Struct(customer); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if customer.Gender != "" && !utils.IsValidValueOfConstant(customer.Gender, constants.GENDER) {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.GENDER_INVALID, errors.New("gender invalid"), "gender")
		}
		c.Locals("inputCreateCustomer", customer)

		// Continue to next handler
		return c.Next()
	}
}
