package repo

import (
	"context"
	"tqd/internal/domain"
)

// DirectorySupplierRepository defines the interface for directory supplier repository
type DirectorySupplierRepository interface {
	// CRUD Operations - sử dụng domain trực tiếp
	Create(ctx context.Context, supplier *domain.DirectorySupplier) error
	GetByID(ctx context.Context, id uint64) (*domain.DirectorySupplier, error)
	Update(ctx context.Context, supplier *domain.DirectorySupplier) error
	Delete(ctx context.Context, id uint64) error

	// List Operations - sử dụng tham số trực tiếp
	List(ctx context.Context, page, size int, search string, categories []string, isActive *bool, minRating, maxRating *float64) ([]domain.DirectorySupplier, int64, error)
}
