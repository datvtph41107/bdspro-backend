package entity

import (
	"time"

	_models "common/models"
)

// Project is the full projects row (admin CRUD).
type Project struct {
	_models.BaseEntity
	Name        string `gorm:"size:255;not null;column:name" json:"name"`
	DeveloperID uint64 `gorm:"column:developer_id" json:"developerId"`
	Description string `gorm:"type:text;column:description" json:"description"`
}

func (Project) TableName() string {
	return "projects"
}

// ProjectItem is a lightweight row for catalog list (table projects).
type ProjectItem struct {
	ID   uint64 `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (ProjectItem) TableName() string {
	return "projects"
}

// ProjectBuild is a building/tower row (table project_build).
type ProjectBuild struct {
	_models.BaseEntity
	Name         string              `gorm:"size:255;not null;column:name" json:"name"`
	NumFloor     uint                `gorm:"column:num_floor" json:"numFloor"`
	NumApartment uint                `gorm:"column:num_apartment" json:"numApartment"`
	Note         string              `gorm:"type:text;column:note" json:"note"`
	ProjectID    uint64              `gorm:"column:project_id;not null" json:"projectId"`
	Apartments   []ApartmentItem     `gorm:"foreignKey:BuildID" json:"apartments,omitempty"`
	Attributes   []ApartmentAttrItem `gorm:"foreignKey:BuildID" json:"attributes,omitempty"`
}

func (ProjectBuild) TableName() string {
	return "project_build"
}

// ApartmentItem is a row in apartment linked to project_build.
type ApartmentItem struct {
	ID              uint64 `gorm:"primaryKey;column:id" json:"id"`
	Name            string `gorm:"size:100;column:name" json:"name"`
	Floor           *int   `gorm:"column:floor" json:"floor"`
	Ordinal         int    `gorm:"column:ordinal;not null" json:"ordinal"`
	BuildID         uint64 `gorm:"column:build_id" json:"buildId"`
	Status          uint   `gorm:"column:status" json:"status"`
	Archived        uint   `gorm:"column:archived" json:"archived"`
	ApartmentAttrID uint64 `gorm:"column:attribute_id" json:"attributeId"`
}

func (ApartmentItem) TableName() string {
	return "apartment"
}

// ApartmentAttrItem is apartment_attribute linked to project_build.
type ApartmentAttrItem struct {
	ID           uint64     `gorm:"primaryKey;column:id" json:"id"`
	BuildID      uint64     `gorm:"column:build_id" json:"buildId"`
	Area         *float64   `gorm:"column:area" json:"area"`
	NumBedroom   *int       `gorm:"column:num_bedroom" json:"numBedroom"`
	NumBathroom  *int       `gorm:"column:num_bathroom" json:"numBathroom"`
	Furniture    string     `gorm:"column:furniture" json:"furniture"`
	BlueprintUrl string     `gorm:"column:blueprint_url" json:"blueprintUrl"`
	Price        float64    `gorm:"column:price" json:"price"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (ApartmentAttrItem) TableName() string {
	return "apartment_attribute"
}
