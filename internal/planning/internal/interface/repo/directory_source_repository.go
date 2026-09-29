package repo

import (
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

// DirectorySourceRepository defines the interface for directory source repository
type DirectorySourceRepository interface {
	// CRUD Operations - sử dụng domain trực tiếp
	Create(ctx context.Context, directorySource *domain.DirectorySource) error
	GetByID(ctx context.Context, id uint64) (*domain.DirectorySource, error)
	Update(ctx context.Context, directorySource *domain.DirectorySource) error
	Delete(ctx context.Context, id uint64) error

	// List Operations - sử dụng ListRequest DTO
	List(ctx context.Context, req *dto.ListDirectorySourcesRequestDTO) ([]dto.DirectorySourceDTO, int64, error)
}
