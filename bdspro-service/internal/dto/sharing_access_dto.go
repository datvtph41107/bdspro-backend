package dto

import (
	"bdspro/internal/enums"
	_dto "common/dto"
)

type SharingAccessSearch struct {
	_dto.Pagable
	DomainID uint64              `form:"domainId"`
	Domain   enums.EDomainAccess `form:"domain"`
	FromType enums.EOwnerType    `form:"fromType"`
	ToType   enums.EOwnerType    `form:"toType"`
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
	Datas   []SharingAccessRequest `json:"datas"`
	Deletes []SharingAccessRequest `json:"deletes"`
}

type ProductAccessUpdate struct {
	TargetID   uint64  `json:"targetId" binding:"required"`
	TargetType uint64  `json:"targetType" binding:"required"`
	ProductID  uint64  `json:"productId" binding:"required"`
	Commission float64 `json:"commission" binding:"required"`
}
