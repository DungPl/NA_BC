package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ErrorResponse(c *fiber.Ctx, status int, message string, err error) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"errors":  err.Error(),
	})
}

func ErrorResponseHaveKey(c *fiber.Ctx, status int, message string, err error, keyError string) error {
	return c.Status(status).JSON(fiber.Map{
		"status":   "error",
		"message":  message,
		"errors":   err.Error(),
		"keyError": keyError,
	})
}

func SuccessResponse(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func ApplyPagination(query *gorm.DB, limit, page *int) *gorm.DB {
	// Kiểm tra nếu có limit thì thêm điều kiện Limit
	if limit != nil && *limit > 0 && page != nil && *page >= 1 {
		query = query.Limit(*limit)
		offset := *limit * (*page - 1)
		query = query.Offset(offset)
	}

	return query
}

func GetLast12Months() []string {
	months := []string{}
	now := time.Now()

	for i := 0; i < 12; i++ {
		month := now.AddDate(0, -i, 0).Format("01-2006") // Định dạng MM-YYYY
		months = append(months, month)
	}

	return months
}
