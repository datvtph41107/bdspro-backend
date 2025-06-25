package entity

import "github.com/tuandq2112/organization-service/internal/enums"

type OrganizationPermission struct {
	Id             uint32
	Name           string
	Key            string
	Description    string
	PermissionType enums.PermissionType      `gorm:"not null;default:10" json:"permission_type"`
	ParentId       uint32                    `json:"parent_id"`
	Childrens      []*OrganizationPermission `gorm:"foreignKey:ParentId" json:"childrens"`
}
