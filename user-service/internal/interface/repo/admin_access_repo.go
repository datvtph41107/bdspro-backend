package repo

import (
	"context"
	"user/internal/domain"
)

type AdminAccessRepository interface {
	// CRUD operations
	Create(c context.Context, access *domain.AdminAccessDomain) (*domain.AdminAccessDomain, error)
	Update(c context.Context, access *domain.AdminAccessDomain) (*domain.AdminAccessDomain, error)
	Delete(c context.Context, id uint64) error
	SoftDelete(c context.Context, id uint64) error
	FindByID(c context.Context, id uint64) (*domain.AdminAccessDomain, error)

	// List operations
	FindAll(c context.Context, page, size int, filters map[string]interface{}) ([]*domain.AdminAccessDomain, int64, error)
	FindByUserID(c context.Context, userID uint64) ([]*domain.AdminAccessDomain, error)
	FindByIPAddress(c context.Context, ipAddress string) ([]*domain.AdminAccessDomain, error)
	FindByDeviceID(c context.Context, deviceID string) ([]*domain.AdminAccessDomain, error)

	// Validation operations
	ValidateAccess(c context.Context, userID uint64, ipAddress, deviceID string) (*domain.AdminAccessValidationResult, error)
	CountByUserID(c context.Context, userID uint64) (int64, error)
	CountByIPAddress(c context.Context, ipAddress string) (int64, error)

	// Log operations
	CreateLog(c context.Context, log *domain.AdminAccessLogDomain) error
	FindLogsByUserID(c context.Context, userID uint64, page, size int) ([]*domain.AdminAccessLogDomain, int64, error)
	FindLogsByIPAddress(c context.Context, ipAddress string, page, size int) ([]*domain.AdminAccessLogDomain, int64, error)
}
