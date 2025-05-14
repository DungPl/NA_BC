package handler

import (
	"errors"
	"fmt"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/helper"
	"order-manager/model"
	"order-manager/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
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
		IsActive:      true,
		ManagerId:     &dataInfo.AccountId,
	}
	if err := database.DB.Create(&customer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create customer",
		})
	}
	db := database.DB
	var account model.Account
	if err := db.First(&account, dataInfo.AccountId).Error; err != nil {
		db.Rollback()
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
func EditCustomer(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	inputCustomerId := c.Locals("inputCustomerId").(int)
	inputUpdateCustomer, ok := c.Locals("inputUpdateCustomer").(model.Customer)
	//return utils.SuccessResponse(c, fiber.StatusOK, dataInfo)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("PARSE DATA TO LOCALS FAIL"))
	}
	db := database.DB
	tx := db.Begin()
	var customer model.Customer

	var account model.Account
	if err := db.First(&account, &dataInfo.AccountId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	if err := tx.Preload("ManageAccount").First(&customer, inputCustomerId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, constants.NOT_FOUND_RECORDS, err)
	}

	if inputUpdateCustomer.Email != "" && inputUpdateCustomer.Email != customer.Email {
		var existingCustomer model.Customer
		if err := tx.Where("email = ? AND id != ?", inputUpdateCustomer.Email, inputCustomerId).First(&existingCustomer).Error; err == nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
	}
	if inputUpdateCustomer.WhatsApp != "" && inputUpdateCustomer.WhatsApp != customer.WhatsApp {
		if len(inputUpdateCustomer.WhatsApp) != 10 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Phone Number must be 10 digits",
			})
		}
		var existingPhone model.Customer
		if err := tx.Where("whats_app = ? AND id != ?", inputUpdateCustomer.WhatsApp, inputCustomerId).First(&existingPhone).Error; err == nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Phone Number already exists",
			})
		}
	}

	// Bảo vệ id và ManagerId
	originalId := customer.ID // Lưu id gốc
	if inputUpdateCustomer.ManagerId == nil {
		inputUpdateCustomer.ManagerId = customer.ManagerId // Giữ ManagerId hiện tại
	}
	copier.Copy(&customer, &inputUpdateCustomer)
	customer.ID = originalId // Khôi phục id để đảm bảo cập nhật

	// Lưu thay đổi
	if err := tx.Save(&customer).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_EDIT, err)
	}
	tx.Commit()
	return utils.SuccessResponse(c, fiber.StatusOK, customer)
}
func DeleteCustomer(c *fiber.Ctx) error {
	// Sau khi xóa, trả về thông tin của khách hàng
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	customerId := c.Locals("customerId").(int)
	db := database.DB

	// Kiểm tra tồn tại
	var customer model.Customer

	var account model.Account
	if err := db.First(&account, dataInfo.AccountId).Error; err != nil {
		db.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}
	if err := db.First(&customer, customerId).Error; err != nil {
		db.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, constants.NOT_FOUND_RECORDS, err)
	}
	// today := time.Now().Truncate(24 * time.Hour)
	// createdAt := customer.CreatedAt.Truncate(24 * time.Hour)
	// if !createdAt.Equal(today) {
	// 	db.Rollback()
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Can only delete customers created today",
	// 	})
	// }
	now := time.Now().In(time.FixedZone("ICT", 7*3600)) // +07
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour).Add(-1 * time.Nanosecond)

	// // Log để kiểm tra
	fmt.Printf("Customer CreatedAt: %v, Today Start: %v, Today End: %v\n", customer.CreatedAt, todayStart, todayEnd)

	if customer.CreatedAt.Before(todayStart) || customer.CreatedAt.After(todayEnd) {
		db.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Can only delete customers created today",
		})
	}

	// Xóa customer
	if err := db.Model(&model.Customer{}).Where("id = ? ", customerId).Update("is_active", false).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to deactivate staff", err)
	}
	db.Commit()
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Customer deleted successfully",
	})
}
