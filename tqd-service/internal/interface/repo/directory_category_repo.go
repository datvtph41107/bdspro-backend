package repo

import (
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

// DirectoryCategoryRepository defines the interface for directory category repository operations
type DirectoryCategoryRepository interface {
	// Create creates a new directory category
	Create(ctx context.Context, category *domain.DirectoryCategory) error

	// GetByID gets a directory category by ID
	GetByID(ctx context.Context, id uint64) (*domain.DirectoryCategory, error)

	// Update updates an existing directory category
	Update(ctx context.Context, category *domain.DirectoryCategory) error

	// Delete deletes a directory category by ID
	Delete(ctx context.Context, id uint64) error

	// List gets a list of directory categories with pagination
	List(ctx context.Context, request *dto.ListDirectoryCategoriesRequest) ([]domain.DirectoryCategory, int64, error)
}
