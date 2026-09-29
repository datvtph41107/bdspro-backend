package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"
)

// IntentRepo - Repository interface cho Intent
type IntentRepo interface {
	// Create tạo intent mới
	Create(ctx context.Context, intent *domain.Intent) error

	// GetByID lấy intent theo ID
	GetByID(ctx context.Context, id uint64) (*domain.Intent, error)

	// GetList lấy danh sách intents với phân trang và filter
	GetList(ctx context.Context, req *dto.IntentSearchRequest) ([]domain.Intent, int64, error)

	// GetBySessionID lấy danh sách intents theo session
	GetBySessionID(ctx context.Context, sessionID string) ([]domain.Intent, error)

	// GetStatistics lấy thống kê intents
	GetStatistics(ctx context.Context, req *dto.IntentStatisticsRequest) ([]domain.IntentStatistics, error)

	// GetTopIntents lấy top intents được sử dụng nhiều nhất
	GetTopIntents(ctx context.Context, category *enums.EIntentCategory, limit int) ([]domain.IntentStatistics, error)

	// Update cập nhật intent
	Update(ctx context.Context, id uint64, intent *domain.Intent) error

	// Delete xóa intent (soft delete)
	Delete(ctx context.Context, id uint64) error

	// GetRecentIntentsByProfile lấy intents gần đây của user
	GetRecentIntentsByProfile(ctx context.Context, profileID uint64, limit int) ([]domain.Intent, error)
}
