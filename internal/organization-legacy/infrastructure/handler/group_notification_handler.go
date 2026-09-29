package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
)

type GroupNotificationHandler struct {
	organizationpb.UnimplementedGroupNotificationServiceServer
	GroupNotificationUsecase     usecase.GroupNotificationUsecase
	GroupNotificationTransformer transformer.GroupNotificationTransformer
	GroupNotificationValidator   validator.GroupNotificationValidator
}

func NewGroupNotificationHandler(
	groupNotificationUsecase usecase.GroupNotificationUsecase,
	groupNotificationTransformer transformer.GroupNotificationTransformer,
	groupNotificationValidator validator.GroupNotificationValidator,
) *GroupNotificationHandler {
	return &GroupNotificationHandler{
		GroupNotificationUsecase:     groupNotificationUsecase,
		GroupNotificationTransformer: groupNotificationTransformer,
		GroupNotificationValidator:   groupNotificationValidator,
	}
}

func (h *GroupNotificationHandler) GetGroupNotifications(ctx context.Context, req *organizationpb.GetGroupNotificationsRequest) (*organizationpb.GetGroupNotificationsResponse, error) {
	if err := h.GroupNotificationValidator.ValidateGetGroupNotificationsRequest(req); err != nil {
		return nil, err
	}

	page := int(*req.Page)
	size := int(*req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	notifications, err := h.GroupNotificationUsecase.GetGroupNotificationsByGroupID(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	return h.GroupNotificationTransformer.EntityToGetGroupNotificationsResponse(notifications), nil
}
