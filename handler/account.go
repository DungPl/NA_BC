package handler

import (
	"order-manager/constants"
	"order-manager/database"
	"order-manager/model"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func AddStaff(c *fiber.Ctx) error {
	type CreateStaffRequest struct {
		Username           string `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50"`
		Password           string `json:"password" validate:"required,min=6,max=50"`
		Name               string `json:"name" validate:"required"`
		PhoneNumber        string `json:"phoneNumber" gorm:"uniqueIndex;not null" validate:"required"`
		Email              string `json:"email" gorm:"uniqueIndex;not null" `
		IdentificationCard string `gorm:"not null;uniqueIndex;require" validate:"required,min=12,max=12" json:"identificationCard"`
	}

	var req CreateStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	// Kiểm tra username có trùng không
	var existingAccount model.Account
	if err := database.DB.Where("username = ?", req.Username).First(&existingAccount).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username already exists",
		})
	}

	// Kiểm tra email có trùng không
	var existingStaff model.Staff
	if err := database.DB.Where("email = ?", req.Email).First(&existingStaff).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}
	var existingPhone model.Staff
	if err := database.DB.Where("phone_number =?", req.PhoneNumber).First(&existingPhone).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone already exists",
		})
	}
	var existingIdentity model.Staff
	if err := database.DB.Where("identification_card =?", req.IdentificationCard).First(&existingIdentity).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Identification Card already exists",
		})
	}
	if len(req.IdentificationCard) != 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Identification Card is not active"})
	}
	// 1. Hash mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	account := model.Account{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     constants.ROLE_SALE, // hoặc ROLE_KETOAN, ROLE_QUANLY tùy bạn
	}
	if err := database.DB.Create(&account).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create account",
		})
	}
	staff := model.Staff{
		Name:               req.Name,
		PhoneNumber:        req.PhoneNumber,
		Email:              req.Email,
		IsActive:           true,
		IdentificationCard: req.IdentificationCard,
		AccountId:          &account.ID,
	}
	database.DB.Create(&account)
	database.DB.Create(&staff)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Account & Staff created successfully",
		"data":    staff,
	})
}
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
func GetStaffByID(c *fiber.Ctx) error {
	// Lấy ID từ URL parameter
	staffID := c.Params("id")

	var staff model.Staff
	if err := database.DB.First(&staff, staffID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": staff,
	})
}
func UpdateStaff(c *fiber.Ctx) error {
	// Lấy ID nhân viên từ URL parameter
	staffID := c.Params("id")

	// Kiểm tra xem nhân viên có tồn tại không
	var staff model.Staff
	if err := database.DB.First(&staff, staffID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff not found",
		})
	}

	// Parse dữ liệu yêu cầu cập nhật từ body request
	var req model.Staff
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Cập nhật các thông tin cần thiết
	staff.Name = req.Name
	staff.PhoneNumber = req.PhoneNumber
	staff.Email = req.Email
	staff.IsActive = req.IsActive
	staff.IdentificationCard = req.IdentificationCard

	// Lưu thông tin đã cập nhật
	if err := database.DB.Save(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update staff",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Staff updated successfully",
		"data":    staff,
	})
}
func DeleteStaff(c *fiber.Ctx) error {
	// Lấy ID nhân viên từ URL parameter
	staffID := c.Params("id")

	// Kiểm tra xem nhân viên có tồn tại không
	var staff model.Staff
	if err := database.DB.First(&staff, staffID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff not found",
		})
	}
	if err := database.DB.Delete(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update staff",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Staff delete successfully",
	})
}
