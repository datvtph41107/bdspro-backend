package usecase

import (
	base_enum "base/enum"
	_db "common/db"
	_enum "common/domain/enum"
	_provider "common/domain/provider"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"errors"
	"strconv"

	"crm/data"
	"crm/internal/domain"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type FriendUsecase struct {
	repo               repo.FriendRepo
	contact            repo.ContactRepo
	blockRepo          repo.BlockRepo
	followRepo         repo.FollowRepo
	blockService       *BlockUsecase
	notificationClient _provider.NotificationProvider
	userClient         provider.UserClient
}

func NewFriendUsecase(repo repo.FriendRepo,
	contact repo.ContactRepo,
	followRepo repo.FollowRepo,
	blockRepo repo.BlockRepo,
	blockService *BlockUsecase,
	notificationClient _provider.NotificationProvider,
	userClient provider.UserClient,
) *FriendUsecase {
	return &FriendUsecase{
		repo:               repo,
		contact:            contact,
		blockRepo:          blockRepo,
		blockService:       blockService,
		followRepo:         followRepo,
		notificationClient: notificationClient,
		userClient:         userClient,
	}
}

func (s *FriendUsecase) GetByID(c context.Context, id uint64) (*domain.FriendEntity, error) {
	return s.repo.FindByID(c, id)
}

// func (s *FriendService) Search(receiverID uint, createdBy uint, isSender bool, page int, size int) ([]domain.FriendEntity, error) {
// 	if isSender {
// 		receiverID = 0
// 	}
// 	return s.Repo.Search(receiverID, createdBy, page, size)
// }

func (s *FriendUsecase) FriendList(c context.Context, dto data.FriendRequest) ([]domain.FriendEntity, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.repo.SearchFriend(c, profileId, dto)
}

func (s *FriendUsecase) Request(c context.Context, receiverId uint64) (*domain.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)

	err := s.blockService.BeforeRequest(c, receiverId)
	if err != nil {
		return nil, err
	}

	existingRequest, _ := s.repo.ExistsByCreatedByAndReceiverId(c, profileId, receiverId)
	if existingRequest {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Đã gửi yêu cầu",
		}
	}

	// Lấy thông tin cả 2 người
	currentUser, _ := s.userClient.GetProfileWithPhoneById(c, profileId)
	receiverUser, _ := s.userClient.GetProfileWithPhoneById(c, receiverId)

	if receiverUser == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Người dùng không tồn tại",
		}
	}

	// Tạo friend request
	e := &domain.FriendEntity{
		Status:     enums.FriendStatusPending,
		ReceiverID: receiverId,
	}
	_db.SaveWithContext(c, &e)

	// Tạo contact cho người gửi (lưu thông tin người nhận)
	s.contact.GetOrCreateByProfileID(c, receiverId, profileId, base_enum.EOwnerOfMember,
		receiverUser.FullName, receiverUser.Phone, receiverUser.Avatar)

	// Tạo contact cho người nhận (lưu thông tin người gửi)
	currentFullName := ""
	currentPhone := ""
	currentAvatar := ""
	if currentUser != nil {
		currentFullName = currentUser.FullName
		currentPhone = currentUser.Phone
		currentAvatar = currentUser.Avatar
	}
	s.contact.GetOrCreateByProfileID(c, profileId, receiverId, base_enum.EOwnerOfMember,
		currentFullName, currentPhone, currentAvatar)

	// Background tasks
	ctxClone := _utils.CloneContext(c)
	go func() {
		profileId := _utils.GetProfileIdWithContext(ctxClone)
		currentUser, _ := s.userClient.GetProfileById(ctxClone, profileId)
		fullName := ""
		if currentUser != nil {
			fullName = currentUser.FullName
		}

		contact, err := s.contact.GetByProfileID(ctxClone, receiverId, profileId, base_enum.EOwnerOfMember)
		if err != nil {
			return
		}

		// Tạo History
		// s.notificationClient.CreateCRMHistory(ctxClone, &base_dto.HistoryDTO{
		// 	Title:      fullName + " đã gửi lời mời kết bạn",
		// 	Note:       []string{"", fullName, " đã gửi lời mời kết bạn tới ", contact.FullName},
		// 	TargetId:   contact.ID,
		// 	TargetType: base_enum.TargetHistoryContact,
		// 	ActionType: base_enum.HistoryContactFriend,
		// 	OwnerID:    &profileId,
		// 	OwnerOf:    base_enum.EOwnerOfMember,
		// })

		// Gửi Notification cho người nhận
		s.notificationClient.CreateNotification(
			ctxClone,
			currentUser.Avatar,
			"Lời mời kết bạn",
			[]string{"", fullName, " đã gửi lời mời kết bạn cho bạn"},
			_enum.NotificationFriendRequest,
			&profileId,
			receiverId,
			_enum.EOwnerOfMember,
			[]string{strconv.FormatUint(contact.ID, 10)},
		)
		// s.notificationClient.SendNotification(ctxClone, &base_dto.NotificationDTO{
		// 	Title:   "Lời mời kết bạn",
		// 	Message: fullName + " đã gửi lời mời kết bạn cho bạn",
		// 	Type:    10, // FriendRequest
		// 	UserID:  receiverId,
		// })
	}()

	return e, nil
}

