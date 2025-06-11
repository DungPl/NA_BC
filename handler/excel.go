package handler

import (
	"fmt"
	"order-manager/utils"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ExportExcelHandler(c *fiber.Ctx) error {
	filename := "Payment_Request_3T.xlsx"

	err := utils.ExportPaymentRequestExcel(filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate Excel file",
		})
	}

	return c.Download(filename, "ĐỀ NGHỊ THANH TOÁN 3T.xlsx")

}
func ExportExcelv2Handler(c *fiber.Ctx) error {
	filename := "Payment_Request_3T.xlsx"

	// Extract JSON data from the request body
	data := c.Body()
	if len(data) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Request body is empty",
		})
	}

	// Generate the Excel file
	err := utils.ExportPaymentRequestExcelv2(filename, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate Excel file: " + err.Error(),
		})
	}

	// Ensure the file is deleted after download
	defer func() {
		if err := os.Remove(filename); err != nil {
			// Log the error (in a real application, use a proper logger)
			fmt.Printf("Failed to delete temporary file %s: %v\n", filename, err)
		}
	}()

	// Set the Content-Disposition header to force download with the desired filename
	c.Set("Content-Disposition", `attachment; filename="ĐỀ NGHỊ THANH TOÁN 3T.xlsx"`)
	return c.SendFile(filename)
}
func ExportExcelv3Handler(c *fiber.Ctx) error {
	filename := "Bảng kê 1.xlsx"
	err := utils.ExportPaymentRequestExcelv3(filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate Excel file",
		})
	}

	return c.Download(filename, "Bảng kê 1.xlsx")
}
