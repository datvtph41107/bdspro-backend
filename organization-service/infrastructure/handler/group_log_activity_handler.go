package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
)

type GroupLogActivityHandler struct {
	organizationpb.UnimplementedGroupLogActivityServiceServer
	GroupLogActivityUsecase     usecase.GroupLogActivityUsecase
	GroupLogActivityTransformer transformer.GroupLogActivityTransformer
	GroupLogActivityValidator   validator.GroupLogActivityValidator
}

func NewGroupLogActivityHandler(
	groupLogActivityUsecase usecase.GroupLogActivityUsecase,
	groupLogActivityTransformer transformer.GroupLogActivityTransformer,
	groupLogActivityValidator validator.GroupLogActivityValidator,
) *GroupLogActivityHandler {
	return &GroupLogActivityHandler{
		GroupLogActivityUsecase:     groupLogActivityUsecase,
		GroupLogActivityTransformer: groupLogActivityTransformer,
		GroupLogActivityValidator:   groupLogActivityValidator,
	}
}

// @Summary Lấy danh sách hoạt động của nhóm
// @Description Lấy danh sách hoạt động của nhóm
// @Tags Hoạt động
// @Accept json
// @Produce json
// @Param group_log_activity query organizationpb.GetGroupLogActivitiesRequest true "Thông tin hoạt động"
// @Param groupId path int true "ID của nhóm"
// @Security BearerAuth
// @Router /group/log-activity/{groupId} [get]
func (h *GroupLogActivityHandler) GetGroupLogActivities(ctx context.Context, req *organizationpb.GetGroupLogActivitiesRequest) (*organizationpb.GetGroupLogActivitiesResponse, error) {
	if err := h.GroupLogActivityValidator.ValidateGetGroupLogActivitiesRequest(req); err != nil {
		return nil, err
	}

	page := 0
	size := 10

	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	logActivities, total, err := h.GroupLogActivityUsecase.GetGroupLogActivitiesByGroupID(ctx, uint32(req.GroupId), page, size)
	if err != nil {
		return nil, err
	}

	return h.GroupLogActivityTransformer.EntityToGetGroupLogActivitiesResponse(logActivities, total), nil
}
