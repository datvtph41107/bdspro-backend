package usecases

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"userv4/internal/interface/repo"
	"userv4/internal/models"
)

type FollowUsecase struct {
	followRepo   repo.IFollowRepo
	contactRepo  repo.IContactRepo
	blockRepo    repo.IBlockRepo
	blockUsecase *BlockUsecase
}

func NewFollowUsecase(repo repo.IFollowRepo,
	contactRepo repo.IContactRepo,
	blockRepo repo.IBlockRepo,
	blockUsecase *BlockUsecase,
) *FollowUsecase {
	return &FollowUsecase{
		followRepo:   repo,
		contactRepo:  contactRepo,
		blockRepo:    blockRepo,
		blockUsecase: blockUsecase,
	}
}

func (u *FollowUsecase) FollowerUser(ctx context.Context) ([]models.Profile, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	return u.followRepo.FollowerUser(profileId, 0, 10)
}

func (u *FollowUsecase) FollowingUser(ctx context.Context) ([]models.Profile, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	return u.followRepo.FollowingUser(profileId, 0, 10)
}

func (u *FollowUsecase) Existed(followId uint64) error {
	ok, err := u.contactRepo.ExistByProfile(followId)
	if !ok {
		return &_routes.Except{
			Code:    400,
			Message: "Người dùng không tồn tại",
		}
	}

	if err != nil {
		return &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}
	return nil
}

func (u *FollowUsecase) GetFollowId(ctx context.Context, followId uint64) (uint64, uint64, error) {
	err := u.contactRepo.Existed(followId)
	if err != nil {
		return 0, 0, err
	}
	profileId := _utils.GetProfileIdWithContext(ctx)

	return profileId, followId, err
}

func (u *FollowUsecase) FollowUser(ctx context.Context, followId uint64) (*models.FollowEntity, error) {
	profileId, followId, err := u.GetFollowId(ctx, followId)
	if err != nil {
		return nil, err
	}

	err = u.blockUsecase.BeforeRequest(ctx, followId)
	if err != nil {
		return nil, err
	}

	e, err := u.followRepo.FollowUser(ctx, profileId, followId)
	return e, err
}

func (u *FollowUsecase) UnfollowUser(ctx context.Context, followId uint64) (*models.FollowEntity, error) {
	profileId, followId, err := u.GetFollowId(ctx, followId)
	if err != nil {
		return nil, err
	}

	err = u.blockUsecase.BeforeRequest(ctx, followId)
	if err != nil {
		return nil, err
	}

	err = u.followRepo.UnfollowUser(ctx, profileId, followId)
	return nil, err
}
