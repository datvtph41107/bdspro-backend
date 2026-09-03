package dto

type CreateLabelRequest struct {
	LayerID         uint64  `json:"layerId"`
	Name            string  `json:"name"`
	DisplayName     string  `json:"displayName"`
	Description     string  `json:"description"`
	Color           string  `json:"color"`
	FillOpacity     float64 `json:"fillOpacity"`
	StrokeColor     string  `json:"strokeColor"`
	StrokeWidth     int     `json:"strokeWidth"`
	StrokeDashArray string  `json:"strokeDashArray"`
	DisplayOrder    int     `json:"displayOrder"`
	IsVisible       bool    `json:"isVisible"`
	MinZoom         int     `json:"minZoom"`
	MaxZoom         int     `json:"maxZoom"`
	Status          int     `json:"status,omitempty"`
}

type UpdateLabelRequest struct {
	Name            *string  `json:"name,omitempty"`
	DisplayName     *string  `json:"displayName,omitempty"`
	Description     *string  `json:"description,omitempty"`
	Color           *string  `json:"color,omitempty"`
	FillOpacity     *float64 `json:"fillOpacity,omitempty"`
	StrokeColor     *string  `json:"strokeColor,omitempty"`
	StrokeWidth     *int     `json:"strokeWidth,omitempty"`
	StrokeDashArray *string  `json:"strokeDashArray,omitempty"`
	DisplayOrder    *int     `json:"displayOrder,omitempty"`
	IsVisible       *bool    `json:"isVisible,omitempty"`
	MinZoom         *int     `json:"minZoom,omitempty"`
	MaxZoom         *int     `json:"maxZoom,omitempty"`
	Status          *int     `json:"status,omitempty"`
}

type LabelResponse struct {
	ID              uint64  `json:"id"`
	LayerID         uint64  `json:"layerId"`
	Name            string  `json:"name"`
	DisplayName     string  `json:"displayName"`
	Description     string  `json:"description"`
	Color           string  `json:"color"`
	FillOpacity     float64 `json:"fillOpacity"`
	StrokeColor     string  `json:"strokeColor"`
	StrokeWidth     int     `json:"strokeWidth"`
	StrokeDashArray string  `json:"strokeDashArray"`
	DisplayOrder    int     `json:"displayOrder"`
	IsVisible       bool    `json:"isVisible"`
	MinZoom         int     `json:"minZoom"`
	MaxZoom         int     `json:"maxZoom"`
	Status          int     `json:"status"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type ListLabelResponse struct {
	Data  []LabelResponse `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}
