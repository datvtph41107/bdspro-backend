package transformer

import (
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
)

type OrganizationLogActivityTransformer interface {
	EntityToGetOrganizationLogActivityByIDResponse(log *entity.OrganizationLogActivity) *organizationpb.GetOrganizationLogActivityByIDResponse
	EntityToGetOrganizationLogActivitiesByOrganizationIDResponse(logs []*entity.OrganizationLogActivity, total uint32) *organizationpb.GetOrganizationLogActivitiesByOrganizationIDResponse
	EntityToGetOrganizationLogActivitiesByActorIDResponse(logs []*entity.OrganizationLogActivity, total uint32) *organizationpb.GetOrganizationLogActivitiesByActorIDResponse
}

type organizationLogActivityTransformer struct{}

func NewOrganizationLogActivityTransformer() OrganizationLogActivityTransformer {
	return &organizationLogActivityTransformer{}
}

func (t *organizationLogActivityTransformer) EntityToGetOrganizationLogActivityByIDResponse(log *entity.OrganizationLogActivity) *organizationpb.GetOrganizationLogActivityByIDResponse {
	if log == nil {
		return nil
	}
	return &organizationpb.GetOrganizationLogActivityByIDResponse{
		Log: &organizationpb.OrganizationLogActivity{
			Id:             log.Id,
			CreatedAt:      log.CreatedAt.Unix(),
			UpdatedAt:      log.UpdatedAt.Unix(),
			CreatedBy:      log.CreatedBy,
			UpdatedBy:      log.UpdatedBy,
			OrganizationId: log.OrganizationId,
			ActorId:        log.ActorId,
			LogType:        log.LogType,
			LogData:        log.LogData,
		},
	}
}

func (t *organizationLogActivityTransformer) EntityToGetOrganizationLogActivitiesByOrganizationIDResponse(logs []*entity.OrganizationLogActivity, total uint32) *organizationpb.GetOrganizationLogActivitiesByOrganizationIDResponse {
	if logs == nil {
		return nil
	}
	response := &organizationpb.GetOrganizationLogActivitiesByOrganizationIDResponse{
		Logs:  make([]*organizationpb.OrganizationLogActivity, 0, len(logs)),
		Total: total,
	}
	for _, log := range logs {
		response.Logs = append(response.Logs, &organizationpb.OrganizationLogActivity{
			Id:             log.Id,
			CreatedAt:      log.CreatedAt.Unix(),
			UpdatedAt:      log.UpdatedAt.Unix(),
			CreatedBy:      log.CreatedBy,
			UpdatedBy:      log.UpdatedBy,
			OrganizationId: log.OrganizationId,
			ActorId:        log.ActorId,
			LogType:        log.LogType,
			LogData:        log.LogData,
		})
	}
	return response
}

func (t *organizationLogActivityTransformer) EntityToGetOrganizationLogActivitiesByActorIDResponse(logs []*entity.OrganizationLogActivity, total uint32) *organizationpb.GetOrganizationLogActivitiesByActorIDResponse {
	if logs == nil {
		return nil
	}
	response := &organizationpb.GetOrganizationLogActivitiesByActorIDResponse{
		Logs:  make([]*organizationpb.OrganizationLogActivity, 0, len(logs)),
		Total: total,
	}
	for _, log := range logs {
		response.Logs = append(response.Logs, &organizationpb.OrganizationLogActivity{
			Id:             log.Id,
			CreatedAt:      log.CreatedAt.Unix(),
			UpdatedAt:      log.UpdatedAt.Unix(),
			CreatedBy:      log.CreatedBy,
			UpdatedBy:      log.UpdatedBy,
			OrganizationId: log.OrganizationId,
			ActorId:        log.ActorId,
			LogType:        log.LogType,
			LogData:        log.LogData,
		})
	}
	return response
}
