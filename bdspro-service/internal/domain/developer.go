package domain

import (
	_models "common/models"
	"time"
)

type Developer struct {
	_models.BaseEntity
	Name          string     `gorm:"size:255" json:"name"`
	Slug          string     `gorm:"size:15" json:"slug"`
	LogoUrl       string     `json:"logoUrl,omitempty"`
	BackgroundUrl string     `json:"backgroundUrl,omitempty"`
	Description   string     `json:"description,omitempty"`
	Address       string     `gorm:"size:15" json:"address,omitempty"`
	Phone         string     `gorm:"size:15" json:"phone,omitempty"`
	Email         string     `gorm:"size:255" json:"email,omitempty"`
	Website       string     `json:"website,omitempty"`
	FoundedAt     *time.Time `json:"foundedAt,omitempty"`
	Status        uint       `json:"status,omitempty"`
}

func (Developer) TableName() string {
	return "developer"
}

type DeveloperItem struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"size:255" json:"name"`
	Status uint   `json:"status"`
}

func (DeveloperItem) TableName() string {
	return "developer"
}
