package org_dto

type ProductShareRequest struct {
	ProductID uint64 `json:"productId" binding:"required"`
	GroupID   uint64 `json:"groupId" binding:"required"`
}
