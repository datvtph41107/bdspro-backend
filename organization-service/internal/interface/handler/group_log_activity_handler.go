package handler

import (
	"context"
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
)

type GroupLogActivityHandler struct {
	grouppb.UnimplementedGroupLogActivityServiceServer
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

func (h *GroupLogActivityHandler) GetGroupLogActivities(ctx context.Context, req *grouppb.GetGroupLogActivitiesRequest) (*grouppb.GetGroupLogActivitiesResponse, error) {
	if err := h.GroupLogActivityValidator.ValidateGetGroupLogActivitiesRequest(req); err != nil {
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

	logActivities, err := h.GroupLogActivityUsecase.GetGroupLogActivitiesByGroupID(ctx, uint32(req.GroupId))
	if err != nil {
		return nil, err
	}

	return h.GroupLogActivityTransformer.EntityToGetGroupLogActivitiesResponse(logActivities), nil
} 