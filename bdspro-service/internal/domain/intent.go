package domain

import (
	"bdspro/internal/enums"
	_models "common/domain/entity"
	"time"
)

// Intent - Domain model cho Intent detection và routing
type Intent struct {
	_models.BaseEntity
	UserInput      string                `gorm:"type:text;not null" json:"userInput"`
	DetectedIntent enums.EIntentAction   `gorm:"type:varchar(100);not null" json:"detectedIntent"`
	Category       enums.EIntentCategory `gorm:"type:varchar(50)" json:"category"`
	Confidence     float32               `gorm:"type:decimal(5,4)" json:"confidence"`
	Entities       string                `gorm:"type:jsonb" json:"entities"` // JSON string chứa extracted entities
	Context        string                `gorm:"type:jsonb" json:"context"`  // JSON string chứa context
	ProfileID      *uint64               `gorm:"index" json:"profileId"`
	OrganizationID *uint64               `gorm:"index" json:"organizationId"`
	SessionID      *string               `gorm:"type:varchar(255)" json:"sessionId"`
	Metadata       string                `gorm:"type:jsonb" json:"metadata"` // JSON string cho additional metadata
	ProcessedAt    *time.Time            `json:"processedAt"`
	ProcessedBy    *string               `gorm:"type:varchar(50)" json:"processedBy"` // "ai" hoặc "manual"
}

// TableName override table name
func (Intent) TableName() string {
	return "intents"
}

// IntentStatistics - Thống kê về intents
type IntentStatistics struct {
	Category      enums.EIntentCategory `json:"category"`
	Action        enums.EIntentAction   `json:"action"`
	TotalCount    int64                 `json:"totalCount"`
	AvgConfidence float32               `json:"avgConfidence"`
	LastUsed      *time.Time            `json:"lastUsed"`
}

// ExtractedEntity - Entity được trích xuất từ user input
type ExtractedEntity struct {
	Type       string  `json:"type"`       // "product_id", "price", "location", "area", etc.
	Value      string  `json:"value"`      // Giá trị thực tế
	Confidence float32 `json:"confidence"` // Độ tin cậy
	StartPos   int     `json:"startPos"`   // Vị trí bắt đầu trong text
	EndPos     int     `json:"endPos"`     // Vị trí kết thúc
}

// IntentContext - Context cho intent processing
type IntentContext struct {
	PreviousIntent *enums.EIntentAction `json:"previousIntent"`
	ConversationID *string              `json:"conversationId"`
	UserLocation   *string              `json:"userLocation"`
	UserRole       *string              `json:"userRole"`
	Timestamp      time.Time            `json:"timestamp"`
}