func (s *FriendUsecase) RequestSent(c context.Context, params data.FriendDTO) ([]domain.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	params.Status = enums.FriendStatusPending
	return s.repo.Search(c, params, &profileId, 0, 10)
}

func (s *FriendUsecase) RequestReceived(c context.Context, params data.FriendDTO) ([]domain.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	params.ReceiverID = profileId
	params.Status = enums.FriendStatusPending

	return s.repo.Search(c, params, nil, 0, 10)
}

func (s *FriendUsecase) ChangeStatus(c context.Context, requestId uint64, status enums.FriendStatus) (int, error) {
	err := s.blockService.BeforeRequest(c, requestId)
	if err != nil {
		return 0, err
	}

	// requestEntity, err := s.repo.FindByID(c, requestId)
	profileId := _utils.GetProfileIdWithContext(c)
	requestEntity, err := s.repo.FindBySenderIdAndReceiverId(c, requestId, profileId)
	if err != nil {
		return 0, errors.New("id không chính xác")
	}

	if requestEntity.ReceiverID != profileId || requestEntity.Status != enums.FriendStatusPending {
		return 0, errors.New("bạn không phải người nhận yêu cầu hoặc không thể thực hiện yêu cầu")
	}

	// if requestEntity.Status != enums.Pending {
	// 	if existingRequest.ReceiverID != request.ID && existingRequest.ID != request.ID {
	// 		return errors.New("Không thể gán nhãn do đây không phải bạn bè của bạn")
	// 	}
	// 	_, err := s.groupRepo.GetByID(*request.GroupID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return s.repo.UpdateGroupId(request.GroupID, id)
	// }

	// if request.Status == 2 {
	// 	friendExists, _ := s.repo.ExistsFriendship(existingRequest.ID, existingRequest.ReceiverID)
	// 	if friendExists {
	// 		return errors.New("Đã là bạn bè")
	// 	}
	// 	t := time.Now()
	// 	return s.repo.UpdateStatusById(request.Status, id, &t)
	// }

	// Update status
	updateErr := s.repo.UpdateStatusById(c, status, requestEntity.ID, nil)
	if updateErr != nil {
		return 0, updateErr
	}

	// Gửi notification và tạo history khi accept hoặc reject
	ctxClone := _utils.CloneContext(c)
	go func() {
		currentUser, _ := s.userClient.GetProfileById(ctxClone, profileId)
		fullName := ""
		if currentUser != nil {
			fullName = currentUser.FullName
		}

		senderID := uint64(0)
		if requestEntity.CreatedBy != nil {
			senderID = *requestEntity.CreatedBy
		}

		if senderID == 0 {
			return
		}

		// contact, err := s.contact.GetByProfileID(ctxClone, senderID, profileId, base_enum.EOwnerOfMember)
		// if err != nil {
		// 	return
		// }

		if status == enums.FriendStatusAccepted {
			// Chấp nhận lời mời
			// s.notificationClient.CreateCRMHistory(ctxClone, &base_dto.HistoryDTO{
			// 	Title:      fullName + " đã chấp nhận lời mời kết bạn",
			// 	Note:       []string{"", fullName, " đã chấp nhận lời mời kết bạn từ ", contact.FullName},
			// 	TargetId:   contact.ID,
			// 	TargetType: base_enum.TargetHistoryContact,
			// 	ActionType: base_enum.HistoryContactFriend,
			// 	OwnerID:    &profileId,
			// 	OwnerOf:    base_enum.EOwnerOfMember,
			// })

			s.notificationClient.CreateNotification(
				ctxClone,
				currentUser.Avatar,
				"Chấp nhận lời mời kết bạn",
				[]string{"", fullName, " đã chấp nhận lời mời kết bạn của bạn"},
				_enum.NotificationFriendResponse,
				&profileId,
				senderID,
				_enum.EOwnerOfMember,
				[]string{strconv.FormatUint(requestEntity.ReceiverID, 10)},
			)
		} else if status == enums.FriendStatusReject {
			// Từ chối lời mời
			// s.notificationClient.CreateCRMHistory(ctxClone, &base_dto.HistoryDTO{
			// 	Title:      fullName + " đã từ chối lời mời kết bạn",
			// 	Note:       []string{"", fullName, " đã từ chối lời mời kết bạn từ ", contact.FullName},
			// 	TargetId:   contact.ID,
			// 	TargetType: base_enum.TargetHistoryContact,
			// 	ActionType: base_enum.HistoryContactFriend,
			// 	OwnerID:    &profileId,
			// 	OwnerOf:    base_enum.EOwnerOfMember,
			// })

			// s.notificationClient.SendNotification(ctxClone, &base_dto.NotificationDTO{
			// 	Title:      "Từ chối lời mời kết bạn",
			// 	Message:    fullName + " đã từ chối lời mời kết bạn của bạn",
			// 	Type:       20, // FriendResponse
			// 	UserID:     senderID,
			// 	AttachData: []string{contact.FullName},
			// })
			s.notificationClient.CreateNotification(
				ctxClone,
				currentUser.Avatar,
				"Từ chối lời mời kết bạn",
				[]string{"", fullName, " đã từ chối lời mời kết bạn của bạn"},
				_enum.NotificationFriendResponse,
				&profileId,
				senderID,
				_enum.EOwnerOfMember,
				[]string{strconv.FormatUint(requestEntity.ReceiverID, 10)},
			)
		}
	}()

	return 0, nil
}

