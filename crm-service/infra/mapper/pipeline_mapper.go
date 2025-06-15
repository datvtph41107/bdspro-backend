package mapper

import (
	_models "common/models"
	"crm/internal/domain"
	"crm/internal/enums"
	crmpb "pb/types/crm"
)

func PipelinePbToDomain(pb *crmpb.PipelineDTO) *domain.PipelineEntity {
	return &domain.PipelineEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		PipelineName: pb.PipelineName,
		Active:       pb.Active,
		Color:        pb.Color,
		OwnerID:      pb.OwnerId,
		OwnerType:    enums.EOwnerType(pb.OwnerType),
		IsDefault:    pb.IsDefault,
		IsCustom:     pb.IsCustom,
	}
}

func PipelineDomainToPb(domain *domain.PipelineEntity) *crmpb.PipelineDTO {
	return &crmpb.PipelineDTO{
		Id:           domain.ID,
		PipelineName: domain.PipelineName,
		Active:       domain.Active,
		Color:        domain.Color,
		OwnerId:      domain.OwnerID,
		OwnerType:    uint32(domain.OwnerType),
	}
}
