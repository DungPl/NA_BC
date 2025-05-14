package validate

import (
	"errors"
	"fmt"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/model"
	"order-manager/utils"
	"strconv"

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
func EditCustomer(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := c.Params(key)
		valueKey, err := strconv.Atoi(params)
		var validate = validator.New()
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.DATA_INPUT_IS_NOT_NUMBER, errors.New("params invalid"))
		}

		var input model.Customer
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid input %s", err.Error()),
			})
		}

		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if input.Gender != "" && !utils.IsValidValueOfConstant(input.Gender, constants.GENDER) {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.GENDER_INVALID, errors.New("gender invalid"), "gender")
		}
		c.Locals("inputUpdateCustomer", input)
		c.Locals("inputCustomerId", valueKey)
		// Continue to next handler
		return c.Next()
	}
}
func DeleteCustomer(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := c.Params(key)

		customerId, err := strconv.Atoi(params)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customer ID", err)
		}

		db := database.DB

		// Kiểm tra tồn tại
		var customer model.Customer
		if err := db.First(&customer, customerId).Error; err != nil {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Staff not found", err)
		}

		c.Locals("customerId", customerId)
		return c.Next()

	}
}
