package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProvinceV2 struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt    *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy    *uint64        `gorm:"column:created_by"`
	UpdatedBy    *uint64        `gorm:"column:updated_by"`
	Name         string         `gorm:"column:name;size:255;not null"`
	// db_v4 hub dump stores code/phone_code as varchar; keep string to match live schema.
	Code         string         `gorm:"column:code;size:20;index;not null"`
	Codename     string         `gorm:"column:codename;size:100"`
	DivisionType string         `gorm:"column:division_type;size:100"`
	PhoneCode    string         `gorm:"column:phone_code;size:10"`
	Wards        []WardV2       `gorm:"foreignKey:ProvinceCode;references:Code"`
}

func (ProvinceV2) TableName() string { return "province_v2" }

type WardV2 struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt     *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy     *uint64        `gorm:"column:created_by"`
	UpdatedBy     *uint64        `gorm:"column:updated_by"`
	Name          string         `gorm:"column:name;size:255;not null"`
	Code          string         `gorm:"column:code;size:20;index;not null"`
	Codename      string         `gorm:"column:codename;size:100"`
	DivisionType  string         `gorm:"column:division_type;size:100"`
	ShortCodename string         `gorm:"column:short_codename;size:100"`
	ProvinceCode  string         `gorm:"column:province_code;size:20;index;not null"`
	ProvinceID    uint64         `gorm:"column:province_id;index;not null"`
}

func (WardV2) TableName() string { return "ward_v2" }

type LocationSearchResultV2 struct {
	WardID       uint64 `gorm:"column:ward_id"`
	WardName     string `gorm:"column:ward_name"`
	WardType     string `gorm:"column:ward_type"`
	ProvinceID   uint64 `gorm:"column:province_id"`
	ProvinceName string `gorm:"column:province_name"`
	ProvinceType string `gorm:"column:province_type"`
	FullAddress  string `gorm:"column:full_address"`
}
