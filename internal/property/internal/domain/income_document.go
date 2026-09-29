package domain

import (
	_models "common/models"
	"time"
)

type IncomeDocument struct {
	_models.BaseEntity
	IncomeID   uint64    `json:"incomeId"`
	FilePath   string    `json:"filePath"`
	FileName   string    `json:"fileName"`
	UploadedAt time.Time `json:"uploadedAt"`
}