func (s *FriendUsecase) ChangeGroup(c context.Context, requestId uint64, groupId uint64) (int, error) {
	err := s.blockService.BeforeRequest(c, requestId)
	if err != nil {
		return 0, err
	}

	requestEntity, err := s.repo.FindByID(c, requestId)
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

	err = s.repo.UpdateGroupById(c, isGroup, groupId, requestId)
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
	}

	return &data.FollowInfoRes{
		NumFollower:  numFollower,
		NumFollowing: numFollowing,
		NumFriend:    friendNumber,
		ContactId:    contactId,
		BlockId:      blockId,
		FollowId:     followId,
	}, nil
}

func (s *FriendUsecase) CancelRequest(c context.Context, receiverId uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)

	requestEntity, err := s.repo.FindBySenderIdAndReceiverId(c, profileId, receiverId)
	if err != nil {
		return &_routes.Except{
			Code:    400,
			Message: "Id không chính xác",
		}
	}
	if requestEntity == nil || requestEntity.ID == 0 {
		return &_routes.Except{
			Code:    400,
			Message: "Không tìm thấy lời mời kết bạn",
		}
	}

	if requestEntity.Status != enums.FriendStatusPending {
		return &_routes.Except{
			Code:    400,
			Message: "Không thể hủy lời mời kết bạn",
		}
	}

	// Xóa friend request
	err = s.repo.CancelRequest(c, requestEntity.ID)
	if err != nil {
		return err
	}

	// Xóa contact của user hiện tại về receiverId
	contact1, err := s.contact.GetByProfileID(c, receiverId, profileId, base_enum.EOwnerOfMember)
	if err == nil && contact1 != nil {
		s.contact.Delete(c, contact1.ID)
	}

	// // Xóa contact của receiverId về user hiện tại
	// contact2, err := s.contact.GetByProfileID(c, profileId, receiverId, base_enum.EOwnerOfMember)
	// if err == nil && contact2 != nil {
	// 	s.contact.Delete(c, contact2.ID)
	// }

	// Background tasks - Gửi notification
	ctxClone := _utils.CloneContext(c)
	go func() {
		currentUser, _ := s.userClient.GetProfileById(ctxClone, profileId)
		fullName := ""
		if currentUser != nil {
			fullName = currentUser.FullName
		}

		s.notificationClient.CreateNotification(
			ctxClone,
			currentUser.Avatar,
			"Hủy lời mời kết bạn",
			[]string{"", fullName, " đã hủy lời mời kết bạn"},
			_enum.NotificationFriendResponse,
			&profileId,
			receiverId,
			_enum.EOwnerOfMember,
			[]string{strconv.FormatUint(requestEntity.ReceiverID, 10)},
		)
	}()

	return nil
}

