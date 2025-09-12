package domain

import (
	"auth/internal/enums"
	_models "common/models"
)

// RoleGroup đại diện cho nhóm vai trò
type RoleGroup struct {
	_models.BaseEntity
	GroupName   string                 `gorm:"size:100;not null" json:"groupName"`
	Description string                 `gorm:"size:500" json:"description"`
	Key         enums.GroupRoleKeyEnum `gorm:"size:50;uniqueIndex" json:"key"`
	Permissions []*Permission          `gorm:"many2many:group_permissions;" json:"permissions,omitempty"`
}

func (RoleGroup) TableName() string {
	return "role_groups"
}
