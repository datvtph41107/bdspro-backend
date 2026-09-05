package access

import (
	_models "common/models"
	"user/internal/enums"
)

// RoleGroup đại diện cho nhóm vai trò
type RoleGroup struct {
	_models.BaseEntity
	GroupName   string                 `gorm:"size:100;not null" json:"groupName"`
	Description string                 `gorm:"size:500" json:"description"`
	Key         enums.GroupRoleKeyEnum `gorm:"size:50;uniqueIndex" json:"key"`
	Code        string                 `gorm:"size:20;uniqueIndex" json:"code"`
	Permissions []*Permission          `gorm:"many2many:group_permissions;" json:"permissions,omitempty"`
}

func (RoleGroup) TableName() string {
	return "role_groups"
}
