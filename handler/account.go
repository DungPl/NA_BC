package handler

import (
	"errors"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/helper"
	"order-manager/model"
	"order-manager/utils"
	"strconv"

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

func AdminChangePassword(c *fiber.Ctx) error {
	staffIdParam := c.Params("staffId")
	staffId, err := strconv.Atoi(staffIdParam)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid staff ID", err)
	}
	db := database.DB
	changePasswordInput, ok := c.Locals("AdminChangePasswordInput").(model.ChangePasswordInput)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_PARSE_DATA_TO_LOCALS, errors.New("PARSE DATA TO LOCALS FAIL"))
	}
	//dataInfo, _, _, _ := helper.GetInfoAccountFromToken(c)
	var account model.Account
	if err := db.First(&account, staffId).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Account not found", err)
	}

	newPasswordHash, err := helper.HashPassword(changePasswordInput.NewPassword)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.CAN_NOT_HASH_PASSWORD, err)
	}
	account.Password = newPasswordHash
	db.Save(&account)

	return utils.SuccessResponse(c, fiber.StatusOK, account)

}
