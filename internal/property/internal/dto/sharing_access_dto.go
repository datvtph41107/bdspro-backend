package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
)

type SharingAccessSearch struct {
	_dto.Pagable
	DomainID uint64              `form:"domainId"`
	Domain   enums.EDomainAccess `form:"domain"`
	FromType enums.EOwnerOf      `form:"fromType"`
	ToType   enums.EOwnerOf      `form:"toType"`
	// TargetType enums.ETargetType `form:"targetType"`
}

type SharingAccessRequest struct {
	// ProductID  *uint64 `json:"parentId" binding:"required"`
	ToID       *uint64 `json:"toId" binding:"required"`
	ToType     uint    `json:"toType"`
	Commission float32 `json:"commission"`
	Fields     string  `json:"fields"`
	// FromType   uint    `json:"fromType" binding:"required"`
	// ProfileIDs []uint64 `json:"profileIds" binding:"required,notemptyarray"`
}

type SharingAccessBulk struct {
	DomainID uint64                  `json:"domainId"`
	Domain   enums.EDomainAccess     `json:"domain"`
	Datas    []*domain.SharingAccess `json:"datas"`
	Deletes  []*domain.SharingAccess `json:"deletes"`
}

type CommissionUpdate struct {
	ToID       uint64  `json:"targetId" binding:"required"`
	ToType     uint64  `json:"targetType" binding:"required"`
	ProductID  uint64  `json:"productId" binding:"required"`
	FromType   uint64  `json:"fromType" binding:"required"`
	Commission float64 `json:"commission" binding:"required"`
}

type SharingAccessDeleted struct {
	DomainID uint64 `json:"domainId" binding:"required"`
	FromType uint64 `json:"fromType" binding:"required"`
	// Owners   []ToOwner `json:"owners" binding:"required"`
}
