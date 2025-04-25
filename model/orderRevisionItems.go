package model

type OrderRevisionItem struct {
	DTO
	OrderRevisionInvoiceId *uint                 `gorm:"not null" json:"orderRevisonInvoiceId"`
	OrderRevisionInvoice   *OrderRevisionInvoice `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:OrderRevisionInvoiceId" json:"orderRevisionInvoice"`

	Content  string `gorm:"type:text;not null" json:"content"`
	ImageURL string `gorm:"type:text" json:"imageUrl"` // Ảnh mô tả nội dung sửa
}
type OrderRevisionItems []OrderRevisionItem
