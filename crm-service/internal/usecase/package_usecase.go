package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
	"fmt"
)

type PackageUsecase interface {
	CreatePackage(ctx context.Context, createDTO *dto.PackageCreateDTO) (*dto.PackageResponseDTO, error)
	UpdatePackage(ctx context.Context, id uint64, updateDTO *dto.PackageUpdateDTO) (*dto.PackageResponseDTO, error)
	DeletePackage(ctx context.Context, id uint64) error
	GetPackage(ctx context.Context, id uint64) (*dto.PackageResponseDTO, error)
	SearchPackages(ctx context.Context, searchDTO dto.PackageSearchDTO) ([]*dto.PackageResponseDTO, int64, error)
	GetPackagesByType(ctx context.Context, packageType domain.PackageType) ([]*dto.PackageResponseDTO, error)

	// D.10.2.2 specific methods
	SelectPackage(ctx context.Context, selectionDTO *dto.PackageSelectionDTO) (*dto.PackageSelectionResponseDTO, error)
	GetPackageRecommendations(ctx context.Context, recommendationDTO *dto.PackageRecommendationDTO) (*dto.PackageRecommendationResponseDTO, error)
	GetPopularPackages(ctx context.Context, limit int) ([]*dto.PackageResponseDTO, error)
	GetRecommendedPackages(ctx context.Context, limit int) ([]*dto.PackageResponseDTO, error)
	GetBestValuePackages(ctx context.Context, limit int) ([]*dto.PackageResponseDTO, error)
	ValidatePackageSelection(ctx context.Context, packageID uint64, customBudget *float64, customDuration *int) error
}

type packageUsecase struct {
	packageRepo repo.PackageRepo
}

func NewPackageUsecase(packageRepo repo.PackageRepo) PackageUsecase {
	return &packageUsecase{
		packageRepo: packageRepo,
	}
}

func (u *packageUsecase) CreatePackage(ctx context.Context, createDTO *dto.PackageCreateDTO) (*dto.PackageResponseDTO, error) {
	// Validate package data
	if err := u.validatePackageData(createDTO); err != nil {
		return nil, err
	}

	pkg := &domain.Package{
		Name:                 createDTO.Name,
		Description:          createDTO.Description,
		Type:                 createDTO.Type,
		Status:               domain.PackageStatusActive,
		Price:                createDTO.Price,
		Currency:             createDTO.Currency,
		Duration:             createDTO.Duration,
		Features:             createDTO.Features,
		MaxImpressions:       createDTO.MaxImpressions,
		MaxClicks:            createDTO.MaxClicks,
		MinBudget:            createDTO.MinBudget,
		MaxBudget:            createDTO.MaxBudget,
		EstimatedReach:       createDTO.EstimatedReach,
		EstimatedClicks:      createDTO.EstimatedClicks,
		IsPopular:            createDTO.IsPopular,
		IsRecommended:        createDTO.IsRecommended,
		Icon:                 createDTO.Icon,
		Color:                createDTO.Color,
		RequiresVerification: createDTO.RequiresVerification,
		RequiresPayment:      createDTO.RequiresPayment,
	}

	createdPackage, err := u.packageRepo.Create(ctx, pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to create package: %w", err)
	}

	return u.toResponseDTO(createdPackage), nil
}

func (u *packageUsecase) UpdatePackage(ctx context.Context, id uint64, updateDTO *dto.PackageUpdateDTO) (*dto.PackageResponseDTO, error) {
	existingPackage, err := u.packageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	// Update fields
	if updateDTO.Name != "" {
		existingPackage.Name = updateDTO.Name
	}
	if updateDTO.Description != "" {
		existingPackage.Description = updateDTO.Description
	}
	if updateDTO.Type != "" {
		existingPackage.Type = updateDTO.Type
	}
	if updateDTO.Status != "" {
		existingPackage.Status = updateDTO.Status
	}
	if updateDTO.Price > 0 {
		existingPackage.Price = updateDTO.Price
	}
	if updateDTO.Currency != "" {
		existingPackage.Currency = updateDTO.Currency
	}
	if updateDTO.Duration > 0 {
		existingPackage.Duration = updateDTO.Duration
	}
	if updateDTO.Features != "" {
		existingPackage.Features = updateDTO.Features
	}

	existingPackage.MaxImpressions = updateDTO.MaxImpressions
	existingPackage.MaxClicks = updateDTO.MaxClicks
	existingPackage.MinBudget = updateDTO.MinBudget
	existingPackage.MaxBudget = updateDTO.MaxBudget
	existingPackage.EstimatedReach = updateDTO.EstimatedReach
	existingPackage.EstimatedClicks = updateDTO.EstimatedClicks
	existingPackage.IsPopular = updateDTO.IsPopular
	existingPackage.IsRecommended = updateDTO.IsRecommended
	existingPackage.Icon = updateDTO.Icon
	existingPackage.Color = updateDTO.Color
	existingPackage.RequiresVerification = updateDTO.RequiresVerification
	existingPackage.RequiresPayment = updateDTO.RequiresPayment

	updatedPackage, err := u.packageRepo.Update(ctx, id, existingPackage)
	if err != nil {
		return nil, fmt.Errorf("failed to update package: %w", err)
	}

	return u.toResponseDTO(updatedPackage), nil
}

