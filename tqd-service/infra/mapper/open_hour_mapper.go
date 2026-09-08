package mapper

import (
	"fmt"
	"time"

	_entity "common/domain/entity"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
)

type OpenHourMapper struct{}

func NewOpenHourMapper() *OpenHourMapper {
	return &OpenHourMapper{}
}

// ==================== DOMAIN -> RESPONSE ====================

func (m *OpenHourMapper) ToResponse(oh *domain.OpenHour) *dto.OpenHourResponse {
	if oh == nil {
		return nil
	}
	resp := &dto.OpenHourResponse{}
	resp.FromDomain(oh)
	return resp
}

// ToResponseList converts slice of domain OpenHour to slice of response DTO
func (m *OpenHourMapper) ToResponseList(openHours []domain.OpenHour) []dto.OpenHourResponse {
	if openHours == nil {
		return nil
	}
	result := make([]dto.OpenHourResponse, len(openHours))
	for i, oh := range openHours {
		result[i].FromDomain(&oh)
	}
	return result
}

// ==================== RESPONSE -> DOMAIN ====================

// ToDomainFromResponse converts response DTO to domain
func (m *OpenHourMapper) ToDomainFromResponse(resp *dto.OpenHourResponse) *domain.OpenHour {
	if resp == nil {
		return nil
	}

	poiID := resp.POIID
	return &domain.OpenHour{
		BaseEntity: _entity.BaseEntity{
			ID:        resp.ID,
			CreatedAt: resp.CreatedAt,
			UpdatedAt: resp.UpdatedAt,
			AuditBase: _entity.AuditBase{
				CreatedBy: &resp.CreatedBy,
				UpdatedBy: &resp.UpdatedBy,
			},
		},
		POIID:       poiID,
		DayOfWeek:   resp.DayOfWeek,
		OpenTime:    resp.OpenTime,
		CloseTime:   resp.CloseTime,
		Note:        resp.Note,
		IsOpen:      resp.IsOpen,
		IsRecurring: resp.IsRecurring,
		StartDate:   resp.StartDate,
		EndDate:     resp.EndDate,
		Status:      enums.OpenHourStatus(resp.Status),
		Type:        enums.OpenHourType(resp.Type),
	}
}

// ==================== REQUEST -> DOMAIN ====================

// ToDomainFromCreate converts create request to domain
func (m *OpenHourMapper) ToDomainFromCreate(req *dto.CreateOpenHourRequest, userID uint64) (*domain.OpenHour, error) {
	if req == nil {
		return nil, nil
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		startDate = _utils.ParseStringToTime(*req.StartDate)
	}
	if req.EndDate != nil && *req.EndDate != "" {
		endDate = _utils.ParseStringToTime(*req.EndDate)
	}

	now := time.Now()
	return &domain.OpenHour{
		BaseEntity: _entity.BaseEntity{
			CreatedAt: &now,
			UpdatedAt: &now,
			AuditBase: _entity.AuditBase{
				CreatedBy: &userID,
				UpdatedBy: &userID,
			},
		},
		POIID:       &req.POIID,
		DayOfWeek:   req.DayOfWeek,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		Note:        req.Note,
		IsOpen:      req.IsOpen,
		IsRecurring: req.IsRecurring,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      enums.OpenHourStatusActive,
		Type:        enums.OpenHourType(req.Type),
	}, nil
}

