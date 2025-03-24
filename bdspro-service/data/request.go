package data

type NewProductRequest struct {
	Name        string  `gorm:"size:255" json:"name"`
	Code        string  `gorm:"size:15" json:"code"`
	Area        float64 `json:"area"`
	Price       int64   `json:"price"`
	Description string  `json:"description"`
	Note        string  `json:"note"`
}

type ProductSearchRequest struct {
	Name        string  `gorm:"size:255" json:"name"`
	Code        string  `gorm:"size:15" json:"code"`
	Area        float64 `json:"area"`
	Price       int64   `json:"price"`
	Description string  `json:"description"`
	Note        string  `json:"note"`
}

type ProductSaveRequest struct {
	Name        string  `gorm:"size:255" json:"name"`
	Code        string  `gorm:"size:15" json:"code"`
	Area        float64 `json:"area"`
	Price       int64   `json:"price"`
	Description string  `json:"description"`
	Note        string  `json:"note"`
}
