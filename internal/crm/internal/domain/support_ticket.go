package domain

import (
	"crm/internal/enums"
	"time"

	_models "common/models"

	"gorm.io/datatypes"
)

type SupportTicket struct {
	_models.BaseEntity
	Code            string                       `gorm:"type:varchar(32);uniqueIndex;not null" json:"code"`
	Title           string                       `gorm:"type:varchar(255);not null" json:"title"`
	Description     string                       `gorm:"type:text" json:"description"`
	IssueType       enums.SupportTicketIssueType `gorm:"not null;default:70" json:"issueType"`
	Status          enums.SupportTicketStatus    `gorm:"not null;default:10;index" json:"status"`
	Priority        enums.SupportTicketPriority  `gorm:"not null;default:20;index" json:"priority"`
	Product         enums.SupportTicketProduct   `gorm:"not null;default:10" json:"product"`
	Source          enums.SupportTicketSource       `gorm:"not null;default:10" json:"source"`
	HandlingTeam    *enums.SupportTicketHandlingTeam `gorm:"column:handling_team;index" json:"handlingTeam"`
	AssigneeID      *uint64                         `gorm:"index" json:"assigneeId"`
	RelatedUserID   *uint64                         `gorm:"index" json:"relatedUserId"`
	RelatedReportID *uint64                         `gorm:"index" json:"relatedReportId"`
	Images          datatypes.JSON                  `gorm:"type:jsonb;default:'[]'" json:"images"`
	ClosedAt        *time.Time                      `json:"closedAt"`
}

func (SupportTicket) TableName() string { return "support_tickets" }

type SupportTicketNote struct {
	_models.BaseEntity
	TicketID uint64 `gorm:"not null;index" json:"ticketId"`
	AuthorID uint64 `gorm:"not null" json:"authorId"`
	Body     string `gorm:"type:text;not null" json:"body"`
}

func (SupportTicketNote) TableName() string { return "support_ticket_notes" }

type SupportTicketEvent struct {
	_models.BaseEntity
	TicketID   uint64 `gorm:"not null;index" json:"ticketId"`
	ActorID    uint64 `gorm:"not null" json:"actorId"`
	Action     string `gorm:"type:varchar(64);not null" json:"action"`
	BeforeJSON string `gorm:"type:text" json:"beforeJson"`
	AfterJSON  string `gorm:"type:text" json:"afterJson"`
	Note       string `gorm:"type:text" json:"note"`
}

func (SupportTicketEvent) TableName() string { return "support_ticket_events" }