// ToDomainFromUpdate converts update request to domain
func (m *OpenHourMapper) ToDomainFromUpdate(req *dto.UpdateOpenHourRequest, existing *domain.OpenHour, userID uint64) (*domain.OpenHour, error) {
	if req == nil {
		return nil, nil
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		startDate = _utils.ParseStringToTime(*req.StartDate)
	}
	if req.EndDate != nil && *req.EndDate != "" {
		endDate = _utils.ParseStringToTime(*req.EndDate)
	}

	now := time.Now()
	updated := &domain.OpenHour{
		BaseEntity: _entity.BaseEntity{
			ID:        existing.ID,
			CreatedAt: existing.CreatedAt,
			UpdatedAt: &now,
			AuditBase: _entity.AuditBase{
				CreatedBy: existing.CreatedBy,
				UpdatedBy: &userID,
			},
		},
		POIID: existing.POIID,
	}

	// Update only provided fields
	if req.DayOfWeek != nil {
		updated.DayOfWeek = *req.DayOfWeek
	} else {
		updated.DayOfWeek = existing.DayOfWeek
	}

	if req.OpenTime != nil {
		updated.OpenTime = *req.OpenTime
	} else {
		updated.OpenTime = existing.OpenTime
	}

	if req.CloseTime != nil {
		updated.CloseTime = *req.CloseTime
	} else {
		updated.CloseTime = existing.CloseTime
	}

	if req.Note != nil {
		updated.Note = *req.Note
	} else {
		updated.Note = existing.Note
	}

	if req.IsOpen != nil {
		updated.IsOpen = *req.IsOpen
	} else {
		updated.IsOpen = existing.IsOpen
	}

	if req.IsRecurring != nil {
		updated.IsRecurring = *req.IsRecurring
	} else {
		updated.IsRecurring = existing.IsRecurring
	}

	if req.Status != nil {
		updated.Status = enums.OpenHourStatus(*req.Status)
	} else {
		updated.Status = existing.Status
	}

	if req.Type != nil {
		updated.Type = enums.OpenHourType(*req.Type)
	} else {
		updated.Type = existing.Type
	}

	if startDate != nil {
		updated.StartDate = startDate
	} else {
		updated.StartDate = existing.StartDate
	}

	if endDate != nil {
		updated.EndDate = endDate
	} else {
		updated.EndDate = existing.EndDate
	}

	return updated, nil
}

// ToDomainFromBulkItem converts bulk item to domain
func (m *OpenHourMapper) ToDomainFromBulkItem(item *dto.BulkSaveOpenHourItem, userID uint64) (*domain.OpenHour, error) {
	if item == nil {
		return nil, nil
	}

	var startDate, endDate *time.Time
	if item.StartDate != nil && *item.StartDate != "" {
		startDate = _utils.ParseStringToTime(*item.StartDate)
	}
	if item.EndDate != nil && *item.EndDate != "" {
		endDate = _utils.ParseStringToTime(*item.EndDate)
	}

	now := time.Now()
	oh := &domain.OpenHour{
		BaseEntity: _entity.BaseEntity{
			CreatedAt: &now,
			UpdatedAt: &now,
			AuditBase: _entity.AuditBase{
				CreatedBy: &userID,
				UpdatedBy: &userID,
			},
		},
		POIID:       &item.POIID,
		DayOfWeek:   item.DayOfWeek,
		OpenTime:    item.OpenTime,
		CloseTime:   item.CloseTime,
		Note:        item.Note,
		IsOpen:      item.IsOpen,
		IsRecurring: item.IsRecurring,
		StartDate:   startDate,
		EndDate:     endDate,
		Type:        enums.OpenHourType(item.Type),
		Status:      enums.OpenHourStatusActive,
	}

	if item.ID != nil {
		oh.ID = *item.ID
	}

	return oh, nil
}

// ==================== DOMAIN -> PROTO ====================

