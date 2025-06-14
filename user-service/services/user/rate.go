package user_service

import (
	_jwt "common/jwt"
	"user/data"
	user_model "user/models/user"
	repo "user/repo"

	"github.com/gin-gonic/gin"
)

type RateService struct {
	RateRepo *repo.RateRepo
}

func NewRateService(rateRepo *repo.RateRepo) *RateService {
	return &RateService{RateRepo: rateRepo}
}

func (s *RateService) Create(c *gin.Context, rate *user_model.Rate) (*user_model.Rate, error) {
	profileId := _jwt.GetProfileId(c)

	rate.CreatedBy = &profileId

	if err := s.RateRepo.Create(c, rate); err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *RateService) List(c *gin.Context, targetID uint64, targetType uint8, dto data.RateSearch) ([]user_model.RateItem, error) {
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
	if err := s.RateRepo.Update(c, id, rate); err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *RateService) Delete(c *gin.Context, id uint64) error {
	profileId := _jwt.GetProfileId(c)
	return s.RateRepo.Delete(c, profileId, id)
}
