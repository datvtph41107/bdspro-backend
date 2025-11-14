package repo

import (
	"context"
	"time"

	"notification/internal/domain"
)

type HistoryAuthSearch struct {
	UserID         uint64
	OrganizationID *uint64
	ActionType     string
	Channel        string
	DeviceID       string
	Success        *bool
	FromDate       *time.Time
	ToDate         *time.Time
	Offset         int
	Limit          int
}

type HistoryAuthRepo interface {
	CreateHistoryAuth(ctx context.Context, entity *domain.HistoryAuthEntity) error
	SearchHistoryAuth(ctx context.Context, params HistoryAuthSearch) ([]domain.HistoryAuthEntity, int64, error)
}
