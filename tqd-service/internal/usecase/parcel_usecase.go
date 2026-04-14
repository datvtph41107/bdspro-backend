package usecase

import (
	"context"
	"fmt"
	"math"

	"tqd/internal/config"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

// ParcelUsecase defines business logic for POI
type ParcelUsecase interface {
	// CRUD operations
	CreateBatch(ctx context.Context, req *dto.ParcelList) (*dto.ParcelList, error)
	FindByLocation(ctx context.Context, lat, lng float64) (*domain.Parcel, error)
	FindParcelsByPolygon(ctx context.Context, points []dto.Point, intersect bool, maxAreaKm2 float64, page, pageSize int) ([]domain.Parcel, int64, error)
}

// parcelUsecase implements ParcelUsecase
type parcelUsecase struct {
	repo repo.IParcelRepo
}

// NewParcelUsecase creates new POI usecase
func NewParcelUsecase(
	repo repo.IParcelRepo,
) ParcelUsecase {
	return &parcelUsecase{
		repo: repo,
	}
}

// Create implements ParcelUsecase.Create
func (u *parcelUsecase) CreateBatch(ctx context.Context, req *dto.ParcelList) (*dto.ParcelList, error) {
	// Get user ID
	// userID := _utils.GetProfileIdWithContext(ctx)
	// if userID == 0 {
	// 	return nil, fmt.Errorf("unauthorized")
	// }
	for _, parcel := range req.Data {
		if parcel.Lat < -90 || parcel.Lat > 90 {
			return nil, fmt.Errorf("latitude must be between -90 and 90")
		}
		if parcel.Lng < -180 || parcel.Lng > 180 {
			return nil, fmt.Errorf("longitude must be between -180 and 180")
		}
	}

	// Create in repository
	if err := u.repo.CreateBatch(ctx, req.Data); err != nil {
		return nil, fmt.Errorf("failed to create POIs: %v", err)
	}

	return req, nil
}

// GetByID implements ParcelUsecase.GetByID
func (u *parcelUsecase) FindByLocation(ctx context.Context, lat, lng float64) (*domain.Parcel, error) {
	parcel, err := u.repo.FindByLocation(ctx, lat, lng)
	if err != nil {
		return nil, fmt.Errorf("failed to get POI: %v", err)
	}
	if parcel == nil {
		return nil, fmt.Errorf("Parcel not found")
	}

	return parcel, nil
}

func (u *parcelUsecase) FindParcelsByPolygon(
	ctx context.Context,
	points []dto.Point,
	intersect bool,
	maxAreaKm2 float64,
	page, pageSize int,
) ([]domain.Parcel, int64, error) {
	if len(points) < config.MinPolygonVertices {
		return nil, 0, fmt.Errorf("polygon must have at least %d points, got %d", config.MinPolygonVertices, len(points))
	}
	if len(points) > config.MaxPolygonVertices {
		return nil, 0, fmt.Errorf("polygon has %d vertices, maximum allowed is %d", len(points), config.MaxPolygonVertices)
	}

	if err := validatePolygon(points); err != nil {
		return nil, 0, fmt.Errorf("invalid polygon: %v", err)
	}

	area := calculatePolygonAreaKm2(points)

	allowedArea := config.MaxPolygonAreaKm2
	if maxAreaKm2 > 0 && maxAreaKm2 < allowedArea {
		allowedArea = maxAreaKm2
	}
	if area > allowedArea {
		return nil, 0, fmt.Errorf("polygon area %.2f km² exceeds limit %.2f km²", area, allowedArea)
	}

	limit, offset := normalizePage(page, pageSize)

	var parcels []domain.Parcel
	var total int64
	var err error

	if intersect {
		parcels, total, err = u.repo.FindIntersectPolygon(ctx, points, limit, offset)
	} else {
		parcels, total, err = u.repo.FindWithinPolygon(ctx, points, 9999, offset)
	}
	if err != nil {
		return nil, 0, err
	}

	return parcels, total, nil
}

func calculatePolygonAreaKm2(points []dto.Point) float64 {
	if len(points) < 3 {
		return 0
	}
	area := 0.0
	n := len(points)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += points[i].Lng * points[j].Lat
		area -= points[j].Lng * points[i].Lat
	}
	area = math.Abs(area) / 2.0
	// 1 độ ~ 111 km
	kmPerDegree := 111.0
	return area * kmPerDegree * kmPerDegree
}

func validatePolygon(points []dto.Point) error {
	n := len(points)
	for i := 0; i < n; i++ {
		for j := i + 2; j < n; j++ {
			// Bỏ qua cạnh kề nhau
			if i == 0 && j == n-1 {
				continue
			}
			if i+1 == j {
				continue
			}
			p1 := points[i]
			p2 := points[(i+1)%n]
			p3 := points[j]
			p4 := points[(j+1)%n]
			if segmentsIntersect(p1, p2, p3, p4) {
				return fmt.Errorf("self-intersection between edges (%d,%d) and (%d,%d)", i, i+1, j, j+1)
			}
		}
	}
	return nil
}

func segmentsIntersect(a, b, c, d dto.Point) bool {
	orient := func(p, q, r dto.Point) int {
		val := (q.Lat-p.Lat)*(r.Lng-q.Lng) - (q.Lng-p.Lng)*(r.Lat-q.Lat)
		if val == 0 {
			return 0
		}
		if val > 0 {
			return 1
		}
		return 2
	}
	o1 := orient(a, b, c)
	o2 := orient(a, b, d)
	o3 := orient(c, d, a)
	o4 := orient(c, d, b)

	if o1 != o2 && o3 != o4 {
		return true
	}
	return false
}

func normalizePage(page, pageSize int) (limit, offset int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = config.DefaultPageSize
	}
	if pageSize > config.MaxPageSize {
		pageSize = config.MaxPageSize
	}
	limit = pageSize
	offset = (page - 1) * pageSize
	return
}
