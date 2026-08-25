package postgres

import (
	"context"
	"strings"
	"time"

	_db "common/db"
	_dto "common/domain/dto"
	"user/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SessionPostgres implementation của Session repository
// @bind: user/internal/interface/repo.SessionRepository
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

func (r *SessionPostgres) GetBySessionIDAndAuthID(ctx context.Context, sessionID uint64, authID uint64) (*domain.UserSessionEntity, error) {
	var session domain.UserSessionEntity
	err := _db.DB.Where("session_id = ? and auth_id = ?", sessionID, authID).First(&session).Error
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

// GetSessionsByAuthID lấy toàn bộ session theo authID
func (r *SessionPostgres) GetSessionsByAuthIDWithoutSessionId(ctx context.Context, authID uint64, sessionID uint64) ([]*domain.UserSessionEntity, error) {
	var sessions []*domain.UserSessionEntity
	if err := _db.DB.WithContext(ctx).
		Where("auth_id = ? and logout_at is null and session_id != ?", authID, sessionID).
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetByDeviceAndAuthID lấy session theo deviceID và authID
func (r *SessionPostgres) GetByDeviceAndAuthID(ctx context.Context, deviceID string, authID uint64) (*domain.UserSessionEntity, error) {
	var session domain.UserSessionEntity
	err := _db.DB.WithContext(ctx).
		Where("device_id = ? AND auth_id = ?", deviceID, authID).
		First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// CreateSession tạo mới session, nếu tồn tại deviceId + authId thì cập nhật, không thì tạo mới
func (r *SessionPostgres) CreateSession(ctx context.Context, session *domain.UserSessionEntity) error {
	var existing domain.UserSessionEntity
	now := time.Now()
	err := _db.DB.
		Where("device_id = ? AND auth_id = ?", session.DeviceID, session.AuthID).
		First(&existing).Error

	if err == nil {
		// Đã tồn tại: cập nhật thông tin session
		existing.Version = session.Version
		existing.Platform = session.Platform
		existing.OS = session.OS
		existing.DeviceName = session.DeviceName
		existing.IPRequest = session.IPRequest
		existing.UserAgent = session.UserAgent
		existing.SessionKey = session.SessionKey
		existing.ClientID = session.ClientID
		existing.TotalRequest = session.TotalRequest
		existing.LastLogin = &now
		existing.LogoutAt = nil
		existing.FinishedDate = session.FinishedDate
		existing.Activate = session.Activate
		existing.LastRequest = session.LastRequest

		session.SessionID = existing.SessionID

		return _db.DB.Save(&existing).Error
	}

	// Nếu không tìm thấy (gorm.ErrRecordNotFound), tạo mới
	if err == gorm.ErrRecordNotFound {
		return _db.DB.Create(session).Error
	}
	// Nếu lỗi khác, trả về lỗi
	return err
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

// GetSessionsByProfile lấy danh sách session theo userID với phân trang
// JOIN với auth_method để lọc theo user_id
func (r *SessionPostgres) GetSessionsByProfile(ctx context.Context, userID uint64, pagable *_dto.Pagable) ([]*domain.UserSessionEntity, int64, error) {
	var sessions []*domain.UserSessionEntity
	var total int64

	if pagable == nil {
		pagable = &_dto.Pagable{}
	}

	baseQuery := _db.DB.WithContext(ctx).
		Model(&domain.UserSessionEntity{}).
		Joins("INNER JOIN auth_method ON user_session.auth_id = auth_method.id").
		Where("auth_method.user_id = ?", userID)

	// Đếm tổng số bản ghi với JOIN
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy dữ liệu với phân trang, JOIN với auth_method, sắp xếp theo thời gian tạo giảm dần (mới nhất lên đầu)
	err := baseQuery.
		Order("user_session.created_date DESC").
		Limit(pagable.GetLimit()).
		Offset(pagable.GetOffset()).
		Find(&sessions).Error

	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

// RecordUserActivity ghi nhận login/refresh theo khung giờ hiện tại (+7).
// Unique (profile_id, slot_at) → cùng giờ chỉ lưu 1 lần.
func (r *SessionPostgres) RecordUserActivity(ctx context.Context, profileID uint64, source string) error {
	if profileID == 0 {
		return nil
	}
	source = strings.TrimSpace(source)
	if source == "" {
		source = "login"
	}

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)
	}
	now := time.Now().In(loc)
	slotAt := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)

	logEntry := &domain.UserActivityLog{
		ProfileID: profileID,
		SlotAt:    slotAt.UTC(),
		Source:    source,
		CreatedAt: time.Now().UTC(),
	}

	return _db.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "profile_id"}, {Name: "slot_at"}},
			DoNothing: true,
		}).
		Create(logEntry).Error
}
