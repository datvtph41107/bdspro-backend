package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
)

type HistoryUsecase struct {
	historyRepo       repo.HistoryRepo
	permissionUsecase PermissionUsecase
}

func NewHistoryUsecase(historyRepo repo.HistoryRepo, permissionUsecase PermissionUsecase) *HistoryUsecase {
	return &HistoryUsecase{
		historyRepo:       historyRepo,
		permissionUsecase: permissionUsecase,
	}
}

func (h *HistoryUsecase) Search(c context.Context, ownerId uint64, ownerType enums.EOwnerType, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	err := h.permissionUsecase.UserInOwner(c, ownerId, ownerType)
	if err != nil {
		return nil, 0, err
	}

	histories, total, err := h.historyRepo.Search(c, ownerId, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}
