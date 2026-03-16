package usecase

import (
	_utils "common/utils"
	"context"
	"fmt"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"
)

// ContactLabelUsecase defines the interface for contact label usecase
type ContactLabelUsecase interface {
	Create(ctx context.Context, req *dto.CreateContactLabelRequestDTO) (*dto.ContactLabelDTO, error)
	GetByID(ctx context.Context, id uint64) (*dto.ContactLabelDTO, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateContactLabelRequestDTO) (*dto.ContactLabelDTO, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, req *dto.ListContactLabelsRequestDTO) (*dto.ListContactLabelsResponseDTO, error)
	UpdateContactCount(ctx context.Context, id uint64, count int) error
}

// contactLabelUsecase implements ContactLabelUsecase interface
type contactLabelUsecase struct {
	contactLabelRepo   repo.ContactLabelRepository
	userProvider       provider.UserProvider
	permissionProvider provider.PermissionProvider
}

// NewContactLabelUsecase creates a new ContactLabelUsecase
func NewContactLabelUsecase(
	contactLabelRepo repo.ContactLabelRepository,
	userProvider provider.UserProvider,
	permissionProvider provider.PermissionProvider,
) ContactLabelUsecase {
	return &contactLabelUsecase{
		contactLabelRepo:   contactLabelRepo,
		userProvider:       userProvider,
		permissionProvider: permissionProvider,
	}
}

func (u *contactLabelUsecase) Create(ctx context.Context, req *dto.CreateContactLabelRequestDTO) (*dto.ContactLabelDTO, error) {
	// Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionContactLabelCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("insufficient permissions")
	}

	// Check if code already exists
	existing, err := u.contactLabelRepo.GetByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("contact label with code %s already exists", req.Code)
	}

	// Create domain object
	contactLabel := &domain.ContactLabel{
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		Color:        req.Color,
		Icon:         req.Icon,
		IsActive:     req.IsActive,
		SortOrder:    req.SortOrder,
		ContactCount: 0,
		IsSystem:     req.IsSystem,
		Status:       0, // Active
		Type:         enums.ContactLabelType(req.Type),
	}

	// Create in repository
	err = u.contactLabelRepo.Create(ctx, contactLabel)
	if err != nil {
		return nil, fmt.Errorf("failed to create contact label: %v", err)
	}

	// Convert to DTO
	dto := &dto.ContactLabelDTO{}
	dto.FromDomain(contactLabel)

	return dto, nil
}

func (u *contactLabelUsecase) GetByID(ctx context.Context, id uint64) (*dto.ContactLabelDTO, error) {
	// Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionContactLabelRead)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("insufficient permissions")
	}

	// Get from repository
	contactLabel, err := u.contactLabelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact label: %v", err)
	}

	// Convert to DTO
	dto := &dto.ContactLabelDTO{}
	dto.FromDomain(contactLabel)

	return dto, nil
}

func (u *contactLabelUsecase) Update(ctx context.Context, id uint64, req *dto.UpdateContactLabelRequestDTO) (*dto.ContactLabelDTO, error) {
	// Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionContactLabelUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("insufficient permissions")
	}

	// Get existing contact label
	existing, err := u.contactLabelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact label: %v", err)
	}

	// Check if code already exists (excluding current record)
	existingByCode, err := u.contactLabelRepo.GetByCode(ctx, req.Code)
	if err == nil && existingByCode != nil && existingByCode.ID != id {
		return nil, fmt.Errorf("contact label with code %s already exists", req.Code)
	}

	// Update fields
	existing.Name = req.Name
	existing.Code = req.Code
	existing.Description = req.Description
	existing.Color = req.Color
	existing.Icon = req.Icon
	existing.IsActive = req.IsActive
	existing.SortOrder = req.SortOrder
	existing.IsSystem = req.IsSystem
	existing.Type = enums.ContactLabelType(req.Type)

	// Update in repository
	err = u.contactLabelRepo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update contact label: %v", err)
	}

	// Convert to DTO
	dto := &dto.ContactLabelDTO{}
	dto.FromDomain(existing)

	return dto, nil
}

func (u *contactLabelUsecase) Delete(ctx context.Context, id uint64) error {
	// Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionContactLabelDelete)
	if err != nil {
		return fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return fmt.Errorf("insufficient permissions")
	}

	// Delete from repository
	err = u.contactLabelRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete contact label: %v", err)
	}

	return nil
}

func (u *contactLabelUsecase) List(ctx context.Context, req *dto.ListContactLabelsRequestDTO) (*dto.ListContactLabelsResponseDTO, error) {
	// Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionContactLabelRead)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("insufficient permissions")
	}

	// Get list from repository
	contactLabels, total, err := u.contactLabelRepo.List(ctx, req.Pagable, req.Search, req.IsActive, req.IsSystem)
	if err != nil {
		return nil, fmt.Errorf("failed to list contact labels: %v", err)
	}

	// Convert to DTOs
	dtoList := make([]dto.ContactLabelDTO, len(contactLabels))
	for i, contactLabel := range contactLabels {
		dto := &dto.ContactLabelDTO{}
		dto.FromDomain(&contactLabel)
		dtoList[i] = *dto
	}

	return &dto.ListContactLabelsResponseDTO{
		Data:  dtoList,
		Total: total,
	}, nil
}

func (u *contactLabelUsecase) UpdateContactCount(ctx context.Context, id uint64, count int) error {
	// Get user info
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Check permission
	hasPermission, err := u.permissionProvider.CheckPermission(ctx, profileID, enums.PermissionContactLabelUpdate)
	if err != nil {
		return fmt.Errorf("failed to check permission: %v", err)
	}
	if !hasPermission {
		return fmt.Errorf("insufficient permissions")
	}

	// Update contact count in repository
	err = u.contactLabelRepo.UpdateContactCount(ctx, id, count)
	if err != nil {
		return fmt.Errorf("failed to update contact count: %v", err)
	}

	return nil
}
