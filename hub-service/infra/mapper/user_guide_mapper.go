package mapper

import (
	_utils "common/utils"
	"hub/internal/domain"
	hubpb "pb/types/hub"
)

type UserGuideMapper struct{}

func NewUserGuideMapper() *UserGuideMapper {
	return &UserGuideMapper{}
}

// StepEntityToProto converts UserGuideStepEntity to UserGuideStep protobuf message
func (m *UserGuideMapper) StepEntityToProto(step *domain.UserGuideStepEntity) *hubpb.UserGuideStep {
	if step == nil {
		return nil
	}

	return &hubpb.UserGuideStep{
		Id:          step.ID,
		UserGuideId: step.UserGuideID,
		StepOrder:   int32(step.StepOrder),
		Image:       step.Image,
		Content:     step.Content,
	}
}

// StepEntitiesToProto converts slice of UserGuideStepEntity to UserGuideStep protobuf messages
func (m *UserGuideMapper) StepEntitiesToProto(steps []*domain.UserGuideStepEntity) []*hubpb.UserGuideStep {
	if steps == nil {
		return []*hubpb.UserGuideStep{}
	}

	result := make([]*hubpb.UserGuideStep, 0, len(steps))
	for _, step := range steps {
		result = append(result, m.StepEntityToProto(step))
	}
	return result
}

// StepInputToEntity converts UserGuideStepInput to UserGuideStepEntity
func (m *UserGuideMapper) StepInputToEntity(input *hubpb.UserGuideStepInput, userGuideID uint64) *domain.UserGuideStepEntity {
	if input == nil {
		return nil
	}

	return &domain.UserGuideStepEntity{
		UserGuideID: userGuideID,
		StepOrder:   int(input.StepOrder),
		Image:       input.Image,
		Content:     input.Content,
	}
}

// EntityToProto converts UserGuideEntity to UserGuide protobuf message
// Description is truncated to 50 characters for list view
func (m *UserGuideMapper) EntityToProto(entity *domain.UserGuideEntity) *hubpb.UserGuide {
	if entity == nil {
		return nil
	}

	// Truncate description to 50 characters for list view
	description := entity.Description
	if len(description) > 50 {
		// Cắt tại 50 ký tự và thêm "..."
		runes := []rune(description)
		if len(runes) > 50 {
			description = string(runes[:50]) + "..."
		}
	}

	return &hubpb.UserGuide{
		Id:          entity.ID,
		Title:       entity.Title,
		Description: description,
		GroupKey:    entity.GroupKey,
		CreatedAt:   _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(entity.UpdatedAt),
		Key:         entity.Key,
		Mode:        hubpb.UserGuideMode(entity.Mode),
	}
}

// EntityToDetailProto converts UserGuideEntity to UserGuideDetail protobuf message
// Returns full description without truncation for detail view
func (m *UserGuideMapper) EntityToDetailProto(entity *domain.UserGuideEntity, steps []*domain.UserGuideStepEntity) *hubpb.UserGuideDetail {
	if entity == nil {
		return nil
	}

	return &hubpb.UserGuideDetail{
		Id:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description, // Full description
		GroupKey:    entity.GroupKey,
		CreatedAt:   _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(entity.UpdatedAt),
		Steps:       m.StepEntitiesToProto(steps),
		Key:         entity.Key,
		Mode:        hubpb.UserGuideMode(entity.Mode),
	}
}

// EntitiesToProto converts slice of UserGuideEntity to UserGuide protobuf messages
func (m *UserGuideMapper) EntitiesToProto(entities []*domain.UserGuideEntity) []*hubpb.UserGuide {
	if entities == nil {
		return []*hubpb.UserGuide{}
	}

	result := make([]*hubpb.UserGuide, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.EntityToProto(entity))
	}
	return result
}

// ProtoToEntity converts CreateUserGuideRequest to UserGuideEntity
func (m *UserGuideMapper) CreateRequestToEntity(req *hubpb.CreateUserGuideRequest) *domain.UserGuideEntity {
	if req == nil {
		return nil
	}

	mode := int32(10) // Default: chế độ đơn
	if req.Mode != hubpb.UserGuideMode_USER_GUIDE_MODE_UNSPECIFIED {
		mode = int32(req.Mode)
	}

	return &domain.UserGuideEntity{
		Title:       req.Title,
		Description: req.Description,
		GroupKey:    req.GroupKey,
		Key:         req.Key,
		Mode:        mode,
	}
}

// UpdateRequestToEntity updates UserGuideEntity with UpdateUserGuideRequest
func (m *UserGuideMapper) UpdateRequestToEntity(entity *domain.UserGuideEntity, req *hubpb.UpdateUserGuideRequest) {
	if entity == nil || req == nil {
		return
	}

	entity.Title = req.Title
	entity.Description = req.Description
	entity.GroupKey = req.GroupKey
	entity.Key = req.Key
	if req.Mode != hubpb.UserGuideMode_USER_GUIDE_MODE_UNSPECIFIED {
		entity.Mode = int32(req.Mode)
	}
}

// GroupMapToProto converts group map to UserGuideGroup protobuf message
func (m *UserGuideMapper) GroupMapToProto(groupMap map[string]interface{}) *hubpb.UserGuideGroup {
	if groupMap == nil {
		return nil
	}

	groupKey := ""
	count := int64(0)

	if val, ok := groupMap["group_key"].(string); ok {
		groupKey = val
	}

	if val, ok := groupMap["count"].(int64); ok {
		count = val
	}

	return &hubpb.UserGuideGroup{
		GroupKey: groupKey,
		Count:    count,
	}
}

// GroupMapsToProto converts slice of group maps to UserGuideGroup protobuf messages
func (m *UserGuideMapper) GroupMapsToProto(groupMaps []map[string]interface{}) []*hubpb.UserGuideGroup {
	if groupMaps == nil {
		return []*hubpb.UserGuideGroup{}
	}

	result := make([]*hubpb.UserGuideGroup, 0, len(groupMaps))
	for _, groupMap := range groupMaps {
		result = append(result, m.GroupMapToProto(groupMap))
	}
	return result
}
