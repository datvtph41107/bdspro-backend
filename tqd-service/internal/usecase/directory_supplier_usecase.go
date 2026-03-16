package usecase

import (
	"context"
	"fmt"
	"tqd/internal/domain"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"
)

// DirectorySupplierUsecase implements DirectorySupplierUsecase interface
type DirectorySupplierUsecase struct {
	directorySupplierRepo repo.DirectorySupplierRepository
	userProvider          provider.UserProvider
	permissionProvider    provider.PermissionProvider
	authUsecase           *AuthorizationUsecase
}

// NewDirectorySupplierUsecase creates a new DirectorySupplierUsecase
func NewDirectorySupplierUsecase(
	directorySupplierRepo repo.DirectorySupplierRepository,
	userProvider provider.UserProvider,
	permissionProvider provider.PermissionProvider,
) *DirectorySupplierUsecase {
	return &DirectorySupplierUsecase{
		directorySupplierRepo: directorySupplierRepo,
		userProvider:          userProvider,
		permissionProvider:    permissionProvider,
		authUsecase: &AuthorizationUsecase{
			permissionProvider: permissionProvider,
			userProvider:       userProvider,
		},
	}
}

// Create creates a new directory supplier with permission check
func (u *DirectorySupplierUsecase) Create(ctx context.Context, supplier *domain.DirectorySupplier) (*domain.DirectorySupplier, error) {
	// 1. Check permission
	if err := u.authUsecase.RequirePermission(ctx, enums.PermissionCreateSupplier); err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	// 2. Validate domain
	if err := supplier.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 3. Set createdBy from context
	userID := u.authUsecase.GetCurrentUserID(ctx)
	supplier.CreatedBy = userID

	// 4. Create in repository
	if err := u.directorySupplierRepo.Create(ctx, supplier); err != nil {
		return nil, fmt.Errorf("failed to create directory supplier: %w", err)
	}

	return supplier, nil
}

// GetByID gets directory supplier by ID
func (u *DirectorySupplierUsecase) GetByID(ctx context.Context, id uint64) (*domain.DirectorySupplier, error) {
	supplier, err := u.directorySupplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory supplier: %w", err)
	}
	return supplier, nil
}

// Update updates directory supplier with permission check
func (u *DirectorySupplierUsecase) Update(ctx context.Context, supplier *domain.DirectorySupplier) (*domain.DirectorySupplier, error) {
	// 1. Check permission
	if err := u.authUsecase.RequirePermission(ctx, enums.PermissionUpdateSupplier); err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	// 2. Check if supplier exists
	existing, err := u.directorySupplierRepo.GetByID(ctx, supplier.ID)
	if err != nil {
		return nil, fmt.Errorf("directory supplier not found: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("directory supplier not found")
	}

	// 3. Validate domain
	if err := supplier.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 4. Set updatedBy from context
	userID := u.authUsecase.GetCurrentUserID(ctx)
	supplier.UpdatedBy = userID

	// 5. Update in repository
	if err := u.directorySupplierRepo.Update(ctx, supplier); err != nil {
		return nil, fmt.Errorf("failed to update directory supplier: %w", err)
	}

	return supplier, nil
}

// Delete deletes directory supplier by ID with permission check
func (u *DirectorySupplierUsecase) Delete(ctx context.Context, id uint64) error {
	// 1. Check permission
	if err := u.authUsecase.RequirePermission(ctx, enums.PermissionDeleteSupplier); err != nil {
		return fmt.Errorf("unauthorized: %w", err)
	}

	// 2. Check if supplier exists
	_, err := u.directorySupplierRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("directory supplier not found: %w", err)
	}

	// 3. Delete from repository
	if err := u.directorySupplierRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete directory supplier: %w", err)
	}

	return nil
}

// List gets list of directory suppliers with filters
func (u *DirectorySupplierUsecase) List(ctx context.Context, page, size int, search string, categories []string, isActive *bool, minRating, maxRating *float64) ([]domain.DirectorySupplier, int64, error) {
	suppliers, total, err := u.directorySupplierRepo.List(ctx, page, size, search, categories, isActive, minRating, maxRating)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list directory suppliers: %w", err)
	}
	return suppliers, total, nil
}
