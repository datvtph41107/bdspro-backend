package dto

import "tqd/internal/enums"

type CreateReportRequestDTO struct {
	ReportType      uint32      `json:"reportType"`
	TargetID        *string     `json:"targetId,omitempty"`
	TargetData      interface{} `json:"targetData,omitempty"`
	Profile         uint32      `json:"profile"`
	Format          uint32      `json:"format"`
	IncludeMapImage bool        `json:"includeMapImage"`
}

type ReportDTO struct {
	ID            uint64
	ProblemReport enums.ProblemReport
	ReportType    enums.ReportType
	Description   string
	TargetId      *uint64
}
