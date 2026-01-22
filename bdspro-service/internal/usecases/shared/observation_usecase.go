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
	// clientID := utils.GetClientID(ctx)
	// sessionID := _utils.GetSessionIdFromContext(ctx)

	o := &domain.Observation{
		EntityID: req.EntityID,
		// EntityType: req.EntityType,
		Action:    req.Action,
		Duration:  req.Duration,
		ProfileID: &profileID,
		// ClientID:   clientID,
		// SessionID: sessionID,
	}

	return u.ObservationRepo.Insert(ctx, o)
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
