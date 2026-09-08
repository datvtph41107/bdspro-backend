package usecase

import (
	_utils "common/utils"
	"context"
	"errors"

	"gorm.io/datatypes"
	"notification/internal/domain"
	"notification/internal/dto"
)

type HistoryAuthUsecase struct {
	historyAuthRepo HistoryAuthStore
}

func NewHistoryAuthUsecase(historyAuthRepo HistoryAuthStore) *HistoryAuthUsecase {
	return &HistoryAuthUsecase{
		historyAuthRepo: historyAuthRepo,
	}
}

func (u *HistoryAuthUsecase) CreateHistoryAuth(ctx context.Context, payload *dto.HistoryAuthCreateDTO) (*domain.HistoryAuthEntity, error) {
	if payload == nil {
		return nil, errors.New("history auth payload is required")
	}

	if payload.UserID == 0 {
		return nil, errors.New("userId is required")
	}

	if payload.ActionType == "" {
		return nil, errors.New("actionType is required")
	}

	if payload.ActionName == "" {
		return nil, errors.New("actionName is required")
	}

	entity := &domain.HistoryAuthEntity{
		UserID:         payload.UserID,
		OrganizationID: payload.OrganizationID,
		ActionType:     payload.ActionType,
		ActionName:     payload.ActionName,
		Description:    payload.Description,
		Success:        true,
		Reason:         payload.Reason,
		IPAddress:      payload.IPAddress,
		UserAgent:      payload.UserAgent,
		PerformedBy:    payload.PerformedBy,
		SessionID:      payload.SessionID,
		Channel:        payload.Channel,
		DeviceID:       payload.DeviceID,
		Location:       payload.Location,
		AdditionalNote: payload.AdditionalNote,
		SourceService:  payload.SourceService,
	}

	if payload.Success != nil {
		entity.Success = *payload.Success
	}

	if payload.Metadata != nil {
		entity.Metadata = datatypes.JSONMap(payload.Metadata)
	}

	if entity.SourceService == "" {
		entity.SourceService = "auth-service"
	}

	if err := u.historyAuthRepo.CreateHistoryAuth(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (u *HistoryAuthUsecase) SearchHistoryAuth(ctx context.Context, payload dto.HistoryAuthSearchDTO) ([]domain.HistoryAuthEntity, int64, error) {
	userID := _utils.GetProfileIdWithContext(ctx)

	searchParams := HistoryAuthSearch{
		UserID:         userID,
		OrganizationID: payload.OrganizationID,
		ActionType:     payload.ActionType,
		Channel:        payload.Channel,
		DeviceID:       payload.DeviceID,
		Success:        payload.Success,
		FromDate:       payload.FromDate,
		ToDate:         payload.ToDate,
		Offset:         payload.GetOffset(),
		Limit:          payload.GetLimit(),
	}

	return u.historyAuthRepo.SearchHistoryAuth(ctx, searchParams)
}
