package qh_domain

import (
	"encoding/json"
	"time"

	"tqd/internal/domain/jsonb"
)

// QHAuditEntry — bản ghi audit log cho mọi thao tác CRUD trên layer/parcel
type QHAuditEntry struct {
	ID         uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	EntityType string      `gorm:"type:varchar(64);not null;index:idx_audit_entity" json:"entityType"`
	EntityID   uint64      `gorm:"not null;index:idx_audit_entity" json:"entityId"`
	Action     string      `gorm:"type:varchar(32);not null;index" json:"action"`
	Actor      string      `gorm:"type:varchar(128)" json:"actor,omitempty"`
	ActorName  string      `gorm:"type:varchar(256)" json:"actorName,omitempty"`
	Timestamp  time.Time   `gorm:"not null;default:now();index:idx_audit_timestamp" json:"timestamp"`
	Changes    jsonb.JSONB `gorm:"type:jsonb" json:"changes,omitempty"`
	Reason     string      `gorm:"type:text" json:"reason,omitempty"`
	RequestID  string      `gorm:"type:varchar(64)" json:"requestId,omitempty"`
	IPAddress  string      `gorm:"type:varchar(45)" json:"ipAddress,omitempty"`
	Metadata   jsonb.JSONB `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"createdAt"`
}

func (QHAuditEntry) TableName() string {
	return "qh_audit_entries"
}

// ToJSONB converts a struct to jsonb.JSONB
func (e *QHAuditEntry) SetChanges(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	e.Changes = jsonb.JSONB(m)
	return nil
}

// SetMetadata sets metadata from any struct
func (e *QHAuditEntry) SetMetadata(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	e.Metadata = jsonb.JSONB(m)
	return nil
}
