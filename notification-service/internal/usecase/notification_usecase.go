package usecase

import (
	_enum "common/domain/enum"
	_errors "common/errors"
	_jwt "common/jwt"
	"common/logging"
	_utils "common/utils"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"notification/infra/cache"
	"notification/internal"
	"notification/internal/domain"
	"notification/internal/dto"
	"notification/internal/enums"
	notificationpb "pb/types/notification"
	"strconv"
	"strings"
)

type NotificationUsecase struct {
	Repo             NotificationStore
	RedisClient      *cache.RedisClient
	AuthClient       PushTokenResolver
	FirebaseProvider PushSender
}

func NewNotificationService(
	repo NotificationStore,
	redisClient *cache.RedisClient,
	authClient PushTokenResolver,
	firebaseProvider PushSender,
) *NotificationUsecase {
	return &NotificationUsecase{
		Repo:             repo,
		RedisClient:      redisClient,
		AuthClient:       authClient,
		FirebaseProvider: firebaseProvider,
	}
}

// --- Gateway ---
// CreateNotification: Tạo thông báo mới
func (s *NotificationUsecase) CreateNotification(c context.Context, request *dto.NotiNewRequest) (*domain.NotificationEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	if profileId == 0 {
		return nil, _errors.ReturnError(service.NotificationOwnerIdentityRequired)
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

	if err := s.Repo.Create(c, notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *NotificationUsecase) CreateBatch(c context.Context, dto *notificationpb.NotiBatchRequest) (*notificationpb.ListResponse, error) {
	metadata := _jwt.MetadataFromContext(c)
	if metadata == nil {
		return nil, _errors.ReturnError(service.NotificationMetadataRequired)
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

	if err := s.Repo.CreateBatch(c, notifications); err != nil {
		return nil, err
	}

	// Publish notification events to Redis for each notification
	for i := range notifications {
		s.publishNotificationEvent(c, &notifications[i])
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

func (s *NotificationUsecase) ReadMany(c context.Context, ids []uint64) (int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.Repo.MarkAsReadMany(c, profileId, ids)
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

	logging.WithComponent(c, "notification.usecase").Debug(
		"notification merge decision",
		slog.Bool("notification.is_merge", dto.IsMerge),
	)
	// if dto.IsMerge {
	// 	// s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, dto.NotificationType, dto.OwnerOf)
	// }
	err := s._removeNotification(c, dto)
	if err != nil {
		return nil, err
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

		Title:            dto.Title,
		NotificationType: _enum.ENotificationType(dto.Type),
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

func (s *NotificationUsecase) RemoveToOwner(c context.Context, dto *dto.NotiNewRequest) error {
	err := s._removeNotification(c, dto)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationUsecase) CreateToOwner(c context.Context, dto *dto.NotiNewRequest) (*domain.NotificationEntity, error) {
	// notifications, err := s.Repo.ListNotiWithTargetIdAndOwnerType(c, dto.TargetID, int32(dto.OwnerType), 4)
	// if err != nil {
	// 	return nil, err
	// }
	if err := s._removeNotification(c, dto); err != nil {
		return nil, err
	}

	logging.WithComponent(c, "notification.usecase").Debug(
		"notification merge decision",
		slog.Bool("notification.is_merge", dto.IsMerge),
	)
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

	// Publish notification event to Redis
	s.publishNotificationEvent(c, &notification)

	// Check user online status and send push notification if offline
	s.checkAndSendPushNotification(c, &notification)

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

// PushByToken gửi push notification đến một token cụ thể (chỉ push, không lưu vào DB)
func (s *NotificationUsecase) PushByToken(c context.Context, dto *dto.PushTokenRequest) (*notificationpb.PushResponse, error) {
	if s.FirebaseProvider == nil {
		return nil, fmt.Errorf("FirebaseProvider is not initialized")
	}

	// Gửi push notification đến token
	err := s.FirebaseProvider.SendPushNotification(c, []string{dto.Token}, dto.Title, dto.Body, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send push notification: %w", err)
	}

	return &notificationpb.PushResponse{
		MessageId: "success",
	}, nil
}

// PushByTopic gửi push notification đến một topic (chỉ push, không lưu vào DB)
// Note: Cần implement method SendToTopic trong FirebaseProvider để hỗ trợ topic
func (s *NotificationUsecase) PushByTopic(c context.Context, dto *dto.PushTopicRequest) (*notificationpb.PushResponse, error) {
	// TODO: Implement SendToTopic trong FirebaseProvider
	return nil, fmt.Errorf("PushByTopic not yet implemented - need SendToTopic method in FirebaseProvider")
}

// PushByUserId gửi push notification đến user theo userId (lấy token tự động, chỉ push, không lưu vào DB)
func (s *NotificationUsecase) PushByUserId(c context.Context, dto *dto.PushByUserIdRequest) (*notificationpb.PushResponse, error) {
	if s.AuthClient == nil {
		return nil, fmt.Errorf("AuthClient is not initialized")
	}

	if s.FirebaseProvider == nil {
		return nil, fmt.Errorf("FirebaseProvider is not initialized")
	}

	// Lấy push tokens từ auth service
	pushTokens, err := s.AuthClient.GetPushTokensByProfileId(c, dto.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get push tokens for user %d: %w", dto.UserID, err)
	}

	if len(pushTokens) == 0 {
		return &notificationpb.PushResponse{
			MessageId: "no_tokens",
		}, nil
	}

	// Gửi push notification
	err = s.FirebaseProvider.SendPushNotification(c, pushTokens, dto.Title, dto.Body, dto.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to send push notification: %w", err)
	}

	return &notificationpb.PushResponse{
		MessageId: "success",
	}, nil
}

func (s *NotificationUsecase) _removeNotification(c context.Context, dto *dto.NotiNewRequest) error {
	// Kiểm tra xem notification type có cần merge (xóa thông báo cũ) không

	if dto.Type == enums.TypeEnum(_enum.NotificationDealMemberWithdrawn) ||
		dto.Type == enums.TypeEnum(_enum.NotificationDealInvitationAccepted) ||
		dto.Type == enums.TypeEnum(_enum.NotificationDealInvitationRejected) {
		s.Repo.RemoveUnique2(c, dto.OwnerID, nil, _enum.NotificationDealInvitation, dto.OwnerOf, []string{"", dto.AttachData[1]})
		s.Repo.RemoveUnique2(c, dto.OwnerID, nil, _enum.NotificationDealInvitationAccepted, dto.OwnerOf, []string{"", dto.AttachData[1]})
		s.Repo.RemoveUnique2(c, dto.OwnerID, nil, _enum.NotificationDealInvitationRejected, dto.OwnerOf, []string{"", dto.AttachData[1]})
	} else if (dto.Type == enums.TypeEnum(_enum.NotificationDealInvitation) ||
		dto.Type == enums.TypeEnum(_enum.NotificationDealMemberRemoved)) &&
		len(dto.AttachData) > 1 {
		s.Repo.RemoveUnique2(c, dto.OwnerID, nil, _enum.ENotificationType(dto.Type), dto.OwnerOf, []string{"", dto.AttachData[1]})
	} else if dto.Type == enums.TypeEnum(_enum.NotificationFriendRequest) ||
		dto.Type == enums.TypeEnum(_enum.NotificationFriendResponse) {
		s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, _enum.NotificationFriendRequest, dto.OwnerOf)
		s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, _enum.NotificationFriendResponse, dto.OwnerOf)
	} else if dto.NotificationType.IsMergeableNotification() {
		s.Repo.RemoveUnique(c, dto.OwnerID, dto.TargetID, dto.NotificationType, dto.OwnerOf)
	}
	return nil
}

// publishNotificationEvent publishes notification event to Redis for relay service
func (s *NotificationUsecase) publishNotificationEvent(ctx context.Context, notification *domain.NotificationEntity) {
	if s.RedisClient == nil {
		logging.WithComponent(ctx, "notification.usecase").Warn(
			"notification event publish skipped",
			slog.String("reason", "redis-client-unavailable"),
		)
		return
	}

	// Create notification event message
	eventData := map[string]interface{}{
		"id":               notification.ID,
		"ownerId":          notification.OwnerID,
		"ownerOf":          int32(notification.OwnerOf),
		"title":            notification.Title,
		"message":          notification.Message,
		"attachData":       notification.AttachData,
		"targetId":         notification.TargetID,
		"notificationType": int32(notification.NotificationType),
		"avatar":           notification.Avatar,
		"isRead":           notification.IsRead,
		"createdAt":        _utils.FormatTimeToString(notification.CreatedAt),
	}

	eventMessage := map[string]interface{}{
		"type": "notification",
		"data": eventData,
	}

	messageBytes, err := json.Marshal(eventMessage)
	if err != nil {
		logging.WithComponent(ctx, "notification.usecase").Error(
			"marshal notification event",
			slog.Uint64("notification.owner_id", notification.OwnerID),
			slog.Any("error", err),
		)
		return
	}

	// Publish to Redis channel: notification:{ownerId}
	channel := "notification:" + strconv.FormatUint(notification.OwnerID, 10)
	if err := s.RedisClient.Publish(ctx, channel, string(messageBytes)); err != nil {
		logging.WithComponent(ctx, "notification.usecase").Error(
			"publish notification event to Redis",
			slog.String("redis.channel", channel),
			slog.Any("error", err),
		)
	}
}

// checkAndSendPushNotification kiểm tra trạng thái online của user và gửi push notification nếu offline
func (s *NotificationUsecase) checkAndSendPushNotification(ctx context.Context, notification *domain.NotificationEntity) {
	if s.RedisClient == nil {
		logging.WithComponent(ctx, "notification.usecase").Warn(
			"notification online-status check skipped",
			slog.String("reason", "redis-client-unavailable"),
		)
		return
	}

	if s.AuthClient == nil {
		logging.WithComponent(ctx, "notification.usecase").Warn(
			"notification push skipped",
			slog.String("reason", "auth-client-unavailable"),
		)
		return
	}

	// Check online status từ Redis
	const USER_ONLINE_STATUS_KEY_PREFIX = "user:online:"
	onlineKey := USER_ONLINE_STATUS_KEY_PREFIX + strconv.FormatUint(notification.OwnerID, 10)
	statusStr, err := s.RedisClient.Get(ctx, onlineKey)
	if err != nil {
		logging.WithComponent(ctx, "notification.usecase").Error(
			"get notification owner online status",
			slog.Uint64("notification.owner_id", notification.OwnerID),
			slog.Any("error", err),
		)
		return
	}

	// Parse status: format "true|timestamp" hoặc "false|timestamp"
	isOnline := false
	if statusStr != "" {
		parts := strings.Split(statusStr, "|")
		if len(parts) > 0 && parts[0] == "true" {
			isOnline = true
		}
	}

	// Nếu user online, không cần gửi push notification (đã nhận qua WebSocket)
	if isOnline {
		logging.WithComponent(ctx, "notification.usecase").Info(
			"notification push skipped for online owner",
			slog.Uint64("notification.owner_id", notification.OwnerID),
		)
		return
	}

	// User offline, lấy push token và gửi push notification
	logging.WithComponent(ctx, "notification.usecase").Info(
		"notification owner offline; resolving push tokens",
		slog.Uint64("notification.owner_id", notification.OwnerID),
	)
	pushTokens, err := s.AuthClient.GetPushTokensByProfileId(ctx, notification.OwnerID)
	if err != nil {
		logging.WithComponent(ctx, "notification.usecase").Error(
			"get notification owner push tokens",
			slog.Uint64("notification.owner_id", notification.OwnerID),
			slog.Any("error", err),
		)
		return
	}

	if len(pushTokens) == 0 {
		logging.WithComponent(ctx, "notification.usecase").Info(
			"notification push skipped; no push tokens",
			slog.Uint64("notification.owner_id", notification.OwnerID),
		)
		return
	}

	// Gửi push notification qua Firebase
	if s.FirebaseProvider == nil {
		logging.WithComponent(ctx, "notification.usecase").Warn(
			"notification push skipped",
			slog.String("reason", "firebase-provider-unavailable"),
			slog.Uint64("notification.owner_id", notification.OwnerID),
		)
		return
	}

	// Tạo title và message từ notification
	title := notification.Title
	if title == "" {
		title = "Thông báo mới"
	}

	message := ""
	if len(notification.Message) > 0 {
		message = notification.Message[0]
		if len(notification.Message) > 1 {
			// Nếu có nhiều message, join lại
			message = strings.Join(notification.Message, " ")
		}
	}

	// Tạo data payload với đầy đủ thông tin
	data := make(map[string]string)
	data["notificationId"] = strconv.FormatUint(notification.ID, 10)
	data["targetId"] = strconv.FormatUint(notification.TargetID, 10)
	data["ownerId"] = strconv.FormatUint(notification.OwnerID, 10)
	data["ownerOf"] = strconv.FormatInt(int64(notification.OwnerOf), 10)
	data["type"] = strconv.FormatInt(int64(notification.NotificationType), 10)
	data["notificationType"] = strconv.FormatInt(int64(notification.NotificationType), 10) // Giữ lại để backward compatibility

	// Avatar - luôn thêm vào data (có thể là empty string)
	data["avatar"] = notification.Avatar

	// AttachData - nếu có thì join lại thành string
	if len(notification.AttachData) > 0 {
		data["attachData"] = strings.Join(notification.AttachData, ",")
	}

	// Title và message cũng thêm vào data để app có thể xử lý
	if title != "" {
		data["title"] = title
	}
	if message != "" {
		data["message"] = message
	}

	// Gửi push notification
	if err := s.FirebaseProvider.SendPushNotification(ctx, pushTokens, title, message, data); err != nil {
		logging.WithComponent(ctx, "notification.usecase").Error(
			"send notification push",
			slog.Uint64("notification.owner_id", notification.OwnerID),
			slog.Any("error", err),
		)
	} else {
		logging.WithComponent(ctx, "notification.usecase").Info(
			"notification push sent",
			slog.Uint64("notification.owner_id", notification.OwnerID),
			slog.Int("push.token_count", len(pushTokens)),
		)
	}
}
