package usecase

import (
	"context"
	"notification/internal/domain"
	"notification/internal/dto"
		shared_dto "pb/dto"
	shared_enum "pb/enums"
)

type HistoryUsecase struct {
	historyRepo HistoryStore
}

func NewHistoryUsecase(historyRepo HistoryStore) *HistoryUsecase {
	return &HistoryUsecase{
		historyRepo: historyRepo,
	}
}

func (h *HistoryUsecase) UserInOwner(c context.Context, ownerId uint64, ownerType shared_enum.EOwnerType) error {

	return nil
}

func (h *HistoryUsecase) Search(c context.Context, ownerId uint64, ownerType shared_enum.EOwnerType, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	err := h.UserInOwner(c, ownerId, ownerType)
	if err != nil {
		return nil, 0, err
	}

	histories, total, err := h.historyRepo.Search(c, ownerId, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (h *HistoryUsecase) CrmHistory(c context.Context, crmId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	histories, total, err := h.historyRepo.CrmHistory(c, crmId, dto)
	if err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}

func (h *HistoryUsecase) ContactHistory(c context.Context, contactId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	// contact, err := h.contactRepo.GetById(c, contactId)
	// if err != nil {
	// 	return nil, 0, err
	// }

	// if contact.OwnerID != dto.UserID {
	// 	return nil, 0, errors.New("user not in contact")
	// }
	// todo: check permission

	histories, total, err := h.historyRepo.ContactHistory(c, contactId, dto)
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (h *HistoryUsecase) AssetHistory(c context.Context, assetId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	histories, total, err := h.historyRepo.AssetHistory(c, assetId, dto)
	if err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}

func (h *HistoryUsecase) ProductHistory(c context.Context, productId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	histories, total, err := h.historyRepo.ProductHistory(c, productId, dto)
	if err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}

func (h *HistoryUsecase) ProductChildHistory(c context.Context, productId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	histories, total, err := h.historyRepo.ProductChildHistory(c, productId, dto)
	if err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}

func (h *HistoryUsecase) RuleEventHistory(c context.Context, ruleEventId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	history, total, err := h.historyRepo.RuleEventHistory(c, ruleEventId, dto)
	if err != nil {
		return nil, 0, err
	}
	return history, total, nil
}

func (h *HistoryUsecase) Create(c context.Context, history *shared_dto.HistoryDTO) (*domain.HistoryEntity, error) {
	historyEntity := &domain.HistoryEntity{
		ActionType: history.ActionType,
		TargetId:   history.TargetId,
		OwnerID:    history.OwnerID,
		Note:       history.Note,
		Title:      history.Title,
		// OwnerType:  history.OwnerType,
	}
	err := h.historyRepo.CreateHistory(c, historyEntity)
	if err != nil {
		return nil, err
	}
	return historyEntity, nil
}

func (h *HistoryUsecase) CreateInternal(c context.Context, history *shared_dto.HistoryDTO) (*domain.HistoryEntity, error) {
	historyEntity := &domain.HistoryEntity{
		ActionType: history.ActionType,
		TargetId:   history.TargetId,
		OwnerID:    history.OwnerID,
		Note:       history.Note,
		Title:      history.Title,
		IsInternal: true,
	}
	err := h.historyRepo.CreateHistory(c, historyEntity)
	if err != nil {
		return nil, err
	}
	return historyEntity, nil
}

func (h *HistoryUsecase) SearchInternal(c context.Context, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	histories, total, err := h.historyRepo.SearchInternal(c, dto)
	if err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}

func (h *HistoryUsecase) CampaignHistory(c context.Context, campaignId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	histories, total, err := h.historyRepo.CampaignHistory(c, campaignId, dto)
	if err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}

func (h *HistoryUsecase) PackageHistory(c context.Context, packageId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	histories, total, err := h.historyRepo.PackageHistory(c, packageId, dto)
	if err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}
