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