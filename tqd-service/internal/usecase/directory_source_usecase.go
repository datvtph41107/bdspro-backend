package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"

	_utils "common/utils"
)

// DirectorySourceUsecase implements DirectorySourceUsecase interface
type DirectorySourceUsecase struct {
	directorySourceRepo repo.DirectorySourceRepository
	userProvider        provider.UserProvider
	permissionProvider  provider.PermissionProvider
}

// NewDirectorySourceUsecase creates a new DirectorySourceUsecase
func NewDirectorySourceUsecase(
	directorySourceRepo repo.DirectorySourceRepository,
	userProvider provider.UserProvider,
	permissionProvider provider.PermissionProvider,
) *DirectorySourceUsecase {
	return &DirectorySourceUsecase{
		directorySourceRepo: directorySourceRepo,
		userProvider:        userProvider,
		permissionProvider:  permissionProvider,
	}
}

func (u *DirectorySourceUsecase) Create(ctx context.Context, req *dto.DirectorySourceDTO) (*dto.DirectorySourceDTO, error) {
	// 1. Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// 2. Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionDirectorySourceCreate)
	if err != nil || !hasPermission {
		return nil, fmt.Errorf("permission denied")
	}

	// 3. Convert DTO to Domain
	domainSource := req.ToDomain()

	// 4. Validate domain
	if err := domainSource.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	// 5. Call repository
	err = u.directorySourceRepo.Create(ctx, domainSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory source: %v", err)
	}

	// 6. Convert Domain to DTO
	result := &dto.DirectorySourceDTO{}
	result.FromDomain(domainSource)

	return result, nil
}

func (u *DirectorySourceUsecase) GetByID(ctx context.Context, id uint64) (*dto.DirectorySourceDTO, error) {
	// 1. Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// 2. Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionDirectorySourceRead)
	if err != nil || !hasPermission {
		return nil, fmt.Errorf("permission denied")
	}

	// 3. Call repository
	directorySource, err := u.directorySourceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory source: %v", err)
	}

	// 4. Convert Domain to DTO
	result := &dto.DirectorySourceDTO{}
	result.FromDomain(directorySource)

	return result, nil
}

func (u *DirectorySourceUsecase) Update(ctx context.Context, req *domain.DirectorySource) (*dto.DirectorySourceDTO, error) {
	// 1. Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// 2. Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionDirectorySourceUpdate)
	if err != nil || !hasPermission {
		return nil, fmt.Errorf("permission denied")
	}

	// 3. Validate domain
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	// 4. Call repository
	err = u.directorySourceRepo.Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update directory source: %v", err)
	}

	// 5. Convert Domain to DTO
	result := &dto.DirectorySourceDTO{}
	result.FromDomain(req)

	return result, nil
}

func (u *DirectorySourceUsecase) Delete(ctx context.Context, id uint64) error {
	// 1. Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// 2. Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionDirectorySourceDelete)
	if err != nil || !hasPermission {
		return fmt.Errorf("permission denied")
	}

	// 3. Call repository directly with domain ID
	err = u.directorySourceRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete directory source: %v", err)
	}

	return nil
}

func (u *DirectorySourceUsecase) List(ctx context.Context, req *dto.ListDirectorySourcesRequestDTO) ([]dto.DirectorySourceDTO, int64, error) {
	// 1. Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// 2. Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionDirectorySourceRead)
	if err != nil || !hasPermission {
		return nil, 0, fmt.Errorf("permission denied")
	}

	// 3. Call repository with ListRequest DTO
	entities, total, err := u.directorySourceRepo.List(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list directory sources: %v", err)
	}

	for i, entity := range entities {
		if entity.RawTags != "" {
			tags := []string{}
			json.Unmarshal([]byte(entity.RawTags), &tags)
			for _, tag := range tags {
				entities[i].Tags = append(entities[i].Tags, tag)
			}
		}
	}

	return entities, total, nil
}
