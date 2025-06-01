package user_service

import (
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	"user/data"
	"user/enums"
	user_model "user/models/user"
	"user/repo"

	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendService struct {
	repo         *repo.FriendRepo
	contact      *repo.ContactRepo
	blockRepo    *repo.BlockRepo
	followRepo   *repo.FollowRepo
	blockService *BlockService
}

func NewFriendService(repo *repo.FriendRepo,
	contact *repo.ContactRepo,
	followRepo *repo.FollowRepo,
	blockRepo *repo.BlockRepo,
	blockService *BlockService,
) *FriendService {
	return &FriendService{
		repo:         repo,
		contact:      contact,
		blockRepo:    blockRepo,
		blockService: blockService,
		followRepo:   followRepo,
	}
}

func (s *FriendService) GetByID(id uint64) (*user_model.FriendEntity, error) {
	return s.repo.FindByID(id)
}

// func (s *FriendService) Search(receiverID uint, createdBy uint, isSender bool, page int, size int) ([]user_model.FriendEntity, error) {
// 	if isSender {
// 		receiverID = 0
// 	}
// 	return s.Repo.Search(receiverID, createdBy, page, size)
// }

func (s *FriendService) FriendList(c *gin.Context) ([]user_model.FriendEntity, error) {
	profileId := _jwt.GetProfileId(c)
	return s.repo.SearchFriend(profileId, data.FriendDTO{}, 0, 10)
}

func (s *FriendService) Request(c *gin.Context) (*user_model.FriendEntity, error) {
	profileId := _jwt.GetProfileId(c)
	receiverIdStr := c.Param("receiverId")
	receiverId, _ := strconv.ParseUint(receiverIdStr, 10, 32)

	err := s.blockService.BeforeRequest(c, receiverId)
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

	e := &user_model.FriendEntity{
		Status:     enums.EFriendStatusPending,
		ReceiverID: receiverId,
	}
	_db.SaveWithAudit(c, &e)
	return e, nil
}

func (s *FriendService) RequestSent(c *gin.Context) ([]user_model.FriendEntity, error) {

	profileId := _jwt.GetProfileId(c)
	var params data.FriendDTO
	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}
	params.Status = enums.EFriendStatusPending

	return s.repo.Search(params, &profileId, 0, 10)
}

func (s *FriendService) RequestReceived(c *gin.Context) ([]user_model.FriendEntity, error) {
	profileId := _jwt.GetProfileId(c)
	var params data.FriendDTO
	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}
	params.ReceiverID = profileId
	params.Status = enums.EFriendStatusPending

	return s.repo.Search(params, nil, 0, 10)
}

func (s *FriendService) ChangeStatus(c *gin.Context, status enums.EFriendStatus) (int, error) {
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

func (s *FriendService) ChangeGroup(c *gin.Context) (int, error) {
	id, _ := strconv.ParseUint(c.Param("requestId"), 10, 32)
	groupId, _ := strconv.ParseUint(c.Param("groupId"), 10, 32)

	err := s.blockService.BeforeRequest(c, id)
	if err != nil {
		return 0, err
	}

	requestEntity, err := s.repo.FindByID(id)
	if err != nil {
		// return 0, errors.New("Id không chính xác")
		return 0, &_routes.Except{
			Code:    400,
			Message: "Thông tin không đúng vui lòng kiểm tra lại",
		}
	}

	profileId := _jwt.GetProfileId(c)

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
	return 0, s.repo.UpdateGroupById(isGroup, groupId, id)
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

func (s *FriendService) Delete(id uint64, userID uint64) error {
	existingRequest, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("Id không chính xác")
	}
	if existingRequest.ID != userID && existingRequest.ReceiverID != userID {
		return errors.New("Bạn không thể xóa yêu cầu này")
	}
	return s.repo.DeleteByID(id)
}

func (s *FriendService) FriendInfo(c *gin.Context) (*data.FollowInfoRes, error) {
	profileId := _jwt.GetProfileId(c)

	numFollower, numFollowing, friendNumber, err := s.followRepo.GetFollowInfoCount(profileId)
	if err != nil {
		return nil, err
	}

	return &data.FollowInfoRes{
		NumFollower:  numFollower,
		NumFollowing: numFollowing,
		NumFriend:    friendNumber,
	}, nil
}
