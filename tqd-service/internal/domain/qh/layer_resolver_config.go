package qh_domain

import (
	"encoding/json"
	"time"
)

type LayerResolverConfig struct {
	ID          uint64          `gorm:"primaryKey" json:"id"`
	ConfigKey   string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"configKey"`
	ConfigValue json.RawMessage `gorm:"type:jsonb;not null" json:"configValue"`
	Description string          `gorm:"type:text" json:"description"`
	UpdatedBy   string          `gorm:"type:varchar(100)" json:"updatedBy"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"createdAt"`
}

func (LayerResolverConfig) TableName() string {
	return "layer_resolver_configs"
}
