package domain

import (
	"time"

	"gorm.io/datatypes"
)

// DealHistoryEntity định nghĩa entity cho bảng deal_history
type DealHistoryEntity struct {
	ID          uint64            `json:"id" gorm:"primaryKey;autoIncrement" db:"id"`
	DealID      uint64            `json:"deal_id" gorm:"not null;index" db:"deal_id"`
	ActorID     uint64            `json:"actor_id" gorm:"not null;index" db:"actor_id"`
	ActorName   string            `json:"actor_name" gorm:"type:varchar(255);default:''" db:"actor_name"`
	ActorAvatar string            `json:"actor_avatar" gorm:"type:varchar(500);default:''" db:"actor_avatar"`
	ActionType  string            `json:"action_type" gorm:"type:varchar(50);not null;index" db:"action_type"`
	ActionName  string            `json:"action_name" gorm:"type:varchar(255);not null" db:"action_name"`
	Content     string            `json:"content" gorm:"type:text;not null" db:"content"`
	Metadata    datatypes.JSONMap `json:"metadata" gorm:"type:jsonb;default:'{}'" db:"metadata"`
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
}

// TableName returns the table name for GORM
func (DealHistoryEntity) TableName() string {
	return "deal_history"
}

// DealHistoryRequest định nghĩa request cho việc tạo deal history
type DealHistoryRequest struct {
	DealID      uint64            `json:"deal_id"`
	ActorID     uint64            `json:"actor_id"`
	ActorName   string            `json:"actor_name"`
	ActorAvatar string            `json:"actor_avatar"`
	ActionType  string            `json:"action_type"`
	ActionName  string            `json:"action_name"`
	Content     string            `json:"content"`
	Metadata    datatypes.JSONMap `json:"metadata"`
}

// DealHistoryResponse định nghĩa response từ deal history service
type DealHistoryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      uint64 `json:"id,omitempty"`
}
