package dto

import "time"

type Layer struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Type        string   `json:"type"`
	Status      string   `json:"status"`
	Order       int      `json:"order"`
	Visible     bool     `json:"visible"`
	Children    []*Layer `json:"children,omitempty"`
}

type CreateLayerRequest struct {
	Name          string                 `json:"name" binding:"required,max=255"`
	DisplayName   string                 `json:"display_name" binding:"required"`
	Type          string                 `json:"type" binding:"required,oneof=vector raster wms wmts"`
	Order         int                    `json:"order"`
	MinZoom       int                    `json:"min_zoom"`
	MaxZoom       int                    `json:"max_zoom"`
	Visible       bool                   `json:"visible"`
	Metadata      map[string]interface{} `json:"metadata"`
	LegalDoc      string                 `json:"legal_doc"`
	EffectiveDate *time.Time             `json:"effective_date"`
}

// UpdateLayerRequest cập nhật
type UpdateLayerRequest struct {
	Name          *string                `json:"name"`
	DisplayName   *string                `json:"display_name"`
	Type          *string                `json:"type" binding:"omitempty,oneof=vector raster wms wmts"`
	Status        *string                `json:"status" binding:"omitempty,oneof=active inactive deleted"`
	ParentID      *uint                  `json:"parent_id"`
	Order         *int                   `json:"order"`
	MinZoom       *int                   `json:"min_zoom"`
	MaxZoom       *int                   `json:"max_zoom"`
	Visible       *bool                  `json:"visible"`
	Metadata      map[string]interface{} `json:"metadata"`
	LegalDoc      *string                `json:"legal_doc"`
	EffectiveDate *time.Time             `json:"effective_date"`
}

// LayerResponse trả về
type LayerResponse struct {
	ID            uint                   `json:"id"`
	Name          string                 `json:"name"`
	DisplayName   string                 `json:"display_name"`
	Type          string                 `json:"type"`
	Status        string                 `json:"status"`
	ParentID      *uint                  `json:"parent_id,omitempty"`
	Order         int                    `json:"order"`
	MinZoom       int                    `json:"min_zoom"`
	MaxZoom       int                    `json:"max_zoom"`
	Visible       bool                   `json:"visible"`
	Metadata      map[string]interface{} `json:"metadata"`
	LegalDoc      string                 `json:"legal_doc"`
	EffectiveDate *time.Time             `json:"effective_date"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}
