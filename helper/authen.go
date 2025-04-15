package helper

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"order-manager/config"
	"order-manager/constants"
	"order-manager/database"
	"order-manager/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Khóa bí mật để ký token
var jwtSecret = []byte(config.Config("SECRET"))

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	log.Println(hash, "super_password_hash")
	return err == nil
}

func GetUserByUsername(u string) (*model.Account, error) {
	db := database.DB
	var account model.Account
	if err := db.Where(&model.Account{Username: u}).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func Valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GenerateAccessToken(tokenClaim model.TokenClaim) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = tokenClaim.Username
	claims["accountId"] = tokenClaim.AccountId
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	t, err := token.SignedString([]byte(jwtSecret))
	return t, err
}

func GenerateRefreshToken(tokenClaim model.TokenClaim) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = tokenClaim.Username
	claims["accountId"] = tokenClaim.AccountId
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	t, err := token.SignedString([]byte(jwtSecret))
	return t, err
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Xác thực thuật toán ký là HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	return token, err
}

func GetInfoAccountFromToken(c *fiber.Ctx) (model.TokenClaim, bool, bool, bool) {
	token := c.Locals("user").(*jwt.Token)
	tokenClaim := token.Claims.(jwt.MapClaims)
	accountId := uint(tokenClaim["accountId"].(float64))
	username := tokenClaim["username"].(string)
	accountInfo := model.TokenClaim{
		AccountId: uint(accountId),
		Username:  username,
	}
	var account model.Account
	db := database.DB
	db.Preload("Role").First(&account, accountId)

	return accountInfo, account.Role == constants.ROLE_ADMIN || account.Role == constants.ROLE_QUANLY, account.Role == constants.ROLE_KETOAN, account.Role == constants.ROLE_SALE
}
