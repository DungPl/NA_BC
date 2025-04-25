package model

import "time"

type OrderRevisionInvoice struct {
	DTO
	OrderId *uint  `gorm:"" json:"orderId"`
	Order   *Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:OrderId" json:"order"`

	Reason      string    `gorm:"type:text;not null" json:"reason"`
	RequestDate time.Time `gorm:"autoCreateTime" json:"requestDate"` // ngày hiện tại, không cho sửa
	Note        string    `gorm:"type:text" json:"note"`

	RevisionItems []OrderRevisionItem
}
type OrderRevisionInvoices []OrderRevisionInvoice
