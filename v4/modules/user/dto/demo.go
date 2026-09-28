package dto

// DemoModeRequest mirrors userpb.DemoModeRequest for POST /v2/user/demo-mode.
type DemoModeRequest struct {
	Enable bool `json:"enable"`
}
