package usecase

import (
	base_enum "base/enum"
	_db "common/db"
	_enum "common/domain/enum"
	_errors "common/errors"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"strconv"
	"time"

	"crm/data"
	"crm/internal/domain"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	sharepb "pb/types/shared"
)

type FriendUsecase struct {
	friendRepo   repo.FriendRepo
	contact      repo.ContactRepo
	blockRepo    repo.BlockRepo
	followRepo   repo.FollowRepo
	blockService *BlockUsecase
	userClient   provider.UserClient

	notificationClient provider.NotificationProvider
	transaction        provider.ITransaction
}

func NewFriendUsecase(
	friendRepo repo.FriendRepo,
	contact repo.ContactRepo,
	followRepo repo.FollowRepo,
	blockRepo repo.BlockRepo,
	blockService *BlockUsecase,
	userClient provider.UserClient,

	notificationClient provider.NotificationProvider,
	transaction provider.ITransaction,
) *FriendUsecase {
	return &FriendUsecase{
		friendRepo:         friendRepo,
		contact:            contact,
		blockRepo:          blockRepo,
		blockService:       blockService,
		followRepo:         followRepo,
		notificationClient: notificationClient,
		userClient:         userClient,
		transaction:        transaction,
	}
}

func (s *FriendUsecase) GetByID(c context.Context, id uint64) (*domain.FriendEntity, error) {
	return s.friendRepo.FindByID(c, id)
}

// func (s *FriendService) Search(receiverID uint, createdBy uint, isSender bool, page int, size int) ([]domain.FriendEntity, error) {
// 	if isSender {
// 		receiverID = 0
// 	}
// 	return s.Repo.Search(receiverID, createdBy, page, size)
// }

func (s *FriendUsecase) FriendList(c context.Context, dto data.FriendRequest) ([]domain.FriendEntity, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.friendRepo.SearchFriend(c, profileId, dto)
}

