package repo

import (
	"context"
	"user/internal/domain/access"
)

type AdminAccessRepository interface {
	// CRUD operations
	Create(c context.Context, access *access.AdminAccessDomain) (*access.AdminAccessDomain, error)
	Update(c context.Context, access *access.AdminAccessDomain) (*access.AdminAccessDomain, error)
	Delete(c context.Context, id uint64) error
	SoftDelete(c context.Context, id uint64) error
	FindByID(c context.Context, id uint64) (*access.AdminAccessDomain, error)

	// List operations
	FindAll(c context.Context, page, size int, filters map[string]interface{}) ([]*access.AdminAccessDomain, int64, error)
	FindByUserID(c context.Context, userID uint64) ([]*access.AdminAccessDomain, error)
	FindByIPAddress(c context.Context, ipAddress string) ([]*access.AdminAccessDomain, error)
	FindByDeviceID(c context.Context, deviceID string) ([]*access.AdminAccessDomain, error)

	// Validation operations
	ValidateAccess(c context.Context, userID uint64, ipAddress, deviceID string) (*access.AdminAccessValidationResult, error)
	CountByUserID(c context.Context, userID uint64) (int64, error)
	CountByIPAddress(c context.Context, ipAddress string) (int64, error)

	// Log operations
	CreateLog(c context.Context, log *access.AdminAccessLogDomain) error
	FindLogsByUserID(c context.Context, userID uint64, page, size int) ([]*access.AdminAccessLogDomain, int64, error)
	FindLogsByIPAddress(c context.Context, ipAddress string, page, size int) ([]*access.AdminAccessLogDomain, int64, error)
}
