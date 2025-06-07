package dto

import (
	_dto "common/dto"
	"crm/internal/enums"
)

type PipelineSearchDTO struct {
	_dto.Page
	PipelineName string        `form:"pipelineName"`
	Steps        []enums.EStep `form:"steps"`
	Colors       []string      `form:"color"`
	Active       bool          `form:"active"`
}

type PipelineDTO struct {
	ID           uint64           `json:"id"`
	PipelineName string           `json:"pipelineName"`
	Color        string           `json:"color"`
	Active       bool             `json:"active"`
	OwnerID      uint64           `json:"ownerId"`
	OwnerType    enums.EOwnerType `json:"ownerType"`
}
