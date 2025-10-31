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

// GetLoginHistoryByAuthID lấy danh sách lịch sử đăng nhập theo authID với phân trang
func (r *SessionPostgres) GetLoginHistoryByAuthID(ctx context.Context, authID uint64, page, size int32) ([]*domain.UserSessionEntity, int64, error) {
	var sessions []*domain.UserSessionEntity
	var total int64

	// Đếm tổng số bản ghi
	if err := _db.DB.Model(&domain.UserSessionEntity{}).Where("auth_id = ?", authID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Default page và size
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	offset := (page - 1) * size

	// Lấy dữ liệu với phân trang, sắp xếp theo thời gian tạo giảm dần (mới nhất lên đầu)
	err := _db.DB.Where("auth_id = ?", authID).
		Order("created_date DESC").
		Limit(int(size)).
		Offset(int(offset)).
		Find(&sessions).Error

	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

// GetLoginHistoryByUserID lấy danh sách lịch sử đăng nhập theo userID với phân trang
// JOIN với auth_method để lọc theo user_id
func (r *SessionPostgres) GetLoginHistoryByUserID(ctx context.Context, userID uint64, page, size int32) ([]*domain.UserSessionEntity, int64, error) {
	var sessions []*domain.UserSessionEntity
	var total int64

	// Đếm tổng số bản ghi với JOIN
	if err := _db.DB.Model(&domain.UserSessionEntity{}).
		Joins("INNER JOIN auth_method ON user_session.auth_id = auth_method.id").
		Where("auth_method.user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Default page và size
	if page < 0 {
		page = 0
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	offset := page * size

	// Lấy dữ liệu với phân trang, JOIN với auth_method, sắp xếp theo thời gian tạo giảm dần (mới nhất lên đầu)
	err := _db.DB.
		Joins("INNER JOIN auth_method ON user_session.auth_id = auth_method.id").
		Where("auth_method.user_id = ?", userID).
		Order("user_session.created_date DESC").
		Limit(int(size)).
		Offset(int(offset)).
		Find(&sessions).Error

	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}
