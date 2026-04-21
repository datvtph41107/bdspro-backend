package mapper

import (
	"fmt"
	"time"

	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
	"tqd/internal/enums"
)

func ToLayerDomainFromCreate(req *tqdpb.CreateLayerRequest) (*domain.Layer, error) {
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

	layer := &domain.Layer{
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
func ToProtoLayer(layer *domain.Layer) (*tqdpb.LayerResponse, error) {
	status := uint32(layer.Status)

	resp := &tqdpb.LayerResponse{
		Id:          uint32(layer.ID),
		Name:        layer.Name,
		DisplayName: layer.DisplayName,
		Type:        uint32(layer.Type),
		Status:      &status,
		Order:       int32(layer.DisplayOrder),
		MinZoom:     int32(layer.MinZoom),
		MaxZoom:     int32(layer.MaxZoom),
		Visible:     layer.Visible,
		LegalDoc:    layer.LegalDoc,
		CreatedAt:   layer.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   layer.UpdatedAt.Format(time.RFC3339),
	}

	if layer.EffectiveDate != nil {
		effectiveDate := layer.EffectiveDate.Format("2006-01-02")
		resp.EffectiveDate = &effectiveDate
	}

	return resp, nil
}
