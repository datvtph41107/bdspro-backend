package domain

import _models "common/models"

// Project struct represents the projects table
// TableName sets the table name explicitly
type ProjectBuild struct {
	_models.BaseEntity
	Name         string              `gorm:"size:255;not null" json:"name"`
	NumFloor     uint                `gorm:"type:int" json:"numFloor"`
	NumApartment uint                `gorm:"type:int" json:"numApartment"`
	Note         string              `gorm:"type:text" json:"note"`
	ProjectID    uint64              `gorm:"not null" json:"projectId"`
	Project      *ProjectItem        `gorm:"foreignKey:ProjectID" json:"project"`
	Apartments   []ApartmentItem     `gorm:"foreignKey:BuildID" json:"apartments"`
	Attributes   []ApartmentAttrItem `gorm:"foreignKey:BuildID" json:"attributes"`
	// ListApartment    []ApartmentItem `gorm:"foreignKey:BuildID" json:"apartments"`
}

func (ProjectBuild) TableName() string {
	return "project_build"
}

type ProjectBuildItem struct {
	ID           uint64 `gorm:"primaryKey"`
	Name         string `gorm:"size:255;not null" json:"name"`
	NumFloor     uint   `gorm:"type:int" json:"numFloor"`
	NumApartment uint   `gorm:"type:int" json:"numApartment"`
	ProjectID    uint64 `gorm:"type:int8" json:"-"`
}

func (ProjectBuildItem) TableName() string {
	return "project_build"
}
