package dto

import (
	_dto "common/domain/dto"
)

type RoleRequest struct {
	Name           string   `json:"name" validate:"required"`
	Key            string   `json:"key" validate:"required"`
	PermissionIDs  []uint64 `json:"permissionIds"`
	OrganizationID uint64   `json:"organizationId" validate:"required"`
	IsDefault      bool     `json:"isDefault"`
	DomainType     string   `json:"domainType"`
	ColorID        *uint64  `json:"colorId"`
}

type RoleResponse struct {
	ID             uint64                `json:"id"`
	Name           string                `json:"name"`
	Key            string                `json:"key"`
	PermissionIDs  []uint64              `json:"permissionIds"`
	OrganizationID uint64                `json:"organizationId"`
	Permissions    []*PermissionResponse `json:"permissions"`
	IsDefault      bool                  `json:"isDefault"`
	DomainType     string                `json:"domainType"`
	ColorID        *uint64               `json:"colorId"`
	Color          *ColorResponse        `json:"color"`
	CreatedAt      string                `json:"createdAt"`
	UpdatedAt      string                `json:"updatedAt"`
}

type RoleListRequest struct {
	_dto.Pagable
	OrganizationID uint64 `form:"organizationId"`
	DomainType     string `form:"domainType"`
	Search         string `form:"search"`
}

type RoleListResponse struct {
	Data          []*RoleResponse `json:"data"`
	TotalElements uint32          `json:"totalElements"`
}

type RoleSearchRequest struct {
	_dto.Pagable
	Scope       uint32  `form:"scope"`
	DomainType  string  `form:"domainType"`
	Name        string  `form:"name"`
	IsDefault   *bool   `form:"isDefault"`
	RoleGroupID *uint64 `form:"roleGroupId"`
}
