package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type SupportTicketPostgres struct {
	db *gorm.DB
}

func NewSupportTicketPostgres(db *gorm.DB) repo.SupportTicketRepo {
	return &SupportTicketPostgres{db: db}
}

func (r *SupportTicketPostgres) Create(ctx context.Context, ticket *domain.SupportTicket) (*domain.SupportTicket, error) {
	if err := r.db.WithContext(ctx).Create(ticket).Error; err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *SupportTicketPostgres) GetByID(ctx context.Context, id uint64) (*domain.SupportTicket, error) {
	var ticket domain.SupportTicket
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&ticket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticket, err
}

func (r *SupportTicketPostgres) List(ctx context.Context, req *dto.SupportTicketListRequest) ([]domain.SupportTicket, int64, error) {
	var tickets []domain.SupportTicket
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.SupportTicket{}).Where("deleted_at IS NULL")
	if req.Status != nil {
		q = q.Where("status = ?", *req.Status)
	}
	if req.Priority != nil {
		q = q.Where("priority = ?", *req.Priority)
	}
	if req.IssueType != nil {
		q = q.Where("issue_type = ?", *req.IssueType)
	}
	if req.AssigneeID != nil {
		q = q.Where("assignee_id = ?", *req.AssigneeID)
	}
	if req.UnassignedOnly {
		q = q.Where("assignee_id IS NULL")
	}
	if req.Product != nil {
		q = q.Where("product = ?", *req.Product)
	}
	if req.Source != nil {
		q = q.Where("source = ?", *req.Source)
	}
	if req.HandlingTeam != nil {
		q = q.Where("handling_team = ?", *req.HandlingTeam)
	}
	if req.Q != "" {
		like := "%" + req.Q + "%"
		q = q.Where("code ILIKE ? OR title ILIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	size := int(req.GetSize())
	offset := req.GetOffset()
	err := q.Order("updated_at DESC").Offset(offset).Limit(size).Find(&tickets).Error
	return tickets, total, err
}

func (r *SupportTicketPostgres) Update(ctx context.Context, ticket *domain.SupportTicket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *SupportTicketPostgres) Summary(ctx context.Context) (*dto.SupportTicketSummaryResponse, error) {
	type row struct {
		Status enums.SupportTicketStatus
		Count  int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&domain.SupportTicket{}).
		Select("status, count(*) as count").
		Where("deleted_at IS NULL").
		Group("status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := &dto.SupportTicketSummaryResponse{}
	for _, r0 := range rows {
		out.Total += r0.Count
		switch r0.Status {
		case enums.SupportTicketStatusNew:
			out.StatusNew = r0.Count
		case enums.SupportTicketStatusInProgress:
			out.StatusInProgress = r0.Count
		case enums.SupportTicketStatusWaitingUser:
			out.StatusWaitingUser = r0.Count
		case enums.SupportTicketStatusResolved:
			out.StatusResolved = r0.Count
		case enums.SupportTicketStatusClosed:
			out.StatusClosed = r0.Count
		case enums.SupportTicketStatusTransferred:
			out.StatusTransferred = r0.Count
		}
	}
	var unassigned int64
	if err := r.db.WithContext(ctx).Model(&domain.SupportTicket{}).
		Where("deleted_at IS NULL AND assignee_id IS NULL AND status <> ?", enums.SupportTicketStatusClosed).
		Count(&unassigned).Error; err != nil {
		return nil, err
	}
	out.Unassigned = unassigned
	var critical int64
	if err := r.db.WithContext(ctx).Model(&domain.SupportTicket{}).
		Where("deleted_at IS NULL AND priority = ? AND status <> ?", enums.SupportTicketPriorityCritical, enums.SupportTicketStatusClosed).
		Count(&critical).Error; err != nil {
		return nil, err
	}
	out.PriorityCritical = critical
	return out, nil
}

func (r *SupportTicketPostgres) AddNote(ctx context.Context, note *domain.SupportTicketNote) (*domain.SupportTicketNote, error) {
	if err := r.db.WithContext(ctx).Create(note).Error; err != nil {
		return nil, err
	}
	return note, nil
}

func (r *SupportTicketPostgres) ListNotes(ctx context.Context, ticketID uint64) ([]domain.SupportTicketNote, error) {
	var notes []domain.SupportTicketNote
	err := r.db.WithContext(ctx).Where("ticket_id = ? AND deleted_at IS NULL", ticketID).
		Order("created_at ASC").Find(&notes).Error
	return notes, err
}

func (r *SupportTicketPostgres) AddEvent(ctx context.Context, event *domain.SupportTicketEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *SupportTicketPostgres) ListEvents(ctx context.Context, ticketID uint64, limit int) ([]domain.SupportTicketEvent, error) {
	var events []domain.SupportTicketEvent
	q := r.db.WithContext(ctx).Where("ticket_id = ? AND deleted_at IS NULL", ticketID).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&events).Error
	return events, err
}

func (r *SupportTicketPostgres) CountByCodePrefix(ctx context.Context, prefix string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.SupportTicket{}).
		Where("code LIKE ?", fmt.Sprintf("%s%%", prefix)).Count(&count).Error
	return count, err
}
