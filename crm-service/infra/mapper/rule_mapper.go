package mapper

import (
	base_enum "base/enum"
	_models "common/models"
	"crm/internal/domain"
	"crm/internal/enums"
	crmpb "pb/types/crm"
)

// ------------------------------------------------------------
func RulePbToDomain(req *crmpb.RuleDTO) *domain.RuleEntity {
	return &domain.RuleEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
		OwnerID:        req.OwnerId,
		OwnerType:      base_enum.EOwnerOf(req.OwnerType),
		Active:         req.Active,
		RuleName:       req.RuleName,
		Note:           req.Note,
		Condition:      enums.ECondition(req.Condition),
		ConditionValue: req.ConditionValue,
		Trigger:        enums.ERuleTrigger(req.Trigger),
		TriggerValue:   req.TriggerValue,
	}
}

func RuleDomainToPb(entity *domain.RuleEntity) *crmpb.RuleDTO {
	return &crmpb.RuleDTO{
		Id:             entity.ID,
		OwnerId:        entity.OwnerID,
		OwnerType:      int32(entity.OwnerType),
		Active:         entity.Active,
		RuleName:       entity.RuleName,
		Note:           entity.Note,
		Condition:      int32(entity.Condition),
		ConditionValue: entity.ConditionValue,
		Trigger:        int32(entity.Trigger),
		TriggerValue:   entity.TriggerValue,
	}
}