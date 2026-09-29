package domain

import (
	_entity "common/domain/entity"
	"fmt"
	"time"
)

// DirectorySupplier represents a directory supplier entity
type DirectorySupplier struct {
	_entity.BaseEntity
	Name            string     `json:"name" gorm:"column:name;not null"`
	Code            string     `json:"code" gorm:"column:code;not null;unique"`
	Description     string     `json:"description" gorm:"column:description"`
	ContactPerson   string     `json:"contactPerson" gorm:"column:contact_person"`
	Phone           string     `json:"phone" gorm:"column:phone"`
	Email           string     `json:"email" gorm:"column:email"`
	Address         string     `json:"address" gorm:"column:address"`
	Website         string     `json:"website" gorm:"column:website"`
	TaxCode         string     `json:"taxCode" gorm:"column:tax_code"`
	BusinessLicense string     `json:"businessLicense" gorm:"column:business_license"`
	Rating          float64    `json:"rating" gorm:"column:rating;default:0"`
	IsActive        bool       `json:"isActive" gorm:"column:is_active;default:true"`
	ServiceCount    uint32     `json:"serviceCount" gorm:"column:service_count;default:0"`
	ContractCount   uint32     `json:"contractCount" gorm:"column:contract_count;default:0"`
	TotalValue      int64      `json:"totalValue" gorm:"column:total_value;default:0"`
	JoinDate        *time.Time `json:"joinDate" gorm:"column:join_date"`
	LastContactDate *time.Time `json:"lastContactDate" gorm:"column:last_contact_date"`
	Notes           string     `json:"notes" gorm:"column:notes"`
	RawTags         string     `json:"-" gorm:"column:tags"`       // JSON string array
	Tags            []string   `json:"tags" gorm:"-"`              // JSON string array
	RawCategories   string     `json:"-" gorm:"column:categories"` // JSON string array
	Categories      []string   `json:"categories" gorm:"-"`        // JSON string array
}

// TableName returns the table name for DirectorySupplier
func (DirectorySupplier) TableName() string {
	return "directory_suppliers"
}

// Validate validates the directory supplier entity
func (ds *DirectorySupplier) Validate() error {
	if ds.Name == "" {
		return fmt.Errorf("name is required")
	}
	if ds.Code == "" {
		return fmt.Errorf("code is required")
	}
	return nil
}
