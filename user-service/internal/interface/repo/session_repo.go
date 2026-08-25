package repo

import (
	_dto "common/domain/dto"
	"context"
	"time"
	"user/internal/domain"
)

// SessionRepository interface cho việc tương tác với dữ liệu session
type SessionRepository interface {
	// Session methods
	GetBySessionID(ctx context.Context, sessionID uint64) (*domain.UserSessionEntity, error)
	GetBySessionIDAndAuthID(ctx context.Context, sessionID uint64, authID uint64) (*domain.UserSessionEntity, error)
	GetByAuthID(ctx context.Context, authID uint64) (*domain.UserSessionEntity, error)
	GetByDeviceAndAuthID(ctx context.Context, deviceID string, authID uint64) (*domain.UserSessionEntity, error)
	GetSessionsByAuthIDWithoutSessionId(ctx context.Context, authID uint64, sessionID uint64) ([]*domain.UserSessionEntity, error)
	CreateSession(ctx context.Context, session *domain.UserSessionEntity) error
	UpdateSession(ctx context.Context, session *domain.UserSessionEntity) error
	DeleteSession(ctx context.Context, sessionID uint64) error

	// Session utilities
	GetActiveSessionsByAuthID(ctx context.Context, authID uint64) ([]*domain.UserSessionEntity, error)
	GetExpiredSessions(ctx context.Context, before time.Time) ([]*domain.UserSessionEntity, error)
	CleanupExpiredSessions(ctx context.Context, before time.Time) error

	// Login history
	GetLoginHistoryByAuthID(ctx context.Context, authID uint64, page, size int32) ([]*domain.UserSessionEntity, int64, error)
	GetSessionsByProfile(ctx context.Context, userID uint64, pagable *_dto.Pagable) ([]*domain.UserSessionEntity, int64, error)

	// User activity (login/refresh) — 1 bản ghi / user / khung giờ
	RecordUserActivity(ctx context.Context, profileID uint64, source string) error
}
