package mapper

import (
	_utils "common/utils"

	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
)

type LabelMapper struct{}

func NewLabelMapper() *LabelMapper {
	return &LabelMapper{}
}

// ToProtoResponse - Convert domain QHLabel to proto QHLabelResponse
func (m *LabelMapper) ToProtoResponse(l *qh_domain.QHLabel) *tqdpb.QHLabelResponse {
	if l == nil {
		return nil
	}

	resp := &tqdpb.QHLabelResponse{
		Id:           l.ID,
		LayerId:      l.LayerID,
		Name:         l.Name,
		DisplayName:  l.DisplayName,
		Description:  l.Description,
		Color:        l.Color,
		FillOpacity:  l.FillOpacity,
		StrokeColor:  l.StrokeColor,
		StrokeWidth:  int32(l.StrokeWidth),
		DisplayOrder: int32(l.DisplayOrder),
		IsVisible:    l.IsVisible,
		MinZoom:      int32(l.MinZoom),
		MaxZoom:      int32(l.MaxZoom),
		Status:       uint32(l.Status),
		CreatedAt:    _utils.FormatTimeToString(&l.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(&l.UpdatedAt),
	}
	if l.StandardAt != nil {
		resp.StandardAt = _utils.FormatTimeToString(l.StandardAt)
	}
	return resp
}

// ToProtoQHLabelResponse - Alias for ToProtoResponse (for embedding)
func (m *LabelMapper) ToProtoQHLabelResponse(l *qh_domain.QHLabel) *tqdpb.QHLabelResponse {
	return m.ToProtoResponse(l)
}

// ToProtoListResponse - Convert list of labels to proto list response
func (m *LabelMapper) ToProtoListResponse(labels []qh_domain.QHLabel, total int64, page, pageSize int32) *tqdpb.ListLabelsResponse {
	data := make([]*tqdpb.QHLabelResponse, len(labels))
	for i, l := range labels {
		data[i] = m.ToProtoResponse(&l)
	}

	return &tqdpb.ListLabelsResponse{
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
