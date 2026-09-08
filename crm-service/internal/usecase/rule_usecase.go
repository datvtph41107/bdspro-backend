package usecase

import (
	base_enum "base/enum"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type RuleUsecase struct {
	ruleRepo           repo.RuleRepo
	permissionUsecase  PermissionUsecase
	notificationClient provider.NotificationProvider
	customerRepo       repo.LeadRepo
}

func NewRuleUsecase(
	ruleRepo repo.RuleRepo,
	permissionUsecase PermissionUsecase,
	notificationClient provider.NotificationProvider,
	customerRepo repo.LeadRepo,
) *RuleUsecase {
	return &RuleUsecase{
		ruleRepo:           ruleRepo,
		permissionUsecase:  permissionUsecase,
		notificationClient: notificationClient,
		customerRepo:       customerRepo,
	}
}

func (u *RuleUsecase) Search(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.RuleSearchDTO) ([]domain.RuleEntity, int64, error) {
	if err := u.permissionUsecase.UserInOwner(c, ownerId, ownerType); err != nil {
		return nil, 0, err
	}
	return u.ruleRepo.Search(c, ownerId, ownerType, dto)
}

func (u *RuleUsecase) Create(c context.Context, entity *domain.RuleEntity) (*domain.RuleEntity, error) {
	if err := u.permissionUsecase.UserInOwner(c, entity.OwnerID, entity.OwnerType); err != nil {
		return nil, err
	}

	// u.notificationClient.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	Title:      "Tạo quy tắc thành công",
	// 	Note:       []string{"Tạo quy tắc", entity.RuleName},
	// 	TargetId:   entity.ID,
	// 	TargetType: base_enum.TargetHistoryRule,
	// 	ActionType: base_enum.HistoryRuleCreate,
	// 	OwnerID:    &entity.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(entity.OwnerType),
	// })

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

	// u.notificationClient.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	Title:      "Cập nhật quy tắc thành công",
	// 	Note:       []string{"Cập nhật quy tắc", entity.RuleName, " - ", entity.Note},
	// 	TargetId:   entity.ID,
	// 	TargetType: base_enum.TargetHistoryRule,
	// 	ActionType: base_enum.HistoryRuleUpdate,
	// 	OwnerID:    &entity.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(entity.OwnerType),
	// })

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

	// u.notificationClient.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	Title:      "Xóa quy tắc thành công",
	// 	Note:       []string{"Xóa quy tắc", entity.RuleName},
	// 	TargetId:   entity.ID,
	// 	TargetType: base_enum.TargetHistoryRule,
	// 	ActionType: base_enum.HistoryRuleDelete,
	// 	OwnerID:    &entity.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(entity.OwnerType),
	// })

	return nil
}

func (u *RuleUsecase) GetByID(c context.Context, id uint64) (*domain.RuleEntity, error) {
	return u.ruleRepo.GetByID(c, id)
}

func (u *RuleUsecase) Active(c context.Context, id uint64, active bool) error {
	return u.ruleRepo.Active(c, id, active)
}
