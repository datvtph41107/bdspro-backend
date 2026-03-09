package usecases

import (
	"context"
	"fmt"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
)

type PropertyIdentifierUsecase struct {
	repo repo.PropertyIdentifierRepo
}

func NewIdentifierUsecase(repo repo.PropertyIdentifierRepo) *PropertyIdentifierUsecase {
	return &PropertyIdentifierUsecase{
		repo: repo,
	}
}

func (u *PropertyIdentifierUsecase) Create(ctx context.Context, identifier *domain.CountryIdentifier, requestedBy string) (*domain.CountryIdentifier, error) {
	if identifier == nil {
		return nil, fmt.Errorf("%w: identifier cannot be nil", domain.ErrInvalidInput)
	}
	if identifier.ProvinceID == "" {
		return nil, fmt.Errorf("%w: province_id is required", domain.ErrInvalidInput)
	}
	if identifier.DistrictID == "" {
		return nil, fmt.Errorf("%w: district_id is required", domain.ErrInvalidInput)
	}
	if identifier.WardID == "" {
		return nil, fmt.Errorf("%w: ward_id is required", domain.ErrInvalidInput)
	}
	if identifier.Type == "" {
		return nil, fmt.Errorf("%w: type is required", domain.ErrInvalidInput)
	}

	// Validate type codes
	// validTypes := map[string]bool{
	// 	"APT": true, "HSE": true, "LND": true,
	// 	"OFF": true, "RTL": true, "IND": true,
	// }
	// if !validTypes[identifier.Type] {
	// 	return nil, fmt.Errorf("%w: invalid type '%s'", domain.ErrInvalidInput, identifier.Type)
	// }

	if identifier.LegalStatus != "" {
		validStatuses := map[string]bool{
			"AVL": true, "MTG": true, "DSP": true,
			"SOLD": true, "RSRV": true,
		}
		if !validStatuses[identifier.LegalStatus] {
			return nil, fmt.Errorf("%w: invalid legal status '%s'", domain.ErrInvalidInput, identifier.LegalStatus)
		}
	} else {
		identifier.LegalStatus = "AVL" // Default
	}

	if identifier.Latitude != 0 || identifier.Longitude != 0 {
		if identifier.Latitude < -90 || identifier.Latitude > 90 {
			return nil, fmt.Errorf("%w: latitude must be between -90 and 90", domain.ErrInvalidInput)
		}
		if identifier.Longitude < -180 || identifier.Longitude > 180 {
			return nil, fmt.Errorf("%w: longitude must be between -180 and 180", domain.ErrInvalidInput)
		}
	}

	if identifier.LandParcelCode != "" {
		exists, err := u.repo.CheckExists(ctx, identifier.LandParcelCode)
		if err != nil {
			return nil, fmt.Errorf("failed to check existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("%w: land parcel code '%s' already exists",
				domain.ErrDuplicate, identifier.LandParcelCode)
		}
	}

	now := time.Now()
	identifier.CreatedAt = now
	identifier.UpdatedAt = now

	err := u.repo.Create(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to create identifier: %w", err)
	}

	return identifier, nil
}

func (u *PropertyIdentifierUsecase) GetByID(ctx context.Context, id string) (*domain.CountryIdentifier, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if len(id) < 10 || len(id) > 40 {
		return nil, fmt.Errorf("%w: invalid id format", domain.ErrInvalidInput)
	}

	identifier, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if identifier == nil {
		return nil, fmt.Errorf("%w: identifier with id '%s' not found", domain.ErrNotFound, id)
	}

	return identifier, nil
}

func (u *PropertyIdentifierUsecase) Update(ctx context.Context, identifier *domain.CountryIdentifier, requestedBy string) (*domain.CountryIdentifier, error) {
	if identifier == nil || identifier.ID == "" {
		return nil, fmt.Errorf("%w: identifier with id is required", domain.ErrInvalidInput)
	}

	existing, err := u.repo.GetByID(ctx, identifier.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("%w: identifier with id '%s' not found",
			domain.ErrNotFound, identifier.ID)
	}

	identifier.CreatedAt = existing.CreatedAt
	identifier.UpdatedAt = time.Now()

	if identifier.Type != "" && identifier.Type != existing.Type {
		if existing.LegalStatus == "SOLD" {
			return nil, fmt.Errorf("%w: cannot change type of sold property",
				domain.ErrInvalidInput)
		}
	}

	err = u.repo.Update(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to update identifier: %w", err)
	}

	return identifier, nil
}

