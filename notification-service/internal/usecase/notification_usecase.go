package usecase

import (
	_db "common/db"
	_enum "common/domain/enum"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"log"
	"notification/internal/domain"
	"notification/internal/dto"
	"notification/internal/enums"
	"notification/internal/repo"
	notificationpb "pb/types/notification"
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

	// var visibleAt *time.Time
	// if request.VisibleAt != nil {
	// 	visibleAt = request.VisibleAt
	// }
	notification := &domain.NotificationEntity{
		OwnerID:          profileId,
		OwnerOf:          request.OwnerOf,
		Avatar:           request.Avatar,
		TargetID:         request.TargetID,
		NotificationType: request.NotificationType,
		Title:            request.Title,
		Message:          request.Message,
		AttachData:       request.AttachData,
	}

	if err := _db.SaveWithContext(c, notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *NotificationUsecase) CreateBatch(c context.Context, dto *notificationpb.NotiBatchRequest) (*notificationpb.ListResponse, error) {
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
		// var visibleAt *time.Time
		// if request.VisibleAt != nil {
		// 	t := request.VisibleAt.AsTime()
		// 	visibleAt = &t
		// }
		notifications = append(notifications, domain.NotificationEntity{
			OwnerID:    request.TargetId,
			OwnerOf:    _enum.EOwnerOf(request.OwnerOf),
			Title:      request.Title,
			Message:    request.Message,
			AttachData: request.AttachData,
			// Link:       request.Link,
			// Type:       enums.TypeEnum(request.Type),
			// VisibleAt:  visibleAt,
		})
	}

	if err := _db.SaveWithContext(c, &notifications); err != nil {
		return nil, err
	}

	// response := s.mapToDTO(notification)
	return &notificationpb.ListResponse{
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
		s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, dto.NotificationType, dto.OwnerOf)
	}

	// message := s.GetMessage(c, enums.TypeEnum(dto.Type), dto.AttachData)
	// str, _ := json.Marshal(dto.AttachData)
	notification := domain.NotificationEntity{
		Avatar:   dto.Avatar,
		TargetID: dto.TargetID,
		OwnerOf:  dto.OwnerOf,
		// Title:     dto.Title,
		OwnerID: dto.OwnerID,
		// PrefixMessage: prefixMessage,
		Message:    dto.Message,
		AttachData: dto.AttachData,
		// Link:       dto.Link,
		// Type:       enums.TypeEnum(dto.Type),
		// VisibleAt:  dto.VisibleAt,
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

func (s *NotificationUsecase) CreateToOwner(c context.Context, dto *dto.NotiNewRequest) (*domain.NotificationEntity, error) {
	// notifications, err := s.Repo.ListNotiWithTargetIdAndOwnerType(c, dto.TargetID, int32(dto.OwnerType), 4)
	// if err != nil {
	// 	return nil, err
	// }
	if err := s._removeNotification(c, dto); err != nil {
		return nil, err
	}

	log.Println("ismerge=", dto.IsMerge)
	// if dto.IsMerge {
	// 	s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, dto.NotificationType, dto.OwnerOf)
	// }

	// message := s.GetMessage(c, enums.TypeEnum(dto.Type), dto.AttachData)
	// str, _ := json.Marshal(dto.AttachData)
	notification := domain.NotificationEntity{
		Avatar:  dto.Avatar,
		Title:   dto.Title,
		Message: dto.Message,
		OwnerOf: dto.OwnerOf,
		// Title:     dto.Title,
		OwnerID: dto.OwnerID,
		// PrefixMessage: prefixMessage,
		AttachData:       dto.AttachData,
		NotificationType: dto.NotificationType,
		TargetID:         dto.TargetID,
		// Link:       dto.Link,
		// Type:       enums.TypeEnum(dto.Type),
		// VisibleAt:  dto.VisibleAt,
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

func (s *NotificationUsecase) SendBatch(c context.Context, notifications []*dto.NotiNewRequest) ([]uint64, error) {
	existedUserIds := map[uint64]bool{}
	notificationIds := make([]uint64, 0)
	for _, notification := range notifications {
		if existedUserIds[notification.OwnerID] {
			continue
		}
		noti, err := s.Create(c, notification)
		if err != nil {
			return nil, err
		}
		notificationIds = append(notificationIds, noti.ID)
		existedUserIds[notification.OwnerID] = true
	}

	return notificationIds, nil
}

func (s *NotificationUsecase) PushByToken(c context.Context, dto *dto.PushTokenRequest) (*notificationpb.PushResponse, error) {
	return nil, nil
}

func (s *NotificationUsecase) PushByTopic(c context.Context, dto *dto.PushTopicRequest) (*notificationpb.PushResponse, error) {
	return nil, nil
}

func (s *NotificationUsecase) _removeNotification(c context.Context, dto *dto.NotiNewRequest) error {
	// Kiểm tra xem notification type có cần merge (xóa thông báo cũ) không

	if dto.NotificationType == _enum.NotificationFriendRequest ||
		dto.NotificationType == _enum.NotificationFriendResponse {
		s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, _enum.NotificationFriendRequest, dto.OwnerOf)
		s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, _enum.NotificationFriendResponse, dto.OwnerOf)
	} else if dto.NotificationType.IsMergeableNotification() {
		s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, dto.NotificationType, dto.OwnerOf)
	}
	return nil
}
