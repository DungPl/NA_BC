package model

import "time"

type TokenData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenClaim struct {
	AccountId uint   `json:"accountId"`
	Username  string `json:"username"`
}

type DTO struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
}

type Pagination struct {
	Limit *int `json:"limit"`
	Page  *int `json:"page"`
}
