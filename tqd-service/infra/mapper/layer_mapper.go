package mapper

import (
	_utils "common/utils"
	"encoding/json"
	"fmt"
	tqdpb "pb/types/tqd"
	"time"
	"tqd/internal/domain"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
)

func ToLayerDomainFromCreate(req *tqdpb.CreateLayerRequest) (*qh_domain.QHLayer, error) {
	// Convert type
	layerType := enums.LayerType(req.Type)
	if !layerType.IsValid() {
		return nil, fmt.Errorf("invalid type: %d", req.Type)
	}

	status := enums.LayerStatusActive
	if req.Status != nil && *req.Status > 0 {
		status = enums.LayerStatus(*req.Status)
		if !status.IsValid() {
			return nil, fmt.Errorf("invalid status: %d", *req.Status)
		}
	}

	layer := &qh_domain.QHLayer{
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Type:         layerType,
		Status:       status,
		DisplayOrder: int(req.Order),
		MinZoom:      int(req.MinZoom),
		MaxZoom:      int(req.MaxZoom),
		Visible:      req.Visible,
		LegalDoc:     req.LegalDoc,
	}

	if req.EffectiveDate != nil && *req.EffectiveDate != "" {
		parsedTime, err := time.Parse("2006-01-02", *req.EffectiveDate)
		if err == nil {
			layer.EffectiveDate = &parsedTime
		}
	}

	return layer, nil
}

// ToProtoLayer chuyển từ domain sang proto response
func ToProtoLayer(layer *qh_domain.QHLayer) (*tqdpb.LayerResponse, error) {
	status := uint32(layer.Status)

	resp := &tqdpb.LayerResponse{
		Id:               uint64(layer.ID),
		Name:             layer.Name,
		DisplayName:      layer.DisplayName,
		Description:      layer.Description,
		Type:             uint32(layer.Type),
		Status:           uint32(layer.Status),
		Visible:          layer.Visible,
		DisplayOrder:     int32(layer.DisplayOrder),
		MinZoom:          int32(layer.MinZoom),
		MaxZoom:          int32(layer.MaxZoom),
		LayerUrl:         layer.LayerURL,
		ImageUrl:         layer.ImageURL,
		ThumbnailUrl:     layer.ThumbnailURL,
		LegalDoc:         layer.LegalDoc,
		LegalDocUrl:      layer.LegalDocURL,
		DocumentNumber:   layer.DocumentNumber,
		IssuingAuthority: layer.IssuingAuthority,
		CreatedAt:        _utils.FormatTimeToString(&layer.CreatedAt),
		UpdatedAt:        _utils.FormatTimeToString(&layer.UpdatedAt),
	}

	// Marshal metadata
	if layer.Metadata != nil {
		metadataBytes, _ := json.Marshal(layer.Metadata)
		resp.Metadata = string(metadataBytes)
	}

	// Marshal style config
	if layer.StyleConfig != nil {
		styleBytes, _ := json.Marshal(layer.StyleConfig)
		resp.StyleConfig = string(styleBytes)
	}

	// Format dates
	if layer.EffectiveDate != nil {
		dateStr := layer.EffectiveDate.Format("2006-01-02")
		resp.EffectiveDate = &dateStr
	}
	if layer.ExpiryDate != nil {
		dateStr := layer.ExpiryDate.Format("2006-01-02")
		resp.ExpiryDate = &dateStr
	}

	return resp, nil
}

func (m *LayerMapper) ToProtoLayerClient(layer *domain.Layer) (*tqdpb.LayerResponse, error) {
	// Client chỉ nhận subset của admin response
	resp := &tqdpb.LayerResponse{
		Id:           uint64(layer.ID),
		Name:         layer.Name,
		DisplayName:  layer.DisplayName,
		Type:         uint32(layer.Type),
		DisplayOrder: int32(layer.DisplayOrder),
		MinZoom:      int32(layer.MinZoom),
		MaxZoom:      int32(layer.MaxZoom),
		LayerUrl:     layer.LayerURL,
		ImageUrl:     layer.ImageURL,
		ThumbnailUrl: layer.ThumbnailURL,
		CreatedAt:    _utils.FormatTimeToString(&layer.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(&layer.UpdatedAt),
	}

	// Client cũng thấy style config để render
	if layer.StyleConfig != nil {
		styleBytes, _ := json.Marshal(layer.StyleConfig)
		resp.StyleConfig = string(styleBytes)
	}

	// Client không thấy: status (vì chỉ có active), legal_doc, metadata chi tiết
	resp.Visible = true
	resp.Status = uint32(enums.LayerStatusActive)

	return resp, nil
}
