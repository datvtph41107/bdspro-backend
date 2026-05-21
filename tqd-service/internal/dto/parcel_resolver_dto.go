package dto

import (
	"errors"
	"regexp"
)

// GetParcelAssessmentRequest - Request cho API đánh giá chuyên sâu (VIP)
type GetParcelAssessmentRequest struct {
	ParcelID uint64 `json:"parcel_id"`
	Mode     string `json:"mode"`
	Preset   string `json:"preset"`
	AsOfDate string `json:"as_of_date"`
}

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func (r *GetParcelAssessmentRequest) Validate() error {
	if r.ParcelID == 0 {
		return errors.New("parcel_id is required")
	}
	if r.Mode == "" {
		r.Mode = "detail"
	}
	if r.Mode != "popup" && r.Mode != "detail" && r.Mode != "compare" {
		return errors.New("invalid mode: must be popup, detail, or compare")
	}
	if r.Preset == "" {
		r.Preset = "all_layers"
	}
	if r.Preset != "active_only" && r.Preset != "all_layers" {
		return errors.New("invalid preset: must be active_only or all_layers")
	}
	if r.AsOfDate != "" && !dateRegex.MatchString(r.AsOfDate) {
		return errors.New("invalid as_of_date: must be YYYY-MM-DD format")
	}
	return nil
}

// ParcelAssessmentResponse - Response DTO cho API đánh giá chuyên sâu (VIP)
type ParcelAssessmentResponse struct {
	ParcelID      uint64          `json:"parcel_id"`
	ParcelAreaSqm float64         `json:"parcel_area_sqm"`
	CurrentPlan   *PlanV2Info     `json:"current_plan,omitempty"`
	OtherPlans    []*PlanV2Info   `json:"other_plans,omitempty"`
	HasConflict   bool            `json:"has_conflict"`
	CanBuild      bool            `json:"can_build"`
	Reason        string          `json:"reason"`
	Summary       string          `json:"summary"`
	Metadata      *ResponseMetaV2 `json:"metadata,omitempty"`
}

type PlanV2Info struct {
	LayerID       uint64  `json:"layer_id"`
	LayerName     string  `json:"layer_name"`
	DisplayName   string  `json:"display_name"`
	LandUseCode   string  `json:"land_use_code"`
	LandUseName   string  `json:"land_use_name"`
	GroupCode     string  `json:"group_code"`
	ImpactPercent float64 `json:"impact_percent"`
	ImpactAreaSqm float64 `json:"impact_area_sqm"`
	CanBuild      bool    `json:"can_build"`
	Note          string  `json:"note,omitempty"`
}

type ResponseMetaV2 struct {
	ResponseTimeMs int64  `json:"response_time_ms"`
	DataVersion    string `json:"data_version"`
}