func (u *packageUsecase) DeletePackage(ctx context.Context, id uint64) error {
	return u.packageRepo.Delete(ctx, id)
}

func (u *packageUsecase) GetPackage(ctx context.Context, id uint64) (*dto.PackageResponseDTO, error) {
	pkg, err := u.packageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	return u.toResponseDTO(pkg), nil
}

func (u *packageUsecase) SearchPackages(ctx context.Context, searchDTO dto.PackageSearchDTO) ([]*dto.PackageResponseDTO, int64, error) {
	// TODO: Get organization ID from context
	organizationID := uint64(1) // Placeholder

	packages, total, err := u.packageRepo.Search(ctx, organizationID, searchDTO)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search packages: %w", err)
	}

	responseDTOs := make([]*dto.PackageResponseDTO, len(packages))
	for i, pkg := range packages {
		responseDTOs[i] = u.toResponseDTO(pkg)
	}

	return responseDTOs, total, nil
}

func (u *packageUsecase) GetPackagesByType(ctx context.Context, packageType domain.PackageType) ([]*dto.PackageResponseDTO, error) {
	packages, err := u.packageRepo.GetByType(ctx, packageType)
	if err != nil {
		return nil, fmt.Errorf("failed to get packages by type: %w", err)
	}

	responseDTOs := make([]*dto.PackageResponseDTO, len(packages))
	for i, pkg := range packages {
		responseDTOs[i] = u.toResponseDTO(pkg)
	}

	return responseDTOs, nil
}

// D.10.2.2 specific methods

func (u *packageUsecase) SelectPackage(ctx context.Context, selectionDTO *dto.PackageSelectionDTO) (*dto.PackageSelectionResponseDTO, error) {
	// Validate package selection
	if err := u.ValidatePackageSelection(ctx, selectionDTO.PackageID, selectionDTO.CustomBudget, selectionDTO.Duration); err != nil {
		return nil, err
	}

	// Get package details
	pkg, err := u.packageRepo.GetByID(ctx, selectionDTO.PackageID)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	// Calculate total price and estimates
	totalPrice := pkg.Price
	var estimatedReach, estimatedClicks *int

	// For VIP packages with custom duration
	if pkg.Type == domain.PackageTypeVIP && selectionDTO.Duration != nil {
		durationRatio := float64(*selectionDTO.Duration) / float64(pkg.Duration)
		totalPrice = pkg.Price * durationRatio
		if pkg.EstimatedReach != nil {
			reach := int(float64(*pkg.EstimatedReach) * durationRatio)
			estimatedReach = &reach
		}
		if pkg.EstimatedClicks != nil {
			clicks := int(float64(*pkg.EstimatedClicks) * durationRatio)
			estimatedClicks = &clicks
		}
	}

	// For Ads packages with custom budget
	if (pkg.Type == domain.PackageTypeMetaAds || pkg.Type == domain.PackageTypeGoogleAds) && selectionDTO.CustomBudget != nil {
		budgetRatio := *selectionDTO.CustomBudget / pkg.Price
		totalPrice = *selectionDTO.CustomBudget
		if pkg.EstimatedReach != nil {
			reach := int(float64(*pkg.EstimatedReach) * budgetRatio)
			estimatedReach = &reach
		}
		if pkg.EstimatedClicks != nil {
			clicks := int(float64(*pkg.EstimatedClicks) * budgetRatio)
			estimatedClicks = &clicks
		}
	}

	// Determine next step
	nextStep := "payment"
	if !pkg.RequiresPayment {
		nextStep = "campaign_creation"
	}

	return &dto.PackageSelectionResponseDTO{
		SelectedPackage: *u.toResponseDTO(pkg),
		CustomBudget:    selectionDTO.CustomBudget,
		CustomDuration:  selectionDTO.Duration,
		TotalPrice:      totalPrice,
		EstimatedReach:  estimatedReach,
		EstimatedClicks: estimatedClicks,
		NextStep:        nextStep,
	}, nil
}

