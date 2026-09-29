package usecase

import (
	"context"
	"fmt"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"

	_utils "common/utils"
)

// DirectorySupplierUsecase implements DirectorySupplierUsecase interface
type DirectorySupplierUsecase struct {
	directorySupplierRepo repo.DirectorySupplierRepository
	userProvider          provider.UserProvider
	permissionProvider    provider.PermissionProvider
	mapper                *mapper.DirectorySupplierMapper
}

// NewDirectorySupplierUsecase creates a new DirectorySupplierUsecase
func NewDirectorySupplierUsecase(
	directorySupplierRepo repo.DirectorySupplierRepository,
	userProvider provider.UserProvider,
	permissionProvider provider.PermissionProvider,
	mapper *mapper.DirectorySupplierMapper,
) *DirectorySupplierUsecase {
	return &DirectorySupplierUsecase{
		directorySupplierRepo: directorySupplierRepo,
		userProvider:          userProvider,
		permissionProvider:    permissionProvider,
		mapper:                mapper,
	}
}

// Create creates a new directory supplier with permission check
func (u *DirectorySupplierUsecase) Create(ctx context.Context, req *dto.CreateDirectorySupplierRequestDTO) (*dto.DirectorySupplierDTO, error) {
	// Get user info
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized: user not found")
	}

	// Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, userID, enums.PermissionDirectorySupplierCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("permission denied: user %d does not have %s permission",
			userID, enums.PermissionDirectorySupplierCreate.String())
	}

	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}

	// Check if code already exists
	existing, err := u.directorySupplierRepo.GetByCode(ctx, req.Code)
	if err != nil {
		// Log error but continue - we only care if we found a record
	}
	if existing != nil {
		return nil, fmt.Errorf("directory supplier with code '%s' already exists", req.Code)
	}

	// Convert request to domain using mapper
	supplier, err := u.mapper.ToDomainFromCreateRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to map request to domain: %v", err)
	}

	// Set metadata
	supplier.CreatedBy = &userID
	supplier.UpdatedBy = &userID

	// Validate domain
	if err := supplier.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	// Create in repository
	if err := u.directorySupplierRepo.Create(ctx, supplier); err != nil {
		return nil, fmt.Errorf("failed to create directory supplier: %v", err)
	}

	// Get created supplier with relations
	created, err := u.directorySupplierRepo.GetByID(ctx, supplier.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created supplier: %v", err)
	}
	if created == nil {
		return nil, fmt.Errorf("created supplier not found")
	}

	// Convert to DTO using mapper
	return u.mapper.ToDTO(created), nil
}

// GetByID gets directory supplier by ID
func (u *DirectorySupplierUsecase) GetByID(ctx context.Context, id uint64) (*dto.DirectorySupplierDTO, error) {
	// Get user info
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized: user not found")
	}

	// Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, userID, enums.PermissionDirectorySupplierRead)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("permission denied: user %d does not have %s permission",
			userID, enums.PermissionDirectorySupplierRead.String())
	}

	// Validate ID
	if id == 0 {
		return nil, fmt.Errorf("invalid supplier ID")
	}

	// Get from repository
	supplier, err := u.directorySupplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory supplier: %v", err)
	}
	if supplier == nil {
		return nil, fmt.Errorf("directory supplier with ID %d not found", id)
	}

	return u.mapper.ToDTO(supplier), nil
}

func (u *DirectorySupplierUsecase) Update(ctx context.Context, id uint64, req *dto.UpdateDirectorySupplierRequestDTO) (*dto.DirectorySupplierDTO, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized: user not found")
	}

	hasPermission, err := u.permissionProvider.CheckPermission(ctx, userID, enums.PermissionDirectorySupplierUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("permission denied: user %d does not have %s permission",
			userID, enums.PermissionDirectorySupplierUpdate.String())
	}

	if id == 0 {
		return nil, fmt.Errorf("invalid supplier ID")
	}

	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}

	existing, err := u.directorySupplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory supplier: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("directory supplier with ID %d not found", id)
	}

	if req.Code != existing.Code {
		existingByCode, err := u.directorySupplierRepo.GetByCode(ctx, req.Code)
		if err != nil {
		}
		if existingByCode != nil && existingByCode.ID != id {
			return nil, fmt.Errorf("directory supplier with code '%s' already exists", req.Code)
		}
	}

	updated, err := u.mapper.ToDomainFromUpdateRequest(req, id)
	if err != nil {
		return nil, fmt.Errorf("failed to map request to domain: %v", err)
	}

	// Preserve existing fields
	updated.CreatedAt = existing.CreatedAt
	updated.CreatedBy = existing.CreatedBy
	updated.ServiceCount = existing.ServiceCount
	updated.ContractCount = existing.ContractCount
	updated.TotalValue = existing.TotalValue
	updated.UpdatedBy = &userID

	if err := updated.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	if err := u.directorySupplierRepo.Update(ctx, updated); err != nil {
		return nil, fmt.Errorf("failed to update directory supplier: %v", err)
	}

	result, err := u.directorySupplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated supplier: %v", err)
	}
	if result == nil {
		return nil, fmt.Errorf("updated supplier not found")
	}

	return u.mapper.ToDTO(result), nil
}

