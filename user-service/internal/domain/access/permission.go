package access

import (
	_models "common/models"
)

type Permission struct {
	_models.BaseEntity
	Name           string        `gorm:"size:100;not null" json:"name"`
	Key            string        `gorm:"size:50;" json:"key"`
	Description    string        `gorm:"size:500" json:"description"`
	PermissionType string        `gorm:"size:20;default:'ORGANIZATION'" json:"permissionType"`
	Module         string        `gorm:"size:50;" json:"module"`
	ParentID       *uint64       `json:"parentId"`
	Children       []*Permission `gorm:"-" json:"children"`
	Priority       uint32        `gorm:"type:int;default:10" json:"priority"`
	RequiredKey    string        `gorm:"size:10;" json:"requiredKey"`
	RoleGroups     []*RoleGroup  `gorm:"many2many:group_permissions;" json:"role_groups,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}
