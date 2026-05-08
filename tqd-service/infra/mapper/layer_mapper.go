package mapper

import (
	"fmt"
	"time"

	tqdpb "pb/types/tqd"
	"tqd/config"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
)

func ToLayerDomainFromCreate(req *tqdpb.CreateLayerRequest) (*qh_domain.QHLayer, error) {
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

	var effectiveDate *time.Time
	if req.EffectiveDate != nil && *req.EffectiveDate != "" {
		t, err := time.Parse("2006-01-02", *req.EffectiveDate)
		if err == nil {
			effectiveDate = &t
		}
	}

	var expiryDate *time.Time
	if req.ExpiryDate != nil && *req.ExpiryDate != "" {
		t, err := time.Parse("2006-01-02", *req.ExpiryDate)
		if err == nil {
			expiryDate = &t
		}
	}

	return &qh_domain.QHLayer{
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		Type:          layerType,
		Status:        status,
		DisplayOrder:  int(req.DisplayOrder),
		SourceType:    req.SourceType,
		MinZoom:       int(req.MinZoom),
		MaxZoom:       int(req.MaxZoom),
		Visible:       req.Visible,
		Avatar:        req.Avatar,
		ImageURL:      req.ImageUrl,
		ThumbnailURL:  req.ThumbnailUrl,
		EffectiveDate: effectiveDate,
		ExpiryDate:    expiryDate,
	}, nil
}

func ToProtoLayer(layer *qh_domain.QHLayer) (*tqdpb.LayerResponse, error) {
	if layer == nil {
		return nil, nil
	}

	sourceCode := config.LayerResolverInstance.GetSourceCode(layer.ID)
	layerUrl := config.LayerResolverInstance.GetLayerURL(layer.ID, sourceCode)

	styleConfig, err := config.LayerResolverInstance.GetStyleConfig(layer.ID, sourceCode, uint32(layer.Type))
	if err != nil {
		styleConfig = "[]" // Default empty array if error
	}

	resp := &tqdpb.LayerResponse{
		Id:           layer.ID,
		Name:         layer.Name,
		DisplayName:  layer.DisplayName,
		Description:  layer.Description,
		Type:         uint32(layer.Type),
		Status:       uint32(layer.Status),
		Visible:      layer.Visible,
		DisplayOrder: int32(layer.DisplayOrder),
		MinZoom:      int32(layer.MinZoom),
		MaxZoom:      int32(layer.MaxZoom),
		LayerUrl:     layerUrl,
		Avatar:       layer.Avatar,
		ImageUrl:     layer.ImageURL,
		ThumbnailUrl: layer.ThumbnailURL,
		StyleConfig:  styleConfig,
		SourceType:   layer.SourceType,
		SourceCode:   sourceCode,
		CreatedAt:    layer.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    layer.UpdatedAt.Format(time.RFC3339),
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
