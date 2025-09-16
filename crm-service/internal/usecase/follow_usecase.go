package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type FollowUsecase struct {
	followRepo   repo.FollowRepo
	contactRepo  repo.ContactRepo
	blockRepo    repo.BlockRepo
	blockService *BlockUsecase
	userClient   provider.UserClient
}

func NewFollowUsecase(followRepo repo.FollowRepo,
	contactRepo repo.ContactRepo,
	blockRepo repo.BlockRepo,
	blockService *BlockUsecase,
	client provider.UserClient,
) *FollowUsecase {
	return &FollowUsecase{
		followRepo:   followRepo,
		contactRepo:  contactRepo,
		blockRepo:    blockRepo,
		blockService: blockService,
		userClient:   client,
	}
}

func (s *FollowUsecase) FollowerUser(c context.Context) ([]domain.Profile, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.followRepo.FollowerUser(c, profileId, 0, 10)
}

func (s *FollowUsecase) FollowingUser(c context.Context) ([]uint64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.followRepo.FollowingUser(c, profileId, 0, 10)
}

func (s *FollowUsecase) Existed(c context.Context, followId uint64) error {
	ok, err := s.contactRepo.ExistByProfile(c, followId)
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

func (s *FollowUsecase) ValidateFollowId(c context.Context, followId uint64) (uint64, uint64, error) {
	err := s.blockService.BeforeRequest(c, followId)
	if err != nil {
		return 0, 0, err
	}
	// err := s.contactRepo.Existed(c, followId)
	// if err != nil {
	// 	return 0, 0, err
	// }
	// profileId := _utils.GetProfileIdWithContext(c)
	profileId := _utils.GetProfileIdWithContext(c)
	return profileId, followId, err
}

func (s *FollowUsecase) FollowUser(c context.Context, followId uint64) (*domain.FollowEntity, error) {
	profileId, followId, err := s.ValidateFollowId(c, followId)
	if err != nil {
		return nil, err
	}

	e, err := s.followRepo.FollowUser(c, profileId, followId)
	return e, err
}

func (s *FollowUsecase) UnfollowUser(c context.Context, followId uint64) error {
	profileId, followId, err := s.ValidateFollowId(c, followId)
	if err != nil {
		return err
	}

	err = s.followRepo.UnfollowUser(c, profileId, followId)
	return err
}
