package mapper

import (
	_utils "common/utils"
	"notification/internal/domain"
	notipb "pb/types/notification"
)

type NotificationMapper struct {
}

func NewNotificationMapper() *NotificationMapper {
	return &NotificationMapper{}
}

func (s *NotificationMapper) MapToPB(noti *domain.NotificationEntity) *notipb.NotificationDTO {
	return &notipb.NotificationDTO{
		Id:               noti.ID,
		Avatar:           noti.Avatar,
		Title:            noti.Title,
		Message:          noti.Message,
		Timestamp:        _utils.FormatTimeToString(noti.CreatedAt),
		NotificationType: uint32(noti.NotificationType),
		TargetId:         &noti.TargetID,
		IsRead:           noti.IsRead,
		// Link:    noti.Link,
		// Type:    int32(noti.Type),
	}
}

// // -- mapping ---
// func (s *NotificationHandler) mapToPB(noti *domain.NotificationEntity) *notipb.NotificationDTO {
// 	return &notipb.NotificationDTO{
// 		Id:      noti.ID,
// 		UserId:  noti.OwnerID,
// 		Title:   noti.Title,
// 		Message: noti.Message,
// 		// Link:      noti.Link,
// 		Type:      int32(noti.NotificationType),
// 		Timestamp: _utils.FormatTimeToString(noti.CreatedAt),
// 		IsRead:    noti.IsRead,
// 		TargetId:  &noti.TargetID,
// 	}
// }
