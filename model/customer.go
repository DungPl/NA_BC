package model

type Customer struct {
	DTO
	Name          string `gorm:"not null" validate:"required" json:"name"`
	Email         string `gorm:"uniqueIndex;not null" validate:"required,email" json:"email"`
	WhatsApp      string `gorm:"uniqueIndex;not null" validate:"required" json:"phone"`
	Address       string `json:"address"`
	AddressOrigin string `json:"addressOrigin"`
	Gender        string `json:"gender"`
	Note          string `json:"note"`

	ManagerId     *uint   ` json:"accountId"` // Tài khoản quản lý (người Sale)
	ManageAccount Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ManagerId" json:"account"`
}
type Customers []Customer
