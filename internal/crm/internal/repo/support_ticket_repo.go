package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type SupportTicketRepo interface {
	Create(ctx context.Context, ticket *domain.SupportTicket) (*domain.SupportTicket, error)
	GetByID(ctx context.Context, id uint64) (*domain.SupportTicket, error)
	List(ctx context.Context, req *dto.SupportTicketListRequest) ([]domain.SupportTicket, int64, error)
	Update(ctx context.Context, ticket *domain.SupportTicket) error
	Summary(ctx context.Context) (*dto.SupportTicketSummaryResponse, error)
	AddNote(ctx context.Context, note *domain.SupportTicketNote) (*domain.SupportTicketNote, error)
	ListNotes(ctx context.Context, ticketID uint64) ([]domain.SupportTicketNote, error)
	AddEvent(ctx context.Context, event *domain.SupportTicketEvent) error
	ListEvents(ctx context.Context, ticketID uint64, limit int) ([]domain.SupportTicketEvent, error)
	CountByCodePrefix(ctx context.Context, prefix string) (int64, error)
}