func (s *FriendUsecase) Request(ctx context.Context, receiverId uint64) (*domain.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, _errors.UnauthorizedException("missing profile id")
	}

	if err := s.blockService.BeforeRequest(ctx, receiverId); err != nil {
		return nil, err
	}

	receiverUser, err := s.userClient.GetProfileById(ctx, receiverId)
	if err != nil || receiverUser == nil {
		return nil, _errors.NotFoundException("receiver not found")
	}

	currentUser, _ := s.userClient.GetProfileById(ctx, profileId)

	var result *domain.FriendEntity

	err = s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		existing, err := s.friendRepo.GetFriendship(txCtx, profileId, receiverId)
		if err != nil {
			return _errors.InternalServerException("failed to check friendship: " + err.Error())
		}

		if existing != nil {
			switch existing.Status {
			case enums.FriendStatusAccepted:
				return _errors.ConflictException("already friends")
			case enums.FriendStatusPending:
				if existing.CreatedBy != nil && *existing.CreatedBy == profileId {
					return _errors.ConflictException("friend request already sent")
				} else {
					return _errors.ConflictException("you have a pending request from this user, please accept or reject it")
				}
			case enums.FriendStatusReject:
				// Xóa mềm bản ghi cũ (hủy bỏ) để tạo mới
				if err := s.friendRepo.CancelRequest(txCtx, existing.ID); err != nil {
					return _errors.InternalServerException("failed to delete old rejected request: " + err.Error())
				}
			}
		}

		// Tạo friend request mới
		friendEntity := &domain.FriendEntity{
			Status:     enums.FriendStatusPending,
			ReceiverID: receiverId,
		}
		if err := _db.SaveWithContext(txCtx, friendEntity); err != nil {
			return _errors.InternalServerException("failed to create friend request: " + err.Error())
		}
		result = friendEntity

		// Tạo contact cho người gửi (lưu thông tin người nhận)
		if _, err := s.contact.GetOrCreateByProfileID(txCtx,
			receiverId, profileId, base_enum.EOwnerOfMember,
			receiverUser.FullName, receiverUser.Phone, receiverUser.Avatar); err != nil {
			return _errors.InternalServerException("failed to create contact for sender: " + err.Error())
		}

		// Tạo contact cho người nhận (lưu thông tin người gửi)
		if currentUser != nil {
			if _, err := s.contact.GetOrCreateByProfileID(txCtx,
				profileId, receiverId, base_enum.EOwnerOfMember,
				currentUser.FullName, currentUser.Phone, currentUser.Avatar); err != nil {
				return _errors.InternalServerException("failed to create contact for receiver: " + err.Error())
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Async notification
	s.sendFriendRequestNotificationAsync(ctx, profileId, receiverId, result.ID)
	return result, nil
}

func (s *FriendUsecase) sendFriendRequestNotificationAsync(ctx context.Context, senderId, receiverId, requestId uint64) {
	go func() {
		asyncCtx, cancel := context.WithTimeout(_utils.CloneContext(ctx), 10*time.Second)
		defer cancel()

		currentUser, _ := s.userClient.GetProfileById(asyncCtx, senderId)
		if currentUser == nil {
			return
		}

		contact, err := s.contact.GetByProfileID(asyncCtx, receiverId, senderId, base_enum.EOwnerOfMember)
		attachData := []string{strconv.FormatUint(receiverId, 10)}
		if err == nil && contact != nil {
			attachData = []string{strconv.FormatUint(contact.ID, 10)}
		}

		_ = s.notificationClient.CreateNotification(
			asyncCtx,
			currentUser.Avatar,
			"Lời mời kết bạn",
			[]string{"", currentUser.FullName, " đã gửi lời mời kết bạn cho bạn"},
			_enum.NotificationFriendRequest,
			&senderId,
			receiverId,
			_enum.EOwnerOfMember,
			attachData,
		)
	}()
}

func (s *FriendUsecase) RequestSent(c context.Context, params data.FriendDTO) ([]domain.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	params.Status = enums.FriendStatusPending
	return s.friendRepo.Search(c, params, &profileId, 0, 10)
}

func (s *FriendUsecase) RequestReceived(c context.Context, params data.FriendDTO) ([]domain.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	params.ReceiverID = profileId
	params.Status = enums.FriendStatusPending

	return s.friendRepo.Search(c, params, nil, 0, 10)
}

func (s *FriendUsecase) ChangeStatus(ctx context.Context, requestId uint64, status enums.FriendStatus) (int, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return 0, _errors.UnauthorizedException("missing profile id")
	}

	// Validate block
	if err := s.blockService.BeforeRequest(ctx, requestId); err != nil { // requestId là id của friend request, cần kiểm tra block với người gửi?
		// Ở đây cần lấy receiverId và senderId từ request để kiểm tra block
		// Tạm thời giữ nguyên logic cũ (chỉ check block với requestId? Không đúng)
		// Nên refactor: lấy requestEntity trước, rồi check block giữa profileId và senderId
	}

	var requestEntity *domain.FriendEntity
	err := s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// Lấy request và kiểm tra quyền
		var err error
		requestEntity, err = s.friendRepo.FindByID(txCtx, requestId)
		if err != nil {
			return _errors.NotFoundException("friend request not found")
		}
		if requestEntity.ReceiverID != profileId {
			return _errors.ForbiddenException("you are not the receiver of this request")
		}
		if requestEntity.Status != enums.FriendStatusPending {
			return _errors.ConflictException("request already processed")
		}

		// Kiểm tra block giữa profileId và người gửi (senderId)
		senderId := uint64(0)
		if requestEntity.CreatedBy != nil {
			senderId = *requestEntity.CreatedBy
		}
		if senderId == 0 {
			return _errors.BadRequestException("invalid sender")
		}
		if err := s.blockService.BeforeRequest(txCtx, senderId); err != nil {
			return err
		}

		// Update status
		respondedAt := time.Now()
		if err := s.friendRepo.UpdateStatusById(txCtx, status, requestId, &respondedAt); err != nil {
			return _errors.InternalServerException("failed to update status: " + err.Error())
		}

		// Nếu accept, có thể cập nhật thêm thông tin contact (ví dụ: đánh dấu là bạn bè) nhưng không cần tạo mới
		// Nếu reject, không làm gì thêm với contact (giữ nguyên contact đã tạo từ lúc request)

		return nil
	})

	if err != nil {
		return 0, err
	}

	// Async notification
	s.sendFriendResponseNotificationAsync(ctx, profileId, requestEntity, status)

	return 0, nil
}

