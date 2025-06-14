package usecase

import (
	_db "common/db"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/data"
	"crm/internal/domain"
	"crm/internal/enums"
	iusecase "crm/internal/interface"
	"crm/internal/repo"
	"errors"
)

type FriendUsecase struct {
	repo         repo.FriendRepo
	contact      repo.ContactRepo
	blockRepo    repo.BlockRepo
	followRepo   repo.FollowRepo
	blockService *BlockUsecase
	historyRepo  repo.HistoryRepo
	userClient   iusecase.UserClient
}

func NewFriendUsecase(repo repo.FriendRepo,
	contact repo.ContactRepo,
	followRepo repo.FollowRepo,
	blockRepo repo.BlockRepo,
	blockService *BlockUsecase,
	historyRepo repo.HistoryRepo,
	userClient iusecase.UserClient,
) *FriendUsecase {
	return &FriendUsecase{
		repo:         repo,
		contact:      contact,
		blockRepo:    blockRepo,
		blockService: blockService,
		followRepo:   followRepo,
		historyRepo:  historyRepo,
		userClient:   userClient,
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

	// ok, _ := s.contact.ExistByProfile(c, receiverId)
	// if !ok {
	// 	return nil, &_routes.Except{
	// 		Code:    400,
	// 		Message: "Người dùng không tồn tại hoặc chưa đăng ký",
	// 	}
	// }

	e := &domain.FriendEntity{
		Status:     enums.FriendStatusPending,
		ReceiverID: receiverId,
	}
	_db.SaveWithContext(c, &e)

	newCtx := _utils.CloneContext(c)
	go func() {
		profileId := _utils.GetProfileIdWithContext(newCtx)
		currentUser, _ := s.userClient.GetProfileById(newCtx, profileId)
		fullName := ""
		if currentUser != nil {
			fullName = currentUser.FullName
		}

		contact, err := s.contact.GetByProfileID(newCtx, receiverId, profileId, enums.EOwnerTypeMember)
		if err != nil {
			return
		}

		s.historyRepo.CreateHistory(newCtx, &domain.HistoryEntity{
			Title:      fullName + " đã gửi lời mời kết bạn",
			Note:       fullName + " đã gửi lời mời kết bạn tới " + contact.FullName,
			TargetId:   contact.ID,
			TargetType: enums.TargetHistoryContact,
			ActionType: enums.HistoryContactFriend,
			OwnerID:    &profileId,
			OwnerType:  enums.EOwnerTypeMember,
		})
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

	requestEntity, err := s.repo.FindByID(c, requestId)
	if err != nil {
		return 0, errors.New("Id không chính xác")
	}

	profileId := _utils.GetProfileIdWithContext(c)

	if requestEntity.ReceiverID != profileId || requestEntity.Status != enums.FriendStatusPending {
		return 0, errors.New("Không thể thực hiện yêu cầu")
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

	return 0, s.repo.UpdateStatusById(c, status, requestId, nil)
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
	return 0, s.repo.UpdateGroupById(c, isGroup, groupId, requestId)
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

func (s *FriendUsecase) Delete(c context.Context, id uint64, userID uint64) error {
	existingRequest, err := s.repo.FindByID(c, id)
	if err != nil {
		return errors.New("Id không chính xác")
	}
	if existingRequest.ID != userID && existingRequest.ReceiverID != userID {
		return errors.New("Bạn không thể xóa yêu cầu này")
	}
	return s.repo.DeleteByID(c, id)
}

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
