package user_service

import (
	_utils "common/utils"
	"errors"
	user_repo "user/infra/repo"
	"user/internal/dto"
	user_model "user/internal/models/user"

	"github.com/gin-gonic/gin"
)

type RateService struct {
	RateRepo *user_repo.RateRepo
}

func NewRateService(rateRepo *user_repo.RateRepo) *RateService {
	return &RateService{RateRepo: rateRepo}
}

func (s *RateService) Create(c *gin.Context, rate *user_model.Rate) (*user_model.Rate, error) {
	ctx := _utils.CreateContextFromGin(c)
	profileId := _utils.GetProfileIdWithContext(ctx)

	rate.CreatedBy = &profileId

	if err := s.RateRepo.Create(c, rate); err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *RateService) List(c *gin.Context, targetID uint64, targetType uint8, dto dto.RateSearch) ([]user_model.RateItem, error) {
	rates, err := s.RateRepo.List(c, targetID, targetType, dto)
	if err != nil {
		return nil, err
	}

	for i, rate := range rates {
		if rate.Anonymous {
			rates[i].CreatedUser = nil
		}
	}
	return rates, nil
}

func (s *RateService) Update(c *gin.Context, id uint64, rate *user_model.Rate) (*user_model.Rate, error) {
	ctx := _utils.CreateContextFromGin(c)
	profileId := _utils.GetProfileIdWithContext(ctx)
	rate.UpdatedBy = &profileId

	if rate.CreatedBy != nil && *rate.CreatedBy != profileId {
		return nil, errors.New("you are not allowed to update this rate")
	}

	if err := s.RateRepo.Update(c, id, rate); err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *RateService) Delete(c *gin.Context, id uint64) error {
	ctx := _utils.CreateContextFromGin(c)
	profileId := _utils.GetProfileIdWithContext(ctx)
	return s.RateRepo.Delete(c, profileId, id)
}
