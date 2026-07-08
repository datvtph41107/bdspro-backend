package domain

// FeedbackRegistry represents a service registry entry
type FeedbackRegistry struct {
	RegistryKey  uint32 `gorm:"column:registry_key;type:int;uniqueIndex;primaryKey" json:"registryKey" binding:"required"`
	RegistryName string `gorm:"column:registry_name;type:varchar(20);not null" json:"registryName" binding:"required,max=20"`
	Port         int    `gorm:"not null" json:"port" binding:"required"`
	Service      string `gorm:"type:varchar(100);not null" json:"service" binding:"required"`
}

// TableName returns the table name for FeedbackRegistry
func (FeedbackRegistry) TableName() string {
	return "feedback_registry"
}

// RegistryListRequest represents request parameters for getting registry list
type RegistryListRequest struct {
	RegistryKey uint32 `json:"registryKey,omitempty"`
	Service     string `json:"service,omitempty"`
	Page        int    `json:"page"`
	Size        int    `json:"size"`
}