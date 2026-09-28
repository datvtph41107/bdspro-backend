package usecases

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"errors"
	"userv4/internal/dto"
	"userv4/internal/enums"
	"userv4/internal/interface/repo"
	models "userv4/internal/models"
)

type FriendUsecase struct {
	repo         repo.IFriendRepo
	contact      repo.IContactRepo
	blockRepo    repo.IBlockRepo
	followRepo   repo.IFollowRepo
	blockService *BlockUsecase
}

func NewFriendUsecase(repo repo.IFriendRepo,
	contact repo.IContactRepo,
	followRepo repo.IFollowRepo,
	blockRepo repo.IBlockRepo,
	blockService *BlockUsecase,
) *FriendUsecase {
	return &FriendUsecase{
		repo:         repo,
		contact:      contact,
		blockRepo:    blockRepo,
		blockService: blockService,
		followRepo:   followRepo,
	}
}

func (s *FriendUsecase) GetByID(id uint64) (*models.FriendEntity, error) {
	return s.repo.FindByID(id)
}

// func (s *FriendUsecase) Search(receiverID uint, createdBy uint, isSender bool, page int, size int) ([]models.FriendEntity, error) {
// 	if isSender {
// 		receiverID = 0
// 	}
// 	return s.Repo.Search(receiverID, createdBy, page, size)
// }

func (s *FriendUsecase) FriendList(ctx context.Context) ([]models.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	return s.repo.SearchFriend(profileId, dto.FriendDTO{}, 0, 10)
}

func (s *FriendUsecase) Request(ctx context.Context, receiverId uint64) (*models.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	err := s.blockService.BeforeRequest(ctx, receiverId)
	if err != nil {
		return nil, err
	}

	if receiverId == profileId {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Bạn bè là chính mình",
		}
	}

	existingRequest, _ := s.repo.ExistsByCreatedByAndReceiverId(profileId, receiverId)
	if existingRequest {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Đã gửi yêu cầu",
		}
	}

	ok, _ := s.contact.ExistByProfile(receiverId)
	if !ok {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Người dùng không tồn tại hoặc chưa đăng ký",
		}
	}

	e := &models.FriendEntity{
		Status:     enums.EFriendStatusPending,
		ReceiverID: receiverId,
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *FriendUsecase) RequestSent(ctx context.Context, params dto.FriendDTO) ([]models.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	params.Status = enums.EFriendStatusPending

	return s.repo.Search(params, &profileId, 0, 10)
}

func (s *FriendUsecase) RequestReceived(ctx context.Context, params dto.FriendDTO) ([]models.FriendEntity, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	params.ReceiverID = profileId
	params.Status = enums.EFriendStatusPending

	return s.repo.Search(params, nil, 0, 10)
}

func (s *FriendUsecase) ChangeStatus(ctx context.Context, requestId uint64, status enums.EFriendStatus) (int, error) {
	err := s.blockService.BeforeRequest(ctx, requestId)
	if err != nil {
		return 0, err
	}

	requestEntity, err := s.repo.FindByID(requestId)
	if err != nil {
		return 0, errors.New("Id không chính xác")
	}

	profileId := _utils.GetProfileIdWithContext(ctx)

	if requestEntity.ReceiverID != profileId || requestEntity.Status != enums.EFriendStatusPending {
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

	return 0, s.repo.UpdateStatusById(status, requestId, nil)
}

func (s *FriendUsecase) ChangeGroup(ctx context.Context, requestId, groupId uint64) (int, error) {
	err := s.blockService.BeforeRequest(ctx, requestId)
	if err != nil {
		return 0, err
	}

	requestEntity, err := s.repo.FindByID(requestId)
	if err != nil {
		// return 0, errors.New("Id không chính xác")
		return 0, &_routes.Except{
			Code:    400,
			Message: "Thông tin không đúng vui lòng kiểm tra lại",
		}
	}

	profileId := _utils.GetProfileIdWithContext(ctx)

	if (requestEntity.ReceiverID != profileId && requestEntity.CreatedBy != &profileId) || requestEntity.Status != enums.EFriendStatusPending {
		return 0, &_routes.Except{
			Code:    400,
			Message: "Không thể thực hiện yêu cầu",
		}
	}

	isGroup := true
	if requestEntity.ReceiverID == profileId {
		isGroup = false
	}
	return 0, s.repo.UpdateGroupById(isGroup, groupId, requestId)
}

// func (s *FriendUsecase) Update(id uint64, request *dto.FriendDTO) error {
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

func (s *FriendUsecase) Delete(id uint64, userID uint64) error {
	existingRequest, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("Id không chính xác")
	}
	if existingRequest.ID != userID && existingRequest.ReceiverID != userID {
		return errors.New("Bạn không thể xóa yêu cầu này")
	}
	return s.repo.DeleteByID(id)
}

func (s *FriendUsecase) FriendInfo(ctx context.Context) (*dto.FollowInfoRes, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	numFollower, numFollowing, friendNumber, err := s.followRepo.GetFollowInfoCount(profileId)
	if err != nil {
		return nil, err
	}

	return &dto.FollowInfoRes{
		NumFollower:  numFollower,
		NumFollowing: numFollowing,
		NumFriend:    friendNumber,
	}, nil
}