// ToProto converts domain OpenHour to proto
func (m *OpenHourMapper) ToProto(oh *domain.OpenHour) *tqdpb.OpenHour {
	if oh == nil {
		return nil
	}

	poiID := uint64(0)
	if oh.POIID != nil {
		poiID = *oh.POIID
	}

	return &tqdpb.OpenHour{
		Id:          oh.ID,
		PoiId:       &poiID,
		DayOfWeek:   dayOfWeekToString(oh.DayOfWeek),
		StartTime:   oh.OpenTime,
		EndTime:     oh.CloseTime,
		IsOpen:      oh.IsOpen,
		IsRecurring: oh.IsRecurring,
		StartDate:   _utils.FormatTimeToString(oh.StartDate),
		EndDate:     _utils.FormatTimeToString(oh.EndDate),
		Notes:       oh.Note,
		Status:      int32(oh.Status),
		Type:        int32(oh.Type),
		CreatedAt:   _utils.FormatTimeToString(oh.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(oh.UpdatedAt),
	}
}

// ToProtoFromResponse converts response DTO to proto
func (m *OpenHourMapper) ToProtoFromResponse(resp *dto.OpenHourResponse) *tqdpb.OpenHour {
	if resp == nil {
		return nil
	}

	return &tqdpb.OpenHour{
		Id:          resp.ID,
		PoiId:       resp.POIID,
		DayOfWeek:   resp.DayName,
		StartTime:   resp.OpenTime,
		EndTime:     resp.CloseTime,
		IsOpen:      resp.IsOpen,
		IsRecurring: resp.IsRecurring,
		StartDate:   _utils.FormatTimeToString(resp.StartDate),
		EndDate:     _utils.FormatTimeToString(resp.EndDate),
		Notes:       resp.Note,
		Status:      int32(resp.Status),
		Type:        int32(resp.Type),
		CreatedAt:   _utils.FormatTimeToString(resp.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(resp.UpdatedAt),
	}
}

// ToProtoList converts slice of domain OpenHour to slice of proto
func (m *OpenHourMapper) ToProtoList(openHours []domain.OpenHour) []*tqdpb.OpenHour {
	if openHours == nil {
		return nil
	}
	result := make([]*tqdpb.OpenHour, len(openHours))
	for i, oh := range openHours {
		result[i] = m.ToProto(&oh)
	}
	return result
}

// ToProtoListFromResponses converts slice of response DTO to slice of proto
func (m *OpenHourMapper) ToProtoListFromResponses(responses []dto.OpenHourResponse) []*tqdpb.OpenHour {
	if responses == nil {
		return nil
	}
	result := make([]*tqdpb.OpenHour, len(responses))
	for i, resp := range responses {
		result[i] = m.ToProtoFromResponse(&resp)
	}
	return result
}

// ToProtoTimeSlot converts domain OpenHour to proto time slot
func (m *OpenHourMapper) ToProtoTimeSlot(oh *domain.OpenHour) *tqdpb.OpenHourTime {
	if oh == nil {
		return nil
	}

	poiID := uint64(0)
	if oh.POIID != nil {
		poiID = *oh.POIID
	}

	return &tqdpb.OpenHourTime{
		Id:        oh.ID,
		PoiId:     poiID,
		StartTime: oh.OpenTime,
		EndTime:   oh.CloseTime,
		IsOpen:    oh.IsOpen,
		Notes:     oh.Note,
	}
}

// ==================== PROTO -> REQUEST ====================

// FromProto converts proto to create request
func (m *OpenHourMapper) FromProto(proto *tqdpb.OpenHour) (*dto.CreateOpenHourRequest, error) {
	if proto == nil {
		return nil, nil
	}
	if proto.PoiId == nil {
		return nil, fmt.Errorf("poi id is required")
	}

	dayOfWeek, err := stringToDayOfWeek(proto.DayOfWeek)
	if err != nil {
		return nil, err
	}

	return &dto.CreateOpenHourRequest{
		POIID:       *proto.PoiId,
		DayOfWeek:   dayOfWeek,
		OpenTime:    proto.StartTime,
		CloseTime:   proto.EndTime,
		Note:        proto.Notes,
		IsOpen:      proto.IsOpen,
		IsRecurring: proto.IsRecurring,
		StartDate:   &proto.StartDate,
		EndDate:     &proto.EndDate,
		Type:        int(proto.Type),
	}, nil
}

// ==================== HELPER FUNCTIONS ====================

func dayOfWeekToString(day int) string {
	switch day {
	case 2:
		return "Monday"
	case 3:
		return "Tuesday"
	case 4:
		return "Wednesday"
	case 5:
		return "Thursday"
	case 6:
		return "Friday"
	case 7:
		return "Saturday"
	case 8:
		return "Sunday"
	default:
		return ""
	}
}

func stringToDayOfWeek(day string) (int, error) {
	switch day {
	case "Monday":
		return 2, nil
	case "Tuesday":
		return 3, nil
	case "Wednesday":
		return 4, nil
	case "Thursday":
		return 5, nil
	case "Friday":
		return 6, nil
	case "Saturday":
		return 7, nil
	case "Sunday":
		return 8, nil
	default:
		return 0, fmt.Errorf("invalid day of week: %s", day)
	}
}
