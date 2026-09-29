package mapper

import (
	_utils "common/utils"
	"notification/internal/domain"
	notificationpb "pb/types/notification"

	"google.golang.org/protobuf/types/known/structpb"
)

type HistoryAuthMapper struct{}

func NewHistoryAuthMapper() *HistoryAuthMapper {
	return &HistoryAuthMapper{}
}

func (m *HistoryAuthMapper) DomainToPb(entity *domain.HistoryAuthEntity) *notificationpb.HistoryAuthDTO {
	if entity == nil {
		return nil
	}

	var metadataStruct *structpb.Struct
	if entity.Metadata != nil {
		metadataStruct, _ = structpb.NewStruct(entity.Metadata)
	}

	response := &notificationpb.HistoryAuthDTO{
		Id:             entity.ID,
		UserId:         entity.UserID,
		ActionType:     entity.ActionType,
		ActionName:     entity.ActionName,
		Description:    entity.Description,
		Success:        entity.Success,
		Reason:         entity.Reason,
		IpAddress:      entity.IPAddress,
		UserAgent:      entity.UserAgent,
		Metadata:       metadataStruct,
		SessionId:      entity.SessionID,
		Channel:        entity.Channel,
		DeviceId:       entity.DeviceID,
		Location:       entity.Location,
		AdditionalNote: entity.AdditionalNote,
		SourceService:  entity.SourceService,
	}

	if entity.OrganizationID != nil {
		response.OrganizationId = entity.OrganizationID
	}
	if entity.PerformedBy != nil {
		response.PerformedBy = entity.PerformedBy
	}
	if entity.CreatedAt != nil {
		response.CreatedAt = _utils.FormatTimeToString(entity.CreatedAt)
	}

	return response
}

func (m *HistoryAuthMapper) DomainsToPb(entities []domain.HistoryAuthEntity) []*notificationpb.HistoryAuthDTO {
	results := make([]*notificationpb.HistoryAuthDTO, 0, len(entities))
	for i := range entities {
		results = append(results, m.DomainToPb(&entities[i]))
	}
	return results
}
