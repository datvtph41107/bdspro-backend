package usecases

import (
	_errors "common/errors"
	_jwt "common/jwt"
	"user/internal"
	"user/internal/dto"
	"user/internal/enums"
	"user/internal/interface/repo"
	models "user/internal/models"

	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (s *FriendUsecase) FriendList(c *gin.Context) ([]models.FriendEntity, error) {
	profileId := _jwt.GetProfileId(c)
	return s.repo.SearchFriend(profileId, dto.FriendDTO{}, 0, 10)
}

func (s *FriendUsecase) Request(c *gin.Context) (*models.FriendEntity, error) {
	profileId := _jwt.GetProfileId(c)
	receiverIdStr := c.Param("receiverId")
	receiverId, _ := strconv.ParseUint(receiverIdStr, 10, 32)

	err := s.blockService.BeforeRequest(c, receiverId)
	if err != nil {
		return nil, err
	}

	if receiverId == profileId {
		return nil, _errors.ReturnError(service.SelfFriendNotAllowed)
	}

	existingRequest, _ := s.repo.ExistsByCreatedByAndReceiverId(profileId, receiverId)
	if existingRequest {
		return nil, _errors.ReturnError(service.FriendRequestAlreadySent)
	}

	ok, _ := s.contact.ExistByProfile(receiverId)
	if !ok {
		return nil, _errors.ReturnError(service.FriendTargetUnavailable)
	}

	e := &models.FriendEntity{
		Status:     enums.EFriendStatusPending,
		ReceiverID: receiverId,
	}
	if err := s.repo.Create(c.Request.Context(), e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *FriendUsecase) RequestSent(c *gin.Context) ([]models.FriendEntity, error) {

	profileId := _jwt.GetProfileId(c)
	var params dto.FriendDTO
	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, _errors.ReturnError(service.RequestValidationFailed, _errors.WithCause(err))
	}
	params.Status = enums.EFriendStatusPending

	return s.repo.Search(params, &profileId, 0, 10)
}

func (s *FriendUsecase) RequestReceived(c *gin.Context) ([]models.FriendEntity, error) {
	profileId := _jwt.GetProfileId(c)
	var params dto.FriendDTO
	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, _errors.ReturnError(service.RequestValidationFailed, _errors.WithCause(err))
	}
	params.ReceiverID = profileId
	params.Status = enums.EFriendStatusPending

	return s.repo.Search(params, nil, 0, 10)
}

func (s *FriendUsecase) ChangeStatus(c *gin.Context, status enums.EFriendStatus) (int, error) {
	idStr := c.Param("requestId")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	err := s.blockService.BeforeRequest(c, id)
	if err != nil {
		return 0, err
	}

	requestEntity, err := s.repo.FindByID(id)
	if err != nil {
		return 0, errors.New("Id không chính xác")
	}

	profileId := _jwt.GetProfileId(c)

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

	return 0, s.repo.UpdateStatusById(status, id, nil)
}

func (s *FriendUsecase) ChangeGroup(c *gin.Context) (int, error) {
	id, _ := strconv.ParseUint(c.Param("requestId"), 10, 32)
	groupId, _ := strconv.ParseUint(c.Param("groupId"), 10, 32)

	err := s.blockService.BeforeRequest(c, id)
	if err != nil {
		return 0, err
	}

	requestEntity, err := s.repo.FindByID(id)
	if err != nil {
		// return 0, errors.New("Id không chính xác")
		return 0, _errors.ReturnError(service.RequestValidationFailed, _errors.WithPublicMessage("Thông tin không đúng vui lòng kiểm tra lại"))
	}

	profileId := _jwt.GetProfileId(c)

	if (requestEntity.ReceiverID != profileId && requestEntity.CreatedBy != &profileId) || requestEntity.Status != enums.EFriendStatusPending {
		return 0, _errors.ReturnError(service.FriendRequestCannotBePerformed)
	}

	isGroup := true
	if requestEntity.ReceiverID == profileId {
		isGroup = false
	}
	return 0, s.repo.UpdateGroupById(isGroup, groupId, id)
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

func (s *FriendUsecase) FriendInfo(c *gin.Context) (*dto.FollowInfoRes, error) {
	profileId := _jwt.GetProfileId(c)

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
