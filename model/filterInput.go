package model

type FilterInput struct {
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
	SearchKey string `query:"search"`
	IsActive  *bool  `query:"is_active"` // con trỏ để phân biệt giữa `false` và `chưa truyền`
}
