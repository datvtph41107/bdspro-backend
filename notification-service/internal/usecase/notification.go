package usecase

import (
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"encoding/json"
	"log"
	"noti/internal/domain"
	"noti/internal/dto"
	"noti/internal/enums"
	"noti/internal/repo"
	pb_noti "pb/types/notification"
	"time"
)

type NotificationUsecase struct {
	Repo repo.NotificationRepo
}

func NewNotificationService(repo repo.NotificationRepo) *NotificationUsecase {
	return &NotificationUsecase{
		Repo: repo,
	}
}

// --- Gateway ---
// CreateNotification: Tạo thông báo mới
func (s *NotificationUsecase) CreateNotification(c context.Context, request *dto.NotiNewRequest) (*domain.NotificationEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	if profileId == 0 {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Không thể tạo thông báo",
		}
	}

	var visibleAt *time.Time
	if request.VisibleAt != nil {
		visibleAt = request.VisibleAt
	}
	notification := domain.NotificationEntity{
		OwnerID:   profileId,
		OwnerType: int32(request.OwnerType),
		TargetID:  request.TargetID,
		Title:     request.Title,
		Message:   request.Message,
		Link:      request.Link,
		Type:      enums.TypeEnum(request.Type),
		VisibleAt: visibleAt,
	}

	if err := s.Repo.Create(c, &notification); err != nil {
		return nil, err
	}

	return &notification, nil
}

func (s *NotificationUsecase) CreateBatch(c context.Context, dto *pb_noti.NotiBatchRequest) (*pb_noti.ListResponse, error) {
	metadata := _jwt.MetadataFromContext(c)
	if metadata == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Thông tin không chính xác",
		}
	}

	notifications := []domain.NotificationEntity{}
	for i := 0; i < len(dto.Datas); i++ {
		request := dto.Datas[i]
		var visibleAt *time.Time
		if request.VisibleAt != nil {
			t := request.VisibleAt.AsTime()
			visibleAt = &t
		}
		notifications = append(notifications, domain.NotificationEntity{
			OwnerID:   request.TargetId,
			OwnerType: request.TargetType,
			Title:     request.Title,
			Message:   []string{request.Message},
			Link:      request.Link,
			Type:      enums.TypeEnum(request.Type),
			VisibleAt: visibleAt,
		})
	}

	if err := _db.SaveWithContext(c, &notifications); err != nil {
		return nil, err
	}

	// response := s.mapToDTO(notification)
	return &pb_noti.ListResponse{
		Data: nil,
	}, nil
}

// Search: Tìm kiếm danh sách thông báo
func (s *NotificationUsecase) Search(c context.Context, dto *dto.SearchNotiRequest) ([]domain.NotificationEntity, int32, error) {
	profileId := _utils.GetProfileIdWithContext(c)

	notifications, total, err := s.Repo.Search(c, profileId, dto)
	if err != nil {
		return nil, 0, err
	}
	return notifications, total, nil
}

// Read: Đánh dấu thông báo là đã đọc
func (s *NotificationUsecase) Read(c context.Context, id uint64) (int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.Repo.MarkAsRead(c, profileId, id)
}

func (s *NotificationUsecase) ReadAll(c context.Context) (int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.Repo.ReadAll(c, profileId)
}

// Count: Đếm số thông báo chưa đọc
func (s *NotificationUsecase) Count(c context.Context) int32 {
	profileId := _utils.GetProfileIdWithContext(c)
	count, err := s.Repo.Count(c, profileId)
	if err != nil {
		return 0
	}
	return int32(count)
}

func (s *NotificationUsecase) Remove(c context.Context, id uint64) (uint64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	err := s.Repo.Delete(c, profileId, id)
	return id, err
}

var typeNotiMerged = []enums.TypeEnum{
	enums.LikeComment,
	enums.LikeNewsFeed,
	enums.CommentChildren,
	enums.CommentNewsFeed,
	enums.Follow,
}

// --- Internal ---
func (s *NotificationUsecase) Create(c context.Context, dto *dto.NotiNewRequest) (*domain.NotificationEntity, error) {
	// notifications, err := s.Repo.ListNotiWithTargetIdAndOwnerType(c, dto.TargetID, int32(dto.OwnerType), 4)
	// if err != nil {
	// 	return nil, err
	// }

	log.Println("ismerge=", dto.IsMerge)
	if dto.IsMerge {
		s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, dto.Type, dto.OwnerType)
	}

	message := s.GetMessage(c, enums.TypeEnum(dto.Type), dto.AttachData)
	str, _ := json.Marshal(dto.AttachData)
	notification := domain.NotificationEntity{
		TargetID:  dto.TargetID,
		OwnerType: int32(dto.OwnerType),
		Title:     dto.Title,
		OwnerID:   dto.OwnerID,
		// PrefixMessage: prefixMessage,
		Message:    message,
		AttachData: string(str),
		Link:       dto.Link,
		Type:       enums.TypeEnum(dto.Type),
		VisibleAt:  dto.VisibleAt,
	}

	// todo: delete old notifications
	// if len(notifications) > 0 {
	// 	ids := []uint64{}
	// 	for _, notification := range notifications {
	// 		ids = append(ids, notification.ID)
	// 	}
	// 	if err := s.Repo.DeleteIds(c, ids); err != nil {
	// 		return nil, err
	// 	}
	// }

	if err := s.Repo.Create(c, &notification); err != nil {
		return nil, err
	}

	return &notification, nil
}

func (s *NotificationUsecase) SendBatch(c context.Context, dto *dto.NotiNewRequest) (*domain.NotificationEntity, error) {
	return nil, nil
}

func (s *NotificationUsecase) PushByToken(c context.Context, dto *dto.PushTokenRequest) (*pb_noti.PushResponse, error) {
	return nil, nil
}

func (s *NotificationUsecase) PushByTopic(c context.Context, dto *dto.PushTopicRequest) (*pb_noti.PushResponse, error) {
	return nil, nil
}

func (s *NotificationUsecase) GetMessage(c context.Context, typeMessage enums.TypeEnum, attachData []string) []string {
	if typeMessage == enums.ShareNewsFeed {
		return []string{"", attachData[0], " đã chia sẻ bài viết của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.CommentNewsFeed {
		return []string{"", attachData[0], " đã bình luận bài viết của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.CommentChildren {
		return []string{"", attachData[0], " đã trả lời bình luận của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.LikeNewsFeed {
		return []string{"", attachData[0], " thích bài viết của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.LikeComment {
		return []string{"", attachData[0], " thích bình luận của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.Follow {
		return []string{"", attachData[0], " đã theo dõi bạn"}
	}
	if typeMessage == enums.FriendRequest {
		return []string{"", attachData[0], " đã gửi lời mời kết bạn cho bạn"}
	}
	if typeMessage == enums.FriendRequest {
		return []string{"", attachData[0], " đã đồng ý kết bạn"}
	}
	if typeMessage == enums.JoinGroup {
		return []string{"", attachData[0], " đã thêm bạn vào nhóm: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.JoinOrganization {
		return []string{"", "Quản trị viên đã thêm bạn vào tổ chức: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.Deposit {
		return []string{"", attachData[0], " đã nạp tiền vào tài khoản của bạn"}
	}
	if typeMessage == enums.Withdraw {
		return []string{"", attachData[0], " đã rút tiền từ tài khoản của bạn"}
	}
	if typeMessage == enums.Payment {
		return []string{"", attachData[0], " đã thanh toán cho bạn"}
	}	
	return []string{}
}
