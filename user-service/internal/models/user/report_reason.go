package user_model

import _models "common/models"

type ReportReason struct {
	_models.BaseEntity
	Reason string `gorm:"type:varchar(255);not null"`
}