func (u *PropertyIdentifierUsecase) Delete(ctx context.Context, id string, requestedBy string) error {
	if id == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}

	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("%w: identifier with id '%s' not found", domain.ErrNotFound, id)
	}

	if existing.LegalStatus == "MTG" {
		return fmt.Errorf("%w: cannot delete mortgaged property", domain.ErrInvalidInput)
	}

	err = u.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete identifier: %w", err)
	}

	return nil
}

func (u *PropertyIdentifierUsecase) Search(ctx context.Context, criteria *dto.SearchCriteriaDTO, pagination *dto.PaginationDTO) ([]*domain.CountryIdentifier, int64, error) {
	if pagination == nil {
		pagination = &dto.PaginationDTO{Page: 1, Limit: 10}
	}
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 || pagination.Limit > 100 {
		pagination.Limit = 10
	}

	filters := make(map[string]interface{})
	if criteria != nil {
		if criteria.ProvinceID != "" {
			filters["province_id"] = criteria.ProvinceID
		}
		if criteria.DistrictID != "" {
			filters["district_id"] = criteria.DistrictID
		}
		if criteria.WardID != "" {
			filters["ward_id"] = criteria.WardID
		}
		if criteria.Type != "" {
			filters["type"] = criteria.Type
		}
		if criteria.LegalStatus != "" {
			filters["legal_status"] = criteria.LegalStatus
		}
		if criteria.OwnerID != "" {
			filters["current_owner_id"] = criteria.OwnerID
		}
	}

	identifiers, total, err := u.repo.Search(ctx, filters, pagination.Page, pagination.Limit)
	if err != nil {
		return nil, 0, fmt.Errorf("search failed: %w", err)
	}

	return identifiers, total, nil
}

func (u *PropertyIdentifierUsecase) TransferOwnership(ctx context.Context, id string, newOwnerID string, requestedBy string) (*domain.CountryIdentifier, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if newOwnerID == "" {
		return nil, fmt.Errorf("%w: new owner id is required", domain.ErrInvalidInput)
	}

	identifier, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if identifier == nil {
		return nil, fmt.Errorf("%w: identifier with id '%s' not found", domain.ErrNotFound, id)
	}

	if identifier.LegalStatus == "MTG" {
		return nil, fmt.Errorf("%w: cannot transfer mortgaged property", domain.ErrInvalidInput)
	}
	if identifier.LegalStatus == "DSP" {
		return nil, fmt.Errorf("%w: cannot transfer disputed property", domain.ErrInvalidInput)
	}
	if identifier.CurrentOwnerID == newOwnerID {
		return nil, fmt.Errorf("%w: new owner is the same as current owner", domain.ErrInvalidInput)
	}

	identifier.CurrentOwnerID = newOwnerID
	identifier.UpdatedAt = time.Now()

	if identifier.LegalStatus == "SOLD" {
		identifier.LegalStatus = "AVL"
	}

	err = u.repo.Update(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer ownership: %w", err)
	}

	return identifier, nil
}

func (u *PropertyIdentifierUsecase) UpdateLegalStatus(ctx context.Context, id string, status string, requestedBy string) (*domain.CountryIdentifier, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if status == "" {
		return nil, fmt.Errorf("%w: status is required", domain.ErrInvalidInput)
	}

	// Validate status
	validStatuses := map[string]bool{
		"AVL": true, "MTG": true, "DSP": true,
		"SOLD": true, "RSRV": true,
	}
	if !validStatuses[status] {
		return nil, fmt.Errorf("%w: invalid legal status '%s'", domain.ErrInvalidInput, status)
	}

	// Get identifier
	identifier, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if identifier == nil {
		return nil, fmt.Errorf("%w: identifier with id '%s' not found", domain.ErrNotFound, id)
	}

	identifier.LegalStatus = status
	identifier.UpdatedAt = time.Now()

	err = u.repo.Update(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to update legal status: %w", err)
	}

	return identifier, nil
}

func (u *PropertyIdentifierUsecase) ValidateIdentifier(ctx context.Context, code string) (*domain.CountryIdentifier, bool, error) {
	if code == "" {
		return nil, false, fmt.Errorf("%w: code is required", domain.ErrInvalidInput)
	}

	identifier, err := u.repo.GetByID(ctx, code)
	if err != nil {
		return nil, false, err
	}
	if identifier == nil {
		return nil, false, nil
	}

	return identifier, true, nil
}
