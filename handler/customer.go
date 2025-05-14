package handler

import (
	"errors"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/helper"
	"order-manager/model"
	"order-manager/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAllCustomer(c *fiber.Ctx) error {
	db := database.DB
	var customers []model.Customer

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
	query := db.Model(&model.Customer{}).Where("is_active = ?", true)

	// Phân trang
	var total int64
	query.Count(&total)
	// Tìm kiếm nếu có
	if filter.SearchKey != "" {
		key := "%" + strings.ToLower(filter.SearchKey) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(identification_card) LIKE ? OR LOWER(phone_number) LIKE ?", key, key, key)
	}
	if err := query.Limit(filter.Limit).Offset(filter.Offset).Find(&customers).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Lỗi lấy danh sách nhân viên", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"data": customers,
		"pagination": fiber.Map{
			"total":  total,
			"limit":  filter.Limit,
			"offset": filter.Offset,
		},
	})
}
func CreateCustomer(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	//return utils.SuccessResponse(c, fiber.StatusOK, dataInfo)
	CustomerInput, ok := c.Locals("inputCreateCustomer").(model.Customer)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("PARSE DATA TO LOCALS FAIL"))
	}
	var existingPhone model.Customer
	if err := database.DB.Where("whats_app =?", CustomerInput.WhatsApp).First(&existingPhone).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone Number  already exists",
		})
	}
	if len(CustomerInput.WhatsApp) != 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Phone Number is not active"})
	}
	var existingEmail model.Customer
	if err := database.DB.Where("email = ?", CustomerInput.Email).First(&existingEmail).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}
	customer := model.Customer{
		Name:          CustomerInput.Name,
		Email:         CustomerInput.Email,
		WhatsApp:      CustomerInput.WhatsApp,
		Address:       CustomerInput.Address,
		AddressOrigin: CustomerInput.AddressOrigin,
		Gender:        CustomerInput.Gender,
		Note:          CustomerInput.Note,
		ManagerId:     &dataInfo.AccountId,
	}
	if err := database.DB.Create(&customer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create customer",
		})
	}
	db := database.DB
	var account model.Account
	if err := db.Preload("Account").First(&account, &dataInfo.AccountId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Customer created successfully",
		"data": fiber.Map{
			"customer": customer,
		},
	})

}
