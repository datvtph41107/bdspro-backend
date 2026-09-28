package entity

import (
	"time"

	"gorm.io/gorm"
)

type Province struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt     *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Name          string         `gorm:"column:name;size:255;not null"`
	Type          int32          `gorm:"column:type;not null"`
	TypeText      string         `gorm:"column:type_text;size:50"`
	Codename      string         `gorm:"column:codename;size:100"`
	ShortCodename string         `gorm:"column:short_codename;size:100"`
	Slug          string         `gorm:"column:slug;size:100;index"`
}

func (Province) TableName() string { return "provinces" }

type District struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt     *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Name          string         `gorm:"column:name;size:255;not null"`
	ProvinceID    string         `gorm:"column:province_id;size:10;index;not null"`
	Type          int32          `gorm:"column:type;not null"`
	TypeText      string         `gorm:"column:type_text;size:50"`
	Codename      string         `gorm:"column:codename;size:100"`
	ShortCodename string         `gorm:"column:short_codename;size:100"`
}

func (District) TableName() string { return "districts" }

type Ward struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt     *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Name          string         `gorm:"column:name;size:255;not null"`
	DistrictID    string         `gorm:"column:district_id;size:10;index;not null"`
	Type          int32          `gorm:"column:type;not null"`
	TypeText      string         `gorm:"column:type_text;size:50"`
	Codename      string         `gorm:"column:codename;size:100"`
	ShortCodename string         `gorm:"column:short_codename;size:100"`
}

func (Ward) TableName() string { return "wards" }

type LocationSearchResult struct {
	DistrictID   uint64 `gorm:"column:district_id"`
	DistrictName string `gorm:"column:district_name"`
	DistrictType string `gorm:"column:district_type"`
	ProvinceID   uint64 `gorm:"column:province_id"`
	ProvinceName string `gorm:"column:province_name"`
	ProvinceType string `gorm:"column:province_type"`
	FullAddress  string `gorm:"column:full_address"`
}
