package validate

import (
	"errors"
	"fmt"
	"order-manager/constants"
	"order-manager/model"
	"order-manager/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func AdminChangePassword(c *fiber.Ctx) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input model.ChangePasswordInput
		// Parse JSON từ request body vào struct
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid input %s", err.Error()),
			})
		}
		var validate = validator.New()
		// Validate input
		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if input.NewPassword != input.ConfirmPassword {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.NEW_PASSWORD_NOT_SAME_REPEAT_PASSWORD, errors.New("repeatPassword invalid"), "repeatPassword")
		}
		if len(input.NewPassword) < 6 || len(input.NewPassword) > 50 {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.NEW_PASSWORD_LENGTH_INVALID, errors.New("Độ dài mật khẩu không phù hợp "), "newPassword")
		}
		// Validate input
		c.Locals("AdminChangePasswordInput", input)

		// Continue to next handler
		return c.Next()

	}

}
