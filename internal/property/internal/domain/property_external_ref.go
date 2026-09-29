package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

// PropertyExternalRef chứa thông tin tham chiếu ngoài (mở rộng của Property)
type PropertyExternalRef struct {
	_models.BaseEntity
	LineageID uint64 `gorm:"type:bigint;not null;index"`
	// đối với bds định danh sẽ có id này
	PropertyIdentifyID *uint64 `gorm:"type:int8;index" json:"propertyIdentifyId"` // Liên kết với Property
	// Để lấy được thông tin của property identify
	PropertyIdentify *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID"`

	// id nguồn tham chiếu
	ExternalRefID uint64 `gorm:"type:bigint;not null;index" json:"externalRefId,omitempty"`
	// loại nguồn tham chiếu
	SourceSystem enums.ERefSourceSys `gorm:"type:smallint;index" json:"sourceSystem,omitempty"`
	// mã nguồn tham chiếu
	SourceCode string `gorm:"size:255" json:"sourceCode,omitempty"`
	// độ tin cậy
	Confidence uint32 `gorm:"type:smallint;default:0" json:"confidence,omitempty"`
	// trạng thái sync
	SyncStatus enums.ESyncStatus `gorm:"type:smallint;default:10" json:"syncStatus,omitempty"`
	// ghi chú
	Note string `gorm:"type:text" json:"note,omitempty"`
	// // ngày sync
	// SyncedAt *time.Time `gorm:"column:synced_at;autoCreateTime;" json:"syncedAt,omitempty"`
}

func (PropertyExternalRef) TableName() string {
	return "property_external_ref"
}
