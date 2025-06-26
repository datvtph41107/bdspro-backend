package transformer

import (
	organizationpb "pb/types/organization"
	"time"

	"organization/internal/domain/entity"
)

type GroupNotificationTransformer interface {
	EntityToGetGroupNotificationsResponse(notifications []*entity.GroupNotification) *organizationpb.GetGroupNotificationsResponse
}

type groupNotificationTransformer struct{}

func NewGroupNotificationTransformer() GroupNotificationTransformer {
	return &groupNotificationTransformer{}
}

func (t *groupNotificationTransformer) EntityToGetGroupNotificationsResponse(notifications []*entity.GroupNotification) *organizationpb.GetGroupNotificationsResponse {
	notificationsResponse := make([]*organizationpb.GroupNotification, len(notifications))
	for i, notification := range notifications {
		notificationsResponse[i] = &organizationpb.GroupNotification{
			Id:        notification.Id,
			GroupId:   notification.GroupId,
			Type:      notification.Type,
			Content:   notification.Content,
			CreatedAt: notification.CreatedAt.Format(time.RFC3339),
			UpdatedAt: notification.UpdatedAt.Format(time.RFC3339),
		}
	}
	return &organizationpb.GetGroupNotificationsResponse{
		Notifications: notificationsResponse,
		Total:         uint32(len(notifications)),
	}
}
