package entity

// ESyncResource identifies which resource type the update_data row belongs to.
type ESyncResource int32

const (
	ESyncResourceProduct ESyncResource = 10
	ESyncResourceAsset   ESyncResource = 20
	ESyncResourceContact ESyncResource = 30
)

// ParseSyncResource maps a resource string to ESyncResource.
func ParseSyncResource(s string) ESyncResource {
	switch s {
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

// UpdateData tracks which resource IDs have changed per owner so clients can sync.
type UpdateData struct {
	ID           uint64        `gorm:"column:id;primaryKey;autoIncrement"`
	OwnerID      uint64        `gorm:"column:owner_id;not null;index:idx_update_data_owner_resource"`
	ResourceType ESyncResource `gorm:"column:resource_type;not null;index:idx_update_data_owner_resource"`
	ResourceID   uint64        `gorm:"column:resource_id;not null"`
	UpdatedAt    int64         `gorm:"column:updated_at;not null;index"`
}

func (UpdateData) TableName() string { return "update_data" }
