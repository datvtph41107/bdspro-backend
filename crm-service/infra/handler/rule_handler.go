package handler

import (
	base_enum "base/enum"
	_dto "common/domain/dto"
	_routes "common/routes"
	"context"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"crm/infra/client"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"
)

type RuleService struct {
	crmpb.UnimplementedRuleServiceServer
	UC         *usecase.RuleUsecase
	UserClient *client.UserClient
}

func NewRuleService(uc *usecase.RuleUsecase, userClient *client.UserClient) *RuleService {
	return &RuleService{
		UC:         uc,
		UserClient: userClient,
	}
}

// @Summary Tìm kiếm quy tắc
// @Description Tìm kiếm quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param query query crmpb.RuleSearchRequest false "Giá trị điều kiện"
// @Security BearerAuth
// @Router /rule/list [get]
func (s *RuleService) Search(ctx context.Context, req *crmpb.RuleSearchRequest) (*crmpb.RuleListDTO, error) {
	dto := dto.RuleSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		RuleName: req.Name,
	}

	rules, total, err := s.UC.Search(ctx, req.OwnerId, base_enum.EOwnerOf(req.OwnerType), dto)
	if err != nil {
		return nil, err
	}

	rulesPb := make([]*crmpb.RuleDTO, len(rules))
	for i, rule := range rules {
		rulesPb[i] = mapper.RuleDomainToPb(&rule)
	}

	s.UserClient.MapToRulePb(ctx, rulesPb)

	return &crmpb.RuleListDTO{
		Data:  rulesPb,
		Total: uint32(total),
	}, nil
}

// @Summary Chi tiết quy tắc
// @Description Chi tiết quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /rule/detail/{id} [get]
func (s *RuleService) Detail(ctx context.Context, req *crmpb.RuleDTO) (*crmpb.RuleDTO, error) {
	result, err := s.UC.GetByID(ctx, req.Id)
	return mapper.RuleDomainToPb(result), err
}

// @Summary Tạo quy tắc
// @Description Tạo quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param entity body crmpb.RuleDTO true "Quy tắc"
// @Security BearerAuth
// @Router /rule [post]
func (s *RuleService) Create(ctx context.Context, req *crmpb.RuleDTO) (*crmpb.RuleDTO, error) {
	entity := mapper.RulePbToDomain(req)
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
	return mapper.RuleDomainToPb(result), err
}

// @Summary Cập nhật quy tắc
// @Description Cập nhật quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param entity body crmpb.RuleDTO true "Quy tắc"
// @Security BearerAuth
// @Router /rule/{id} [put]
func (s *RuleService) Update(ctx context.Context, req *crmpb.RuleDTO) (*crmpb.RuleDTO, error) {
	entity := mapper.RulePbToDomain(req)
	if entity.ID == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "id is required",
		}
	}
	result, err := s.UC.Update(ctx, entity.ID, entity)
	return mapper.RuleDomainToPb(result), err
}

// @Summary Xóa quy tắc
// @Description Xóa quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /rule/{id} [delete]
func (s *RuleService) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "id is required",
		}
	}
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Delete rule success",
	}, err
}

// @Summary Active quy tắc
// @Description Active quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param active body crmpb.ActiveRuleRequest false "Active"
// @Security BearerAuth
// @Router /rule/{id}/active [put]
func (s *RuleService) Active(ctx context.Context, req *crmpb.ActiveRuleRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "id is required",
		}
	}
	err := s.UC.Active(ctx, req.Id, req.Active)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Active rule success",
	}, err
}