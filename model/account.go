package model

type Account struct {
	DTO
	Username     string `gorm:"uniqueIndex;not null" validate:"required,min=3,max=50" json:"username"`
	Password     string `gorm:"not null" validate:"required,min=6,max=50" json:"password"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Active       bool   `gorm:"not null;default:true" json:"active"`
	Role         string `json:"role"`
	Staff        *Staff `gorm:"foreignKey:AccountId" json:"staff"`
}

type Accounts []Account
