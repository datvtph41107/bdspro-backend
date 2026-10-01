package dto

// SystemConfigItem is the API representation of a system config row.
type SystemConfigItem struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	GroupKey  string `json:"groupKey"`
	GroupName string `json:"groupName"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// UpsertConfigItem is one entry in a bulk-upsert request.
type UpsertConfigItem struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	GroupKey string `json:"groupKey"`
}

// BulkUpsertSystemConfigRequest is the JSON body for POST /v2/hub/system-config/bulk-upsert.
type BulkUpsertSystemConfigRequest struct {
	Configs []UpsertConfigItem `json:"configs"`
}

// GetSystemConfigByGroupQuery binds the {groupKey} path param.
type GetSystemConfigByGroupQuery struct {
	GroupKey string `path:"groupKey"`
}

// GetSystemConfigByKeyQuery binds the {key} path param.
type GetSystemConfigByKeyQuery struct {
	Key string `path:"key"`
}

// ResetDefaultByGroupQuery binds the {groupKey} path param for reset.
type ResetDefaultByGroupQuery struct {
	GroupKey string `path:"groupKey"`
}

// GetUserSettingsByGroupQuery binds the {groupKey} path param for settings.
type GetUserSettingsByGroupQuery struct {
	GroupKey string `path:"groupKey"`
}
