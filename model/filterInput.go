package model

type FilterInput struct {
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
	SearchKey string `query:"search"`
	IsActive  *bool  `query:"is_active"` // con trỏ để phân biệt giữa `false` và `chưa truyền`
}
type RevisionInvoiceFilter struct {
	OrderCode     *string `query:"orderCode"`
	CustomerName  *string `query:"customerName"`
	CustomerPhone *string `query:"customerPhone"`
	TimeFilter    *string `query:"timeFilter"` // Format: "month:1", "quarter:1", "year:2025", "lastYear"
}