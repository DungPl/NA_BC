package handler

import (
	"order-manager/utils"

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
	// file, err := os.Open("test_logo.png")
	// if err != nil {
	// 	fmt.Println("Error opening file:", err)
	// 	return err
	// }
	// defer file.Close()

	// _, format, err := image.DecodeConfig(file)
	// if err != nil {
	// 	fmt.Println("Error decoding image:", err)
	// 	return err
	// }
	// fmt.Println("Image format:", format)
	// return nil
}
