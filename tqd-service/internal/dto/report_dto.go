package dto

type CreateReportRequestDTO struct {
	ReportType      uint32      `json:"reportType"`
	TargetID        *string     `json:"targetId,omitempty"`
	TargetData      interface{} `json:"targetData,omitempty"`
	Profile         uint32      `json:"profile"`
	Format          uint32      `json:"format"`
	IncludeMapImage bool        `json:"includeMapImage"`
}
