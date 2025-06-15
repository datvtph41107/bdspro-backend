package grpc_service

import (
	_dto "common/dto"
	_models "common/models"
	_routes "common/routes"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
)

type RuleService struct {
	crmpb.UnimplementedRuleServiceServer
	UC *usecase.RuleUsecase
}

func NewRuleService(uc *usecase.RuleUsecase) *RuleService {
	return &RuleService{
		UC: uc,
	}
}

func (s *RuleService) Search(ctx context.Context, req *crmpb.RuleSearchRequest) (*crmpb.RuleListDTO, error) {
	dto := dto.RuleSearchDTO{
		Page: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		RuleName: req.Name,
	}

	rules, total, err := s.UC.Search(ctx, req.OwnerId, enums.EOwnerType(req.OwnerType), dto)
	if err != nil {
		return nil, err
	}

	rulesPb := make([]*crmpb.RuleDTO, len(rules))
	for i, rule := range rules {
		rulesPb[i] = s.EntityToPb(&rule)
	}

	return &crmpb.RuleListDTO{
		Data:          rulesPb,
		TotalElements: uint32(total),
	}, nil
}
func (s *RuleService) Detail(ctx context.Context, req *crmpb.RuleDTO) (*crmpb.RuleDTO, error) {
	// result, err := s.UC.GetByID(ctx, req.Id)
	// return s.EntityToPb(result), err
	return nil, nil
}
func (s *RuleService) Create(ctx context.Context, req *crmpb.RuleDTO) (*crmpb.RuleDTO, error) {
	entity := s.PbToEntity(req)
	if entity.OwnerID == 0 ||
		entity.OwnerType == 0 ||
		entity.RuleName == "" ||
		entity.Condition == 0 ||
		entity.Trigger == 0 ||
		entity.ConditionValue == "" ||
		entity.TriggerValue == "" {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ownerId, ownerType, ruleName, condition, trigger, conditionValue and triggerValue are required",
		}
	}
	result, err := s.UC.Create(ctx, entity)
	return s.EntityToPb(result), err
}
func (s *RuleService) Update(ctx context.Context, req *crmpb.RuleDTO) (*crmpb.RuleDTO, error) {
	entity := s.PbToEntity(req)
	if entity.ID == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "id is required",
		}
	}
	result, err := s.UC.Update(ctx, entity.ID, entity)
	return s.EntityToPb(result), err
}
func (s *RuleService) Delete(ctx context.Context, req *crmpb.RuleDTO) (*crmpb.RuleDTO, error) {
	err := s.UC.Delete(ctx, req.Id)
	return nil, err
}

// ------------------------------------------------------------
func (s *RuleService) PbToEntity(req *crmpb.RuleDTO) *domain.RuleEntity {
	return &domain.RuleEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
		OwnerID:        req.OwnerId,
		OwnerType:      enums.EOwnerType(req.OwnerType),
		Active:         req.Active,
		RuleName:       req.RuleName,
		Note:           req.Note,
		Condition:      enums.ECondition(req.Condition),
		ConditionValue: req.ConditionValue,
		Trigger:        enums.ERule(req.Trigger),
		TriggerValue:   req.TriggerValue,
	}
}

func (s *RuleService) EntityToPb(entity *domain.RuleEntity) *crmpb.RuleDTO {
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
