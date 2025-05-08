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

func GetAllStaff(c *fiber.Ctx) error {
	var staff []model.Staff

	// Truy vấn tất cả nhân viên từ cơ sở dữ liệu
	if err := database.DB.Find(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to fetch staff data",
		})
	}

	// Trả về danh sách nhân viên dưới dạng JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Staff list fetched successfully",
		"data":    staff,
	})
}

func CreateStaff(c *fiber.Ctx) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Lấy thông tin nhân viên từ request body
		var input model.CreateStaffInput
		var validate = validator.New()
		// Parse JSON từ request body vào struct
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid input %s", err.Error()),
			})
		}

		// Validate input
		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if input.StatusWorking != "" && !utils.IsValidValueOfConstant(input.StatusWorking, constants.TYPE_WORKING) {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.TYPE_WORKING_INVALID, errors.New("typeWorking invalid"), "typeWorking")
		}
		if input.Position != "" && !utils.IsValidValueOfConstant(input.Position, constants.TYPE_POSITION) {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.TYPE_POSITION_INVALID, errors.New("typeWorking invalid"), "typePosition")
		}
		if input.Gender != "" && !utils.IsValidValueOfConstant(input.Gender, constants.GENDER) {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.GENDER_INVALID, errors.New("gender invalid"), "gender")
		}
		// Save input to context locals
		c.Locals("inputCreateStaff", input)

		// Continue to next handler
		return c.Next()
	}

}
func EditStaff(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := c.Params(key)
		valueKey, err := strconv.Atoi(params)
		var validate = validator.New()
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.DATA_INPUT_IS_NOT_NUMBER, errors.New("params invalid"))
		}
		var input model.UpdateStaffInput

		// Parse JSON từ request body vào struct
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid input %s", err.Error()),
			})
		}

		// Validate input
		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// check type dữ liệu
		if input.StatusWorking != "" && !utils.IsValidValueOfConstant(input.StatusWorking, constants.TYPE_WORKING) {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.TYPE_WORKING_INVALID, errors.New("typeWorking invalid"), "typeWorking")
		}
		if input.Position != "" && !utils.IsValidValueOfConstant(input.Position, constants.TYPE_POSITION) {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.TYPE_POSITION_INVALID, errors.New("typeWorking invalid"), "typePosition")
		}
		if input.Gender != "" && !utils.IsValidValueOfConstant(input.Gender, constants.GENDER) {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.GENDER_INVALID, errors.New("gender invalid"), "gender")
		}
		// Save input to context locals
		c.Locals("UpdateStaffInput", input)
		c.Locals("inputStaffId", valueKey)

		// Continue to next handler
		return c.Next()
	}
}
func DeleteStaff(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := c.Params(key)

		staffId, err := strconv.Atoi(params)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid staff ID", err)
		}

		db := database.DB

		// Kiểm tra tồn tại
		var staff model.Staff
		if err := db.First(&staff, staffId).Error; err != nil {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Staff not found", err)
		}

		c.Locals("staffId", staffId)
		return c.Next()

	}
}
func ActiveStaff(c *fiber.Ctx) fiber.Handler {
	return func(c *fiber.Ctx) error {
		isActive := c.Query("active")
		staffId := c.Params("staffId")
		valueKeyIsActive, err := strconv.ParseBool(isActive)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.DATA_INPUT_IS_NOT_BOOL, errors.New("params invalid"))
		}
		valueKeyStaffId, err := strconv.Atoi(staffId)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.DATA_INPUT_IS_NOT_NUMBER, errors.New("params invalid"))
		}

		// Save input to context locals
		c.Locals("isActive", valueKeyIsActive)
		c.Locals("staffId", valueKeyStaffId)

		// Continue to next handler
		return c.Next()
	}
}
func StaffChangePassword(c *fiber.Ctx) fiber.Handler {
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

		if input.OldPassword == input.NewPassword {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.NEW_PASSWORD_SAME_CURRENT_PASSWORD, errors.New("newPassword invalid"), "newPassword")
		}
		if input.NewPassword != input.ConfirmPassword {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.NEW_PASSWORD_NOT_SAME_REPEAT_PASSWORD, errors.New("repeatPassword invalid"), "repeatPassword")
		}
		if len(input.NewPassword) < 6 || len(input.NewPassword) > 50 {
			return utils.ErrorResponseHaveKey(c, fiber.StatusBadRequest, constants.NEW_PASSWORD_LENGTH_INVALID, errors.New("Độ dài mật khẩu không phù hợp "), "newPassword")
		}
		// Save input to context locals
		c.Locals("staffChangePasswordInput", input)

		// Continue to next handler
		return c.Next()
	}
}
