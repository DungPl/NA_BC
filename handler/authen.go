package handler

import (
	"errors"
	"fmt"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/helper"

	"order-manager/model"
	"order-manager/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	// kiểm tra tài khoản
	// 1. kiểm tra nhập toàn khoản
	if loginInput.UserName == "" || loginInput.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, constants.MISSING_LOGIN_INPUT, errors.New("username and password are required"))
	}

	username := loginInput.UserName
	password := loginInput.Password
	accountModel, err := new(model.Account), *new(error)
	// Kiểm tra tài  khoản trong cơ sơr dữ liệu

	accountModel, err = helper.GetUserByUsername(username) // Tìm kiếm theo username

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_INTERNAL_ERROR, err)
	}
	if accountModel == nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, constants.INVALID_USERNAME, errors.New("username not exists"))
	}

	if !helper.CheckPasswordHash(password, accountModel.Password) {
		return utils.ErrorResponse(c, fiber.StatusNotFound, constants.INVALID_PASSWORD, errors.New("password does not match username"))
	}
	// kiểm tra tài khoản đã được kích hoạt chưa
	if !accountModel.Active {
		return utils.ErrorResponse(c, fiber.StatusForbidden, constants.ACCOUNT_NOT_ACTIVE, errors.New("active false"))
	}
	// tạo token nhúng vào JWT
	tokenClaim := model.TokenClaim{
		AccountId: accountModel.ID,
		Username:  accountModel.Username,
	}
	token, err := helper.GenerateAccessToken(tokenClaim) // goij hàm tự động tạo  token
	//
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, constants.ERROR_INTERNAL_ERROR, err)
	}

	refreshToken, err := helper.GenerateRefreshToken(tokenClaim) // gọi hàm tự động tạo token mới khi hết hạn mà không cần login lại

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
	// Nhận refreshtoken từ body
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

	// Tạo mới access token và refresh token kiểm soát phiên đăng nhập
	// Giúp ngăn chặn các token cũ đã hết phiên hạn
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

func ChangePassword(c *fiber.Ctx) error {
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
	//return utils.SuccessResponse(c, fiber.StatusOK, token)
	// Lấy thông tin tài khoản từ token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	accountIdFloat, ok := claims["accountId"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid accountId in token"})
	}
	accountId := uint(accountIdFloat)
	//return utils.SuccessResponse(c, fiber.StatusOK, accountId)

	// Lấy dữ liệu đầu vào (mật khẩu cũ và mật khẩu mới)
	type ChangePasswordInput struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	var input ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if len(input.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "New password too short"})
	}

	// 5. Lấy account từ DB
	var account model.Account
	if err := database.DB.First(&account, accountId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Account not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	var staff model.Staff
	if staff.Position != constants.ROLE_ADMIN {
		if input.OldPassword == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Old password is required"})
		}
		if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.OldPassword)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Old password is incorrect"})
		}
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	account.Password = string(hashed)
	if err := database.DB.Save(&account).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
	}

	// Trả về phản hồi thành công
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Password changed successfully"})
}
