package handler

import (
	"errors"
	"fmt"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/helper"
	"order-manager/model"
	"order-manager/utils"
	"strconv"
	"strings"
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
	currentDate := time.Now().In(time.FixedZone("ICT", 7*60*60)) // 2025-05-23 00:00:00 +07:00

	var factoryReceiveRevisionAt *time.Time
	if input.RevisionStatus == "Đã nhận hàng cần sửa" {
		if input.FactoryReceiveRevisionAt != nil && *input.FactoryReceiveRevisionAt != "" {
			parsedDate, err := time.ParseInLocation("2006-01-02", *input.FactoryReceiveRevisionAt, time.FixedZone("ICT", 7*60*60))
			if err != nil {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid factory receive revision date format, expected YYYY-MM-DD", err)
			}
			if parsedDate.Before(currentDate) {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factory receive revision date cannot be in the past", errors.New("Ngày xưởng nhận đơn sửa không được ở quá khứ"))
			}
			factoryReceiveRevisionAt = &parsedDate
		} else {
			// Auto-set to current date if no input provided
			factoryReceiveRevisionAt = &currentDate
		}
	} else if input.FactoryReceiveRevisionAt != nil && *input.FactoryReceiveRevisionAt != "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factory receive revision date requires revision status to be 'Đã nhận hàng cần sửa'", errors.New("Ngày xưởng nhận đơn sửa chỉ được cập nhật khi trạng thái là 'Đã nhận hàng cần sửa'"))
	}

	var factoryShipRevisionAt *time.Time
	if input.RevisionStatus == "Đang giao hàng" || input.RevisionStatus == "Đã đóng gói - chờ giao" {
		if input.FactoryShipRevisionAt != nil && *input.FactoryShipRevisionAt != "" {
			parsedDate, err := time.ParseInLocation("2006-01-02", *input.FactoryShipRevisionAt, time.FixedZone("ICT", 7*60*60))
			if err != nil {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid factory ship revision date format, expected YYYY-MM-DD", err)
			}
			if parsedDate.Before(currentDate) {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factory ship revision date cannot be in the past", errors.New("Ngày xưởng giao lại không được ở quá khứ"))
			}
			factoryShipRevisionAt = &parsedDate
		} else if input.RevisionStatus == "Đang giao hàng" {
			// Auto-set to current date if no input provided and status is "Đang giao hàng"
			factoryShipRevisionAt = &currentDate
		}
	} else if input.FactoryShipRevisionAt != nil && *input.FactoryShipRevisionAt != "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factory ship revision date requires revision status to be 'Đang giao hàng' or 'Đã đóng gói - chờ giao'", errors.New("Ngày xưởng giao lại chỉ được cập nhật khi trạng thái sửa đơn là 'Đang giao hàng' hoặc 'Đã đóng gói - chờ giao'"))
	}

	if factoryReceiveRevisionAt != nil {
		revisionInvoice.FactoryReceiveRevisionAt = factoryReceiveRevisionAt
	}
	revisionInvoice.RevisionStatus = input.RevisionStatus
	if input.RevisionProductionStatus != nil && *input.RevisionProductionStatus != "" {
		revisionInvoice.RevisionProductStatus = *input.RevisionProductionStatus
	} else {
		revisionInvoice.RevisionProductStatus = ""
	}
	if factoryReceiveRevisionAt != nil {
		revisionInvoice.FactoryReceiveRevisionAt = factoryReceiveRevisionAt
	}
	if factoryShipRevisionAt != nil {
		revisionInvoice.FactoryRevisionShipAt = factoryShipRevisionAt
	} else if input.RevisionStatus != "Đang giao hàng" && input.RevisionStatus != "Đã đóng gói - chờ giao" {
		// Only clear FactoryShipRevisionAt if status is not eligible
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
func ListOrder(c *fiber.Ctx) error {
	filter, ok := c.Locals("filter").(model.OrderFilter)
	if !ok {
		fmt.Printf("EditDraftOrder: Failed to parse inputOrderDraft\n")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("failed to parse inputOrderDraft from locals"))
	}
	db := database.DB

	query := db.Model(&model.Order{})

	// Apply filters
	if filter.OrderCode != nil && *filter.OrderCode != "" {
		query = query.Where("order_code ILIKE ?", "%"+*filter.OrderCode+"%")
	}
	if filter.CustomerName != nil && *filter.CustomerName != "" {
		query = query.Where("customer_name ILIKE ?", "%"+*filter.CustomerName+"%")
	}
	if filter.CustomerPhone != nil && *filter.CustomerPhone != "" {
		query = query.Where("customer_phone ILIKE ?", "%"+*filter.CustomerPhone+"%")
	}
	if filter.TimeFilter != nil {
		var startDate, endDate time.Time
		now := time.Now()
		timeFilter := *filter.TimeFilter

		if strings.HasPrefix(timeFilter, "month:") {
			monthStr := strings.TrimPrefix(timeFilter, "month:")
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid month format", err)
			}

			startDate = time.Date(now.Year(), time.Month(month), 1, 0, 0, 0, 0, now.Location())

			endDate = startDate.AddDate(0, 1, 0)
			// return c.JSON(fiber.Map{
			// 	"end":   endDate,
			// 	"start": startDate,
			// })
		} else if strings.HasPrefix(timeFilter, "quarter:") {
			qStr := strings.TrimPrefix(timeFilter, "quarter:")
			quarter, err := strconv.Atoi(qStr)
			if err != nil || quarter < 1 || quarter > 4 {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid quarter format", err)
			}
			month := (quarter-1)*3 + 1
			startDate = time.Date(now.Year(), time.Month(month), 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(0, 3, 0)
			// return c.JSON(fiber.Map{
			// 	"end":   endDate,
			// 	"start": startDate,
			// })
		} else if strings.HasPrefix(timeFilter, "year:") {
			yearStr := strings.TrimPrefix(timeFilter, "year:")
			year, err := strconv.Atoi(yearStr)
			if err != nil {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid year format", err)
			}
			startDate = time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
			endDate = time.Date(year+1, 1, 1, 0, 0, 0, 0, now.Location())
			// return c.JSON(fiber.Map{
			// 	"end":   endDate,
			// 	"start": startDate,
			// })
		} else if timeFilter == "lastYear" {
			lastYear := now.Year() - 1
			startDate = time.Date(lastYear, 1, 1, 0, 0, 0, 0, now.Location())
			endDate = time.Date(lastYear+1, 1, 1, 0, 0, 0, 0, now.Location())

		} else if timeFilter == "thisYear" {
			year := now.Year()
			startDate = time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
			endDate = time.Date(year+1, 1, 1, 0, 0, 0, 0, now.Location())

		} else {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid time filter", nil)
		}

		query = query.Where("order_date >= ? AND order_date < ?", startDate, endDate)
	}

	var orders []model.Order
	if err := query.Find(&orders).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch orders", err)
	}

	//var OrderResponse model.OrderResponse
	// Prepare response
	response := make([]model.OrderResponse, len(orders))
	for i, order := range orders {
		response[i] = model.OrderResponse{
			ID:               order.ID,
			OrderCode:        order.OrderCode,
			CustomerName:     order.CustomerName,
			CustomerPhone:    order.PhoneNumber,
			Address:          order.Address,
			OrderDate:        order.OrderDate,
			Status:           order.Status,
			ProductionStatus: order.ProductionStatus,
		}
	}

	// Calculate statistics
	var stats model.OrderStatisticsResponse
	db.Model(&model.Order{}).Count(&stats.TotalOrders)
	db.Model(&model.Order{}).Where("status = ?", "Đã huỷ").Count(&stats.CanceledOrders)
	db.Model(&model.Order{}).Where("production_status = ?", "Đang sản xuất").Count(&stats.ProducingOrders)
	db.Model(&model.Order{}).Where("status IN ?", []string{"Đang giao hàng", "Đã đóng gói"}).Count(&stats.ShippedOrders)

	return c.JSON(fiber.Map{
		"orders":     response,
		"statistics": stats,
	})
}
