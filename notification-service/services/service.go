package services

import (
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"log"
	"noti/config"
	"noti/data"
	"noti/enums"
	"noti/models"
	"noti/repo"
	pb_noti "pb/types/notification"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationService2 struct {
	Repo        *repo.NotificationRepository
	firebaseApp *firebase.App
}

func NewNotificationService2(
	repo *repo.NotificationRepository,
	firebaseConfig *config.FirebaseConfig,
) *NotificationService2 {
	return &NotificationService2{
		Repo:        repo,
		firebaseApp: firebaseConfig.App,
	}
}

// CreateNotification: Tạo thông báo mới
func (s *NotificationService2) CreateNotification(c context.Context, request *pb_noti.NotiNewRequest) (*pb_noti.NotiNewResponse, error) {
	metadata := _jwt.MetadataFromContext(c)
	if metadata == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Thông tin không chính xác",
		}
	}

	notification := models.NotificationEntity{
		UserID:  request.UserID,
		Title:   request.Title,
		Message: request.Message,
		Link:    request.Link,
		Type:    enums.TypeEnum(request.Type),
	}

	if err := _db.SaveWithContext(c, &notification); err != nil {
		return nil, err
	}

	// response := s.mapToDTO(notification)
	return &pb_noti.NotiNewResponse{
		Success: true,
	}, nil
}

func (s *NotificationService2) CreateBatch(c context.Context, dto *pb_noti.NotiBatchRequest) (*pb_noti.ListResponse, error) {
	metadata := _jwt.MetadataFromContext(c)
	if metadata == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Thông tin không chính xác",
		}
	}

	notifications := []models.NotificationEntity{}
	for i := 0; i < len(dto.Datas); i++ {
		request := dto.Datas[i]
		notifications = append(notifications, models.NotificationEntity{
			UserID:  request.UserID,
			Title:   request.Title,
			Message: request.Message,
			Link:    request.Link,
			Type:    enums.TypeEnum(request.Type),
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
func (s *NotificationService2) Search(c context.Context, dto *pb_noti.ListRequest) (*pb_noti.ListResponse, error) {
	metadata := _jwt.MetadataFromContext(c)
	if metadata.ProfileID == nil {
		return nil, nil
	}
	// dto.UserID = metadata.ProfileID
	query := data.SearchNotiRequest{
		UserID: metadata.ProfileID,
	}
	pageable := _utils.Page{
		Skip:  int(dto.Page * dto.Size),
		Limit: int(dto.Size),
	}
	notifications, err := s.Repo.Search(query, pageable)
	if err != nil {
		return nil, err
	}

	var responses []*pb_noti.Notification
	for _, noti := range notifications {
		responses = append(responses, s.mapToDTO(noti))
	}
	return &pb_noti.ListResponse{
		Data: responses,
	}, nil
}

// Read: Đánh dấu thông báo là đã đọc
func (s *NotificationService2) Read(c context.Context, id uint64) (int64, error) {
	metadata := _jwt.MetadataFromContext(c)
	if metadata.ProfileID == nil {
		return 0, nil
	}
	return s.Repo.MarkNotificationAsRead(id, *metadata.ProfileID)
}

// Count: Đếm số thông báo chưa đọc
func (s *NotificationService2) Count(c context.Context) (int64, error) {
	metadata := _jwt.MetadataFromContext(c)
	if metadata.ProfileID == nil {
		return 0, nil
	}
	return s.Repo.Count(*metadata.ProfileID)
}

// mapToDTO: Chuyển đổi NotificationEntity thành NotificationResponse
func (s *NotificationService2) mapToDTO(entity models.NotificationEntity) *pb_noti.Notification {
	return &pb_noti.Notification{
		Id:        entity.ID,
		UserID:    entity.UserID,
		Title:     entity.Title,
		Message:   entity.Message,
		Link:      entity.Link,
		Type:      int32(entity.Type),
		IsRead:    entity.IsRead,
		Timestamp: timestamppb.New(*entity.CreatedAt),
	}
}

func (s *NotificationService2) Delete(c context.Context, dto *pb_noti.IdRequest) (int, error) {
	metadata := _jwt.MetadataFromContext(c)
	if metadata.ProfileID == nil {
		return 0, nil
	}
	return 1, s.Repo.DeleteByID(*metadata.ProfileID, dto.Id)
}

// SendToTopic gửi thông báo đến một topic
func (s *NotificationService2) SendToTopic(c context.Context, request *pb_noti.PushTopicRequest) (*pb_noti.PushResponse, error) {
	var entity models.HistoryEntity
	copier.Copy(&entity, &request)
	ctx := context.Background()

	notification := &messaging.Notification{
		Title: request.Title,
		Body:  request.Body,
		// ImageURL: request.Image,
	}
	messageBody := &messaging.Message{
		Topic:        request.Topic,
		Notification: notification,
	}
	client, err := s.firebaseApp.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	result, err := client.Send(ctx, messageBody)
	if err != nil {
		log.Println("Error sending message:", err)
		entity.Error = err.Error()
	} else {
		entity.Response = result
	}

	_db.SaveWithContext(c, entity)
	return &pb_noti.PushResponse{
		MessageId: result,
	}, nil
}

// SendToToken gửi thông báo đến một thiết bị dựa vào token
func (s *NotificationService2) SendToToken(c context.Context, request *pb_noti.PushTokenRequest) (*pb_noti.PushResponse, error) {
	var entity models.HistoryEntity
	copier.Copy(&entity, &request)
	ctx := context.Background()

	notification := &messaging.Notification{
		Title: request.Title,
		Body:  request.Body,
		// ImageURL: request.Image,
	}
	message := &messaging.Message{
		Token:        request.Token,
		Notification: notification,
	}

	client, err := s.firebaseApp.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	result, err := client.Send(ctx, message)
	if err != nil {
		log.Println("Error sending message:", err)
		entity.Error = err.Error()
	} else {
		entity.Response = result
	}

	_db.SaveWithContext(c, &entity)
	return &pb_noti.PushResponse{
		MessageId: result,
	}, nil
}
