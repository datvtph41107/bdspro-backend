package dto

import (
	"errors"
	"regexp"
	"strings"
)

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// GetParcelPlanningRequest - Request canonical cho parcel planning.
// View quyết định projection trả về, mode quyết định cách engine xử lý.
type GetParcelPlanningRequest struct {
	ParcelID          uint64 `json:"parcel_id"`
	View              string `json:"view"`
	Mode              string `json:"mode"`
	Preset            string `json:"preset"`
	AsOfDate          string `json:"as_of_date"`
	IncludeOverview   bool   `json:"include_overview"`
	IncludeAssessment bool   `json:"include_assessment"`
	Timestamp         int64  `json:"timestamp"`
}

func (r *GetParcelPlanningRequest) Validate() error {
	if r.ParcelID == 0 {
		return errors.New("parcel_id is required")
	}

	r.View = strings.TrimSpace(strings.ToLower(r.View))
	if r.View == "" {
		r.View = "overview"
	}

	switch r.View {
	case "overview":
		if !r.IncludeOverview && !r.IncludeAssessment {
			r.IncludeOverview = true
		}
	case "detail":
		if !r.IncludeOverview && !r.IncludeAssessment {
			r.IncludeOverview = true
			r.IncludeAssessment = true
		}
	case "assessment":
		if !r.IncludeOverview && !r.IncludeAssessment {
			r.IncludeAssessment = true
		}
	case "compare":
		if !r.IncludeOverview && !r.IncludeAssessment {
			r.IncludeAssessment = true
		}
	default:
		return errors.New("invalid view: must be overview, detail, assessment, or compare")
	}

	if r.Mode == "" {
		switch r.View {
		case "overview":
			r.Mode = "overview"
		case "compare":
			r.Mode = "compare"
		default:
			r.Mode = "detail"
		}
	}

	if r.Mode != "popup" && r.Mode != "overview" && r.Mode != "detail" && r.Mode != "compare" {
		return errors.New("invalid mode: must be popup, overview, detail, or compare")
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
