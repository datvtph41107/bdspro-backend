package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_dto "common/domain/dto"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type BdsDomainMapper struct{}

func NewBdsDomainMapper() *BdsDomainMapper {
	return &BdsDomainMapper{}
}

// PbToBdsDomainListRequest converts protobuf request to DTO
func (m *BdsDomainMapper) PbToBdsDomainListRequest(req *bdspropb.BdsDomainListRequest) *dto.BdsDomainListRequest {
	result := &dto.BdsDomainListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
	}

	if req.Sort != "" {
		result.Sort = req.Sort
	}

	if req.Text != nil {
		text := *req.Text
		result.Text = &text
	}

	if req.ProductId != nil {
		result.ProductID = req.ProductId
	}

	if req.AssetId != nil {
		result.AssetID = req.AssetId
	}

	if req.ProvinceId != nil {
		result.ProvinceID = req.ProvinceId
	}

	if req.WardId != nil {
		result.WardID = req.WardId
	}

	if req.DistrictId != nil {
		result.DistrictID = req.DistrictId
	}

	return result
}

// BdsDomainToProto converts domain BDSDomain to protobuf BdsDomain
func (m *BdsDomainMapper) BdsDomainToProto(bds *domain.BDSDomain) *bdspropb.BdsDomain {
	if bds == nil {
		return nil
	}

	result := &bdspropb.BdsDomain{
		Id:           bds.ID,
		Address:      bds.Address,
		Street:       bds.Street,
		StreetNumber: bds.StreetNumber,
		Note:         bds.Note,
	}

	// Optional fields
	if bds.ProductID != nil {
		result.ProductId = bds.ProductID
	}
	if bds.AssetID != nil {
		result.AssetId = bds.AssetID
	}
	if bds.Area != nil {
		result.Area = bds.Area
	}
	if bds.AreaUse != nil {
		result.AreaUse = bds.AreaUse
	}
	if bds.AreaBuild != nil {
		result.AreaBuild = bds.AreaBuild
	}
	if bds.AreaLand != nil {
		result.AreaLand = bds.AreaLand
	}
	if bds.AreaFloor != nil {
		result.AreaFloor = bds.AreaFloor
	}
	if bds.AreaGreen != nil {
		result.AreaGreen = bds.AreaGreen
	}
	if bds.AreaParking != nil {
		result.AreaParking = bds.AreaParking
	}
	if bds.AreaBalcony != nil {
		result.AreaBalcony = bds.AreaBalcony
	}
	if bds.AreaTerrace != nil {
		result.AreaTerrace = bds.AreaTerrace
	}
	if bds.Length != nil {
		result.Length = bds.Length
	}
	if bds.Width != nil {
		result.Width = bds.Width
	}
	if bds.Height != nil {
		result.Height = bds.Height
	}
	if bds.FrontWidth != nil {
		result.FrontWidth = bds.FrontWidth
	}
	if bds.BackWidth != nil {
		result.BackWidth = bds.BackWidth
	}
	if bds.LeftWidth != nil {
		result.LeftWidth = bds.LeftWidth
	}
	if bds.RightWidth != nil {
		result.RightWidth = bds.RightWidth
	}
	if bds.Latitude != nil {
		result.Latitude = bds.Latitude
	}
	if bds.Longitude != nil {
		result.Longitude = bds.Longitude
	}
	if bds.WardID != nil {
		result.WardId = bds.WardID
	}
	if bds.DistrictID != nil {
		result.DistrictId = bds.DistrictID
	}
	if bds.ProvinceID != nil {
		result.ProvinceId = bds.ProvinceID
	}
	if bds.RoadWidth != nil {
		result.RoadWidth = bds.RoadWidth
	}
	if bds.RoadType != nil {
		result.RoadType = *bds.RoadType
	}
	if bds.RoadAccess != nil {
		result.RoadAccess = *bds.RoadAccess
	}
	if bds.DistanceToRoad != nil {
		result.DistanceToRoad = bds.DistanceToRoad
	}
	if bds.Orientation != nil {
		result.Orientation = *bds.Orientation
	}
	if bds.MainDirection != nil {
		result.MainDirection = *bds.MainDirection
	}
	if bds.CornerLot != nil {
		result.CornerLot = bds.CornerLot
	}
	if bds.AlleyAccess != nil {
		result.AlleyAccess = bds.AlleyAccess
	}
	if bds.AlleyWidth != nil {
		result.AlleyWidth = bds.AlleyWidth
	}
	if bds.Shape != nil {
		result.Shape = *bds.Shape
	}
	if bds.Topography != nil {
		result.Topography = *bds.Topography
	}
	if bds.Elevation != nil {
		result.Elevation = bds.Elevation
	}
	if bds.ProfileID != nil {
		result.ProfileId = bds.ProfileID
	}

	// Timestamps
	if !bds.CreatedAt.IsZero() {
		result.CreatedAt = _utils.FormatTimeToString(bds.CreatedAt)
	}
	if !bds.UpdatedAt.IsZero() {
		result.UpdatedAt = _utils.FormatTimeToString(bds.UpdatedAt)
	}

	return result
}

// BdsDomainListToProto converts list of domain BDSDomain to protobuf BdsDomain list
func (m *BdsDomainMapper) BdsDomainListToProto(bdsList []domain.BDSDomain) []*bdspropb.BdsDomain {
	result := make([]*bdspropb.BdsDomain, 0, len(bdsList))
	for _, bds := range bdsList {
		result = append(result, m.BdsDomainToProto(&bds))
	}
	return result
}
