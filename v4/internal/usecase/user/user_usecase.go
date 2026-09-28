package usecases

import (
	"context"
	user_infra_postgres "userv4/infra/postgres"
	"userv4/internal/domain/user"
)

type UserUsecase interface {
	CountCurrent(ctx context.Context) (int64, error)
}

type userUsecase struct {
	profileRepo *user_infra_postgres.ProfilePostgres
}

func NewUserUsecase(profileRepo *user_infra_postgres.ProfilePostgres) UserUsecase {
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
