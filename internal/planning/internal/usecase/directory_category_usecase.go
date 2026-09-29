package usecase

import (
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"
)

type DirectoryCategoryUsecase struct {
	*AuthorizationUsecase
	directoryCategoryRepo repo.DirectoryCategoryRepository
}

func NewDirectoryCategoryUsecase(directoryCategoryRepo repo.DirectoryCategoryRepository,
	userProvider provider.UserProvider,
	permissionProvider provider.PermissionProvider) *DirectoryCategoryUsecase {
	return &DirectoryCategoryUsecase{
		AuthorizationUsecase:  NewBaseUsecase(userProvider, permissionProvider),
		directoryCategoryRepo: directoryCategoryRepo,
	}
}

func (u *DirectoryCategoryUsecase) CreateDirectoryCategory(ctx context.Context, category *domain.DirectoryCategory) error {
	if err := u.RequirePermission(ctx, enums.PermissionDirectoryCategoryCreate); err != nil {
		return err
	}
	return u.directoryCategoryRepo.Create(ctx, category)
}

func (u *DirectoryCategoryUsecase) GetDirectoryCategoryByID(ctx context.Context, id uint64) (*domain.DirectoryCategory, error) {
	if err := u.RequirePermission(ctx, enums.PermissionDirectoryCategoryRead); err != nil {
		return nil, err
	}
	return u.directoryCategoryRepo.GetByID(ctx, id)
}

func (u *DirectoryCategoryUsecase) UpdateDirectoryCategory(ctx context.Context, id uint64, category *domain.DirectoryCategory) (*domain.DirectoryCategory, error) {
	if err := u.RequirePermission(ctx, enums.PermissionDirectoryCategoryUpdate); err != nil {
		return nil, err
	}

	// Set ID for update
	category.ID = id
	err := u.directoryCategoryRepo.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	// Return updated category
	return u.directoryCategoryRepo.GetByID(ctx, id)
}

func (u *DirectoryCategoryUsecase) DeleteDirectoryCategory(ctx context.Context, id uint64) error {
	if err := u.RequirePermission(ctx, enums.PermissionDirectoryCategoryDelete); err != nil {
		return err
	}
	return u.directoryCategoryRepo.Delete(ctx, id)
}

func (u *DirectoryCategoryUsecase) ListDirectoryCategories(ctx context.Context, request *dto.ListDirectoryCategoriesRequest) ([]domain.DirectoryCategory, int64, error) {
	if err := u.RequirePermission(ctx, enums.PermissionDirectoryCategoryRead); err != nil {
		return nil, 0, err
	}

	categories, total, err := u.directoryCategoryRepo.List(ctx, request)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