func (s *FriendUsecase) sendFriendResponseNotificationAsync(ctx context.Context, responderId uint64, requestEntity *domain.FriendEntity, status enums.FriendStatus) {
	go func() {
		asyncCtx, cancel := context.WithTimeout(_utils.CloneContext(ctx), 10*time.Second)
		defer cancel()

		currentUser, _ := s.userClient.GetProfileById(asyncCtx, responderId)
		if currentUser == nil {
			return
		}

		senderId := uint64(0)
		if requestEntity.CreatedBy != nil {
			senderId = *requestEntity.CreatedBy
		}
		if senderId == 0 {
			return
		}

		title := "Chấp nhận lời mời kết bạn"
		content := []string{"", currentUser.FullName, " đã chấp nhận lời mời kết bạn của bạn"}
		if status == enums.FriendStatusReject {
			title = "Từ chối lời mời kết bạn"
			content = []string{"", currentUser.FullName, " đã từ chối lời mời kết bạn của bạn"}
		}

		_ = s.notificationClient.CreateNotification(
			asyncCtx,
			currentUser.Avatar,
			title,
			content,
			_enum.NotificationFriendResponse,
			&responderId,
			senderId,
			_enum.EOwnerOfMember,
			[]string{strconv.FormatUint(requestEntity.ReceiverID, 10)},
		)
	}()
}

func (s *FriendUsecase) ChangeGroup(c context.Context, requestId uint64, groupId uint64) (int, error) {
	err := s.blockService.BeforeRequest(c, requestId)
	if err != nil {
		return 0, err
	}

	requestEntity, err := s.friendRepo.FindByID(c, requestId)
	if err != nil {
		// return 0, errors.New("Id không chính xác")
		return 0, &_routes.Except{
			Code:    400,
			Message: "Thông tin không đúng vui lòng kiểm tra lại",
		}
	}

	profileId := _utils.GetProfileIdWithContext(c)

	if (requestEntity.ReceiverID != profileId && requestEntity.CreatedBy != &profileId) || requestEntity.Status != enums.FriendStatusPending {
		return 0, &_routes.Except{
			Code:    400,
			Message: "Không thể thực hiện yêu cầu",
		}
	}

	isGroup := true
	if requestEntity.ReceiverID == profileId {
		isGroup = false
	}

	err = s.friendRepo.UpdateGroupById(c, isGroup, groupId, requestId)
	if err != nil {
		return 0, err
	}

	// Background tasks - Tạo history khi thay đổi nhóm
	// ctxClone := _utils.CloneContext(c)
	// go func() {
	// 	currentUser, _ := s.userClient.GetProfileById(ctxClone, profileId)
	// 	fullName := ""
	// 	if currentUser != nil {
	// 		fullName = currentUser.FullName
	// 	}

	// 	// Xác định người bạn (người còn lại trong mối quan hệ)
	// 	friendId := requestEntity.ReceiverID
	// 	if requestEntity.ReceiverID == profileId {
	// 		if requestEntity.CreatedBy != nil {
	// 			friendId = *requestEntity.CreatedBy
	// 		}
	// 	}

	// 	contact, err := s.contact.GetByProfileID(ctxClone, friendId, profileId, base_enum.EOwnerOfMember)
	// 	if err != nil {
	// 		return
	// 	}

	// 	// Tạo History
	// 	s.notificationClient.CreateCRMHistory(ctxClone, &base_dto.HistoryDTO{
	// 		Title:      fullName + " đã thay đổi nhóm bạn bè",
	// 		Note:       []string{"", fullName, " đã thay đổi nhóm của ", contact.FullName},
	// 		TargetId:   contact.ID,
	// 		TargetType: base_enum.TargetHistoryContact,
	// 		ActionType: base_enum.HistoryContactFriend,
	// 		OwnerID:    &profileId,
	// 		OwnerOf:    base_enum.EOwnerOfMember,
	// 	})
	// }()

	return 0, nil
}

