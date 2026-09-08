package repo

import (
	"context"

	"bdspro/internal/domain"
)

// IdentifierRepository defines ONLY the data operations contract
// This is just an interface - no implementation, no business logic
type PropertyIdentifierRepo interface {
	// Create inserts a new identifier
	Create(ctx context.Context, identifier *domain.CountryIdentifier) error

	// GetByID retrieves an identifier by its ID
	GetByID(ctx context.Context, id string) (*domain.CountryIdentifier, error)

	// Update updates an existing identifier
	Update(ctx context.Context, identifier *domain.CountryIdentifier) error

	// Delete soft deletes an identifier
	Delete(ctx context.Context, id string) error

	// Search finds identifiers based on filters
	Search(ctx context.Context, filters map[string]interface{}, page, limit int32) ([]*domain.CountryIdentifier, int64, error)

	// GetByLocation finds identifiers within a radius (purely spatial query)
	GetByLocation(ctx context.Context, lat, lng float64, radiusKm float64) ([]*domain.CountryIdentifier, error)

	// GetByOwner retrieves identifiers by owner ID
	GetByOwner(ctx context.Context, ownerID string) ([]*domain.CountryIdentifier, error)

	// CheckExists checks if identifier with given land parcel code exists
	CheckExists(ctx context.Context, landParcelCode string) (bool, error)
}
