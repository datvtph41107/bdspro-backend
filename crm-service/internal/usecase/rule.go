package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
)

type RuleUsecase struct {
	ruleRepo          repo.RuleRepo
	permissionUsecase PermissionUsecase
	historyRepo       repo.HistoryRepo
	customerRepo      repo.LeadRepo
}

func NewRuleUsecase(
	ruleRepo repo.RuleRepo,
	permissionUsecase PermissionUsecase,
	historyRepo repo.HistoryRepo,
	customerRepo repo.LeadRepo,
) *RuleUsecase {
	return &RuleUsecase{
		ruleRepo:          ruleRepo,
		permissionUsecase: permissionUsecase,
		historyRepo:       historyRepo,
		customerRepo:      customerRepo,
	}
}

func (u *RuleUsecase) Search(c context.Context, ownerId uint64, ownerType enums.EOwnerType, dto dto.RuleSearchDTO) ([]domain.RuleEntity, int64, error) {
	if err := u.permissionUsecase.UserInOwner(c, ownerId, ownerType); err != nil {
		return nil, 0, err
	}
	return u.ruleRepo.Search(c, ownerId, ownerType, dto)
}

func (u *RuleUsecase) Create(c context.Context, entity *domain.RuleEntity) (*domain.RuleEntity, error) {
	if err := u.permissionUsecase.UserInOwner(c, entity.OwnerID, entity.OwnerType); err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(c, &domain.HistoryEntity{
		Note:       "Tạo quy tắc " + entity.RuleName,
		TargetId:   entity.ID,
		TargetType: enums.TargetHistoryRule,
		ActionType: enums.HistoryRuleCreate,
		OwnerID:    &entity.OwnerID,
		OwnerType:  entity.OwnerType,
	})

	return u.ruleRepo.Create(c, entity)
}

func (u *RuleUsecase) Update(c context.Context, id uint64, entity *domain.RuleEntity) (*domain.RuleEntity, error) {
	if err := u.permissionUsecase.UserInOwner(c, entity.OwnerID, entity.OwnerType); err != nil {
		return nil, err
	}

	entity, err := u.ruleRepo.Update(c, entity)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(c, &domain.HistoryEntity{
		Note:       "Cập nhật quy tắc " + entity.RuleName + " - " + entity.Note,
		TargetId:   entity.ID,
		TargetType: enums.TargetHistoryRule,
		ActionType: enums.HistoryRuleUpdate,
		OwnerID:    &entity.OwnerID,
		OwnerType:  entity.OwnerType,
	})

	return entity, nil
}

func (u *RuleUsecase) Delete(c context.Context, id uint64) error {
	entity, err := u.ruleRepo.GetByID(c, id)
	if err != nil {
		return err
	}
	if err := u.permissionUsecase.UserInOwner(c, entity.OwnerID, entity.OwnerType); err != nil {
		return err
	}

	err = u.ruleRepo.Delete(c, id)
	if err != nil {
		return err
	}

	u.historyRepo.CreateHistory(c, &domain.HistoryEntity{
		Note:       "Xóa quy tắc " + entity.RuleName,
		TargetId:   entity.ID,
		TargetType: enums.TargetHistoryRule,
		ActionType: enums.HistoryRuleDelete,
		OwnerID:    &entity.OwnerID,
		OwnerType:  entity.OwnerType,
	})

	return nil
}

func (u *RuleUsecase) GetByID(c context.Context, id uint64) (*domain.RuleEntity, error) {
	return u.ruleRepo.GetByID(c, id)
}
