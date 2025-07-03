package mapper

import (
	"noti/internal/domain"
	notipb "pb/types/notification"
)

type NotificationMapper struct {
}

func NewNotificationMapper() *NotificationMapper {
	return &NotificationMapper{}
}

func (s *NotificationMapper) MapToPB(noti *domain.NotificationEntity) *notipb.NotificationDTO {
	return &notipb.NotificationDTO{
		Id:      noti.ID,
		Title:   noti.Title,
		Message: noti.Message,
		Link:    noti.Link,
		Type:    int32(noti.Type),
	}
}
