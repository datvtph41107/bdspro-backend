package dto

import (
	_dto "common/domain/dto"
	"crm/internal/enums"
	"time"
)

type SupportTicketCreateRequest struct {
	Title           string   `json:"title" binding:"required"`
	Description     string   `json:"description"`
	IssueType       uint32   `json:"issueType"`
	Priority        uint32   `json:"priority"`
	RelatedUserID   *uint64  `json:"relatedUserId"`
	RelatedReportID *uint64  `json:"relatedReportId"`
	Product         uint32   `json:"product"`
	Source          uint32   `json:"source"`
	Images          []string `json:"images"`
}

type SupportTicketImagesRequest struct {
	Images []string `json:"images" binding:"required"`
}

type SupportTicketListRequest struct {
	_dto.Pagable
	Status         *uint32 `json:"status"`
	Priority       *uint32 `json:"priority"`
	IssueType      *uint32 `json:"issueType"`
	AssigneeID     *uint64 `json:"assigneeId"`
	Product        *uint32 `json:"product"`
	Source         *uint32 `json:"source"`
	HandlingTeam   *uint32 `json:"handlingTeam"`
	UnassignedOnly bool    `json:"unassignedOnly"`
	Q              string  `json:"q"`
}

type SupportTicketAssignRequest struct {
	AssigneeID uint64 `json:"assigneeId" binding:"required"`
	Note       string `json:"note"`
}

type SupportTicketTransferRequest struct {
	HandlingTeam uint32 `json:"handlingTeam" binding:"required"`
	Reason       string `json:"reason" binding:"required"`
}

type SupportTicketPriorityRequest struct {
	Priority uint32 `json:"priority" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
}

type SupportTicketStatusRequest struct {
	Status uint32 `json:"status" binding:"required"`
	Note   string `json:"note"`
}

type SupportTicketNoteRequest struct {
	Body string `json:"body" binding:"required"`
}

type SupportTicketCloseRequest struct {
	Note string `json:"note"`
}

type SupportTicketUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type SupportTicketItemResponse struct {
	ID              uint64     `json:"id"`
	Code            string     `json:"code"`
	Title           string     `json:"title"`
	IssueType       uint32     `json:"issueType"`
	Status          uint32     `json:"status"`
	Priority        uint32     `json:"priority"`
	AssigneeID      *uint64    `json:"assigneeId"`
	RelatedUserID   *uint64    `json:"relatedUserId"`
	RelatedReportID *uint64    `json:"relatedReportId"`
	Product         uint32     `json:"product"`
	Source          uint32     `json:"source"`
	HandlingTeam    *uint32    `json:"handlingTeam"`
	Images          []string   `json:"images,omitempty"`
	CreatedBy       uint64     `json:"createdBy"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	CreatedAt       *time.Time `json:"createdAt"`
}

type SupportTicketNoteResponse struct {
	ID        uint64     `json:"id"`
	AuthorID  uint64     `json:"authorId"`
	Body      string     `json:"body"`
	CreatedAt *time.Time `json:"createdAt"`
}

type SupportTicketEventResponse struct {
	ID         uint64     `json:"id"`
	ActorID    uint64     `json:"actorId"`
	Action     string     `json:"action"`
	Note       string     `json:"note"`
	CreatedAt  *time.Time `json:"createdAt"`
	ActorName  string     `json:"actorName,omitempty"`
	BeforeJSON string     `json:"beforeJson,omitempty"`
	AfterJSON  string     `json:"afterJson,omitempty"`
}

type SupportTicketDetailResponse struct {
	ID              uint64                       `json:"id"`
	Code            string                       `json:"code"`
	Title           string                       `json:"title"`
	Description     string                       `json:"description"`
	IssueType       uint32                       `json:"issueType"`
	Status          uint32                       `json:"status"`
	Priority        uint32                       `json:"priority"`
	AssigneeID      *uint64                      `json:"assigneeId"`
	RelatedUserID   *uint64                      `json:"relatedUserId"`
	RelatedReportID *uint64                      `json:"relatedReportId"`
	Product         uint32                       `json:"product"`
	Source          uint32                       `json:"source"`
	HandlingTeam    *uint32                      `json:"handlingTeam"`
	Images          []string                     `json:"images"`
	CreatedBy       uint64                       `json:"createdBy"`
	CreatedAt       *time.Time                   `json:"createdAt"`
	UpdatedAt       *time.Time                   `json:"updatedAt"`
	ClosedAt        *time.Time                   `json:"closedAt"`
	Notes           []SupportTicketNoteResponse  `json:"notes"`
	Events          []SupportTicketEventResponse `json:"events"`
}

type SupportTicketListResponse struct {
	Data  []SupportTicketItemResponse `json:"data"`
	Total int64                       `json:"total"`
}

type SupportTicketSummaryResponse struct {
	Total              int64 `json:"total"`
	StatusNew          int64 `json:"statusNew"`
	StatusInProgress   int64 `json:"statusInProgress"`
	StatusWaitingUser  int64 `json:"statusWaitingUser"`
	StatusResolved     int64 `json:"statusResolved"`
	StatusClosed       int64 `json:"statusClosed"`
	StatusTransferred  int64 `json:"statusTransferred"`
	Unassigned         int64 `json:"unassigned"`
	PriorityCritical   int64 `json:"priorityCritical"`
}

// Defaults helpers
func DefaultIssueType(v uint32) enums.SupportTicketIssueType {
	if v == 0 {
		return enums.SupportTicketIssueOther
	}
	return enums.SupportTicketIssueType(v)
}

func DefaultPriority(v uint32) enums.SupportTicketPriority {
	if v == 0 {
		return enums.SupportTicketPriorityNormal
	}
	return enums.SupportTicketPriority(v)
}

func DefaultProduct(v uint32) enums.SupportTicketProduct {
	if v == 0 {
		return enums.SupportTicketProductQHPro
	}
	return enums.SupportTicketProduct(v)
}

func DefaultSource(v uint32) enums.SupportTicketSource {
	if v == 0 {
		return enums.SupportTicketSourceAdmin
	}
	return enums.SupportTicketSource(v)
}