// Unfriend hủy kết bạn và xóa contact của cả 2 user
func (s *FriendUsecase) Unfriend(c context.Context, friendId uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)

	// Tìm friend relationship giữa 2 user
	requestEntity, err := s.repo.FindBySenderIdAndReceiverId(c, profileId, friendId)
	if err != nil {
		// Thử tìm ngược lại
		requestEntity, err = s.repo.FindBySenderIdAndReceiverId(c, friendId, profileId)
		if err != nil {
			return &_routes.Except{
				Code:    400,
				Message: "Không tìm thấy quan hệ bạn bè",
			}
		}
	}

	// Kiểm tra xem có phải là bạn bè không (status == Accepted)
	if requestEntity.Status != enums.FriendStatusAccepted {
		return &_routes.Except{
			Code:    400,
			Message: "Hai người chưa là bạn bè",
		}
	}

	// Xóa friend relationship
	err = s.repo.DeleteByID(c, requestEntity.ID)
	if err != nil {
		return err
	}

	// Xóa contact của user hiện tại về friendId
	contact1, err := s.contact.GetByProfileID(c, friendId, profileId, base_enum.EOwnerOfMember)
	if err == nil && contact1 != nil {
		s.contact.Delete(c, contact1.ID)
	}

	// Xóa contact của friendId về user hiện tại
	contact2, err := s.contact.GetByProfileID(c, profileId, friendId, base_enum.EOwnerOfMember)
	if err == nil && contact2 != nil {
		s.contact.Delete(c, contact2.ID)
	}

	// Background tasks - Gửi notification
	ctxClone := _utils.CloneContext(c)
	go func() {
		currentUser, _ := s.userClient.GetProfileById(ctxClone, profileId)
		fullName := ""
		if currentUser != nil {
			fullName = currentUser.FullName
		}

		s.notificationClient.CreateNotification(
			ctxClone,
			currentUser.Avatar,
			"Hủy kết bạn",
			[]string{"", fullName, " đã hủy kết bạn với bạn"},
			_enum.NotificationFriendResponse,
			&profileId,
			friendId,
			_enum.EOwnerOfMember,
			[]string{},
		)
	}()

	return nil
}
