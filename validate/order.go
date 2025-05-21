package validate

import (
	"errors"
	"fmt"
	"order-manager/constants"
	"order-manager/model"
	"order-manager/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateOrder(c *fiber.Ctx) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input model.InputDraftOrder

		// Parse JSON từ request body vào struct
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid input: %s", err.Error()),
			})
		}

		validate := validator.New()
		if err := validate.Struct(input); err != nil {
			errors := []string{}
			for _, e := range err.(validator.ValidationErrors) {
				errors = append(errors, fmt.Sprintf("Key: '%s' Error: %s", e.Namespace(), e.Error()))
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": errors,
			})
		}
		c.Locals("inputOrderDraft", input)

		// Continue to next handler
		return c.Next()

	}
}
func SendOrder(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := c.Params(key)
		valueKey, err := strconv.Atoi(params)

		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.DATA_INPUT_IS_NOT_NUMBER, errors.New("params invalid"))
		}

		c.Locals("orderId", valueKey)
		return c.Next()

	}
}
func EditDraftOrder(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := c.Params(key)
		valueKey, err := strconv.Atoi(params)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.DATA_INPUT_IS_NOT_NUMBER, errors.New("params invalid"))
		}
		var input model.InputDraftOrder
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid input %s", err.Error()),
			})
		}

		c.Locals("inputUpdateOrder", input)
		c.Locals("orderId", valueKey)
		return c.Next()
	}
}
func DeleteDraftOrder(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := c.Params(key)
		valueKey, err := strconv.Atoi(params)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.DATA_INPUT_IS_NOT_NUMBER, errors.New("params invalid"))
		}
		c.Locals("orderId", valueKey)
		return c.Next()
	}
}
