package model

import (
	"time"
)

type OrderRevisionInvoice struct {
	DTO
	OrderId                  *uint               `gorm:"" json:"orderId"`
	Order                    *Order              `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:OrderId" json:"order"`
	RevisionInvoiceCode      string              `gorm:"type:text;not null" json:"revicionInvoiceCode"`
	Reason                   string              `gorm:"type:text;not null" json:"reason"`
	RequestDate              time.Time           `gorm:"autoCreateTime" json:"requestDate"` // ngày hiện tại, không cho sửa
	Note                     string              `gorm:"type:text" json:"note"`
	FactoryReceiveRevisionAt *time.Time          ` json:"factoryReceiveRevisionAt"`                      // ngày nhận lại
	RevisionStatus           string              `gorm:"type:varchar(255)" json:"revisionStatus"`        // trạng  thái Đã nhận hàng cần sửa, Đang sản xuất, Đã đóng gói - chờ giao, Đang giao hàng
	RevisionProductStatus    string              `gorm:"type:varchar(255)" json:"revisionProductStatus"` // trạng thái sản phẩm sửa lại , chỉ hiển thị khi revision_status = 'Đang sản xuất'
	FactoryRevisionShipAt    *time.Time          ` json:"factoryRevisionShipAt"`                         //  Lưu ngày xưởng giao lại, chỉ cập nhật khi revision_status là Đang giao hàng hoặc Đã đóng gói - chờ giao.
	RevisionItems            []OrderRevisionItem `json:"revisionItems"`
}
type OrderRevisionInvoices []OrderRevisionInvoice
type InputRevisionStatus struct {
	FactoryReceiveRevisionAt *string `json:"factoryReceiveRevisionAt"`
	RevisionStatus           string  `json:"revisionStatus" validate:"required"`
	RevisionProductionStatus *string `json:"revisionProductionStatus"`
	FactoryShipRevisionAt    *string `json:"factoryShipRevisionAt"`
}

type RevisionHistory struct {
	DTO
	OrderId   uint                    `gorm:"not null;foreignKey:OrderId;references:orders(id);constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"orderId"`
	AccountId uint                    `gorm:"not null;foreignKey:AccountId;references:accounts(id);constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"accountId"`
	Action    string                  `gorm:"type:varchar(50);not null" json:"action"`
	Details   []RevisionHistoryDetail ` json:"details"`
}
type RevisionHistories []RevisionHistory
type RevisionHistoryDetail struct {
	DTO
	RevisionHistoryId        uint       `gorm:"foreignKey:RevisionHistoryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"revisionHistoryId"`
	Note                     string     `gorm:"type:text" json:"note"`
	RevisionStatus           *string    `json:"revisionStatus"`
	RevisionProductionStatus *string    `json:"revisionProductionStatus"`
	FactoryReceiveRevisionAt *time.Time `json:"factoryReceiveRevisionAt"`
	FactoryShipRevisionAt    *time.Time `json:"factoryShipRevisionAt"`
}
type RevisionHistoryDetails []RevisionHistoryDetail
type RevisionInvoiceResponse struct {
	RevisionCode             string              `json:"revisionCode"`
	OrderCode                string              `json:"orderCode"`
	Note                     string              `json:"note"`
	CreatedAt                time.Time           `json:"createdAt"`
	RevisionHistory          []RevisionHistory   `json:"revisionHistory"`
	RevisionStatus           string              `json:"revisionStatus"`
	FactoryReceiveRevisionAt *time.Time          `json:"factoryReceiveRevisionAt"`
	FactoryRevisionShipAt    *time.Time          `json:"factoryRevisionShipAt"`
	Image                    []OrderRevisionItem `json:"images"`
}