func (u *packageUsecase) GetPackageRecommendations(ctx context.Context, recommendationDTO *dto.PackageRecommendationDTO) (*dto.PackageRecommendationResponseDTO, error) {
	// Get popular packages
	popularPackages, err := u.GetPopularPackages(ctx, 5)
	if err != nil {
		return nil, err
	}

	// Get recommended packages
	recommendedPackages, err := u.GetRecommendedPackages(ctx, 5)
	if err != nil {
		return nil, err
	}

	// Get best value packages
	bestValuePackages, err := u.GetBestValuePackages(ctx, 5)
	if err != nil {
		return nil, err
	}

	// TODO: Implement combo suggestions logic
	var comboSuggestions []dto.ComboPackageDTO

	return &dto.PackageRecommendationResponseDTO{
		RecommendedPackages: recommendedPackages,
		PopularPackages:     popularPackages,
		BestValuePackages:   bestValuePackages,
		ComboSuggestions:    comboSuggestions,
	}, nil
}

func (u *packageUsecase) GetPopularPackages(ctx context.Context, limit int) ([]*dto.PackageResponseDTO, error) {
	packages, err := u.packageRepo.GetPopularPackages(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular packages: %w", err)
	}

	responseDTOs := make([]*dto.PackageResponseDTO, len(packages))
	for i, pkg := range packages {
		responseDTOs[i] = u.toResponseDTO(pkg)
	}

	return responseDTOs, nil
}

func (u *packageUsecase) GetRecommendedPackages(ctx context.Context, limit int) ([]*dto.PackageResponseDTO, error) {
	packages, err := u.packageRepo.GetRecommendedPackages(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommended packages: %w", err)
	}

	responseDTOs := make([]*dto.PackageResponseDTO, len(packages))
	for i, pkg := range packages {
		responseDTOs[i] = u.toResponseDTO(pkg)
	}

	return responseDTOs, nil
}

func (u *packageUsecase) GetBestValuePackages(ctx context.Context, limit int) ([]*dto.PackageResponseDTO, error) {
	packages, err := u.packageRepo.GetBestValuePackages(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get best value packages: %w", err)
	}

	responseDTOs := make([]*dto.PackageResponseDTO, len(packages))
	for i, pkg := range packages {
		responseDTOs[i] = u.toResponseDTO(pkg)
	}

	return responseDTOs, nil
}

func (u *packageUsecase) ValidatePackageSelection(ctx context.Context, packageID uint64, customBudget *float64, customDuration *int) error {
	pkg, err := u.packageRepo.GetByID(ctx, packageID)
	if err != nil {
		return fmt.Errorf("package not found: %w", err)
	}

	// Validate custom budget for Ads packages
	if (pkg.Type == domain.PackageTypeMetaAds || pkg.Type == domain.PackageTypeGoogleAds) && customBudget != nil {
		if pkg.MinBudget != nil && *customBudget < *pkg.MinBudget {
			return fmt.Errorf("budget must be at least %f %s", *pkg.MinBudget, pkg.Currency)
		}
		if pkg.MaxBudget != nil && *customBudget > *pkg.MaxBudget {
			return fmt.Errorf("budget cannot exceed %f %s", *pkg.MaxBudget, pkg.Currency)
		}
	}

	// Validate custom duration for VIP packages
	if pkg.Type == domain.PackageTypeVIP && customDuration != nil {
		if *customDuration < 1 {
			return fmt.Errorf("duration must be at least 1 day")
		}
		if *customDuration > 365 {
			return fmt.Errorf("duration cannot exceed 365 days")
		}
	}

	return nil
}

// Helper methods

func (u *packageUsecase) validatePackageData(createDTO *dto.PackageCreateDTO) error {
	if createDTO.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	if createDTO.Duration < 1 {
		return fmt.Errorf("duration must be at least 1 day")
	}
	if createDTO.MinBudget != nil && createDTO.MaxBudget != nil && *createDTO.MinBudget > *createDTO.MaxBudget {
		return fmt.Errorf("min budget cannot be greater than max budget")
	}
	return nil
}

func (u *packageUsecase) toResponseDTO(pkg *domain.Package) *dto.PackageResponseDTO {
	return &dto.PackageResponseDTO{
		ID:                   pkg.ID,
		Name:                 pkg.Name,
		Description:          pkg.Description,
		Type:                 pkg.Type,
		Status:               pkg.Status,
		Price:                pkg.Price,
		Currency:             pkg.Currency,
		Duration:             pkg.Duration,
		Features:             pkg.Features,
		MaxImpressions:       pkg.MaxImpressions,
		MaxClicks:            pkg.MaxClicks,
		MinBudget:            pkg.MinBudget,
		MaxBudget:            pkg.MaxBudget,
		EstimatedReach:       pkg.EstimatedReach,
		EstimatedClicks:      pkg.EstimatedClicks,
		IsPopular:            pkg.IsPopular,
		IsRecommended:        pkg.IsRecommended,
		Icon:                 pkg.Icon,
		Color:                pkg.Color,
		RequiresVerification: pkg.RequiresVerification,
		RequiresPayment:      pkg.RequiresPayment,
		CreatedBy:            pkg.CreatedBy,
		UpdatedBy:            pkg.UpdatedBy,
		OrganizationID:       pkg.OrganizationID,
		CreatedAt:            pkg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:            pkg.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}