package repo

import (
	"context"
	"map/internal/domain"
)

// DictCommonRepository defines the interface for dictionary common repository operations
type DictCommonRepository interface {
	// GetByCode gets a dictionary item by code
	GetByCode(ctx context.Context, code string) (*domain.DictCommon, error)

	// GetByCategory gets all dictionary items by category
	GetByCategory(ctx context.Context, category string) ([]domain.DictCommon, error)

	// GetByCategoryAndParent gets dictionary items by category and parent
	GetByCategoryAndParent(ctx context.Context, category string, parentID *uint) ([]domain.DictCommon, error)

	// GetHierarchical gets dictionary items in hierarchical structure
	GetHierarchical(ctx context.Context, category string) ([]domain.DictCommonWithChildren, error)

	// Create creates a new dictionary item
	Create(ctx context.Context, item *domain.DictCommon) (*domain.DictCommon, error)

	// Update updates an existing dictionary item
	Update(ctx context.Context, id uint, item *domain.DictCommon) (*domain.DictCommon, error)

	// Delete deletes a dictionary item by ID
	Delete(ctx context.Context, id uint) error

	// GetByID gets a dictionary item by ID
	GetByID(ctx context.Context, id uint) (*domain.DictCommon, error)

	// List gets dictionary items with filtering and pagination
	List(ctx context.Context, filter *domain.DictCommonFilter, page, size int) ([]domain.DictCommon, int64, error)

	// GetSummary gets summary statistics of dictionary data
	GetSummary(ctx context.Context) (*domain.DictCommonSummary, error)

	// BulkCreate creates multiple dictionary items
	BulkCreate(ctx context.Context, items []domain.DictCommon) error

	// BulkUpdate updates multiple dictionary items
	BulkUpdate(ctx context.Context, items []domain.DictCommon) error

	// BulkDelete deletes multiple dictionary items
	BulkDelete(ctx context.Context, ids []uint) error

	// GetActiveByCategory gets active dictionary items by category
	GetActiveByCategory(ctx context.Context, category string) ([]domain.DictCommon, error)

	// Search searches dictionary items by name or description
	Search(ctx context.Context, query string, category *string) ([]domain.DictCommon, error)
}
