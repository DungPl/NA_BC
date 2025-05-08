package model

import (
	"time"
)

type Staff struct {
	DTO
	Name               string    `gorm:"not null" validate:"required" json:"name"`
	BirthDay           time.Time `gorm:"not null;type:timestamp" validate:"required" json:"birthDay"`
	Address            string    `json:"address"`
	AddressOrigin      string    `json:"addressOrigin"`
	PhoneNumber        string    `gorm:"not null" json:"phoneNumber"`
	Email              string    `json:"email"`
	Gender             string    `json:"gender"`
	IsActive           bool      `gorm:"not null;default:true" json:"isActive"`
	IdentificationCard string    `gorm:"not null;uniqueIndex;require" validate:"required,min=12,max=12" json:"identificationCard"`
	Position           string    `json:"position"`
	StatusWorking      string    `json:"statusWorking"`
	Note               string    `json:"note"`
	AccountId          *uint     `json:"accountId"`
	Account            Account   `gorm:"references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"account"`
}
type CreateStaffInput struct {
	Name               string `json:"name" validate:"required"`
	PhoneNumber        string `json:"phoneNumber" gorm:"uniqueIndex;not null" validate:"required"`
	Email              string `json:"email" gorm:"uniqueIndex;not null" `
	IdentificationCard string `gorm:"not null;uniqueIndex;require" validate:"required,min=12,max=12" json:"identificationCard"`
	Password           string `json:"password" validate:"required,min=6,max=50"`
	Position           string `json:"position"`
	StatusWorking      string `json:"statusWorking"`
	Gender             string `json:"gender"`
	SearchKeyStr       string `json:"searchKey"`
	Act                bool   `json:"active"`
}
type CreateStaffInputs []CreateStaffInput

func (s CreateStaffInput) SearchKey() string {
	return s.SearchKeyStr

}
func (s CreateStaffInput) Active() bool {
	return s.Act

}

type UpdateStaffInput struct {
	Name               string    `json:"name" validate:"required"`
	PhoneNumber        string    `json:"phoneNumber"`
	Email              string    `json:"email"`
	IdentificationCard string    `json:"identificationCard" validate:"min=12,max=12"`
	Address            string    `json:"address"`
	AddressOrigin      string    `json:"addressOrigin"`
	Gender             string    `json:"gender"`
	Role               string    `json:"role"`
	StatusWorking      string    `json:"statusWorking"`
	Note               string    `json:"note"`
	BirthDay           time.Time `json:"birthDay"`
	Username           string    `json:"username"` // cho account
	Position           string    `json:"position"` // role của account
}
type UpdateStaffInputs []UpdateStaffInput

type ChangePasswordInput struct {
	OldPassword     string `gorm:"not null" json:"old_password" `
	NewPassword     string `gorm:"not null" json:"new_password" validate:"required"`
	ConfirmPassword string `gorm:"not null" json:"confirm_password" validate:"required"`
}
type ChangePasswords []ChangePasswordInput
type Staffs []Staff
