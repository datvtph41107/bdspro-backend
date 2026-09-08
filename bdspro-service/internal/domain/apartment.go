package domain

import _models "common/models"

type Apartment struct {
	_models.BaseEntity
	Name            string `gorm:"size:100" json:"name"`
	Note            string `gorm:"type:text" json:"note"`
	Floor           *int   `gorm:"type:int" json:"floor"`
	Ordinal         int    `gorm:"type:int;not null" json:"ordinal"`
	ApartmentAttrID uint64 `gorm:"column:attribute_id;type:int8;foreignKey:ID" json:"attributeId,omitempty"`
	Status          uint   `gorm:"type:smallint" json:"status,omitempty"`
	Archived        uint   `gorm:"type:smallint" json:"archived,omitempty"`
	BuildID         uint64 `gorm:"type:int8" json:"buildId,omitempty"`
}

func (Apartment) TableName() string {
	return "apartment"
}

type ApartmentItem struct {
	ID              uint64 `gorm:"primaryKey" json:"id"`
	Name            string `gorm:"size:100" json:"name"`
	Floor           *int   `gorm:"type:int" json:"floor"`
	Ordinal         int    `gorm:"type:int;not null" json:"ordinal"`
	BuildID         uint64 `gorm:"type:int8" json:"-"`
	Status          uint   `gorm:"type:smallint" json:"status,omitempty"`
	Archived        uint   `gorm:"type:smallint" json:"archived,omitempty"`
	ApartmentAttrID uint64 `gorm:"column:attribute_id;type:int8;foreignKey:ID" json:"attributeId,omitempty"`
}

func (ApartmentItem) TableName() string {
	return "apartment"
}
