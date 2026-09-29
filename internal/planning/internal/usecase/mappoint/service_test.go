package mappoint

import (
	"context"
	"errors"
	"testing"

	domain "tqd/internal/domain/mappoint"
)

type stubRepository struct{}

func (stubRepository) Create(context.Context, domain.Point) (domain.Point, error) {
	return domain.Point{ID: 1}, nil
}
func (stubRepository) Get(context.Context, uint64) (domain.Point, error) {
	return domain.Point{ID: 1}, nil
}
func (stubRepository) Update(context.Context, domain.Point) (domain.Point, error) {
	return domain.Point{ID: 1}, nil
}
func (stubRepository) Delete(context.Context, uint64) error { return nil }
func (stubRepository) FindNearby(context.Context, domain.Coordinate, float64, int) ([]domain.Point, error) {
	return nil, nil
}
func (stubRepository) FindWithinPolygon(context.Context, []domain.Coordinate, int) ([]domain.Point, error) {
	return nil, nil
}

func TestCreateRejectsInvalidPoint(t *testing.T) {
	service := NewService(stubRepository{})
	_, err := service.Create(context.Background(), domain.Point{Latitude: 91, Name: "invalid"})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestFindPolygonRequiresThreeCoordinates(t *testing.T) {
	service := NewService(stubRepository{})
	_, err := service.FindWithinPolygon(context.Background(), []domain.Coordinate{{}, {}})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}
