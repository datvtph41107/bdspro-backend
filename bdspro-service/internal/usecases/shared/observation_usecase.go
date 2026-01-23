package shared_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	_utils "common/utils"
	"context"
)

type ObservationUsecase struct {
	ObservationRepo repo.ObservationRepository
}

func NewObservationUsecase(
	observationRepo repo.ObservationRepository,
) *ObservationUsecase {
	return &ObservationUsecase{
		ObservationRepo: observationRepo,
	}
}

func (u *ObservationUsecase) Track(
	ctx context.Context,
	req *dto.TrackObservationRequest,
) error {
	profileID := _utils.GetProfileIdWithContext(ctx)
	// clientID := _utils.GetClientIDFromContext(ctx)
	// sessionID := _utils.GetSessionIdFromContext(ctx)

	for _, event := range req.Events {
		duration := event.Duration
		if duration == 0 && event.EndTime > 0 && event.StartTime > 0 {
			duration = uint32((event.EndTime - event.StartTime) / 1000)
		}

		o := &domain.Observation{
			EntityType: &event.Event,
			EntityID:   event.RefId,
			ProfileID:  &profileID,
			Action:     event.Event,
			Duration:   duration,
			StartTime:  event.StartTime,
			EndTime:    event.EndTime,
			Screen:     &event.Screen,
		}

		err := u.ObservationRepo.Insert(ctx, o)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *ObservationUsecase) GetProductStats(
	ctx context.Context,
	productID uint64,
) (*dto.ProductStats, error) {
	// profileID := _utils.GetProfileIdWithContext(ctx)

	// ok, _ := u.productUserRepo.CheckAccess(
	// 	ctx, profileID, productID,
	// )

	// if !ok {
	// 	return nil, _errors.ForbiddenException()
	// }

	return u.ObservationRepo.GetProductStats(ctx, productID)
}
