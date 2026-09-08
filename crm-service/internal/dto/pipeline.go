package dto

import (
	_dto "common/domain/dto"
	"crm/internal/domain"
	"crm/internal/enums"
)

type PipelineSearchDTO struct {
	_dto.Pagable
	PipelineName string        `form:"pipelineName"`
	Steps        []enums.EStep `form:"steps"`
	Colors       []uint32      `form:"color"`
	Active       bool          `form:"active"`
}

type PipelineDTO struct {
	ID           uint64         `json:"id"`
	PipelineName string         `json:"pipelineName"`
	Color        uint32         `json:"color"`
	Active       bool           `json:"active"`
	OwnerID      uint64         `json:"ownerId"`
	OwnerType    enums.EOwnerOf `json:"ownerType"`
}

type PipelineSaveDTO struct {
	ID             uint64                `json:"id"`
	PipelineName   string                `json:"pipelineName"`
	Color          uint32                `json:"color"`
	Active         bool                  `json:"active"`
	IsDefault      bool                  `json:"isDefault"`
	Stages         []*domain.StageEntity `json:"stages"`
	RemoveStageIds []uint64              `json:"removeStageIds"`
}