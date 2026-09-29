package domain

import _models "common/models"

// ProjectInventory struct cho bảng project_inventory
type ProjectInventory struct {
	_models.BaseEntity
	// ID        uint64    `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	ProjectID uint64  `gorm:"not null" json:"project_id"`
	Block     string  `gorm:"type:VARCHAR(50);default:null" json:"block"`
	Floor     int     `gorm:"default:null" json:"floor"`
	UnitCode  string  `gorm:"type:VARCHAR(50);default:null" json:"unit_code"`
	Area      float64 `gorm:"type:DECIMAL(10,2);default:0;check:area >= 0" json:"area"`
	Price     float64 `gorm:"type:DECIMAL(15,0);default:0;check:price >= 0" json:"price"`
	Status    int     `gorm:"type:smallint;check:status IN (1,2,3)" json:"status"`
	ProductID *uint64 `gorm:"default:null" json:"product_id"` // nullable field
	// CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	// UpdatedAt time.Time `gorm:"default:current_timestamp on update current_timestamp" json:"updated_at"`
}

// GORM table name override
func (ProjectInventory) TableName() string {
	return "project_inventory"
}
