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

type AdminCreateReportRequest struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	ProblemReport   uint32   `json:"problemReport"`
	ReportType      uint32   `json:"reportType"`
	Severity        uint32   `json:"severity"`
	TargetID        *uint64  `json:"targetId"`
	SupportTicketID *uint64  `json:"supportTicketId"`
	Images          []string `json:"images"`
}

type AdminImagesRequest struct {
	Images []string `json:"images" binding:"required"`
}

type AdminAssignRequest struct {
	AssigneeID uint64 `json:"assigneeId"`
}

type AdminSeverityRequest struct {
	Severity uint32 `json:"severity"`
	Reason   string `json:"reason"`
}

type AdminQaStatusRequest struct {
	QaStatus uint32 `json:"qaStatus"`
	Note     string `json:"note"`
}

type AdminCloseRequest struct {
	Reject bool   `json:"reject"`
	Note   string `json:"note"`
}
