package handler

import (
	"errors"
	"fmt"
	"order-manager/constants"
	"order-manager/helper"
	"order-manager/model"
	"order-manager/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	loginInput := new(LoginInput)

	if err := c.BodyParser(loginInput); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.MISSING_LOGIN_INPUT, err)
	}

	// Manual validation
	if loginInput.UserName == "" || loginInput.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.MISSING_LOGIN_INPUT, errors.New("username and password are required"))
	}

	username := loginInput.UserName
	password := loginInput.Password
	accountModel, err := new(model.Account), *new(error)

	accountModel, err = helper.GetUserByUsername(username)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_INTERNAL_ERROR, err)
	}
	if accountModel == nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, constants.INVALID_USERNAME, errors.New("username not exists"))
	}

	if !helper.CheckPasswordHash(password, accountModel.Password) {
		return utils.ErrorResponse(c, fiber.StatusNotFound, constants.INVALID_PASSWORD, errors.New("password does not match username"))
	}

	if !accountModel.Active {
		return utils.ErrorResponse(c, fiber.StatusForbidden, constants.ACCOUNT_NOT_ACTIVE, errors.New("active false"))
	}

	tokenClaim := model.TokenClaim{
		AccountId: accountModel.ID,
		Username:  accountModel.Username,
	}
	token, err := helper.GenerateAccessToken(tokenClaim)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_INTERNAL_ERROR, err)
	}

	refreshToken, err := helper.GenerateRefreshToken(tokenClaim)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_INTERNAL_ERROR, err)
	}

	tokenData := model.TokenData{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}

	return utils.SuccessResponse(c, fiber.StatusOK, tokenData)
}

func RefreshToken(c *fiber.Ctx) error {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refreshToken"`
	}

	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Xác thực refresh token
	token, err := helper.ParseToken(req.RefreshToken)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	var tokenClaim model.TokenClaim

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Xác thực thành công, truy xuất thông tin từ payload
		accountId := claims["accountId"].(float64)
		username := claims["username"].(string)

		tokenClaim = model.TokenClaim{
			AccountId: uint(accountId),
			Username:  username,
		}
	} else {
		// Xác thực thất bại
		fmt.Println("Invalid Token:", err)
		return err
	}

	// Tạo mới access token và refresh token
	newAccessToken, err := helper.GenerateAccessToken(tokenClaim)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate access token"})
	}

	newRefreshToken, err := helper.GenerateRefreshToken(tokenClaim)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token"})
	}

	tokenData := model.TokenData{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return utils.SuccessResponse(c, fiber.StatusOK, tokenData)
}
