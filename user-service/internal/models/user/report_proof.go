package user_model

import _models "common/models"

type ReportProof struct {
	_models.BaseEntity
	ReportID uint64 `gorm:"not null"`
	FileName string `gorm:"type:varchar(255);not null"`
	FileUrl  string `gorm:"type:varchar(255);not null"`
	FileType string `gorm:"type:varchar(255);not null"`
}
