package domain

import (
	_models "common/models"
	"time"

	"gorm.io/gorm"
)

// Bảng tỉnh/thành phố
type Province struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Code      int            `gorm:"unique;not null" json:"code"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Bảng quận/huyện
type District struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Code       int            `gorm:"unique;not null" json:"code"`
	ProvinceID uint64         `gorm:"not null" json:"province_id"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Province   *Province      `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
}

// Bảng phường/xã
type Ward struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Code       int            `gorm:"unique;not null" json:"code"`
	DistrictID uint64         `gorm:"not null" json:"district_id"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	District   *District      `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
}

// Bảng phường/xã
type Region struct {
	_models.BaseEntity
	Name     string `gorm:"size:100;not null" json:"name"`
	Code     uint   `json:"code"`
	CodeName string `json:"codeName"`
	Level    uint   `json:"level"`
	ParentID uint64 `json:"parentId"`
	Unit     string `json:"unit"`
	// Parent   *Region `gorm:"foreignKey:ParentID"`
}

func (Region) TableName() string {
	return "region"
}
