package repo

import (
	"bdspro/internal/domain"
	"context"
)

type ApartmentRepo interface {
	// Create(c context.Context, entity *domain.Apartment) error
	// Update(c context.Context, id uint64, entity *domain.Apartment) error
	// Delete(c context.Context, id uint64) error
	// GetAll() ([]domain.Apartment, error)
	// GetByID(id uint64) (*domain.Apartment, error)

	UpdateApartments(c context.Context, apartments []domain.Apartment, attributeID uint64) error
	UpdateStatusApartments(c context.Context, apartments []domain.Apartment, status int, archived int) error
	UpdateApartment2(c context.Context, apartments []domain.Apartment, attributeID uint64) error
}
