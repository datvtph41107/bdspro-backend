package usecase

import (
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"noti/enums"
	"noti/internal/domain"
	"noti/internal/dto"
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
	// metadata := _jwt.MetadataFromContext(c)
	// if metadata == nil {
	// 	return nil, &_routes.Except{
	// 		Code:    400,
	// 		Message: "Thông tin không chính xác",
	// 	}
	// }

	var visibleAt *time.Time
	if request.VisibleAt != nil {
		visibleAt = request.VisibleAt
	}
	notification := domain.NotificationEntity{
		TargetID:  profileId,
		OwnerType: request.OwnerType,
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
			TargetID:  request.TargetId,
			OwnerType: request.TargetType,
			Title:     request.Title,
			Message:   request.Message,
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

// --- Internal ---
func (s *NotificationUsecase) Create(c context.Context, dto *dto.NotiNewRequest) (*domain.NotificationEntity, error) {
	return nil, nil
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
