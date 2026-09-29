package usecase

import (
	_errors "common/errors"
	"context"
	"fmt"
	"notification/internal"
	"notification/internal/domain"
	"notification/internal/dto"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"time"
)

type PropertyHistoryUseCase struct {
	PropertyHistoryRepo PropertyHistoryStore
	UserProvider        UserProfileReader
}

func NewPropertyHistoryUseCase(
	propertyHistoryRepo PropertyHistoryStore,
	userProvider UserProfileReader,
) *PropertyHistoryUseCase {
	return &PropertyHistoryUseCase{
		PropertyHistoryRepo: propertyHistoryRepo,
		UserProvider:        userProvider,
	}
}

func (uc *PropertyHistoryUseCase) Create(
	ctx context.Context,
	dto *dto.CreatePropertyHistoryDTO,
) (*domain.PropertyHistory, error) {
	now := time.Now()
	entity := &domain.PropertyHistory{
		SubjectID: dto.SubjectID,
		// SubjectType: dto.SubjectType,
		Action:      dto.Action,
		Description: dto.Description,
		Metadata:    dto.Metadata,
		ActorID:     dto.ActorID,
		OccurredAt:  &now,
	}

	return uc.PropertyHistoryRepo.Create(ctx, entity)
}

func (uc *PropertyHistoryUseCase) GetDetail(
	ctx context.Context,
	id uint64,
) (*domain.PropertyHistory, *userpb.ProfileResponse, error) {
	history, err := uc.PropertyHistoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if history == nil {
		return nil, nil, _errors.ReturnError(service.PropertyHistoryNotFound)
	}

	actor, err := uc.UserProvider.GetProfileByID(
		ctx,
		&userpb.GetProfileByIDRequest{ProfileId: history.ActorID},
	)
	if err != nil {
		return nil, nil,
			fmt.Errorf("get property history actor profile: %w", err)
	}

	return history, actor, nil
}

func (uc *PropertyHistoryUseCase) SoftDelete(
	ctx context.Context,
	id uint64,
	actorID uint64,
) error {

	history, err := uc.PropertyHistoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if history == nil {
		return _errors.ReturnError(service.PropertyHistoryNotFound)
	}

	return uc.PropertyHistoryRepo.Delete(ctx, id)
}

func (uc *PropertyHistoryUseCase) Search(
	ctx context.Context,
	q *dto.PropertyHistoryQueryDTO,
) ([]*domain.PropertyHistory, map[uint64]*sharepb.ProfileItem, int64, error) {
	histories, total, err := uc.PropertyHistoryRepo.Search(ctx, q)
	if err != nil {
		return nil, nil, 0,
			fmt.Errorf("search property history: %w", err)
	}

	if len(histories) == 0 {
		return histories, map[uint64]*sharepb.ProfileItem{}, total, nil
	}

	userIDs := make([]uint64, 0, len(histories))
	for _, h := range histories {
		userIDs = append(userIDs, h.ActorID)
	}

	userMap, err := uc.UserProvider.GetMapProfileByIds(
		ctx,
		&sharepb.GetProfileByIdsRequest{Ids: userIDs},
	)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("get property history profile map: %w", err)
	}
	if len(userMap) == 0 {
		userMap = map[uint64]*sharepb.ProfileItem{}
	}

	return histories, userMap, total, nil
}
