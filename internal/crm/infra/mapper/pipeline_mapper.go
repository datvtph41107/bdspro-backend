package mapper

import (
	base_enum "base/enum"
	_models "common/models"
	"crm/internal/domain"
	"crm/internal/dto"
	crmpb "pb/types/crm"
)

func PipelineSavePbToDomain(pb *crmpb.PipelineSaveRequest) *dto.PipelineSaveDTO {
	result := &dto.PipelineSaveDTO{
		ID:             pb.Id,
		PipelineName:   pb.PipelineName,
		Active:         pb.Active,
		Color:          pb.Color,
		IsDefault:      pb.IsDefault,
		RemoveStageIds: pb.RemoveStageIds,
	}

	if pb.Stages != nil {
		result.Stages = make([]*domain.StageEntity, len(pb.Stages))
		for i, stage := range pb.Stages {
			result.Stages[i] = StagePbToDomain(stage)
		}
	}

	return result
}

func PipelinePbToDomain(pb *crmpb.PipelineDTO) *domain.PipelineEntity {
	result := &domain.PipelineEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		PipelineName: pb.PipelineName,
		Active:       pb.Active,
		Color:        pb.Color,
		OwnerID:      pb.OwnerId,
		OwnerType:    base_enum.EOwnerOf(pb.OwnerType),
		IsDefault:    pb.IsDefault,
		IsCustom:     pb.IsCustom,
	}

	if pb.Stages != nil {
		result.Stages = make([]domain.StageEntity, len(pb.Stages))
		for i, stage := range pb.Stages {
			result.Stages[i] = *StagePbToDomain(stage)
		}
	}

	return result
}

func PipelineDomainToPb(domain *domain.PipelineEntity) *crmpb.PipelineDTO {
	pb := &crmpb.PipelineDTO{
		Id:           domain.ID,
		PipelineName: domain.PipelineName,
		Active:       domain.Active,
		Color:        domain.Color,
		OwnerId:      domain.OwnerID,
		OwnerType:    uint32(domain.OwnerType),
		IsDefault:    domain.IsDefault,
		IsCustom:     domain.IsCustom,
	}

	if domain.Stages != nil {
		pbStages := make([]*crmpb.StageDTO, len(domain.Stages))
		for i, stage := range domain.Stages {
			pbStages[i] = StageDomainToPb(&stage)
		}
		pb.Stages = pbStages
	}

	if domain.DefaultStage != nil {
		pb.DefaultStage = StageDomainToPb(domain.DefaultStage)
	}

	return pb
}