func (u *DirectorySupplierUsecase) Delete(ctx context.Context, id uint64) error {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return fmt.Errorf("unauthorized: user not found")
	}
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, userID, enums.PermissionDirectorySupplierDelete)
	if err != nil {
		return fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return fmt.Errorf("permission denied: user %d does not have %s permission",
			userID, enums.PermissionDirectorySupplierDelete.String())
	}

	// Validate ID
	if id == 0 {
		return fmt.Errorf("invalid supplier ID")
	}

	// Check if exists
	existing, err := u.directorySupplierRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get directory supplier: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("directory supplier with ID %d not found", id)
	}

	// TODO: Check if supplier has any dependencies (services, contracts, etc.)
	// if existing.ServiceCount > 0 {
	//     return fmt.Errorf("cannot delete supplier with existing services")
	// }
	// if existing.ContractCount > 0 {
	//     return fmt.Errorf("cannot delete supplier with existing contracts")
	// }

	if err := u.directorySupplierRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete directory supplier: %v", err)
	}

	return nil
}

func (u *DirectorySupplierUsecase) List(ctx context.Context, req *dto.DirectorySupplierFilterDTO) (*dto.ListDirectorySuppliersResponseDTO, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized: user not found")
	}

	hasPermission, err := u.permissionProvider.CheckPermission(ctx, userID, enums.PermissionDirectorySupplierList)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		hasReadPermission, err := u.permissionProvider.CheckPermission(ctx, userID, enums.PermissionDirectorySupplierRead)
		if err != nil || !hasReadPermission {
			return nil, fmt.Errorf("permission denied: user %d does not have %s or %s permission",
				userID, enums.PermissionDirectorySupplierList.String(), enums.PermissionDirectorySupplierRead.String())
		}
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}

	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	suppliers, total, err := u.directorySupplierRepo.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory suppliers: %v", err)
	}

	dtos := u.mapper.ToDTOList(suppliers)

	return &dto.ListDirectorySuppliersResponseDTO{
		Data:  dtos,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

func (u *DirectorySupplierUsecase) GetByCode(ctx context.Context, code string) (*dto.DirectorySupplierDTO, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized: user not found")
	}

	hasPermission, err := u.permissionProvider.CheckPermission(ctx, userID, enums.PermissionDirectorySupplierRead)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("permission denied: user %d does not have %s permission",
			userID, enums.PermissionDirectorySupplierRead.String())
	}

	// Validate code
	if code == "" {
		return nil, fmt.Errorf("code is required")
	}

	supplier, err := u.directorySupplierRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory supplier by code: %v", err)
	}
	if supplier == nil {
		return nil, fmt.Errorf("directory supplier with code '%s' not found", code)
	}

	// Convert to DTO using mapper
	return u.mapper.ToDTO(supplier), nil
}

func (u *DirectorySupplierUsecase) UpdateStats(ctx context.Context, id uint64, serviceCount, contractCount uint32, totalValue int64) (*dto.DirectorySupplierDTO, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized: user not found")
	}

	hasPermission, err := u.permissionProvider.CheckPermission(ctx, userID, enums.PermissionDirectorySupplierUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("permission denied: user %d does not have %s permission",
			userID, enums.PermissionDirectorySupplierUpdate.String())
	}

	existing, err := u.directorySupplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory supplier: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("directory supplier with ID %d not found", id)
	}

	// Update stats
	existing.ServiceCount = serviceCount
	existing.ContractCount = contractCount
	existing.TotalValue = totalValue
	existing.UpdatedBy = &userID

	// Update in repository
	if err := u.directorySupplierRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update supplier stats: %v", err)
	}

	// Get updated supplier
	result, err := u.directorySupplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated supplier: %v", err)
	}

	// Convert to DTO using mapper
	return u.mapper.ToDTO(result), nil
}
