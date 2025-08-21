package domain

import (
	_models "common/models"
)

type Role struct {
	_models.BaseEntity
	RoleName        string        `gorm:"size:100;not null" json:"roleName"`
	RoleDescription string        `gorm:"size:255" json:"roleDescription"`
	Key             string        `gorm:"size:50" json:"key"`
	PermissionIDs   []uint64      `gorm:"-" json:"permissionIds"`
	Permissions     []*Permission `gorm:"many2many:role_permissions;"`
	IsDefault       bool          `gorm:"default:false" json:"isDefault"`
	DomainType      string        `gorm:"size:20;default:'ORGANIZATION'" json:"domainType"`
	RoleGroupID     *uint64       `gorm:"type:int" json:"roleGroupId"`
	RoleGroup       *RoleGroup    `gorm:"foreignKey:RoleGroupID" json:"roleGroup"`
	ColorID         *uint64       `json:"colorId"`
	Color           *Color        `gorm:"foreignKey:ColorID" json:"color"`
	OrganizationID  uint64        `gorm:"not null" json:"organizationId"`
}

func (Role) TableName() string {
	return "roles"
}
