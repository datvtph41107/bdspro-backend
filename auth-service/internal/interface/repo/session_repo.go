package repo

import (
	"auth/internal/domain"
	"context"
	"time"
)

// SessionRepository interface cho việc tương tác với dữ liệu session
type SessionRepository interface {
	// Session methods
	GetBySessionID(ctx context.Context, sessionID uint64) (*domain.UserSessionEntity, error)
	GetByAuthID(ctx context.Context, authID uint64) (*domain.UserSessionEntity, error)
	CreateSession(ctx context.Context, session *domain.UserSessionEntity) error
	UpdateSession(ctx context.Context, session *domain.UserSessionEntity) error
	DeleteSession(ctx context.Context, sessionID uint64) error

	// Session utilities
	GetActiveSessionsByAuthID(ctx context.Context, authID uint64) ([]*domain.UserSessionEntity, error)
	GetExpiredSessions(ctx context.Context, before time.Time) ([]*domain.UserSessionEntity, error)
	CleanupExpiredSessions(ctx context.Context, before time.Time) error

	// Login history
	GetLoginHistoryByAuthID(ctx context.Context, authID uint64, page, size int32) ([]*domain.UserSessionEntity, int64, error)
	GetLoginHistoryByUserID(ctx context.Context, userID uint64, page, size int32) ([]*domain.UserSessionEntity, int64, error)
}
