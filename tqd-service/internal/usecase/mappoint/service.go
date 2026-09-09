package mappoint

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domain "tqd/internal/domain/mappoint"
)

var (
	ErrInvalidInput = errors.New("invalid map point input")
	ErrNotFound     = errors.New("map point not found")
)

type Repository interface {
	Create(context.Context, domain.Point) (domain.Point, error)
	Get(context.Context, uint64) (domain.Point, error)
	Update(context.Context, domain.Point) (domain.Point, error)
	Delete(context.Context, uint64) error
	FindNearby(context.Context, domain.Coordinate, float64, int) ([]domain.Point, error)
	FindWithinPolygon(context.Context, []domain.Coordinate, int) ([]domain.Point, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func validateCoordinate(coordinate domain.Coordinate) error {
	if coordinate.Latitude < -90 || coordinate.Latitude > 90 {
		return fmt.Errorf("%w: latitude must be between -90 and 90", ErrInvalidInput)
	}
	if coordinate.Longitude < -180 || coordinate.Longitude > 180 {
		return fmt.Errorf("%w: longitude must be between -180 and 180", ErrInvalidInput)
	}
	return nil
}

func validatePoint(point domain.Point) error {
	if err := validateCoordinate(domain.Coordinate{Latitude: point.Latitude, Longitude: point.Longitude}); err != nil {
		return err
	}
	if strings.TrimSpace(point.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	return nil
}

func (s *Service) Create(ctx context.Context, point domain.Point) (domain.Point, error) {
	if s == nil || s.repository == nil {
		return domain.Point{}, errors.New("map point repository is unavailable")
	}
	point.Name = strings.TrimSpace(point.Name)
	if err := validatePoint(point); err != nil {
		return domain.Point{}, err
	}
	return s.repository.Create(ctx, point)
}

func (s *Service) Get(ctx context.Context, id uint64) (domain.Point, error) {
	if s == nil || s.repository == nil {
		return domain.Point{}, errors.New("map point repository is unavailable")
	}
	if id == 0 {
		return domain.Point{}, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}
	return s.repository.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, point domain.Point) (domain.Point, error) {
	if s == nil || s.repository == nil {
		return domain.Point{}, errors.New("map point repository is unavailable")
	}
	if point.ID == 0 {
		return domain.Point{}, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}
	point.Name = strings.TrimSpace(point.Name)
	if err := validatePoint(point); err != nil {
		return domain.Point{}, err
	}
	return s.repository.Update(ctx, point)
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	if s == nil || s.repository == nil {
		return errors.New("map point repository is unavailable")
	}
	if id == 0 {
		return fmt.Errorf("%w: id is required", ErrInvalidInput)
	}
	return s.repository.Delete(ctx, id)
}

func (s *Service) FindNearby(ctx context.Context, center domain.Coordinate, radius float64) ([]domain.Point, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("map point repository is unavailable")
	}
	if err := validateCoordinate(center); err != nil {
		return nil, err
	}
	if radius <= 0 {
		radius = 5000
	}
	if radius > 100000 {
		return nil, fmt.Errorf("%w: radius cannot exceed 100000 meters", ErrInvalidInput)
	}
	return s.repository.FindNearby(ctx, center, radius, 200)
}

func (s *Service) FindWithinPolygon(ctx context.Context, coordinates []domain.Coordinate) ([]domain.Point, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("map point repository is unavailable")
	}
	if len(coordinates) < 3 {
		return nil, fmt.Errorf("%w: polygon requires at least three coordinates", ErrInvalidInput)
	}
	if len(coordinates) > 1000 {
		return nil, fmt.Errorf("%w: polygon cannot exceed 1000 coordinates", ErrInvalidInput)
	}
	for _, coordinate := range coordinates {
		if err := validateCoordinate(coordinate); err != nil {
			return nil, err
		}
	}
	return s.repository.FindWithinPolygon(ctx, coordinates, 200)
}
