package domain

import (
	_entity "common/domain/entity"
	"fmt"
	"time"
)

// DirectorySource represents a directory source entity
type DirectorySource struct {
	_entity.BaseEntity
	Name           string     `json:"name" gorm:"column:name;not null"`
	Code           string     `json:"code" gorm:"column:code;not null;unique"`
	Description    string     `json:"description" gorm:"column:description"`
	CategoryID     *uint64    `json:"categoryId" gorm:"column:category_id"`
	Type           string     `json:"type" gorm:"column:type;not null"`
	Icon           string     `json:"icon" gorm:"column:icon"`
	Color          string     `json:"color" gorm:"column:color"`
	IsActive       bool       `json:"isActive" gorm:"column:is_active;default:true"`
	SortOrder      uint32     `json:"sortOrder" gorm:"column:sort_order;default:0"`
	ExpectedAmount int64      `json:"expectedAmount" gorm:"column:expected_amount;default:0"`
	ActualAmount   int64      `json:"actualAmount" gorm:"column:actual_amount;default:0"`
	Frequency      string     `json:"frequency" gorm:"column:frequency"`
	StartDate      *time.Time `json:"startDate" gorm:"column:start_date"`
	IsRecurring    bool       `json:"isRecurring" gorm:"column:is_recurring;default:false"`
	PaymentMethod  string     `json:"paymentMethod" gorm:"column:payment_method"`
	Notes          string     `json:"notes" gorm:"column:notes"`
	RawTags        string     `json:"-" gorm:"column:tags"` // JSON string array
	Tags           []string   `json:"tags" gorm:"-"`        // JSON string array

	Category DirectoryCategory `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
}

// TableName returns the table name for DirectorySource
func (DirectorySource) TableName() string {
	return "directory_sources"
}

// Validate validates the directory source entity
func (ds *DirectorySource) Validate() error {
	if ds.Name == "" {
		return fmt.Errorf("name is required")
	}
	if ds.Code == "" {
		return fmt.Errorf("code is required")
	}
	// if ds.Category == "" {
	// 	return fmt.Errorf("category is required")
	// }
	if ds.Type == "" {
		return fmt.Errorf("type is required")
	}
	// if ds.CategoryID == 0 {
	// 	return fmt.Errorf("category ID is required")
	// }
	return nil
}
