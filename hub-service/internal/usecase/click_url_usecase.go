package usecase

import (
	"context"

	_err "common/domain/err"
	_utils "common/utils"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
)

type IClickURLUsecase interface {
	TrackBatch(ctx context.Context, req *dto.TrackClickURLBatchRequest) *_err.ErrorDTO
}

type ClickURLUsecase struct {
	ClickURLRepo repo.IClickURLRepo
}

// @bind: internal/usecase/IClickURLUsecase
func NewClickURLUsecase(
	clickURLRepo repo.IClickURLRepo,
) IClickURLUsecase {
	return &ClickURLUsecase{
		ClickURLRepo: clickURLRepo,
	}
}

func (u *ClickURLUsecase) TrackBatch(
	ctx context.Context,
	req *dto.TrackClickURLBatchRequest,
) *_err.ErrorDTO {
	if req == nil || len(req.Events) == 0 {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "Request không hợp lệ hoặc danh sách events rỗng",
		}
	}

	var profileID *uint64
	if profileIDVal := _utils.GetProfileIdWithContext(ctx); profileIDVal != 0 {
		profileID = &profileIDVal
	}

	events := make([]*domain.ClickURL, 0, len(req.Events))
	for _, dtoReq := range req.Events {
		if !dtoReq.DomainType.Validate() {
			return &_err.ErrorDTO{
				Code:    400,
				Message: "domain_type không hợp lệ",
			}
		}

		events = append(events, &domain.ClickURL{
			EventTime:  dtoReq.EventTime,
			DomainID:   dtoReq.DomainID,
			DomainType: uint32(dtoReq.DomainType),
			ProfileID:  profileID,
		})
	}

	if err := u.ClickURLRepo.InsertBatch(ctx, events); err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lưu batch click url: " + err.Error(),
		}
	}

	return nil
}
