package usecase

import (
	"context"

	_err "common/domain/err"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

// IPOIUsecase defines the interface for POI usecase operations
type IPOIUsecase interface {
	CreatePOI(ctx context.Context, req *dto.CreatePOIRequestDTO) (*dto.POIDTO, *_err.ErrorDTO)
	UpdatePOI(ctx context.Context, id uint64, req *dto.UpdatePOIRequestDTO) (*dto.POIDTO, *_err.ErrorDTO)
	DeletePOI(ctx context.Context, id uint64) *_err.ErrorDTO
	GetPOIByID(ctx context.Context, id uint64) (*dto.POIDTO, *_err.ErrorDTO)
	GetPOIByCode(ctx context.Context, code string) (*dto.POIDTO, *_err.ErrorDTO)
	GetPOIList(ctx context.Context, req *dto.ListPOIsRequestDTO) (*dto.ListPOIsResponseDTO, *_err.ErrorDTO)
	GetPOIsByCategory(ctx context.Context, categoryID uint64, limit int) ([]dto.POIDTO, *_err.ErrorDTO)
	GetNearbyPOIs(ctx context.Context, lat, lng, radius float64, limit int) ([]dto.POIDTO, *_err.ErrorDTO)
	UpdatePOIRating(ctx context.Context, id uint64, ratingPoint float64, ratingCount int) *_err.ErrorDTO
}

// POIUsecase implements the POI usecase interface
type POIUsecase struct {
	poiRepo repo.IPOIRepo
}

// NewPOIUsecase creates a new POI usecase instance
// @bind: internal/usecase/poi_usecase.go
func NewPOIUsecase(poiRepo repo.IPOIRepo) IPOIUsecase {
	return &POIUsecase{
		poiRepo: poiRepo,
	}
}

// CreatePOI creates a new POI
func (u *POIUsecase) CreatePOI(ctx context.Context, req *dto.CreatePOIRequestDTO) (*dto.POIDTO, *_err.ErrorDTO) {
	// Validate category if provided
	if req.CategoryID != nil {
		if _, err := u.poiRepo.GetByID(ctx, *req.CategoryID); err != nil {
			return nil, &_err.ErrorDTO{
				Code:    400,
				Message: "Danh mục không tồn tại",
			}
		}
	}

	// Convert request to domain
	poi := mapper.CreatePOIRequestToDomain(req)
	if poi == nil {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Dữ liệu không hợp lệ",
		}
	}

	// Create POI
	err := u.poiRepo.Create(ctx, poi)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể tạo điểm quan tâm",
		}
	}

	// Convert to DTO
	poiDTO := mapper.POIToDTO(poi)
	return poiDTO, nil
}

// UpdatePOI updates an existing POI
func (u *POIUsecase) UpdatePOI(ctx context.Context, id uint64, req *dto.UpdatePOIRequestDTO) (*dto.POIDTO, *_err.ErrorDTO) {
	// Get existing POI
	existingPOI, err := u.poiRepo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Điểm quan tâm không tồn tại",
		}
	}

	// Validate category if provided
	if req.CategoryID != nil {
		if _, err := u.poiRepo.GetByID(ctx, *req.CategoryID); err != nil {
			return nil, &_err.ErrorDTO{
				Code:    400,
				Message: "Danh mục không tồn tại",
			}
		}
	}

	// Convert request to domain
	updatedPOI := mapper.UpdatePOIRequestToDomain(req, existingPOI)
	if updatedPOI == nil {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Dữ liệu không hợp lệ",
		}
	}

	// Update POI
	err = u.poiRepo.Update(ctx, id, updatedPOI)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật điểm quan tâm",
		}
	}

	// Convert to DTO
	poiDTO := mapper.POIToDTO(updatedPOI)
	return poiDTO, nil
}

