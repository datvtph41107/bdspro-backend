package dto

// CheckAppUpdateRequest mirrors hubpb.CheckAppUpdateRequest query params.
type CheckAppUpdateRequest struct {
	Platform      string `form:"platform"`
	TargetVersion string `form:"targetVersion"`
}
