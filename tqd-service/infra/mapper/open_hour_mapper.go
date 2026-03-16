package mapper

import (
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
)

// OpenHourMapper handles mapping between domain and proto
type OpenHourMapper struct{}

// NewOpenHourMapper creates a new OpenHourMapper
func NewOpenHourMapper() *OpenHourMapper {
	return &OpenHourMapper{}
}

// ToProto converts domain OpenHour to proto OpenHour
func (m *OpenHourMapper) ToProto(openHour *domain.OpenHour) *tqdpb.OpenHour {
	if openHour == nil {
		return nil
	}

	return &tqdpb.OpenHour{
		Id:        openHour.ID,
		PoiId:     openHour.POIID,
		DayOfWeek: openHour.DayOfWeek,
		StartTime: openHour.OpenTime,
		EndTime:   openHour.CloseTime,
		IsOpen:    openHour.IsOpen,
		Notes:     openHour.Note,
		Status:    int32(openHour.Status),
		Type:      int32(openHour.Type),
		CreatedAt: _utils.FormatTimeToString(openHour.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(openHour.UpdatedAt),
	}
}

// ToProtoList converts slice of domain OpenHour to slice of proto OpenHour
func (m *OpenHourMapper) ToProtoList(openHours []*domain.OpenHour) []*tqdpb.OpenHour {
	result := make([]*tqdpb.OpenHour, len(openHours))
	for i, openHour := range openHours {
		result[i] = m.ToProto(openHour)
	}
	return result
}
