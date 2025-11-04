package usecase

import (
	base_enum "base/enum"
	_enum "common/domain/enum"
	_provider "common/domain/provider"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"strconv"
)

type FollowUsecase struct {
	followRepo         repo.FollowRepo
	contactRepo        repo.ContactRepo
	blockRepo          repo.BlockRepo
	blockService       *BlockUsecase
	userClient         provider.UserClient
	notificationClient _provider.NotificationProvider
}

func NewFollowUsecase(followRepo repo.FollowRepo,
	contactRepo repo.ContactRepo,
	blockRepo repo.BlockRepo,
	blockService *BlockUsecase,
	client provider.UserClient,
	notificationClient _provider.NotificationProvider,
) *FollowUsecase {
	return &FollowUsecase{
		followRepo:         followRepo,
		contactRepo:        contactRepo,
		blockRepo:          blockRepo,
		blockService:       blockService,
		userClient:         client,
		notificationClient: notificationClient,
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
	if err != nil {
		return nil, err
	}

	// Background tasks - Gửi notification cho người được theo dõi
	ctxClone := _utils.CloneContext(c)
	go func() {
		currentUser, _ := s.userClient.GetProfileById(ctxClone, profileId)
		if currentUser == nil {
			return
		}

		// Lấy thông tin contact để có targetId
		_, err := s.contactRepo.GetByProfileID(ctxClone, followId, profileId, base_enum.EOwnerOfMember)
		if err != nil {
			// Nếu chưa có contact thì tạo mới
			targetUser, _ := s.userClient.GetProfileById(ctxClone, followId)
			if targetUser != nil {
				_, _ = s.contactRepo.GetOrCreateByProfileID(
					ctxClone,
					followId,
					profileId,
					base_enum.EOwnerOfMember,
					targetUser.FullName,
					targetUser.Phone,
					targetUser.Avatar,
				)
			}
		}

		s.notificationClient.CreateNotification(
			ctxClone,
			currentUser.Avatar,
			"Theo dõi",
			[]string{"", currentUser.FullName, " đã theo dõi bạn"},
			_enum.NotificationContactFollow,
			&profileId,
			followId,
			_enum.EOwnerOfMember,
			[]string{strconv.FormatUint(followId, 10)},
		)
	}()

	return e, err
}

func (s *FollowUsecase) UnfollowUser(c context.Context, followId uint64) error {
	profileId, followId, err := s.ValidateFollowId(c, followId)
	if err != nil {
		return err
	}

	err = s.followRepo.UnfollowUser(c, profileId, followId)
	if err != nil {
		return err
	}
	go func() {
		ctxClone := _utils.CloneContext(c)
		s.notificationClient.RemoveNotification(
			ctxClone,
			_enum.EOwnerOfMember,
			followId,
			profileId,
			_enum.NotificationContactFollow,
		)
	}()

	return nil
}
