package pending_data

// ToResponse chuyển entity PendingData sang DTO trả về.
func ToResponse(entity *PendingData) *PendingDataResponse {
	if entity == nil {
		return nil
	}

	return &PendingDataResponse{
		ID:          entity.ID,
		ServiceName: entity.ServiceName,
		Port:        entity.Port,
		Path:        entity.Path,
		Data:        entity.Data,
	}
}

func ToResponses(entities []PendingData) []*PendingDataResponse {
	if len(entities) == 0 {
		return []*PendingDataResponse{}
	}

	result := make([]*PendingDataResponse, 0, len(entities))
	for i := range entities {
		result = append(result, ToResponse(&entities[i]))
	}
	return result
}
