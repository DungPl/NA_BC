package model

import "time"

type OrderRevisionInvoice struct {
	DTO
	OrderId *uint  `gorm:"" json:"orderId"`
	Order   *Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:OrderId" json:"order"`

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
type RevisionInvoiceResponse struct {
	ID                       uint      `json:"id"`
	OrderCode                string    `json:"orderCode"`
	CustomerName             string    `json:"customerName"`
	CustomerPhone            string    `json:"customerPhone"`
	RequestDate              time.Time `json:"requestDate"`
	Reason                   string    `json:"reason"`
	RevisionStatus           *string   `json:"revisionStatus"`
	RevisionProductionStatus *string   `json:"revisionProductionStatus"`
}
type StatisticsResponse struct {
	TotalOrders     int `json:"totalOrders"`
	CanceledOrders  int `json:"canceledOrders"`
	ProducingOrders int `json:"producingOrders"`
	ShippedOrders   int `json:"shippedOrders"`
}
