package mapper

import (
	_models "common/models"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/crm/enums"
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
	out := &crmpb.StageDTO{
		Id:          domain.ID,
		StageName:   domain.StageName,
		OrderNumber: int32(domain.OrderNumber),
		PipelineId:  domain.PipelineID,
		Step:        uint32(domain.Step),
		Active:      domain.Active,
		Color:       domain.Color,
		ColorRGB:    domain.ColorRGB,
	}
	if domain.RuleID != nil {
		out.RuleId = *domain.RuleID
	}
	return out
}
