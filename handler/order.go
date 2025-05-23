package handler

import (
	"errors"
	"fmt"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/helper"
	"order-manager/model"
	"order-manager/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateOrder(c *fiber.Ctx) error {
	// Get the request body
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	input, ok := c.Locals("inputOrderDraft").(model.InputDraftOrder)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("PARSE DATA TO LOCALS FAIL"))
	}
	db := database.DB
	tx := db.Begin()
	var account model.Account
	if err := tx.First(&account, dataInfo.AccountId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}
	if account.Role != "SALE" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only Sale account",
		})
	}
	var customer model.Customer
	if err := tx.First(&customer, *input.CustomerId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Customer not found", err)
	}
	var orderDate time.Time
	if input.OrderDate == "" {
		orderDate = time.Now().In(time.FixedZone("ICT", 7*3600)).Truncate(24 * time.Hour)
	} else {
		parsedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", input.OrderDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid orderDate format: %s", err.Error()),
			})
		}
		orderDate = parsedDate
	}
	orderDatePtr := &orderDate
	// Tạo mã đơn
	today := time.Now().In(time.FixedZone("ICT", 7*3600)).Format("20060102")
	var count int64
	tx.Model(&model.Order{}).Where("order_code LIKE ?", fmt.Sprintf("ORD-%s%%", today)).Count(&count)
	orderCode := fmt.Sprintf("ORD-%s-%03d", today, count+1)

	// Tạo đơn nháp
	order := model.Order{
		OrderCode:    orderCode,
		OrderDate:    orderDatePtr,
		CustomerId:   input.CustomerId,
		CustomerName: input.CustomerName,
		PhoneNumber:  input.PhoneNumber,
		Address:      input.Address,
		CreatedById:  &dataInfo.AccountId,
		Discount:     input.Discount,
		Status:       "bản nháp",
	}
	var totalAmount float64

	for _, item := range *input.OrderItems {
		total := float64(*item.Quantity) * *item.UnitPrice
		totalAmount += total
		order.OrderItems = append(order.OrderItems, model.OrderItem{
			ProductName: item.Name,
			SizeInch:    &item.Size,
			Quantity:    item.Quantity,
			Unit:        &item.Unit,
			UnitPrice:   item.UnitPrice,
			Subtotal:    &totalAmount,
		})
	}
	order.TotalAmount = totalAmount - input.Discount
	if order.TotalAmount < 0 {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Total amount cannot be negative",
		})
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create order", err)
	}

	tx.Commit()
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Draft order created successfully",
		"data":    order,
	})
}
func SendOrder(c *fiber.Ctx) error {
	// Lấy thông tin đơn hàng
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	orderId, ok := c.Locals("orderId").(int)
	if !ok || orderId <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.DATA_INPUT_IS_NOT_NUMBER, errors.New("invalid customer ID"))
	}
	db := database.DB
	tx := db.Begin()
	var account model.Account
	if err := tx.First(&account, dataInfo.AccountId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}
	var order model.Order
	// Tìm đơn hàng

	if err := tx.Preload("OrderItems").First(&order, orderId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found", err)
	}
	if order.CreatedById == nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Order has no creator assigned",
		})
	}
	if *order.CreatedById != dataInfo.AccountId {
		tx.Rollback()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fmt.Sprintf("You are not authorized to send this order (Creator: %d, Your ID: %d)", *order.CreatedById, dataInfo.AccountId),
		})
	}
	if order.Status != "bản nháp" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only draft orders can be sent",
		})
	}
	//now := time.Now().In(time.FixedZone("ICT", 7*3600))
	//estimatedShipAt := now.Add(7 * 24 * time.Hour)
	order.Status = "Đã gửi và chưa sản xuất"
	//order.FactoryReceiveAt = &now
	//order.EstimatedShipAt = &estimatedShipAt

	// Lưu thay đổi
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		fmt.Printf("Send Order Error: %v\n", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send order", err)
	}

	// Commit transaction
	tx.Commit()
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Order sent successfully",
		"data":    order,
	})

}
func EditDraftOrder(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	orderId := c.Locals("orderId").(int)
	//fmt.Printf("EditDraftOrder: Order ID: %d\n", orderId)
	input, ok := c.Locals("inputUpdateOrder").(model.InputDraftOrder)
	if !ok {
		fmt.Printf("EditDraftOrder: Failed to parse inputOrderDraft\n")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("failed to parse inputOrderDraft from locals"))
	}
	db := database.DB

	tx := db.Begin()
	var account model.Account
	if err := tx.First(&account, dataInfo.AccountId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}
	var order model.Order
	if err := tx.Preload("OrderItems").First(&order, orderId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found", err)
	}
	if order.Status != "bản nháp" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only draft orders can be edited",
		})
	}
	var customer model.Customer
	if err := tx.First(&customer, *input.CustomerId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Customer not found", err)
	}

	var orderDate time.Time
	if input.OrderDate == "" {
		orderDate = time.Now().In(time.FixedZone("ICT", 7*3600)).Truncate(24 * time.Hour)
	} else {
		parsedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", input.OrderDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid orderDate format: %s", err.Error()),
			})
		}
		orderDate = parsedDate
	}
	orderDatePtr := &orderDate
	// Tạo mã đơn
	today := time.Now().In(time.FixedZone("ICT", 7*3600)).Format("20060102")
	var count int64
	tx.Model(&model.Order{}).Where("order_code LIKE ?", fmt.Sprintf("ORD-%s%%", today)).Count(&count)

	order.OrderDate = orderDatePtr
	order.CustomerId = input.CustomerId
	order.CustomerName = input.CustomerName
	order.PhoneNumber = input.PhoneNumber
	order.Address = input.Address
	order.Discount = input.Discount
	// Delete existing OrderItems
	if err := tx.Where("order_id = ?", order.ID).Delete(&model.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update order items", err)
	}
	var totalAmount float64
	order.OrderItems = nil // Reset slice
	for _, item := range *input.OrderItems {
		if item.Quantity == nil || item.UnitPrice == nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Quantity and UnitPrice cannot be null",
			})
		}
		total := float64(*item.Quantity) * *item.UnitPrice
		totalAmount += total
		order.OrderItems = append(order.OrderItems, model.OrderItem{
			ProductName: item.Name,
			SizeInch:    &item.Size,
			Quantity:    item.Quantity,
			Unit:        &item.Unit,
			UnitPrice:   item.UnitPrice,
			Subtotal:    &totalAmount,
		})
	}
	order.TotalAmount = totalAmount - input.Discount
	if order.TotalAmount < 0 {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Total amount cannot be negative",
		})
	}
	// Save order
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		fmt.Printf("EditDraftOrder Error: %v\n", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update order", err)
	}
	tx.Commit()
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Draft order updated successfully",
		"data":    order,
	})
}
func DeleteDraftOrder(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	orderId := c.Locals("orderId").(int)
	db := database.DB
	tx := db.Begin()
	if tx.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Transaction error", tx.Error)
	}
	// Verify account
	var account model.Account
	if err := tx.First(&account, dataInfo.AccountId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}
	// Find order
	var order model.Order
	if err := tx.First(&order, orderId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found", err)
	}
	// Check authorization
	if order.CreatedById == nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Order has no creator assigned",
		})
	}
	if *order.CreatedById != dataInfo.AccountId {
		tx.Rollback()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fmt.Sprintf("You are not authorized to delete this order (Creator: %d, Your ID: %d)", *order.CreatedById, dataInfo.AccountId),
		})
	}

	// Check draft status
	if order.Status != "bản nháp" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only draft orders can be deleted",
		})
	}

	// Delete OrderItems
	// Soft-delete OrderItems (assuming OrderItem has DeletedAt)
	// if err := tx.Model(&model.OrderItem{}).
	// 	Where("order_id = ? AND deleted_at =?", order.ID, "0001-01-01 06:42:04+06:42:04").
	// 	Update("deleted_at", time.Now()).Error; err != nil {
	// 	tx.Rollback()
	// 	fmt.Printf("DeleteDraftOrder: Failed to soft-delete order items: %v\n", err)
	// 	return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to soft-delete order items", err)
	// }

	// // Soft-delete order
	// if err := tx.Model(&order).
	// 	Updates(map[string]interface{}{
	// 		"deleted_at": time.Now(),
	// 		"status":     "deleted",
	// 	}).Error; err != nil {
	// 	tx.Rollback()
	// 	fmt.Printf("DeleteDraftOrder Error: %v\n", err)
	// 	return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to soft-delete order", err)
	// }
	if err := tx.Where("order_id = ?", order.ID).Delete(&model.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete order items", err)
	}
	if err := tx.Delete(&order).Error; err != nil {
		tx.Rollback()
		fmt.Printf("DeleteDraftOrder Error: %v\n", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete order", err)
	}
	tx.Commit()
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Draft order soft-deleted successfully",
	})
}
func PreviewDraftOrder(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	orderId := c.Locals("orderId").(int)

	db := database.DB
	var order model.Order
	if err := db.Where("deleted_at =?", "0001-01-01 06:42:04+06:42:04").Preload("OrderItems").First(&order, orderId).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found or deleted", err)
	}

	// Check authorization
	if order.CreatedById == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Order has no creator assigned",
		})
	}
	if *order.CreatedById != dataInfo.AccountId {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fmt.Sprintf("You are not authorized to view this order (Creator: %d, Your ID: %d)", *order.CreatedById, dataInfo.AccountId),
		})
	}

	// Check draft status
	if order.Status != "bản nháp" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only draft orders can be previewed",
		})
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Draft order retrieved successfully",
		"data":    order,
	})
}
func DownloadDraftOrder(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	orderId := c.Locals("orderId").(int)

	db := database.DB
	var order model.Order
	if err := db.Where("deleted_at = ?", "0001-01-01 06:42:04+06:42:04").Preload("OrderItems").First(&order, orderId).Error; err != nil {
		fmt.Printf("DownloadDraftOrder: Order not found: %v\n", err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found or deleted", err)
	}

	// Check authorization
	if order.CreatedById == nil {
		fmt.Printf("DownloadDraftOrder: Order has no creator assigned\n")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Order has no creator assigned",
		})
	}
	if *order.CreatedById != dataInfo.AccountId {
		fmt.Printf("DownloadDraftOrder: Unauthorized - Creator: %d, Account: %d\n", *order.CreatedById, dataInfo.AccountId)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fmt.Sprintf("You are not authorized to download this order (Creator: %d, Your ID: %d)", *order.CreatedById, dataInfo.AccountId),
		})
	}

	// Check draft status
	if order.Status != "bản nháp" {
		fmt.Printf("DownloadDraftOrder: Order is not a draft: %s\n", order.Status)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only draft orders can be downloaded",
		})
	}

	// Set headers for JSON download
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=order_%d.json", orderId))
	c.Set("Content-Type", "application/json")

	return c.JSON(fiber.Map{
		"order": order,
	})
}
func CancelDraftOrder(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	orderId := c.Locals("orderId").(int)

	db := database.DB
	tx := db.Begin()
	if tx.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Transaction error", tx.Error)
	}

	// Verify account
	var account model.Account
	if err := tx.First(&account, dataInfo.AccountId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	// Find order (exclude soft-deleted)
	var order model.Order

	if err := tx.Preload("OrderItems").First(&order, orderId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found", err)
	}
	// Check authorization
	if order.CreatedById == nil {
		fmt.Printf("CancelDraftOrder: Order has no creator assigned\n")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Order has no creator assigned",
		})
	}
	if *order.CreatedById != dataInfo.AccountId {
		fmt.Printf("CancelDraftOrder: Unauthorized - Creator: %d, Account: %d\n", *order.CreatedById, dataInfo.AccountId)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fmt.Sprintf("You are not authorized to cancel this order (Creator: %d, Your ID: %d)", *order.CreatedById, dataInfo.AccountId),
		})
	}

	// Check draft status
	if order.Status != "Đã gửi và chưa sản xuất" {
		fmt.Printf("CancelDraftOrder: Order is not a draft: %s\n", order.Status)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only draft orders can be cancelled",
		})
	}

	// Update status to cancelled
	order.Status = "yêu cầu cần hủy"
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		fmt.Printf("CancelDraftOrder Error: %v\n", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to cancel order", err)
	}

	tx.Commit()
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Draft order cancelled successfully",
		"data":    order,
	})
}
func UpdateStatusOrder(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	orderId := c.Locals("orderId").(int)
	//fmt.Printf("EditDraftOrder: Order ID: %d\n", orderId)
	var input model.CancelOrder
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid input: %s", err.Error()),
		})
	}
	db := database.DB

	tx := db.Begin()
	var account model.Account
	if err := tx.First(&account, dataInfo.AccountId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}
	var order model.Order
	if err := tx.Preload("OrderItems").First(&order, orderId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found", err)
	}
	if order.Status != "Đã giao hàng" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Chỉ có thể cập nhật invoice sửa đơn cho đơn hàng ở trạng thái 'Đã giao hàng'",
		})
	}
	validStatuses := []string{"Nhận hàng", "Hủy đơn", "Yêu cầu sửa đơn"}
	if order.Status == "Đã giao hàng" {
		// Validate new status
		isValidStatus := false
		for _, status := range validStatuses {
			if input.Status == status {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Trạng thái không hợp lệ. Phải là một trong: %v", validStatuses),
			})
		}
	}
	//var CancelOrder model.CancelOrder
	if input.Status == "Nhận hàng" {
		currentTime := time.Now()
		order.FinalizedAt = &currentTime
		order.Status = "Hoàn thành" // "Nhận hàng" leads to "Hoàn thành"
	} else if input.Status == "Hủy đơn" {
		if input.CancelReason == nil || *input.CancelReason == "" {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Lý do hủy là bắt buộc khi hủy đơn",
			})
		}
		order.CancelReason = input.CancelReason
		if input.CancelImageURLs != nil {
			order.CancelImageURL = input.CancelImageURLs
		}

	} else if input.Status == "Yêu cầu sửa đơn" {
		// No additional logic here; invoice creation is handled by CreateEditInvoice
		order.Status = input.Status
	}

	order.Status = input.Status
	if err := db.Save(&order).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update order", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order updated successfully",
		"order":   order,
	})
}
func CreateEditInvoice(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	orderId := c.Locals("orderId").(int)
	//fmt.Printf("EditDraftOrder: Order ID: %d\n", orderId)
	input, ok := c.Locals("InputEditInvoice").(model.OrderRevisionInvoice)
	if !ok {
		fmt.Printf("EditDraftOrder: Failed to parse inputOrder\n")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("failed to parse inputOrderDraft from locals"))
	}
	db := database.DB

	tx := db.Begin()
	var account model.Account
	if err := tx.First(&account, dataInfo.AccountId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}
	var order model.Order
	if err := tx.Preload("OrderItems").First(&order, orderId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found", err)
	}
	if order.Status != "Yêu cầu sửa đơn" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Đơn đã hoàn thành không thể sửa",
		})
	}
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format", err)
	}

	if input.Reason == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Reason required", errors.New("Lý do sửa đơn là bắt buộc"))
	}
	if input.Note == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Content required", errors.New("Nội dung sửa là bắt buộc"))
	}

	revisionInvoice := model.OrderRevisionInvoice{
		OrderId:     &order.ID,
		Reason:      input.Reason,
		RequestDate: time.Now().In(time.FixedZone("ICT", 7*60*60)), // Set to current date in ICT timezone
		Note:        input.Note,
	}

	// Create RevisionItems
	for _, item := range input.RevisionItems {
		revisionItem := model.OrderRevisionItem{
			Content:  item.Content,
			ImageURL: item.ImageURL,
		}
		revisionInvoice.RevisionItems = append(revisionInvoice.RevisionItems, revisionItem)
	}

	if err := tx.Create(&revisionInvoice).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create edit invoice", err)
	}

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":         "Edit invoice created successfully",
		"revisionInvoice": revisionInvoice,
	})

}
