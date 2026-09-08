package dto

import qh_domain "tqd/internal/domain/qh"

type PolygonQueryRequest struct {
	GeoJSON         string  `json:"geo_json"`
	MaxAreaKm2      float64 `json:"max_area_km2"` // gợi ý, không bắt buộc
	Page            int     `json:"page"`
	PageSize        int     `json:"page_size"`
	AllowAutoShrink bool    `json:"allow_auto_shrink"` // client cho phép server tự scale
}

type PolygonQueryResponse struct {
	Parcels           []qh_domain.Parcel `json:"parcels"`
	Total             int64              `json:"total"`
	Page              int                `json:"page"`
	PageSize          int                `json:"page_size"`
	AdjustedGeoJSON   string             `json:"adjusted_geo_json"` // polygon thực tế đã dùng để query
	AdjustmentReason  string             `json:"adjustment_reason"` // "NONE" | "SHRUNK_AREA" | "SHRUNK_RADIUS"
	OriginalAreaKm2   float64            `json:"original_area_km2"`
	AdjustedAreaKm2   float64            `json:"adjusted_area_km2"`
	MaxAllowedAreaKm2 float64            `json:"max_allowed_area_km2"`
	Warning           string             `json:"warning,omitempty"`
}

type Point struct {
	Lat float64
	Lng float64
}
