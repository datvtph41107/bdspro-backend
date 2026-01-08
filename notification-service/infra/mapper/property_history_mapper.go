package mapper

import (
	_utils "common/utils"
	"encoding/json"
	"time"

	"notification/internal/domain"
	"notification/internal/enums"
	notificationpb "pb/types/notification"

	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
)

func RequestToPrototypeHistoryEntity(
	req *notificationpb.CreatePropertyHistoryRequest,
) (*domain.PropertyHistoryEntity, error) {

	metaMap := map[string]interface{}{}
	if req.Metadata != nil {
		metaMap = req.Metadata.AsMap()
	}

	metaBytes, err := json.Marshal(metaMap)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	return &domain.PropertyHistoryEntity{
		SubjectID:   req.SubjectId,
		SubjectType: req.SubjectType,
		Description: req.Description,
		ActorID:     req.ActorId,
		Metadata:    datatypes.JSON(metaBytes),
		Action:      enums.PropertyAction(req.Action),
		ActorType:   enums.ActorType(req.ActorType),
		OccurredAt:  &now,
	}, nil
}

func EntityToActivityHistoryItem(
	entity *domain.PropertyHistoryEntity,
) (*notificationpb.PropertyHistoryItem, error) {

	// ---- metadata: datatypes.JSON -> structpb.Struct
	metaMap := map[string]interface{}{}
	if len(entity.Metadata) > 0 {
		if err := json.Unmarshal(entity.Metadata, &metaMap); err != nil {
			return nil, err
		}
	}

	metaStruct, err := structpb.NewStruct(metaMap)
	if err != nil {
		return nil, err
	}

	return &notificationpb.PropertyHistoryItem{
		Id:          entity.ID,
		SubjectId:   entity.SubjectID,
		SubjectType: entity.SubjectType,
		Action:      string(entity.Action),
		Description: entity.Description,
		Metadata:    metaStruct,

		ActorId:    entity.ActorID,
		ActorType:  string(entity.ActorType),
		OccurredAt: _utils.FormatTimeToString(entity.OccurredAt),
		CreatedAt:  _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:  _utils.FormatTimeToString(entity.UpdatedAt),
	}, nil
}

func EntitiesToActivityHistoryResponse(
	entities []*domain.PropertyHistoryEntity,
	total uint32,
	page, size int32,
) (*notificationpb.GetPropertyHistoryResponse, error) {

	items := make([]*notificationpb.PropertyHistoryItem, 0, len(entities))

	for _, e := range entities {
		item, err := EntityToActivityHistoryItem(e)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &notificationpb.GetPropertyHistoryResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

func EntityToActivityHistoryByIDResponse(
	entity *domain.PropertyHistoryEntity,
) (*notificationpb.GetPropertyHistoryByIDResponse, error) {

	item, err := EntityToActivityHistoryItem(entity)
	if err != nil {
		return nil, err
	}

	return &notificationpb.GetPropertyHistoryByIDResponse{
		Data: item,
	}, nil
}
