package repo

import (
	_crud "common/domain/crud"
	"tqd/internal/domain"
)

// IAmenityRepo defines the interface for amenity repository operations
type IAmenityRepo interface {
	_crud.ICrudRepo[domain.Amenity]
}
