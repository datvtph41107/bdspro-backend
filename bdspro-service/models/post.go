package models

import (
	"bdspro/enums"
	_models "common/models"
	"time"
)

type Post struct {
	_models.BaseEntity
	ID         uint64               `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductId  uint64               `json:"productId"`
	FinishedAt *time.Time           `json:"finishedAt"`
	Status     enums.PostStatus     `json:"status"`
	Visibility enums.PostVisibility `json:"visibility"`
}
