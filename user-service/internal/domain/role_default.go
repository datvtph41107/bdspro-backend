package domain

import (
	_models "common/models"
)

type RoleDefault struct {
	_models.BaseEntity
	RoleName        string        `gorm:"size:100;not null" json:"roleName"`
	RoleDescription string        `gorm:"size:255" json:"roleDescription"`
	PermissionIDs   []uint64      `gorm:"-" json:"permissionIds"`
	Permissions     []*Permission `gorm:"many2many:role_default_permissions;"`
	RoleGroupID     *uint64       `gorm:"type:int" json:"roleGroupId"`
	RoleGroup       *RoleGroup    `gorm:"foreignKey:RoleGroupID" json:"roleGroup"`
	ColorID         *uint64       `json:"colorId"`
	Color           *Color        `gorm:"foreignKey:ColorID" json:"color"`
}

func (RoleDefault) TableName() string {
	return "role_defaults"
}
