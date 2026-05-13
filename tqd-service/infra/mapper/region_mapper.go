package mapper

import (
	"encoding/json"

	_utils "common/utils"

	"google.golang.org/protobuf/types/known/structpb"

	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
)

type RegionMapper struct{}

func NewRegionMapper() *RegionMapper {
	return &RegionMapper{}
}

// ToProtoResponse - Convert domain QHRegion to proto RegionResponse (full version for admin)
func (m *RegionMapper) ToProtoResponse(region *qh_domain.QHRegion, label *qh_domain.QHLabel) *tqdpb.RegionResponse {
	if region == nil {
		return nil
	}

	// Convert original properties to Struct
	var props *structpb.Struct
	if len(region.OriginalProperties) > 0 {
		var m map[string]interface{}
		if err := json.Unmarshal(region.OriginalProperties, &m); err == nil {
			props, _ = structpb.NewStruct(m)
		}
	}

	// Get processing status (default to 10 if not set)
	processingStatus := uint32(10)
	if region.ProcessingStatus != 0 {
		processingStatus = uint32(region.ProcessingStatus)
	}

	resp := &tqdpb.RegionResponse{
		Id:          region.ID,
		LayerId:     region.LayerID,
		Name:        region.Name,
		DisplayName: region.DisplayName,
		Description: region.Description,
		AreaHa:      region.AreaHa,
		AreaSqm:     region.AreaSqm,
		// ColorRed:           uint32(region.ColorRed),
		// ColorGreen:         uint32(region.ColorGreen),
		// ColorBlue:          uint32(region.ColorBlue),
		// ColorHex:           region.CalculateColorHex(),
		OriginalProperties: props,
		SourceFile:         region.SourceFile,
		ImportBatchId:      region.ImportBatchID,
		Status:             uint32(region.Status),
		ProcessingStatus:   processingStatus,
		Version:            uint32(region.Version),
		IsLatest:           region.IsLatest,
		CreatedAt:          _utils.FormatTimeToString(&region.CreatedAt),
		UpdatedAt:          _utils.FormatTimeToString(&region.UpdatedAt),
		Geometry:           region.Geometry.Raw,
	}

	// Add style config
	// styleConfig := region.GetStyleConfig()
	// styleJSON, _ := json.Marshal(styleConfig)
	// resp.StyleConfig = string(styleJSON)

	// Add labels
	// for _, label := range labels {
	// 	resp.Labels = append(resp.Labels, NewLabelMapper().ToProtoQHLabelResponse(&label))
	// }
	resp.Label = NewLabelMapper().ToProtoQHLabelResponse(label)

	return resp
}

// ToProtoClientResponse - Convert domain to proto (simplified version for client)
func (m *RegionMapper) ToProtoClientResponse(region *qh_domain.QHRegion) *tqdpb.RegionResponse {
	if region == nil {
		return nil
	}

	resp := &tqdpb.RegionResponse{
		Id:          region.ID,
		LayerId:     region.LayerID,
		Name:        region.Name,
		DisplayName: region.DisplayName,
		Description: region.Description,
		AreaHa:      region.AreaHa,
		AreaSqm:     region.AreaSqm,
		ColorHex:    region.CalculateColorHex(),
		Status:      0,
		Version:     0,
		CreatedAt:   _utils.FormatTimeToString(&region.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(&region.UpdatedAt),
		Geometry:    region.Geometry.Raw,
	}

	// Add style config
	// styleConfig := region.GetStyleConfig()
	// styleJSON, _ := json.Marshal(styleConfig)
	// resp.StyleConfig = string(styleJSON)

	return resp
}
