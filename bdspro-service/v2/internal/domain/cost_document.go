package domain

import (
	_models "common/models"
	"time"
)

type CostDocument struct {
	_models.BaseEntity
	CostID     uint64    `json:"costId"`
	FilePath   string    `json:"filePath"`
	FileName   string    `json:"fileName"`
	UploadedAt time.Time `json:"uploadedAt"`
}
