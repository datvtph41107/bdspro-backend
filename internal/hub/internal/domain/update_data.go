package domain

import (
	"strings"
)

type ESyncResource int32

const (
	ESyncResourceProduct ESyncResource = 10
	ESyncResourceAsset   ESyncResource = 20
	ESyncResourceContact ESyncResource = 30
)

// ParseResource maps string to ESyncResource
func ParseResource(s string) ESyncResource {
	switch strings.ToLower(s) {
	case "product":
		return ESyncResourceProduct
	case "asset":
		return ESyncResourceAsset
	case "contact":
		return ESyncResourceContact
	default:
		return ESyncResourceProduct
	}
}

// String returns resource name for ESyncResource
func (e ESyncResource) String() string {
	switch e {
	case ESyncResourceProduct:
		return "product"
	case ESyncResourceAsset:
		return "asset"
	case ESyncResourceContact:
		return "contact"
	default:
		return "product"
	}
}

type UpdateData struct {
	ID           uint64        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OwnerID      uint64        `gorm:"column:owner_id;not null;index:idx_update_data_owner_resource" json:"ownerId"`
	ResourceType ESyncResource `gorm:"column:resource_type;not null;index:idx_update_data_owner_resource" json:"resourceType"`
	ResourceID   uint64        `gorm:"column:resource_id;not null" json:"resourceId"`
	UpdatedAt    int64         `gorm:"column:updated_at;not null;index" json:"updatedAt"`
}

func (UpdateData) TableName() string {
	return "update_data"
}
