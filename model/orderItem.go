package model

type OrderItem struct {
	DTO
	OrderId *uint  ` json:"oderId"`
	Order   *Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:OrderId" json:"order"`

	ProductName string   `gorm:"type:varchar(255);not null" json:"productName"`
	SizeInch    *float64 ` json:"sizeInch"`                     // kích thước (inch)
	Quantity    *int     ` json:"quantity"`                     // số lượng
	Unit        *string  `gorm:"type:varchar(20);" json:"unit"` // đơn vị: wig, bundles, pieces, kg
	UnitPrice   *float64 ` json:"unitPrice"`                    // USD
	Subtotal    *float64 ` json:"subtotal"`                     // = Quantity * UnitPrice
}
