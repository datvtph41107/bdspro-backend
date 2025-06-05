package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
)

type RuleUsecase struct {
	ruleRepo repo.RuleRepo
}

func NewRuleUsecase(ruleRepo repo.RuleRepo) *RuleUsecase {
	return &RuleUsecase{ruleRepo: ruleRepo}
}

func (u *RuleUsecase) Search(c context.Context, ownerId uint64, ownerType enums.EOwnerType, dto dto.RuleSearchDTO) ([]domain.RuleEntity, int64, error) {
	return u.ruleRepo.Search(c, ownerId, ownerType, dto)
}

func (u *RuleUsecase) Create(c context.Context, entity *domain.RuleEntity) (*domain.RuleEntity, error) {
	return u.ruleRepo.Create(c, entity)
}

func (u *RuleUsecase) Update(c context.Context, id uint64, entity *domain.RuleEntity) (*domain.RuleEntity, error) {
	return u.ruleRepo.Update(c, entity)
}

func (u *RuleUsecase) Delete(c context.Context, id uint64) error {
	return u.ruleRepo.Delete(c, id)
}
