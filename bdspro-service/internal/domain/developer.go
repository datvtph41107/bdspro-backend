package domain

import (
	_models "common/models"
	"time"
)

type Developer struct {
	_models.BaseEntity
	Name          string     `gorm:"size:255" json:"name"`
	Slug          string     `gorm:"size:15" json:"slug"`
	LogoUrl       string     `json:"logoUrl"`
	BackgroundUrl string     `json:"backgroundUrl"`
	Description   string     `json:"description"`
	Address       string     `gorm:"size:15" json:"address"`
	Phone         string     `gorm:"size:15" json:"phone"`
	Email         string     `gorm:"size:255" json:"email"`
	Website       string     `json:"website"`
	FoundedAt     *time.Time `json:"foundedAt"`
	Status        uint       `json:"status"`
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
