package transformer

import (
	organizationpb "pb/types/organization"
	"time"

	"organization/internal/domain/entity"
)

type GroupLogActivityTransformer interface {
	EntityToGetGroupLogActivitiesResponse(logActivities []*entity.GroupLogActivity, total uint32) *organizationpb.GetGroupLogActivitiesResponse
}

type groupLogActivityTransformer struct{}

func NewGroupLogActivityTransformer() GroupLogActivityTransformer {
	return &groupLogActivityTransformer{}
}

func (t *groupLogActivityTransformer) EntityToGetGroupLogActivitiesResponse(logActivities []*entity.GroupLogActivity, total uint32) *organizationpb.GetGroupLogActivitiesResponse {
	response := &organizationpb.GetGroupLogActivitiesResponse{
		LogActivities: make([]*organizationpb.GroupLogActivity, 0, len(logActivities)),
		Total:         total,
	}

	for _, logActivity := range logActivities {
		response.LogActivities = append(response.LogActivities, &organizationpb.GroupLogActivity{
			Id:      logActivity.Id,
			GroupId: logActivity.GroupId,
			// ActorId:   logActivity.ActorId,
			LogType:   logActivity.LogType,
			LogData:   logActivity.LogData,
			CreatedAt: logActivity.CreatedAt.Format(time.RFC3339),
			UpdatedAt: logActivity.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response
}
