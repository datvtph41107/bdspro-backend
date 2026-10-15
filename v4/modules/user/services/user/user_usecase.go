package services

import (
	"context"
	userpostgres "userv4/modules/user/postgres"
	"userv4/modules/user/domain/user"
)

type UserUsecase interface {
	CountCurrent(ctx context.Context) (int64, error)
}

type userUsecase struct {
	profileRepo *userpostgres.ProfilePostgres
}

func NewUserUsecase(profileRepo *userpostgres.ProfilePostgres) UserUsecase {
	return &userUsecase{
		profileRepo: profileRepo,
	}
}

func (uc *userUsecase) CountCurrent(ctx context.Context) (int64, error) {
	var count int64
	err := uc.profileRepo.DB.WithContext(ctx).
		Model(&user.UserProfileEntity{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
