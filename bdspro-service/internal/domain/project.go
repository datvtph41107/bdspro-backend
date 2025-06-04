package domain

import _models "common/models"

// Project struct represents the projects table
// TableName sets the table name explicitly
type Project struct {
	_models.BaseEntity
	Name        string             `gorm:"size:255;not null" json:"name"`
	DeveloperID uint64             `gorm:"type:uint64" json:"developerId"`
	Developer   *DeveloperItem     `gorm:"foreignKey:DeveloperID;references:ID" json:"developer"`
	Description string             `gorm:"type:text" json:"description"`
	Builds      []ProjectBuildItem `gorm:"foreignKey:ProjectID" json:"builds,omitempty"`
	// DeveloperName string             `gorm:"type:text" json:"developerName"`
}

func (Project) TableName() string {
	return "projects"
}

type ProjectItem struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255" json:"name"`
}

func (ProjectItem) TableName() string {
	return "projects"
}
