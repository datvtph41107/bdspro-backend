package dto

// CheckAppUpdateRequest mirrors hubpb.CheckAppUpdateRequest path and query params.
type CheckAppUpdateRequest struct {
	AppName       string `path:"appName"`
	Platform      string `form:"platform"`
	TargetVersion string `form:"targetVersion"`
}
