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
