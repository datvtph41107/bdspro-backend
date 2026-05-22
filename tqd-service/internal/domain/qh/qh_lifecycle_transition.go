package qh_domain

import "time"

// QHLifecycleTransition — bản ghi chuyển trạng thái trong vòng đời layer
type QHLifecycleTransition struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	LayerID   uint64    `gorm:"not null;index:idx_lifecycle_layer" json:"layerId"`
	FromState string    `gorm:"type:varchar(32);not null" json:"fromState"`
	ToState   string    `gorm:"type:varchar(32);not null" json:"toState"`
	Actor     string    `gorm:"type:varchar(128)" json:"actor,omitempty"`
	Reason    string    `gorm:"type:text" json:"reason,omitempty"`
	LegalDoc  string    `gorm:"type:varchar(512)" json:"legalDoc,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_lifecycle_created" json:"createdAt"`
}

func (QHLifecycleTransition) TableName() string {
	return "qh_lifecycle_transitions"
}
