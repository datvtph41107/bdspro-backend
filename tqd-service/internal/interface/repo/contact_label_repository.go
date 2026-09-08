package repo

import (
	"context"
	"tqd/internal/domain"

	_dto "common/domain/dto"
)

// ContactLabelRepository defines the interface for contact label repository
type ContactLabelRepository interface {
	// Create creates a new contact label
	Create(ctx context.Context, contactLabel *domain.ContactLabel) error

	// GetByID gets contact label by ID
	GetByID(ctx context.Context, id uint64) (*domain.ContactLabel, error)

	// GetByCode gets contact label by code
	GetByCode(ctx context.Context, code string) (*domain.ContactLabel, error)

	// Update updates an existing contact label
	Update(ctx context.Context, contactLabel *domain.ContactLabel) error

	// Delete deletes a contact label by ID
	Delete(ctx context.Context, id uint64) error

	// List gets list of contact labels with pagination
	List(ctx context.Context, pagable _dto.Pagable, search string, isActive *bool, isSystem *bool) ([]domain.ContactLabel, int64, error)

	// UpdateContactCount updates the contact count for a label
	UpdateContactCount(ctx context.Context, id uint64, count int) error
}
