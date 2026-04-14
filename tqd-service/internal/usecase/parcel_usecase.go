package usecase

import (
	"context"
	"fmt"
	"math"

	"common/pkg/geoutils"
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
	FindByPolygon(ctx context.Context, req *dto.PolygonQueryRequest) (*dto.PolygonQueryResponse, error)
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

func (u *parcelUsecase) FindByPolygon(ctx context.Context, req *dto.PolygonQueryRequest) (*dto.PolygonQueryResponse, error) {
	// 1. Validate GeoJSON cơ bản
	if req.GeoJSON == "" {
		return nil, fmt.Errorf("geo_json is required")
	}
	if err := geoutils.ValidateGeoJSON(req.GeoJSON); err != nil {
		return nil, fmt.Errorf("invalid polygon: %v", err)
	}

	// 2. Tính diện tích gốc và bán kính
	originalArea, err := geoutils.CalcAreaKm2(req.GeoJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate area: %v", err)
	}
	originalRadius, err := geoutils.CalcRadiusKm(req.GeoJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate radius: %v", err)
	}

	// 3. Xác định giới hạn diện tích (lấy min của client gợi ý và server max)
	maxAllowedArea := config.MaxPolygonAreaKm2
	if req.MaxAreaKm2 > 0 && req.MaxAreaKm2 < maxAllowedArea {
		maxAllowedArea = req.MaxAreaKm2
	}

	// 4. Quyết định có cần scale hay không
	adjustedGeoJSON := req.GeoJSON
	adjustedArea := originalArea
	adjustmentReason := "NONE"

	// Kiểm tra diện tích
	if originalArea > maxAllowedArea {
		if !req.AllowAutoShrink {
			return nil, fmt.Errorf("polygon area %.2f km² exceeds limit %.2f km² (auto-shrink disabled)", originalArea, maxAllowedArea)
		}
		// Scale để diện tích = maxAllowedArea
		factor := math.Sqrt(maxAllowedArea / originalArea)
		scaled, err := geoutils.ScalePolygon(req.GeoJSON, factor)
		if err != nil {
			return nil, fmt.Errorf("auto-shrink failed: %v", err)
		}
		adjustedGeoJSON = scaled
		adjustedArea, _ = geoutils.CalcAreaKm2(adjustedGeoJSON)
		adjustmentReason = fmt.Sprintf("SHRUNK_AREA from %.2f to %.2f km²", originalArea, adjustedArea)
	} else if originalRadius > config.MaxPolygonRadiusKm {
		// Kiểm tra bán kính (có thể áp dụng scale nếu cho phép)
		if !req.AllowAutoShrink {
			return nil, fmt.Errorf("polygon radius %.2f km exceeds limit %.2f km (auto-shrink disabled)", originalRadius, config.MaxPolygonRadiusKm)
		}
		factor := config.MaxPolygonRadiusKm / originalRadius
		scaled, err := geoutils.ScalePolygon(req.GeoJSON, factor)
		if err != nil {
			return nil, fmt.Errorf("auto-shrink (radius) failed: %v", err)
		}
		adjustedGeoJSON = scaled
		adjustedArea, _ = geoutils.CalcAreaKm2(adjustedGeoJSON)
		adjustmentReason = fmt.Sprintf("SHRUNK_RADIUS from %.2f to %.2f km", originalRadius, config.MaxPolygonRadiusKm)
	}

	// 5. Chuẩn hóa phân trang
	page, pageSize := pageParams(req.Page, req.PageSize)
	limit := pageSize
	offset := (page - 1) * pageSize

	// 6. Gọi repository
	parcels, total, err := u.repo.FindByPolygon(ctx, adjustedGeoJSON, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	var warning string
	if total > 1000 {
		warning = fmt.Sprintf("Có %d parcel trong vùng, hiển thị có thể chậm. Nên thu nhỏ vùng để cải thiện tốc độ.", total)
	}
	if total > config.MaxParcelReturn {
		warning = fmt.Sprintf("Vùng chứa quá nhiều parcel (%d > %d). Chỉ trả về %d parcel đầu tiên.", total, config.MaxParcelReturn, limit)
	}

	return &dto.PolygonQueryResponse{
		Parcels:           parcels,
		Total:             total,
		Page:              page,
		PageSize:          pageSize,
		AdjustedGeoJSON:   adjustedGeoJSON,
		AdjustmentReason:  adjustmentReason,
		OriginalAreaKm2:   originalArea,
		AdjustedAreaKm2:   adjustedArea,
		MaxAllowedAreaKm2: maxAllowedArea,
		Warning:           warning,
	}, nil
}

func pageParams(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = config.DefaultPageSize
	}
	if pageSize > config.MaxPageSize {
		pageSize = config.MaxPageSize
	}
	return page, pageSize
}