// func (s *FriendService) Update(id uint64, request *data.FriendDTO) error {
// 	existingRequest, err := s.repo.FindByID(id)
// 	if err != nil {
// 		return errors.New("Id không chính xác")
// 	}

// 	if request.GroupID != nil && existingRequest.Status == 2 {
// 		if existingRequest.ReceiverID != request.ID && existingRequest.ID != request.ID {
// 			return errors.New("Không thể gán nhãn do đây không phải bạn bè của bạn")
// 		}
// 		_, err := s.groupRepo.GetByID(*request.GroupID)
// 		if err != nil {
// 			return err
// 		}
// 		return s.repo.UpdateGroupId(request.GroupID, id)
// 	}

// 	if request.Status == 2 {
// 		friendExists, _ := s.repo.ExistsFriendship(existingRequest.ID, existingRequest.ReceiverID)
// 		if friendExists {
// 			return errors.New("Đã là bạn bè")
// 		}
// 		t := time.Now()
// 		return s.repo.UpdateStatusById(request.Status, id, &t)
// 	}

// 	return s.repo.UpdateStatusById(request.Status, id, nil)
// }

// func (s *FriendUsecase) Delete(c context.Context, id uint64, userID uint64) error {
// 	existingRequest, err := s.repo.FindByID(c, id)
// 	if err != nil {
// 		return errors.New("Id không chính xác")
// 	}
// 	if existingRequest.ID != userID && existingRequest.ReceiverID != userID {
// 		return errors.New("Bạn không thể xóa yêu cầu này")
// 	}
// 	return s.repo.DeleteByID(c, id)
// }

func (s *FriendUsecase) FriendInfo(c context.Context) (*data.FollowInfoRes, error) {
	profileId := _utils.GetProfileIdWithContext(c)

	numFollower, numFollowing, friendNumber, err := s.followRepo.GetFollowInfoCount(c, profileId)
	if err != nil {
		return nil, err
	}

	return &data.FollowInfoRes{
		NumFollower:  numFollower,
		NumFollowing: numFollowing,
		NumFriend:    friendNumber,
	}, nil
}

// FriendInfoByProfileId lấy friend info theo profileId cụ thể
func (s *FriendUsecase) FriendInfoByProfileId(c context.Context, profileId uint64) (*data.FollowInfoRes, error) {
	numFollower, numFollowing, friendNumber, err := s.followRepo.GetFollowInfoCount(c, profileId)
	if err != nil {
		return nil, err
	}

	// Lấy contactId nếu có (contact của current user về profileId này)
	var contactId *uint64
	var blockId *uint64
	var followId *uint64
	var friendCommons []*sharepb.UserV3Proto
	var numFriendCommons uint64
	currentUserId := _utils.GetProfileIdWithContext(c)
	if currentUserId > 0 {
		contact, err := s.contact.GetByProfileID(c, profileId, currentUserId, base_enum.EOwnerOfMember)
		if err == nil && contact != nil {
			contactId = &contact.ID
		}

		// Kiểm tra xem currentUserId có block profileId hay không
		if s.blockRepo.IsBlocked(c, currentUserId, profileId) {
			blockId = &profileId
		}

		// Lấy followId nếu current user đang follow profileId này
		followId, _ = s.followRepo.GetFollowId(c, currentUserId, profileId)

		// Lấy bạn chung giữa currentUserId và profileId
		friendCommons, numFriendCommons, err = s.GetCommonFriends(c, currentUserId, profileId, 2)
		if err != nil {
			// Log error nhưng không fail toàn bộ request
			friendCommons = []*sharepb.UserV3Proto{}
			numFriendCommons = 0
		}
	}

	return &data.FollowInfoRes{
		NumFollower:      numFollower,
		NumFollowing:     numFollowing,
		NumFriend:        friendNumber,
		ContactId:        contactId,
		BlockId:          blockId,
		FollowId:         followId,
		FriendCommons:    friendCommons,
		NumFriendCommons: numFriendCommons,
	}, nil
}

