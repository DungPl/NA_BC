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

func AdminEditOrder(c *fiber.Ctx) error {
	// Get order ID from URL parameter
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	orderId := c.Locals("orderId").(int)
	input, ok := c.Locals("inputUpdateOrder").(model.InputDraftOrder)
	if !ok {
		fmt.Printf("EditDraftOrder: Failed to parse inputOrderDraft\n")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("failed to parse inputOrderDraft from locals"))
	}
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

	var order model.Order
	if err := tx.Preload("OrderItems").First(&order, orderId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Order not found", err)
	}

	if order.Status == "Hoàn thành" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Không thể chỉnh sửa đơn hàng đã hoàn thành",
		})
	}
	if input.Status == "Đã sản xuất" {
		if order.Status != "Đang sản xuất" || order.ProductionStatus == nil || *order.ProductionStatus != "Đang xử lý mềm mượt" {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Chỉ có thể cập nhật trạng thái 'Đã sản xuất' khi đơn hàng đang ở trạng thái 'Đang sản xuất' và trạng thái sản xuất là 'Đang xử lý mềm mượt'",
			})
		}
	}
	if input.Status == "Đã giao hàng" {
		if order.Status != "Đã sản xuất" && order.ProductionStatus != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Chỉ có thể cập nhật trạng thái 'Đã giao hàng' khi đơn hàng đang ở trạng thái 'Đã sản xuất'",
			})
		}
	}
	// Update FactoryReceiveAt (default to current date if not provided)
	if input.FactoryReceiveAt != nil && *input.FactoryReceiveAt != "" {
		parsedDate, err := time.Parse("2006-01-02", *input.FactoryReceiveAt)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Định dạng ngày xưởng nhận đơn không hợp lệ",
			})
		}
		order.FactoryReceiveAt = &parsedDate
	} else {
		currentTime := time.Now()
		order.FactoryReceiveAt = &currentTime
	}

	// Update EstimatedShipAt
	if input.ActualShipAt != nil && *input.ActualShipAt != "" {
		parsedDate, err := time.Parse("2006-01-02", *input.ActualShipAt)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Định dạng ngày xưởng nhận đơn không hợp lệ",
			})
		}
		order.ActualShipAt = &parsedDate
	} else {
		currentTime := time.Now()
		order.ActualShipAt = &currentTime
	}

	// Update ActualShipAt (only allowed when Status is "Đang giao hàng")
	if input.Status == "Đã giao hàng" {
		if input.ActualShipAt != nil && *input.ActualShipAt != "" {
			parsedDate, err := time.Parse("2006-01-02", *input.ActualShipAt)
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Định dạng ngày giao thực tế không hợp lệ",
				})
			}
			order.ActualShipAt = &parsedDate
		} else {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ngày giao thực tế là bắt buộc khi trạng thái là 'Đã giao'",
			})
		}
	} else {
		order.ActualShipAt = nil // Reset ActualShipAt if status is not "Đã giao"
	}

	// Update Status
	order.Status = input.Status

	// Update ProductionStatus (only allowed when Status is "Đang sản xuất")
	validProductionStatuses := []string{
		"Đang chia hàng",
		"Đã gửi lace",
		"Đang làm màu",
		"Đang tẩy màu",
		"Đang xử lý mềm mượt",
	}
	if input.Status == "Đang sản xuất" && input.ProductionStatus != nil {
		isValid := false
		for _, status := range validProductionStatuses {
			if *input.ProductionStatus == status {
				isValid = true
				break
			}
		}
		if !isValid {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Trạng thái sản xuất không hợp lệ. Phải là một trong: %v", validProductionStatuses),
			})
		}
		order.ProductionStatus = input.ProductionStatus
	} else {
		order.ProductionStatus = nil // Reset ProductionStatus if Status is not "Đang sản xuất"
	}
	if err := db.Save(&order).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update order", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order updated successfully",
		"order":   order,
	})
}
func UpdateRevisionStatus(c *fiber.Ctx) error {
	// Get account info from token
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	revisionInvoiceId := c.Locals("revisionInvoiceId").(int)

	// Parse input
	input, ok := c.Locals("InputUpdateRevisionStatus").(model.InputRevisionStatus)
	if !ok {
		fmt.Printf("EditDraftOrder: Failed to parse inputOrder\n")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("failed to parse inputOrderDraft from locals"))
	}
	// Validate input

	validRevisionStatuses := []string{
		"Đã nhận hàng cần sửa",
		"Đang sản xuất",
		"Đã đóng gói - chờ giao",
		"Đang giao hàng",
	}
	isValidRevisionStatus := false
	for _, status := range validRevisionStatuses {
		if input.RevisionStatus == status {
			isValidRevisionStatus = true
			break
		}
	}
	if !isValidRevisionStatus {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid revision status", fmt.Errorf("Trạng thái sửa đơn phải là một trong: %v", validRevisionStatuses))
	}

	if input.RevisionProductionStatus != nil && *input.RevisionProductionStatus != "" {
		if input.RevisionStatus != "Đang sản xuất" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Revision production status is only allowed when revision status is 'Đang sản xuất'", errors.New("Trạng thái sản xuất chỉ được cập nhật khi trạng thái sửa đơn là 'Đang sản xuất'"))
		}
		if input.RevisionProductionStatus != nil {
			validProductionStatuses := []string{
				"Đang chia hàng",
				"Đã gửi lace",
				"Đang làm màu",
				"Đang tẩy màu",
				"Đang xử lý mềm mượt",
			}
			isValidProductionStatus := false
			for _, status := range validProductionStatuses {
				if *input.RevisionProductionStatus == status {
					isValidProductionStatus = true
					break
				}
			}
			if !isValidProductionStatus {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid revision production status", fmt.Errorf("Trạng thái sản xuất phải là một trong: %v", validProductionStatuses))
			}
		}
	}

	if input.FactoryShipRevisionAt != nil && input.RevisionStatus != "Đang giao hàng" && input.RevisionStatus != "Đã đóng gói - chờ giao" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factory ship revision date is only allowed when revision status is 'Đang giao hàng' or 'Đã đóng gói - chờ giao'", errors.New("Ngày xưởng giao lại chỉ được cập nhật khi trạng thái sửa đơn là 'Đang giao hàng' hoặc 'Đã đóng gói - chờ giao'"))
	}
	if input.FactoryShipRevisionAt != nil && *input.FactoryShipRevisionAt != "" {
		if input.RevisionStatus != "Đang giao hàng" && input.RevisionStatus != "Đã đóng gói - chờ giao" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factory ship revision date requires revision status to be 'Đang giao hàng' or 'Đã đóng gói - chờ giao'", errors.New("Ngày xưởng giao lại chỉ được cập nhật khi trạng thái sửa đơn là 'Đang giao hàng' hoặc 'Đã đóng gói - chờ giao'"))
		}
	}
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

	// Verify revision invoice exists
	var revisionInvoice model.OrderRevisionInvoice
	if err := tx.Preload("Order").Preload("RevisionItems").First(&revisionInvoice, revisionInvoiceId).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Revision invoice not found", err)
	}

	// Update fields
	// Update revision status
	revisionInvoice.RevisionStatus = input.RevisionStatus

	// Update production status nếu có
	if input.RevisionProductionStatus != nil && *input.RevisionProductionStatus != "" {
		revisionInvoice.RevisionProductStatus = *input.RevisionProductionStatus
	} else {
		revisionInvoice.RevisionProductStatus = ""
	}

	// Cập nhật FactoryRevisionShipAt chỉ khi trạng thái đúng
	if input.FactoryShipRevisionAt != nil && *input.FactoryShipRevisionAt != "" {
		if input.RevisionStatus == "Đang giao hàng" || input.RevisionStatus == "Đã đóng gói - chờ giao" {
			parsedTime, err := time.Parse(time.RFC3339, *input.FactoryShipRevisionAt)
			if err != nil {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid FactoryShipRevisionAt format", err)
			}
			revisionInvoice.FactoryRevisionShipAt = &parsedTime
		} else {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Không thể cập nhật ngày giao nếu trạng thái không phải 'Đang giao hàng' hoặc 'Đã đóng gói - chờ giao'", errors.New("Trạng thái đơn không hợp lệ để cập nhật ngày giao"))
		}
	} else {
		revisionInvoice.FactoryRevisionShipAt = nil
	}
	// Save revision invoice
	if err := tx.Save(&revisionInvoice).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update revision invoice", err)
	}

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":         "Revision status updated successfully",
		"revisionInvoice": revisionInvoice,
	})
}
