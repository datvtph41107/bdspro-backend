package qh

import (
	"time"

	_models "common/domain/entity"

	"gorm.io/datatypes"
)

// ParcelHistoryEventType — loại sự kiện lịch sử thửa đất (string, giống Đồ án QH).
const (
	ParcelHistoryEventCreate         = "CREATE"
	ParcelHistoryEventUpdateInfo     = "UPDATE_INFO"
	ParcelHistoryEventUpdateGeometry = "UPDATE_GEOMETRY"
	ParcelHistoryEventOwnership      = "OWNERSHIP"
	ParcelHistoryEventStatus         = "STATUS"
	ParcelHistoryEventNote           = "NOTE"
	ParcelHistoryEventOther          = "OTHER"
)

// ParcelHistory — lịch sử thay đổi của thửa đất (cùng shape với QHPlanningEvent).
type ParcelHistory struct {
	_models.BaseEntity

	ParcelID    uint64         `gorm:"column:parcel_id;not null;index:idx_parcel_history_parcel_timeline" json:"parcelId"`
	EventType   string         `gorm:"column:event_type;type:varchar(80)" json:"eventType"`
	EventName   string         `gorm:"column:event_name;type:varchar(255);not null" json:"eventName"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	EventDate   *time.Time     `gorm:"column:event_date" json:"eventDate"`
	DocumentID  *uint64        `gorm:"column:document_id" json:"documentId"`
	SourceURL   string         `gorm:"column:source_url;type:text" json:"sourceUrl"`
	Metadata    datatypes.JSON `gorm:"column:metadata;type:jsonb;not null;default:'{}'" json:"metadata"`
}

func (ParcelHistory) TableName() string {
	return "parcel_history"
}
