package mapper

import (
	"fmt"
	"time"

	_utils "common/utils"
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

	out := &qh_domain.QHLayer{
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		Type:          layerType,
		Status:        status,
		DisplayOrder:  int(req.DisplayOrder),
		SourceType:    req.SourceType,
		MinZoom:       req.MinZoom,
		MaxZoom:       req.MaxZoom,
		Visible:       req.Visible,
		Avatar:        req.Avatar,
		ImageURL:      req.ImageUrl,
		ThumbnailURL:  req.ThumbnailUrl,
		EffectiveDate: effectiveDate,
		ExpiryDate:    expiryDate,
	}
	if req.FamilyId != nil && *req.FamilyId != 0 {
		fid := *req.FamilyId
		out.FamilyID = &fid
	}
	return out, nil
}

func ToProtoLayer(layer *qh_domain.QHLayer) (*tqdpb.LayerResponse, error) {
	if layer == nil {
		return nil, nil
	}

	sourceCode := config.LayerResolverInstance.GetSourceCode(layer.ID)
	tileName := config.LayerResolverInstance.GetTileName(layer.ID)
	layerUrl := config.LayerResolverInstance.GetLayerURL(layer.ID, tileName)

	styleConfig, err := config.LayerResolverInstance.GetStyleConfig(layer.ID, sourceCode, uint32(layer.Type), layer.MinZoom, layer.MaxZoom)
	if err != nil {
		styleConfig = "[]" // Default empty array if error
	}

	resp := &tqdpb.LayerResponse{
		Id:              layer.ID,
		Name:            layer.Name,
		DisplayName:     layer.DisplayName,
		Description:     layer.Description,
		Type:            uint32(layer.Type),
		Status:          uint32(layer.Status),
		Visible:         layer.Visible,
		DisplayOrder:    int32(layer.DisplayOrder),
		MinZoom:         layer.MinZoom,
		MaxZoom:         layer.MaxZoom,
		LayerUrl:        layerUrl,
		Avatar:          layer.Avatar,
		ImageUrl:        layer.ImageURL,
		ThumbnailUrl:    layer.ThumbnailURL,
		StyleConfig:     styleConfig,
		SourceType:      layer.SourceType,
		SourceCode:      sourceCode,
		CreatedAt:       layer.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       layer.UpdatedAt.Format(time.RFC3339),
		LegalStatus:     uint32(layer.LegalStatus),
		LegalStatusName: enums.LegalStatusMap[layer.LegalStatus],
		TrustValue:      float64(layer.TrustValue),
		FamilyId:        layer.FamilyID,
		PublishScopes:   layer.PublishScopes,
		LegalDoc:        layer.LegalDoc,
		DefaultVisible:  layer.DefaultVisible,
	}
	if layer.ReplacedByID != nil {
		rid := *layer.ReplacedByID
		resp.ReplacedByLayerId = &rid
	}
	if layer.FamilyID != nil {
		fid := *layer.FamilyID
		resp.FamilyId = &fid
	}
	if layer.Family != nil {
		resp.Family = ToProtoQHLayerFamily(layer.Family)
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

// ToProtoQHLayerFamily map domain họ lớp → proto (dùng chung cho LayerResponse.family).
func ToProtoQHLayerFamily(f *qh_domain.QHLayerFamily) *tqdpb.QHLayerFamilyResponse {
	if f == nil {
		return nil
	}
	var createdAt, updatedAt string
	if f.CreatedAt != nil {
		createdAt = _utils.FormatTimeToString(f.CreatedAt)
	}
	if f.UpdatedAt != nil {
		updatedAt = _utils.FormatTimeToString(f.UpdatedAt)
	}
	return &tqdpb.QHLayerFamilyResponse{
		Id:         f.ID,
		Name:       f.Name,
		SortNumber: f.SortNumber,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

// ToProtoQHLayerFamilyClient map domain họ lớp → proto kèm URL tile và style cho client/admin list.
func ToProtoQHLayerFamilyClient(f *qh_domain.QHLayerFamily, minZoom, maxZoom uint32) *tqdpb.QHLayerFamilyResponse {
	if f == nil {
		return nil
	}
	resp := ToProtoQHLayerFamily(f)
	sourceCode := config.LayerResolverInstance.GetFamilySourceCode(f.ID)
	tileName := config.LayerResolverInstance.GetFamilyTileName(f.ID)
	layerUrl := config.LayerResolverInstance.GetFamilyURL(f.ID, tileName)

	styleConfig, err := config.LayerResolverInstance.GetFamilyStyleConfig(f.ID, minZoom, maxZoom)
	if err != nil {
		styleConfig = "[]"
	}

	resp.FamilyUrl = layerUrl
	resp.LayerUrl = layerUrl
	resp.MinZoom = minZoom
	resp.MaxZoom = maxZoom
	resp.StyleConfig = styleConfig
	resp.SourceType = "vector"
	resp.SourceCode = sourceCode
	return resp
}
