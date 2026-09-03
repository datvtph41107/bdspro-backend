package types

import "time"

// AuditAction — loại hành động được audit
type AuditAction string

const (
	ActionCreated      AuditAction = "created"
	ActionUpdated      AuditAction = "updated"
	ActionDeleted      AuditAction = "deleted"
	ActionReplaced     AuditAction = "replaced"
	ActionStateChanged AuditAction = "state_changed"
)

// AuditEntry — một bản ghi audit
type AuditEntry struct {
	ID         uint64
	EntityType string
	EntityID   uint64
	Action     AuditAction
	Actor      string
	ActorName  string
	Timestamp  time.Time
	Changes    []FieldDiff
	Reason     string
	RequestID  string
	IPAddress  string
	Metadata   map[string]any
}

// FieldDiff — diff 1 field
type FieldDiff struct {
	Field    string
	OldValue any
	NewValue any
}

// AuditQuery — query audit log
type AuditQuery struct {
	EntityType string
	EntityID   uint64
	Action     AuditAction
	Actor      string
	From       time.Time
	To         time.Time
	Page       int
	PageSize   int
}

// AuditConfig — cấu hình cho Audit Engine
type AuditConfig struct {
	RetentionMonths int
	PartitionBy     string
	EnabledActions  []AuditAction
}

// DefaultAuditConfig — cấu hình audit mặc định
func DefaultAuditConfig() *AuditConfig {
	return &AuditConfig{
		RetentionMonths: 24,
		PartitionBy:     "month",
		EnabledActions:  []AuditAction{ActionCreated, ActionUpdated, ActionDeleted, ActionReplaced, ActionStateChanged},
	}
}