// GetCommonFriends lấy danh sách bạn chung giữa 2 user
// Trả về: danh sách UserV3Proto (tối đa limit), tổng số bạn chung
func (s *FriendUsecase) GetCommonFriends(c context.Context, user1Id uint64, user2Id uint64, limit int) ([]*sharepb.UserV3Proto, uint64, error) {
	// Lấy danh sách profileId của bạn chung
	commonFriendIds, total, err := s.friendRepo.GetCommonFriends(c, user1Id, user2Id, limit)
	if err != nil {
		return nil, 0, err
	}

	if len(commonFriendIds) == 0 {
		return []*sharepb.UserV3Proto{}, uint64(total), nil
	}

	// Lấy thông tin profile của các bạn chung bằng cách gọi từng profile
	// Vì UserClient interface không có method GetProfileByIds, ta sẽ gọi từng cái
	friendCommons := make([]*sharepb.UserV3Proto, 0, len(commonFriendIds))
	for _, profileId := range commonFriendIds {
		profile, err := s.userClient.GetProfileById(c, profileId)
		if err != nil || profile == nil {
			continue
		}
		friendCommons = append(friendCommons, &sharepb.UserV3Proto{
			ProfileId: profile.ProfileId,
			FullName:  profile.FullName,
			Avatar:    profile.Avatar,
		})
	}

	return friendCommons, uint64(total), nil
}

func (s *FriendUsecase) CancelRequest(ctx context.Context, receiverId uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return _errors.UnauthorizedException("missing profile id")
	}

	return s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// Tìm request pending từ profileId đến receiverId
		request, err := s.friendRepo.FindBySenderIdAndReceiverId(txCtx, profileId, receiverId)
		if err != nil || request == nil {
			return _errors.NotFoundException("friend request not found")
		}
		if request.Status != enums.FriendStatusPending {
			return _errors.ConflictException("cannot cancel non-pending request")
		}

		// Xóa request (soft delete)
		if err := s.friendRepo.CancelRequest(txCtx, request.ID); err != nil {
			return _errors.InternalServerException("failed to cancel request: " + err.Error())
		}

		// Xóa contact của current user về receiver
		contact1, err := s.contact.GetByProfileID(txCtx, receiverId, profileId, base_enum.EOwnerOfMember)
		if err == nil && contact1 != nil {
			if err := s.contact.Delete(txCtx, contact1.ID); err != nil {
				return _errors.InternalServerException("failed to delete contact for sender: " + err.Error())
			}
		}

		// Xóa contact của receiver về current user
		contact2, err := s.contact.GetByProfileID(txCtx, profileId, receiverId, base_enum.EOwnerOfMember)
		if err == nil && contact2 != nil {
			if err := s.contact.Delete(txCtx, contact2.ID); err != nil {
				return _errors.InternalServerException("failed to delete contact for receiver: " + err.Error())
			}
		}

		return nil
	})
}

func (s *FriendUsecase) Unfriend(ctx context.Context, friendId uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return _errors.UnauthorizedException("missing profile id")
	}

	return s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// Tìm friendship giữa hai người (status = accepted)
		friendship, err := s.friendRepo.GetFriendship(txCtx, profileId, friendId)
		if err != nil {
			return _errors.InternalServerException("failed to get friendship: " + err.Error())
		}
		if friendship == nil || friendship.Status != enums.FriendStatusAccepted {
			return _errors.NotFoundException("friendship not found")
		}

		// Xóa friendship (soft delete)
		if err := s.friendRepo.DeleteByID(txCtx, friendship.ID); err != nil {
			return _errors.InternalServerException("failed to delete friendship: " + err.Error())
		}

		// Xóa contact của profileId về friendId
		contact1, err := s.contact.GetByProfileID(txCtx, friendId, profileId, base_enum.EOwnerOfMember)
		if err == nil && contact1 != nil {
			if err := s.contact.Delete(txCtx, contact1.ID); err != nil {
				return _errors.InternalServerException("failed to delete contact for profile: " + err.Error())
			}
		}

		// Xóa contact của friendId về profileId
		contact2, err := s.contact.GetByProfileID(txCtx, profileId, friendId, base_enum.EOwnerOfMember)
		if err == nil && contact2 != nil {
			if err := s.contact.Delete(txCtx, contact2.ID); err != nil {
				return _errors.InternalServerException("failed to delete contact for friend: " + err.Error())
			}
		}

		return nil
	})
}
