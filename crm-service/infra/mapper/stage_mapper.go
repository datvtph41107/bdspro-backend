package mapper

import (
	_models "common/models"
	"crm/internal/domain"
	"crm/internal/enums"
	crmpb "pb/types/crm"
)

func StagePbToDomain(pb *crmpb.StageDTO) *domain.StageEntity {
	return &domain.StageEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		StageName:   pb.StageName,
		PipelineID:  pb.PipelineId,
		OrderNumber: int(pb.OrderNumber),
		RuleID:      &pb.RuleId,
		Step:        enums.EStep(pb.Step),
		Active:      pb.Active,
		Color:       pb.Color,
		ColorRGB:    pb.ColorRGB,
	}
}

func StageDomainToPb(domain *domain.StageEntity) *crmpb.StageDTO {
	return &crmpb.StageDTO{
		Id:          domain.ID,
		StageName:   domain.StageName,
		OrderNumber: int32(domain.OrderNumber),
		PipelineId:  domain.PipelineID,
		RuleId:      *domain.RuleID,
		Step:        uint32(domain.Step),
		Active:      domain.Active,
		Color:       domain.Color,
		ColorRGB:    domain.ColorRGB,
	}
}