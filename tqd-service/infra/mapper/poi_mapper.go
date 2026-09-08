package mapper

import (
	"encoding/json"
	"time"

	_entity "common/domain/entity"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

// PoiMapper handles mapping between domain, DTO and proto
type PoiMapper struct{}

func NewPoiMapper() *PoiMapper {
	return &PoiMapper{}
}

// ==================== DOMAIN -> RESPONSE ====================

// ToResponse converts domain POI to response DTO
func (m *PoiMapper) ToResponse(p *domain.POI) *dto.PoiResponse {
	if p == nil {
		return nil
	}
	resp := &dto.PoiResponse{}
	resp.FromDomain(p)
	return resp
}

// ToResponseList converts slice of domain POI to slice of response DTO
func (m *PoiMapper) ToResponseList(pois []domain.POI) []dto.PoiResponse {
	if pois == nil {
		return nil
	}
	result := make([]dto.PoiResponse, len(pois))
	for i, p := range pois {
		result[i].FromDomain(&p)
	}
	return result
}

// ==================== REQUEST -> DOMAIN ====================

// ToDomainFromCreate converts create request to domain
func (m *PoiMapper) ToDomainFromCreate(req *dto.CreatePoiRequest, userID uint64) (*domain.POI, error) {
	if req == nil {
		return nil, nil
	}

	// Generate code if not provided
	code := req.Code
	if code == "" {
		code = generatePoiCode() // Will be handled by repository
	}

	// Convert images to JSON
	imagesJSON := ""
	if len(req.Images) > 0 {
		imagesBytes, _ := json.Marshal(req.Images)
		imagesJSON = string(imagesBytes)
	}

	// Convert tags to JSON
	tagsJSON := ""
	if len(req.Tags) > 0 {
		tagsBytes, _ := json.Marshal(req.Tags)
		tagsJSON = string(tagsBytes)
	}

	// Convert amenity IDs to JSON
	amenitiesJSON := ""
	if len(req.AmenityIDs) > 0 {
		amenitiesBytes, _ := json.Marshal(req.AmenityIDs)
		amenitiesJSON = string(amenitiesBytes)
	}

	now := time.Now()
	return &domain.POI{
		BaseEntity: _entity.BaseEntity{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Code:          code,
		Name:          req.Name,
		Description:   req.Description,
		Address:       req.Address,
		Phone:         req.Phone,
		Email:         req.Email,
		Website:       req.Website,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		CategoryID:    req.CategoryID,
		IsActive:      req.IsActive,
		IsFeatured:    req.IsFeatured,
		CoverImage:    req.CoverImage,
		ImagesJSON:    imagesJSON,
		TagsJSON:      tagsJSON,
		AmenitiesJSON: amenitiesJSON,
		Rating:        0,
		ReviewCount:   0,
		IsVerified:    false,
		ViewCount:     0,
		LikeCount:     0,
		CreatedBy:     userID,
		UpdatedBy:     userID,
	}, nil
}

// ToDomainFromUpdate converts update request to domain
func (m *PoiMapper) ToDomainFromUpdate(req *dto.UpdatePoiRequest, existing *domain.POI, userID uint64) (*domain.POI, error) {
	if req == nil {
		return nil, nil
	}

	now := time.Now()
	updated := &domain.POI{
		BaseEntity: _entity.BaseEntity{
			ID:        existing.ID,
			CreatedAt: existing.CreatedAt,
			UpdatedAt: &now,
		},
		Code:          existing.Code,
		Name:          existing.Name,
		Description:   existing.Description,
		Address:       existing.Address,
		Phone:         existing.Phone,
		Email:         existing.Email,
		Website:       existing.Website,
		Latitude:      existing.Latitude,
		Longitude:     existing.Longitude,
		CategoryID:    existing.CategoryID,
		Rating:        existing.Rating,
		ReviewCount:   existing.ReviewCount,
		IsVerified:    existing.IsVerified,
		IsActive:      existing.IsActive,
		IsFeatured:    existing.IsFeatured,
		CoverImage:    existing.CoverImage,
		ImagesJSON:    existing.ImagesJSON,
		TagsJSON:      existing.TagsJSON,
		AmenitiesJSON: existing.AmenitiesJSON,
		ViewCount:     existing.ViewCount,
		LikeCount:     existing.LikeCount,
		CreatedBy:     existing.CreatedBy,
		UpdatedBy:     userID,
	}

	// Update only provided fields
	if req.Name != nil {
		updated.Name = *req.Name
	}
	if req.Code != nil {
		updated.Code = *req.Code
	}
	if req.Description != nil {
		updated.Description = *req.Description
	}
	if req.Address != nil {
		updated.Address = *req.Address
	}
	if req.Phone != nil {
		updated.Phone = *req.Phone
	}
	if req.Email != nil {
		updated.Email = *req.Email
	}
	if req.Website != nil {
		updated.Website = *req.Website
	}
	if req.Latitude != nil {
		updated.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		updated.Longitude = *req.Longitude
	}
	if req.CategoryID != nil {
		updated.CategoryID = *req.CategoryID
	}
	if req.IsActive != nil {
		updated.IsActive = *req.IsActive
	}
	if req.IsFeatured != nil {
		updated.IsFeatured = *req.IsFeatured
	}
	if req.CoverImage != nil {
		updated.CoverImage = *req.CoverImage
	}
	if req.Images != nil {
		imagesBytes, _ := json.Marshal(req.Images)
		updated.ImagesJSON = string(imagesBytes)
	}
	if req.Tags != nil {
		tagsBytes, _ := json.Marshal(req.Tags)
		updated.TagsJSON = string(tagsBytes)
	}
	if req.AmenityIDs != nil {
		amenitiesBytes, _ := json.Marshal(req.AmenityIDs)
		updated.AmenitiesJSON = string(amenitiesBytes)
	}

	return updated, nil
}

// ==================== DOMAIN -> PROTO ====================

// ToProto converts domain POI to proto
func (m *PoiMapper) ToProto(p *domain.POI) *tqdpb.Poi {
	if p == nil {
		return nil
	}

	// Parse images
	var images []string
	if p.ImagesJSON != "" {
		json.Unmarshal([]byte(p.ImagesJSON), &images)
	}

	// Parse tags
	var tags []string
	if p.TagsJSON != "" {
		json.Unmarshal([]byte(p.TagsJSON), &tags)
	}

	return &tqdpb.Poi{
		Id:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Description: p.Description,
		Address:     p.Address,
		Phone:       p.Phone,
		Email:       p.Email,
		Website:     p.Website,
		Latitude:    p.Latitude,
		Longitude:   p.Longitude,
		CategoryId:  p.CategoryID,
		Rating:      p.Rating,
		ReviewCount: p.ReviewCount,
		IsVerified:  p.IsVerified,
		IsActive:    p.IsActive,
		IsFeatured:  p.IsFeatured,
		CoverImage:  p.CoverImage,
		Images:      images,
		Tags:        tags,
		ViewCount:   p.ViewCount,
		LikeCount:   p.LikeCount,
		CreatedAt:   _utils.FormatTimeToString(p.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(p.UpdatedAt),
		CreatedBy:   p.CreatedBy,
		UpdatedBy:   p.UpdatedBy,
	}
}

// ToProtoFromResponse converts response DTO to proto
func (m *PoiMapper) ToProtoFromResponse(resp *dto.PoiResponse) *tqdpb.Poi {
	if resp == nil {
		return nil
	}

	return &tqdpb.Poi{
		Id:          resp.ID,
		Name:        resp.Name,
		Code:        resp.Code,
		Description: resp.Description,
		Address:     resp.Address,
		Phone:       resp.Phone,
		Email:       resp.Email,
		Website:     resp.Website,
		Latitude:    resp.Latitude,
		Longitude:   resp.Longitude,
		CategoryId:  resp.CategoryID,
		Rating:      resp.Rating,
		ReviewCount: resp.ReviewCount,
		IsVerified:  resp.IsVerified,
		IsActive:    resp.IsActive,
		IsFeatured:  resp.IsFeatured,
		CoverImage:  resp.CoverImage,
		Images:      resp.Images,
		Tags:        resp.Tags,
		ViewCount:   resp.ViewCount,
		LikeCount:   resp.LikeCount,
		CreatedAt:   _utils.FormatTimeToString(&resp.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(&resp.UpdatedAt),
		CreatedBy:   resp.CreatedBy,
		UpdatedBy:   resp.UpdatedBy,
	}
}

// ToProtoList converts slice of domain POI to slice of proto
func (m *PoiMapper) ToProtoList(pois []domain.POI) []*tqdpb.Poi {
	if pois == nil {
		return nil
	}
	result := make([]*tqdpb.Poi, len(pois))
	for i, p := range pois {
		result[i] = m.ToProto(&p)
	}
	return result
}

// ToProtoListFromResponses converts slice of response DTO to slice of proto
func (m *PoiMapper) ToProtoListFromResponses(responses []dto.PoiResponse) []*tqdpb.Poi {
	if responses == nil {
		return nil
	}
	result := make([]*tqdpb.Poi, len(responses))
	for i, resp := range responses {
		result[i] = m.ToProtoFromResponse(&resp)
	}
	return result
}

// ==================== PROTO -> REQUEST ====================

// FromProto converts proto to create request
func (m *PoiMapper) FromProto(proto *tqdpb.Poi) (*dto.CreatePoiRequest, error) {
	if proto == nil {
		return nil, nil
	}

	return &dto.CreatePoiRequest{
		Name:        proto.Name,
		Code:        proto.Code,
		Description: proto.Description,
		Address:     proto.Address,
		Phone:       proto.Phone,
		Email:       proto.Email,
		Website:     proto.Website,
		Latitude:    proto.Latitude,
		Longitude:   proto.Longitude,
		CategoryID:  proto.CategoryId,
		IsActive:    proto.IsActive,
		IsFeatured:  proto.IsFeatured,
		CoverImage:  proto.CoverImage,
		Images:      proto.Images,
		Tags:        proto.Tags,
	}, nil
}

func generatePoiCode() string {
	return "ABCD"
}
