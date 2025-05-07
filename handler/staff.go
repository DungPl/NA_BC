package handler

import (
	"errors"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/model"
	"order-manager/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateStaff(c *fiber.Ctx) error {
	staffInput, ok := c.Locals("inputCreateStaff").(model.CreateStaffInput)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("PARSE DATA TO LOCALS FAIL"))
	}
	var existingIdentity model.Staff
	if err := database.DB.Where("identification_card =?", staffInput.IdentificationCard).First(&existingIdentity).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Identification Card already exists",
		})
	}
	if len(staffInput.IdentificationCard) != 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Identification Card is not active"})
	}
	// Kiểm tra email có trùng không
	var existingStaff model.Staff
	if err := database.DB.Where("email = ?", staffInput.Email).First(&existingStaff).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}
	var existingPhone model.Staff
	if err := database.DB.Where("phone_number =?", staffInput.PhoneNumber).First(&existingPhone).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone already exists",
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(staffInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	// 1. Tạo account trước
	account := model.Account{
		Username: staffInput.IdentificationCard,
		Password: string(hashedPassword),
		Active:   true,
		Role:     staffInput.Position,
	}
	if err := database.DB.Create(&account).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create account",
		})
	}

	// 2. Tạo staff và gán AccountId
	staff := model.Staff{
		Name:               staffInput.Name,
		PhoneNumber:        staffInput.PhoneNumber,
		Email:              staffInput.Email,
		IdentificationCard: staffInput.IdentificationCard,
		AccountId:          &account.ID,
	}

	if err := database.DB.Create(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create staff",
		})
	}

	// 3. Trả về dữ liệu
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Account & Staff created successfully",
		"data": fiber.Map{
			"account": account,
			"staff":   staff,
		},
	})

}
func EditStaff(c *fiber.Ctx) error {
	inputStaffId := c.Locals("inputStaffId").(int)

	staffInput, ok := c.Locals("UpdateStaffInput").(model.UpdateStaffInput)
	//fmt.Println("Locals: ", c.Locals("UpdateStaffInput"))
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("PARSE DATA TO LOCALS FAIL"))
	}
	db := database.DB
	tx := db.Begin()
	var staff model.Staff
	if err := tx.Preload("Account").First(&staff, inputStaffId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, constants.NOT_FOUND_RECORDS, err)
	}
	staff.Name = staffInput.Name
	staff.PhoneNumber = staffInput.PhoneNumber
	staff.Email = staffInput.Email
	staff.IdentificationCard = staffInput.IdentificationCard
	staff.Gender = staffInput.Gender
	staff.Position = staffInput.Position
	staff.AddressOrigin = staffInput.AddressOrigin
	staff.Address = staffInput.Address
	staff.BirthDay = staffInput.BirthDay

	staff.StatusWorking = staffInput.StatusWorking

	staff.Account.Username = staffInput.IdentificationCard
	staff.Account.Role = staffInput.Position
	if err := tx.Save(staff.Account).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_EDIT, err)
	}

	// Lưu staff
	if err := tx.Save(&staff).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_EDIT, err)
	}

	tx.Commit()
	return utils.SuccessResponse(c, fiber.StatusOK, staff)
}
func GetStaffById(c *fiber.Ctx) error {
	idParam := c.Params("staffId") // Lấy tham số id từ URL
	//fmt.Println("Locals: ", c.Params("staffId"))
	// Chuyển đổi id sang int
	staffId, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid staff ID", err)
	}

	db := database.DB
	var staff model.Staff

	if err := db.Preload("Account").First(&staff, staffId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff not found",
		})
	}
	return utils.SuccessResponse(c, fiber.StatusOK, staff)
}
func DeleteStaff(c *fiber.Ctx) error {
	staffId := c.Locals("staffId").(int)
	db := database.DB

	// Kiểm tra tồn tại
	var staff model.Staff
	if err := db.First(&staff, staffId).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Staff not found", err)
	}

	if err := db.Model(&model.Staff{}).Where("id = ? ", staffId).Update("is_active", false).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to deactivate staff", err)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Staff deleted successfully",
	})
}

func ActiveStaff(c *fiber.Ctx) error {
	staffId, ok := c.Locals("staffId").(int)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("PARSE DATA TO LOCALS FAIL"))
	}
	db := database.DB
	var staff model.Staff
	if err := db.First(&staff, int(staffId)).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, constants.NOT_FOUND_RECORDS, err)
	}
	staff.IsActive = true
	db.Save(&staff).Scan(&staff)

	return utils.SuccessResponse(c, fiber.StatusOK, staff)
}

func GetAllStaff(c *fiber.Ctx) error {
	db := database.DB
	var staffs []model.Staff

	// Lấy filter input từ query params
	var filter model.FilterInput
	if err := c.QueryParser(&filter); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid pagination input", err)
	}

	// Thiết lập mặc định nếu chưa có
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	// Query staff
	query := db.Model(&model.Staff{})

	// Phân trang
	var total int64
	query.Count(&total)
	// Tìm kiếm nếu có
	if filter.SearchKey != "" {
		key := "%" + strings.ToLower(filter.SearchKey) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(identification_card) LIKE ? OR LOWER(phone_number) LIKE ?", key, key, key)
	}

	// Lọc theo trạng thái is_active nếu có
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if err := query.Limit(filter.Limit).Offset(filter.Offset).Find(&staffs).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Lỗi lấy danh sách nhân viên", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"data": staffs,
		"pagination": fiber.Map{
			"total":  total,
			"limit":  filter.Limit,
			"offset": filter.Offset,
		},
	})
}
