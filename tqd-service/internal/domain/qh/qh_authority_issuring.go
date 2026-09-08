package qh_domain

import _models "common/domain/entity"

type QHAuthorityIssuring struct {
	_models.BaseEntity
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Code        string `gorm:"type:varchar(100);not null" json:"code"`
	Description string `gorm:"type:text" json:"description,omitempty"`
}

func (QHAuthorityIssuring) TableName() string {
	return "qh_authority_issuring"
}
