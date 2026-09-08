package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type PropertyAuditLog struct {
	ID               uint64           `gorm:"primaryKey;autoIncrement"`
	PropertyID       uint64           `gorm:"not null;index:idx_property_audit_property"`
	UserID           uint64           `gorm:"not null"`
	Action           string           `gorm:"type:varchar(50);not null"`
	ImpactLevel      int32            `gorm:"column:impact_level"`
	AffectedProducts int              `gorm:"column:affected_products"`
	AffectedListings int              `gorm:"column:affected_listings"`
	AffectedAssets   int              `gorm:"column:affected_assets"`
	Changes          AuditChangesJSON `gorm:"type:jsonb"`
	ConfirmedAt      *time.Time       `gorm:"column:confirmed_at"`
	SessionID        string           `gorm:"type:varchar(100);column:session_id"`
	BlockedReason    string           `gorm:"type:text;column:blocked_reason"`
	ProposeLineageID *uint64          `gorm:"column:propose_lineage_id"`
	CreatedAt        time.Time        `gorm:"autoCreateTime"`
}

func (PropertyAuditLog) TableName() string {
	return "property_audit_logs"
}

type AuditChangesJSON map[string]interface{}

func (j AuditChangesJSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *AuditChangesJSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}
