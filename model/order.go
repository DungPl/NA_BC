package model

import "time"

type Order struct {
	DTO
	OrderCode    string     `gorm:"type:varchar(50);unique;not null" json:"orderCode"` // Mã đơn
	OrderDate    *time.Time `gorm:"type:timestamp" validate:"required" json:"orderDate"`
	CustomerId   *uint      `gorm:"" json:"customerId"`
	Customer     *Customer  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CustomerId" json:"customer"`
	CustomerName string     `gorm:"type:varchar(100);not null" json:"customerName"` // cho phép sửa
	Address      string     `json:"address"`
	PhoneNumber  string     `gorm:"not null" json:"phoneNumber"`

	CreatedById *uint    `json:"accountId"` // tài khoản Sale
	CreatedBy   *Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CreatedById" json:"account"`

	Discount    float64 `gorm:"default:0"` // USD
	TotalAmount float64 `gorm:"default:0"` // tính từ hàng hóa

	Status           string     `json:"orderUpdateStatus"` // Đã sản xuất, Đang giao hàng, Nhận hàng, Hoàn thành, Yêu cầu sửa, Hủy đơn
	CancelReason     *string    `gorm:"type:text" json:"cancelReason"`
	CancelImageURL   *string    `gorm:"type:text" json:"cancelImageUrl"`
	FinalizedAt      *time.Time // Ngày hoàn thành (set khi nhận hàng hoặc hoàn thành)
	FactoryReceiveAt *time.Time `gorm:"type:timestamp" validate:"required" json:"factoryReceiveAt"` //Ngày xưởng nhận đơn

	EstimatedShipAt *time.Time `gorm:"type:timestamp" validate:"required" json:"estimatedShipAt"` //Ngày dự kiến xuất xưởng
	ActualShipAt    *time.Time `gorm:"type:timestamp" json:"actualShipAt"`                        //Ngày thực tế xuất xưởng

	ProductionStatus *string `json:"prodStatus"` // Trạng thái sản xuất
	OrderItems       []OrderItem
}
