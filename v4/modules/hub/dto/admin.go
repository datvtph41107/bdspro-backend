package dto

// ---- FAQ admin DTOs ----

// CreateFAQRequest is the JSON body for admin POST /v2/hub/admin/faqs.
type CreateFAQRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	GroupKey string `json:"groupKey,omitempty"`
}

// UpdateFAQRequest is the JSON body for admin PUT /v2/hub/admin/faqs/{id}.
type UpdateFAQRequest struct {
	ID       uint64 `path:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	GroupKey string `json:"groupKey,omitempty"`
}

// DeleteFAQPath binds the {id} for admin DELETE.
type DeleteFAQPath struct {
	ID uint64 `path:"id"`
}

// ---- UserGuide admin DTOs ----

// CreateUserGuideRequest is the JSON body for admin POST /v2/hub/user-guides.
type CreateUserGuideRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	GroupKey    string `json:"groupKey"`
	Key         string `json:"key,omitempty"`
	Mode        int32  `json:"mode,omitempty"`
}

// UpdateUserGuideRequest is the JSON body for admin PUT /v2/hub/user-guides/{id}.
type UpdateUserGuideRequest struct {
	ID          uint64 `path:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	GroupKey    string `json:"groupKey"`
	Key         string `json:"key,omitempty"`
	Mode        int32  `json:"mode,omitempty"`
}

// DeleteUserGuidePath binds the {id} for admin DELETE.
type DeleteUserGuidePath struct {
	ID uint64 `path:"id"`
}

// AdminListQuery is a generic pagination + search query for admin list endpoints.
type AdminListQuery struct {
	Page int32  `form:"page"`
	Size int32  `form:"size"`
	Text string `form:"text"`
}

// ---- Version admin DTOs ----

// CreateVersionRequest is the JSON body for admin POST /v2/hub/admin/versions.
type CreateVersionAdminRequest struct {
	AppName      string `json:"appName"`
	Platform     string `json:"platform"`
	VersionName  string `json:"versionName"`
	BuildNumber  uint64 `json:"buildNumber"`
	ForceUpdate  bool   `json:"forceUpdate"`
	Active       bool   `json:"active"`
	DownloadURL  string `json:"downloadUrl,omitempty"`
	ReleaseNotes string `json:"releaseNotes,omitempty"`
	BundleName   string `json:"bundleName,omitempty"`
	BundleSize   uint64 `json:"bundleSize,omitempty"`
}

// UpdateVersionRequest is the JSON body for admin PUT /v2/hub/admin/versions/{id}.
type UpdateVersionRequest struct {
	ID           uint64 `path:"id"`
	AppName      string `json:"appName"`
	Platform     string `json:"platform"`
	VersionName  string `json:"versionName"`
	BuildNumber  uint64 `json:"buildNumber"`
	ForceUpdate  bool   `json:"forceUpdate"`
	Active       bool   `json:"active"`
	DownloadURL  string `json:"downloadUrl,omitempty"`
	ReleaseNotes string `json:"releaseNotes,omitempty"`
	BundleName   string `json:"bundleName,omitempty"`
	BundleSize   uint64 `json:"bundleSize,omitempty"`
}

// DeleteVersionPath binds the {id} for admin DELETE.
type DeleteVersionPath struct {
	ID uint64 `path:"id"`
}

// VersionListQuery binds list query params for admin version list.
type VersionListQuery struct {
	Page    int32  `form:"page"`
	Size    int32  `form:"size"`
	AppName string `form:"appName"`
	Platform string `form:"platform"`
}
