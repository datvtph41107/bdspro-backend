package transformer

import (
	grouppb "pb/types/group"
	"time"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupNotificationTransformer interface {
	EntityToGetGroupNotificationsResponse(notifications []*entity.GroupNotification) *grouppb.GetGroupNotificationsResponse
}

type groupNotificationTransformer struct{}

func NewGroupNotificationTransformer() GroupNotificationTransformer {
	return &groupNotificationTransformer{}
}

func (t *groupNotificationTransformer) EntityToGetGroupNotificationsResponse(notifications []*entity.GroupNotification) *grouppb.GetGroupNotificationsResponse {
	notificationsResponse := make([]*grouppb.GroupNotification, len(notifications))
	for i, notification := range notifications {
		notificationsResponse[i] = &grouppb.GroupNotification{
			Id:        notification.Id,
			GroupId:   notification.GroupId,
			Type:      notification.Type,
			Content:   notification.Content,
			CreatedAt: notification.CreatedAt.Format(time.RFC3339),
			UpdatedAt: notification.UpdatedAt.Format(time.RFC3339),
		}
	}
	return &grouppb.GetGroupNotificationsResponse{
		Notifications: notificationsResponse,
		Total:        uint32(len(notifications)),
	}
} 