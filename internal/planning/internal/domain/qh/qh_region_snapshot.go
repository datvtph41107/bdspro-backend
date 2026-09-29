package qh_domain

import "encoding/json"

// RegionSnapshot — dữ liệu tối giản để log/retry import (khi CreateBatch lỗi).
// Geometry là chuỗi GeoJSON (MultiPolygon) để có thể insert lại qua ST_GeomFromGeoJSON.
type RegionSnapshot struct {
	ID                 uint64          `json:"id,omitempty"`
	LayerID            uint64          `json:"layerId"`
	Name               string          `json:"name"`
	DisplayName        string          `json:"displayName"`
	ImportBatchID      string          `json:"importBatchId"`
	SourceFile         string          `json:"sourceFile"`
	OriginalProperties json.RawMessage `json:"originalProperties,omitempty"`
	Geometry           string          `json:"geometry,omitempty"`
}
