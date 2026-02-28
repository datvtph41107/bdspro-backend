// file/internal/infrastructure/persistence/gorm/session_repository.go
package gorm

import (
	"context"
	"file/internal/core/domain"
	"file/internal/ports"
	"time"

	"gorm.io/gorm"
)

type uploadSessionRepository struct {
	db *gorm.DB
}

func NewUploadSessionRepository(db *gorm.DB) ports.UploadSessionRepository {
	return &uploadSessionRepository{db: db}
}

func (r *uploadSessionRepository) CreateSession(ctx context.Context, session *domain.UploadSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *uploadSessionRepository) GetSession(ctx context.Context, sessionID string) (*domain.UploadSession, error) {
	var session domain.UploadSession
	err := r.db.WithContext(ctx).
		Where("id = ? AND expires_at > ?", sessionID, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *uploadSessionRepository) UpdateSession(ctx context.Context, session *domain.UploadSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *uploadSessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", sessionID).
		Delete(&domain.UploadSession{}).Error
}

func (r *uploadSessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&domain.UploadSession{}).Error
}

func (r *uploadSessionRepository) UpdateChunkStatus(ctx context.Context, sessionID string, chunkNumber int, status string) error {
	// Update chunk metadata in session
	// This would typically involve updating a JSON field or a separate chunks table
	return nil
}
