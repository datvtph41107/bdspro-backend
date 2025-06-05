package transformer

import (
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupLogActivityTransformer interface {
	EntityToGetGroupLogActivitiesResponse(logActivities []*entity.GroupLogActivity) *grouppb.GetGroupLogActivitiesResponse
}

type groupLogActivityTransformer struct{}

func NewGroupLogActivityTransformer() GroupLogActivityTransformer {
	return &groupLogActivityTransformer{}
}

func (t *groupLogActivityTransformer) EntityToGetGroupLogActivitiesResponse(logActivities []*entity.GroupLogActivity) *grouppb.GetGroupLogActivitiesResponse {
	response := &grouppb.GetGroupLogActivitiesResponse{
		LogActivities: make([]*grouppb.GroupLogActivity, 0, len(logActivities)),
		Total:         uint32(len(logActivities)),
	}

	for _, logActivity := range logActivities {
		response.LogActivities = append(response.LogActivities, &grouppb.GroupLogActivity{
			Id: logActivity.Id,//todo
			// GroupId:   logActivity.GroupId,
			// Action:    logActivity.Action,
			// CreatedAt: logActivity.CreatedAt.Format(time.RFC3339),
			// UpdatedAt: logActivity.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response
}
