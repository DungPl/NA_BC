package handler

import (
	"order-manager/constants"
	"order-manager/database"
	"order-manager/helper"
	"order-manager/model"
	"order-manager/utils"

	"github.com/gofiber/fiber/v2"
)

func Me(c *fiber.Ctx) error {
	dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	accountId := dataInfo.AccountId

	db := database.DB
	var account model.Account
	if err := db.Preload("Staff").First(&account, accountId).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, constants.NOT_FOUND_RECORDS, err)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, account)
}
