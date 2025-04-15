package database

import (
	"log"
	"order-manager/constants"
	"order-manager/model"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	bytes, err := bcrypt.GenerateFromPassword([]byte("123456cn"), 10)
	HashPassword := string(bytes)
	if err != nil {
		HashPassword = "123456cn"
	}
	accounts := []model.Account{
		{Username: "Administration", Password: HashPassword, Active: true},
	}

	for _, account := range accounts {
		// Tạo mới nếu không tồn tại
		if err := db.Where(model.Account{Username: account.Username}).FirstOrCreate(&account).Error; err != nil {
			log.Println("failed to seed data for account:", account.Username, "error:", err)
		}
	}

	var account model.Account
	db.Where(model.Account{Username: "Administration"}).First(&account)

	birthDay, _ := time.Parse("2006-01-02", "1994-04-17")
	staffs := []model.Staff{
		{Name: "Admin", BirthDay: birthDay, PhoneNumber: "0969013457", IdentificationCard: "027094000624", AccountId: &account.ID, Role: constants.ROLE_ADMIN},
	}

	for _, staff := range staffs {
		// Tạo mới nếu không tồn tại
		if err := db.Where(model.Staff{IdentificationCard: staff.IdentificationCard}).FirstOrCreate(&staff).Error; err != nil {
			db.Where(model.Staff{IdentificationCard: staff.IdentificationCard}).Updates(&staff)
			log.Println("failed to seed data for staff:", staff.IdentificationCard, "error:", err)
		}
	}
}
