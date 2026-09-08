package usecase

import (
	base_enum "base/enum"
	_enum "common/domain/enum"
	_errors "common/errors"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"database/sql"
	"log"
	"strconv"
	"time"
)

type FollowUsecase struct {
	followRepo         repo.FollowRepo
	contactRepo        repo.ContactRepo
	blockRepo          repo.BlockRepo
	blockService       *BlockUsecase
	userClient         provider.UserClient
	notificationClient provider.NotificationProvider
	transaction        provider.ITransaction
}

func NewFollowUsecase(followRepo repo.FollowRepo,
	contactRepo repo.ContactRepo,
	blockRepo repo.BlockRepo,
	blockService *BlockUsecase,
	client provider.UserClient,
	transaction provider.ITransaction,
	notificationClient provider.NotificationProvider,
) *FollowUsecase {
	return &FollowUsecase{
		followRepo:         followRepo,
		contactRepo:        contactRepo,
		blockRepo:          blockRepo,
		blockService:       blockService,
		userClient:         client,
		notificationClient: notificationClient,
		transaction:        transaction,
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

func (s *FollowUsecase) ValidateFollowId(ctx context.Context, followId uint64) (uint64, uint64, error) {
	err := s.blockService.BeforeRequest(ctx, followId)
	if err != nil {
		return 0, 0, err
	}

	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return 0, 0, _errors.UnauthorizedException("User not authenticated")
	}

	if followId == 0 {
		return 0, 0, _errors.BadRequestException("Following ID is required")
	}

	return profileId, followId, nil
}

func (s *FollowUsecase) FollowUser(ctx context.Context, followId uint64) (*domain.FollowEntity, error) {
	var result *domain.FollowEntity
	var profileId uint64

	err := s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// Lấy profileId của người follow từ context
		profileId = _utils.GetProfileIdWithContext(txCtx)
		if profileId == 0 {
			return _errors.UnauthorizedException("missing profile id")
		}

		// Kiểm tra followId hợp lệ (người được follow tồn tại)
		if err := s.validateUserExists(txCtx, followId); err != nil {
			return err
		}

		// Không tự follow chính mình
		if profileId == followId {
			return _errors.BadRequestException("cannot follow yourself")
		}

		// Kiểm tra follow đã tồn tại chưa
		existingFollow, err := s.followRepo.GetFollow(txCtx, profileId, followId)
		if err != nil && err != sql.ErrNoRows {
			return _errors.InternalServerException("failed to check follow: " + err.Error())
		}

		// Nếu đã follow active -> return luôn
		if existingFollow != nil && existingFollow.Status == 1 {
			result = existingFollow
			return nil
		}

		// Nếu đã follow nhưng đã unfollow (status = 2) -> cập nhật lại status
		if existingFollow != nil && existingFollow.Status == 2 {
			existingFollow.Status = 1
			if err := s.followRepo.Update(txCtx, existingFollow); err != nil {
				return _errors.InternalServerException("failed to update follow: " + err.Error())
			}
			result = existingFollow
			return nil
		}

		// Tạo follow mới
		newFollow, err := s.followRepo.CreateFollow(txCtx, profileId, followId)
		if err != nil {
			return _errors.InternalServerException("failed to create follow: " + err.Error())
		}
		result = newFollow

		if err := s.hasContact(txCtx, profileId, followId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sau khi commit thành công, chạy async tasks
	s.sendNotificationAsync(ctx, profileId, followId)

	return result, nil
}

func (s *FollowUsecase) hasContact(ctx context.Context, ownerID, profileID uint64) error {
	targetUser, err := s.userClient.GetProfileById(ctx, profileID)
	if err != nil {
		return _errors.InternalServerException("failed to get target user: " + err.Error())
	}
	if targetUser == nil {
		return _errors.NotFoundException("target user not found")
	}

	existingContact, err := s.contactRepo.GetByProfileID(ctx, profileID, ownerID, base_enum.EOwnerOfMember)
	if err != nil && err != sql.ErrNoRows {
		return _errors.InternalServerException("failed to check contact: " + err.Error())
	}

	if existingContact != nil {
		return nil // đã có contact
	}

	// Tạo contact mới
	_, err = s.contactRepo.Create(ctx, &domain.ContactEntity{
		ProfileID: &profileID,
		OwnerID:   ownerID,
		OwnerOf:   base_enum.EOwnerOfMember,
		FullName:  targetUser.FullName,
		Phone:     targetUser.Phone,
		Avatar:    targetUser.Avatar,
		// Các trường khác có thể set mặc định
	})
	return err
}

func (s *FollowUsecase) validateUserExists(ctx context.Context, userID uint64) error {
	// Gọi user service để kiểm tra user tồn tại
	user, err := s.userClient.GetProfileById(ctx, userID)
	if err != nil {
		return _errors.InternalServerException("failed to validate user: " + err.Error())
	}
	if user == nil {
		return _errors.NotFoundException("user not found")
	}
	return nil
}

func (s *FollowUsecase) sendNotificationAsync(ctx context.Context, profileId, followId uint64) {
	go func() {
		asyncCtx, cancel := context.WithTimeout(_utils.CloneContext(ctx), 10*time.Second)
		defer cancel()

		currentUser, err := s.userClient.GetProfileById(asyncCtx, profileId)
		if err != nil || currentUser == nil {
			log.Printf("follow notification: get current user failed profile_id=%d: %v", profileId, err)
			return
		}

		err = s.notificationClient.CreateNotification(
			asyncCtx,
			currentUser.Avatar,
			"Theo dõi mới",
			[]string{"", currentUser.FullName, " đã theo dõi bạn"},
			_enum.NotificationContactFollow,
			&profileId,
			followId,
			_enum.EOwnerOfMember,
			[]string{strconv.FormatUint(followId, 10)},
		)
		if err != nil {
			log.Printf("follow notification: create failed profile_id=%d follow_id=%d: %v", profileId, followId, err)
		}
	}()
}

func (s *FollowUsecase) UnfollowUser(ctx context.Context, followId uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return _errors.UnauthorizedException("missing profile id")
	}

	// Có thể dùng transaction nếu cần xóa contact hoặc các logic khác
	err := s.followRepo.UnfollowUser(ctx, profileId, followId)
	if err != nil {
		return err
	}

	// Async remove notification
	go func() {
		asyncCtx, cancel := context.WithTimeout(_utils.CloneContext(ctx), 5*time.Second)
		defer cancel()
		_ = s.notificationClient.RemoveNotification(
			asyncCtx,
			_enum.EOwnerOfMember,
			followId,
			profileId,
			_enum.NotificationContactFollow,
		)
	}()

	return nil
}
