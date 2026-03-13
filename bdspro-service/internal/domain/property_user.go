package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

type PropertyUser struct {
	_models.BaseEntity
	OwnerOriginID uint64 `gorm:"column:owner_origin_id;not null;index:idx_property_user_user_profile,priority:1"`

	// đối với bds định danh sẽ có id này
	PropertyLineageID *uint64 `gorm:"column:property_lineage_id;not null;index:idx_property_user_user_profile,priority:2"`
	// Để lấy được thông tin của property lineage
	PropertyLineage *PropertyLineage `gorm:"foreignKey:PropertyLineageID;references:ID" json:"propertyLineage,omitempty"`

	OwnerAt time.Time `gorm:"column:ownered_at;"`
	RoleID  *uint64   `gorm:"column:role_id;"`

	OriginProfileID *uint64    `gorm:"column:origin_id;type:bigint;index:idx_property_origin"`
	ArchivedAt      *time.Time `gorm:"column:archived_at" json:"archivedAt,omitempty"`
	HiddenAt        *time.Time `gorm:"column:hidden_at" json:"hiddenAt,omitempty"`

	RecordStatus enums.EPropertyStatus `gorm:"-" json:"recordStatus,omitempty"`
}

func (p *PropertyUser) ResolveRecordStatus() enums.EPropertyStatus {
	if p == nil {
		return enums.EPropertyStatusActive
	}
	if p.ArchivedAt != nil {
		return enums.EPropertyStatusArchived
	}
	if p.HiddenAt != nil {
		return enums.EPropertyStatusHide
	}
	return enums.EPropertyStatusActive
}

func (PropertyUser) TableName() string {
	return "property_user"
}
