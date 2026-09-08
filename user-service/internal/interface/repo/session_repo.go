package repo

import (
	"context"

	_dto "common/domain/dto"
	"user/internal/domain/auth"
)

// SessionRepository contains only session operations consumed by serving usecases.
type SessionRepository interface {
	GetBySessionID(ctx context.Context, sessionID uint64) (*auth.UserSessionEntity, error)
	GetBySessionIDAndAuthID(ctx context.Context, sessionID uint64, authID uint64) (*auth.UserSessionEntity, error)
	GetSessionsByAuthIDWithoutSessionId(ctx context.Context, authID uint64, sessionID uint64) ([]*auth.UserSessionEntity, error)
	CreateSession(ctx context.Context, session *auth.UserSessionEntity) error
	UpdateSession(ctx context.Context, session *auth.UserSessionEntity) error

	GetLoginHistoryByAuthID(ctx context.Context, authID uint64, page, size int32) ([]*auth.UserSessionEntity, int64, error)
	GetSessionsByProfile(ctx context.Context, userID uint64, pagable *_dto.Pagable) ([]*auth.UserSessionEntity, int64, error)
}
