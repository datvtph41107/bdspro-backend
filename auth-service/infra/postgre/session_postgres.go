package postgre

import (
	"context"
	"time"

	"auth/internal/domain"
	_db "common/db"
)

// SessionPostgres implementation của Session repository
// @bind: auth/internal/interface/repo.SessionRepository
type SessionPostgres struct {
	// TODO: Add database connection
}

// NewSessionRepository tạo mới SessionRepository
func NewSessionRepository() *SessionPostgres {
	return &SessionPostgres{}
}

// GetBySessionID lấy session theo sessionID
func (r *SessionPostgres) GetBySessionID(ctx context.Context, sessionID uint64) (*domain.UserSessionEntity, error) {
	var session domain.UserSessionEntity
	err := _db.DB.Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByAuthID lấy session theo authID
func (r *SessionPostgres) GetByAuthID(ctx context.Context, authID uint64) (*domain.UserSessionEntity, error) {
	var session domain.UserSessionEntity
	err := _db.DB.Where("auth_id = ?", authID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// CreateSession tạo mới session
func (r *SessionPostgres) CreateSession(ctx context.Context, session *domain.UserSessionEntity) error {
	return _db.DB.Create(session).Error
}

// UpdateSession cập nhật session
func (r *SessionPostgres) UpdateSession(ctx context.Context, session *domain.UserSessionEntity) error {
	return _db.DB.Save(session).Error
}

// DeleteSession xóa session
func (r *SessionPostgres) DeleteSession(ctx context.Context, sessionID uint64) error {
	return _db.DB.Where("session_id = ?", sessionID).Delete(&domain.UserSessionEntity{}).Error
}

// GetActiveSessionsByAuthID lấy danh sách session đang hoạt động theo authID
func (r *SessionPostgres) GetActiveSessionsByAuthID(ctx context.Context, authID uint64) ([]*domain.UserSessionEntity, error) {
	var sessions []*domain.UserSessionEntity
	err := _db.DB.Where("auth_id = ? AND finished_date IS NULL", authID).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetExpiredSessions lấy danh sách session đã hết hạn
func (r *SessionPostgres) GetExpiredSessions(ctx context.Context, before time.Time) ([]*domain.UserSessionEntity, error) {
	var sessions []*domain.UserSessionEntity
	err := _db.DB.Where("created_at < ?", before).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// CleanupExpiredSessions xóa các session đã hết hạn
func (r *SessionPostgres) CleanupExpiredSessions(ctx context.Context, before time.Time) error {
	return _db.DB.Where("created_at < ?", before).Delete(&domain.UserSessionEntity{}).Error
}
