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

	Status string `json:"orderUpdateStatus"` // Đã sản xuất, Đang giao hàng, Nhận hàng, Hoàn thành, Yêu cầu sửa, Hủy đơn

	CancelReason     *string    `gorm:"type:text" json:"cancelReason"`
	CancelImageURL   *string    `gorm:"type:text" json:"cancelImageUrl"`
	FinalizedAt      *time.Time // Ngày hoàn thành (set khi nhận hàng hoặc hoàn thành)
	FactoryReceiveAt *time.Time `gorm:"type:timestamp" validate:"required" json:"factoryReceiveAt"` //Ngày xưởng nhận đơn

	EstimatedShipAt *time.Time `gorm:"type:timestamp" validate:"required" json:"estimatedShipAt"` //Ngày dự kiến xuất xưởng
	ActualShipAt    *time.Time `gorm:"type:timestamp" json:"actualShipAt"`                        //Ngày thực tế xuất xưởng

	ProductionStatus *string `json:"prodStatus"` //Đang sản xuất: Đang chia hàng, Đã gửi lace, Đang làm màu, Đang tẩy màu, Đang xử lý mềm mượt
	OrderItems       []OrderItem
}
type Orders []Order
type InputDraftOrder struct {
	OrderCode        string            `gorm:"type:varchar(50);unique;not null" json:"orderCode"` // Mã đơn
	OrderDate        string            `json:" orderDate" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00" `
	CustomerId       *uint             `json:"customerId" validate:"required"`
	CustomerName     string            `json:"customerName" validate:"required"`
	PhoneNumber      string            `json:"phoneNumber" validate:"required"`
	Address          string            `json:"address"`
	Discount         float64           `json:"discount" validate:"gte=0"`
	FactoryReceiveAt *string           `gorm:"type:timestamp"  json:"factoryReceiveAt"` //Ngày xưởng nhận đơn
	Status           string            `json:"orderUpdateStatus"`
	ProductionStatus *string           `json:"prodStatus"`
	EstimatedShipAt  *string           `gorm:"type:timestamp"  json:"estimatedShipAt"` //Ngày dự kiến xuất xưởng
	ActualShipAt     *string           `gorm:"type:timestamp" json:"actualShipAt"`     //Ngày thực tế xuất xưởng
	OrderItems       *[]OrderItemInput `json:"orderItems" validate:"required,dive"`
}

type OrderItemInput struct {
	Name      string   `json:"name" validate:"required"`
	Size      float64  `json:"size" validate:"gte=0"`
	Quantity  *int     `json:"quantity" validate:"required,gt=0"`
	Unit      string   `json:"unit" validate:"required,oneof=wig bundles pieces kg"`
	UnitPrice *float64 `json:"unitPrice" validate:"required,gt=0"`
}
type InputDraftOrders []InputDraftOrder
type InputEditInvoice struct {
	Reason  string   `json:"reason"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
	Notes   string   `json:"notes"`
}

type EditInvoice struct {
	ID          int       `gorm:"primaryKey"`
	OrderID     int       `gorm:"index"`
	Reason      string    `json:"reason"`
	RequestDate time.Time `json:"requestDate"`
	Content     string    `json:"content"`
	Images      string    `json:"images" gorm:"type:json"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type CancelOrder struct {
	Status          string  `json:"orderUpdateStatus" validate:"required"`
	CancelReason    *string `json:"reason"`
	CancelImageURLs *string `json:"images"` // Mảng URL hình ảnh
}