// DeletePOI deletes a POI
func (u *POIUsecase) DeletePOI(ctx context.Context, id uint64) *_err.ErrorDTO {
	// Check if POI exists
	_, err := u.poiRepo.GetByID(ctx, id)
	if err != nil {
		return &_err.ErrorDTO{
			Code:    404,
			Message: "Điểm quan tâm không tồn tại",
		}
	}

	// Delete POI
	err = u.poiRepo.Delete(ctx, id)
	if err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa điểm quan tâm",
		}
	}

	return nil
}

// GetPOIByID retrieves a POI by ID
func (u *POIUsecase) GetPOIByID(ctx context.Context, id uint64) (*dto.POIDTO, *_err.ErrorDTO) {
	poi, err := u.poiRepo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Điểm quan tâm không tồn tại",
		}
	}

	poiDTO := mapper.POIToDTO(poi)
	return poiDTO, nil
}

// GetPOIByCode retrieves a POI by code
func (u *POIUsecase) GetPOIByCode(ctx context.Context, code string) (*dto.POIDTO, *_err.ErrorDTO) {
	poi, err := u.poiRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Điểm quan tâm không tồn tại",
		}
	}

	poiDTO := mapper.POIToDTO(poi)
	return poiDTO, nil
}

// GetPOIList retrieves a paginated list of POIs
func (u *POIUsecase) GetPOIList(ctx context.Context, req *dto.ListPOIsRequestDTO) (*dto.ListPOIsResponseDTO, *_err.ErrorDTO) {
	// Convert DTO to repo request
	repoReq := mapper.ListPOIsRequestToRepoRequest(req)

	// Get POIs from repository
	resp, err := u.poiRepo.GetListInterface(ctx, repoReq)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể lấy danh sách điểm quan tâm",
		}
	}

	// Convert to DTO
	poiDTOs := mapper.POIListToDTO(resp.Data)
	return &dto.ListPOIsResponseDTO{
		Data:  poiDTOs,
		Total: resp.Total,
		Code:  0,
	}, nil
}

// GetPOIsByCategory retrieves POIs by category
func (u *POIUsecase) GetPOIsByCategory(ctx context.Context, categoryID uint64, limit int) ([]dto.POIDTO, *_err.ErrorDTO) {
	pois, err := u.poiRepo.GetListByCategory(ctx, categoryID, limit)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể lấy danh sách điểm quan tâm theo danh mục",
		}
	}

	poiDTOs := mapper.POIListToDTO(pois)
	return poiDTOs, nil
}

// GetNearbyPOIs retrieves POIs within a certain radius
func (u *POIUsecase) GetNearbyPOIs(ctx context.Context, lat, lng, radius float64, limit int) ([]dto.POIDTO, *_err.ErrorDTO) {
	// Validate coordinates
	if lat < -90 || lat > 90 {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Vĩ độ không hợp lệ",
		}
	}
	if lng < -180 || lng > 180 {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Kinh độ không hợp lệ",
		}
	}

	pois, err := u.poiRepo.GetNearby(ctx, lat, lng, radius, limit)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể lấy danh sách điểm quan tâm gần đây",
		}
	}

	poiDTOs := mapper.POIListToDTO(pois)
	return poiDTOs, nil
}

// UpdatePOIRating updates the rating information for a POI
func (u *POIUsecase) UpdatePOIRating(ctx context.Context, id uint64, ratingPoint float64, ratingCount int) *_err.ErrorDTO {
	// Validate rating
	if ratingPoint < 0 || ratingPoint > 5 {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "Điểm đánh giá phải từ 0 đến 5",
		}
	}

	if ratingCount < 0 {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "Số lượng đánh giá không hợp lệ",
		}
	}

	// Check if POI exists
	_, err := u.poiRepo.GetByID(ctx, id)
	if err != nil {
		return &_err.ErrorDTO{
			Code:    404,
			Message: "Điểm quan tâm không tồn tại",
		}
	}

	// Update rating
	err = u.poiRepo.UpdateRating(ctx, id, ratingPoint, ratingCount)
	if err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật đánh giá",
		}
	}

	return nil
}
