package qh_domain

import (
	"encoding/json"
	"time"
)

// QHLayerResolverConfig - Thiết lập thay đổi cách đánh giá tính toán cho Layer Resolver
// Lưu trong bảng qh_layer_resolver_configs, mỗi row là một key-value config (JSONB)
// Dùng để override các giá trị mặc định từ YAML/config.go
type QHLayerResolverConfig struct {
	ID          uint64          `gorm:"primaryKey" json:"id"`
	ConfigKey   string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"configKey"`
	ConfigValue json.RawMessage `gorm:"type:jsonb;not null" json:"configValue"`
	Description string          `gorm:"type:text" json:"description"`
	UpdatedBy   string          `gorm:"type:varchar(100)" json:"updatedBy"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"createdAt"`
}

func (QHLayerResolverConfig) TableName() string {
	return "qh_layer_resolver_configs"
}
