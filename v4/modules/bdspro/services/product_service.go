package services

import (
	"context"
	"strconv"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"
	"userv4/modules/bdspro/repository"

	"gorm.io/gorm"
)

// ProductService owns product list/detail/create/delete/timestamps flows for modules/bdspro.
type ProductService interface {
	ProductMe(ctx context.Context, req dto.SyncListRequest) (*bdspropb.SearchResponse, error)
	// GetDetail returns product payload and sync timestamp (ms) for HTTP envelope root.
	GetDetail(ctx context.Context, req dto.SyncListRequest) (*bdspropb.ProductResponse, int64, error)
	ProductNew(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.Response, error)
	DeleteProduct(ctx context.Context, id uint64) error
	GetTimestamps(ctx context.Context, id uint64) (*bdspropb.ProductTimestampsResponse, error)
}

const productMeMaxRecords uint32 = 5

// LocalProductService owns product flows against the v4 Postgres DSN + SyncUtil.
type LocalProductService struct {
	sync *_utils.SyncUtil
	repo repository.ProductRepository
}

func NewLocalProductService(sync *_utils.SyncUtil, repo repository.ProductRepository) *LocalProductService {
	return &LocalProductService{sync: sync, repo: repo}
}

func (s *LocalProductService) GetTimestamps(ctx context.Context, id uint64) (*bdspropb.ProductTimestampsResponse, error) {
	if id == 0 {
		return nil, _errors.BadRequestException("product id is required")
	}
	if s.sync == nil {
		return &bdspropb.ProductTimestampsResponse{Timestamps: map[string]int64{}}, nil
	}

	redisKeys := []string{
		s.sync.GetKey(ctx, _utils.SyncKeyProductDetail, id),
		s.sync.GetKey(ctx, _utils.SyncKeyProductPriceHist, id),
		s.sync.GetKey(ctx, _utils.SyncKeyProductDistrs, id),
		s.sync.GetKey(ctx, _utils.SyncKeyProductPosts, id),
		s.sync.GetKey(ctx, _utils.SyncKeyProductDeals, id),
		s.sync.GetKey(ctx, _utils.SyncKeyProductNotes, id),
		s.sync.GetKey(ctx, _utils.SyncKeyProductAppointments, id),
		s.sync.GetKey(ctx, _utils.SyncKeyProductContacts, id),
		s.sync.GetKey(ctx, _utils.SyncKeyPropertyHistories, id),
	}
	names := []string{
		"detail", "price_history", "distributions", "posts", "deals",
		"notes", "appointments", "contacts", "histories",
	}

	vals, err := s.sync.MGet(ctx, redisKeys)
	if err != nil {
		return nil, _errors.InternalServerException(err.Error())
	}

	timestamps := make(map[string]int64, len(names))
	for i, v := range vals {
		timestamps[names[i]] = parseSyncTimestamp(v)
	}
	return &bdspropb.ProductTimestampsResponse{Timestamps: timestamps}, nil
}

func (s *LocalProductService) ProductMe(ctx context.Context, req dto.SyncListRequest) (*bdspropb.SearchResponse, error) {
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyProductMe, req.ID)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Page == 0 {
			return &bdspropb.SearchResponse{Timestamp: s.sync.GetTimestamp(ctx, key)}, nil
		}
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	size := productMeMaxRecords

	list, total, err := s.repo.ListMine(ctx, profileID, req.Page, size, req.Archived)
	if err != nil {
		return nil, err
	}

	products := make([]*bdspropb.ProductResponse, 0, len(list))
	for _, p := range list {
		products = append(products, mapProductToPB(p))
	}

	now := time.Now().UnixMilli()
	if s.sync != nil && req.Page == 0 {
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductMe, 0), now)
	}

	totalPages := int32(0)
	if size > 0 {
		totalPages = int32((total + int64(size) - 1) / int64(size))
	}
	return &bdspropb.SearchResponse{
		Data:          products,
		TotalElements: total,
		Page:          int32(req.Page),
		Size:          int32(size),
		TotalPages:    totalPages,
		Timestamp:     now,
	}, nil
}

func (s *LocalProductService) GetDetail(ctx context.Context, req dto.SyncListRequest) (*bdspropb.ProductResponse, int64, error) {
	if req.ID == 0 {
		return nil, 0, _errors.BadRequestException("id is required")
	}
	var (
		key string
		ts  int64
	)
	if s.sync != nil {
		key = s.sync.GetKey(ctx, _utils.SyncKeyProductDetail, req.ID)
		ts = s.sync.GetTimestamp(ctx, key)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Timestamp > 0 {
			if ts <= 0 {
				ts = req.Timestamp
			}
			return &bdspropb.ProductResponse{Id: req.ID}, ts, nil
		}
	}
	p, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, 0, err
	}
	if p == nil {
		return nil, 0, _errors.NotFoundException("product not found")
	}
	if ts <= 0 {
		if s.sync != nil && key != "" {
			ts = s.sync.EnsureTimestamp(ctx, key)
		} else {
			ts = time.Now().UnixMilli()
		}
	}
	return mapProductToPB(p), ts, nil
}

func (s *LocalProductService) ProductNew(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.Response, error) {
	if req == nil {
		return nil, _errors.BadRequestException("body is required")
	}
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.UnauthorizedException("profileId không hợp lệ")
	}

	ownerOf := req.GetOwnerType()
	if ownerOf == 0 {
		ownerOf = enums.EOwnerOfMember
	}

	saleStatus := req.GetSaleStatus()
	if saleStatus == 0 {
		saleStatus = uint32(enums.EProductNotSold)
	}
	rentStatus := req.GetRentStatus()
	if rentStatus == 0 {
		rentStatus = uint32(enums.EProductNotRent)
	}
	sourceStatus := req.GetSourceStatus()
	if sourceStatus == 0 {
		sourceStatus = uint32(enums.SourceStatusNew)
	}
	saleVis := req.GetSaleVisibility()
	if saleVis == 0 {
		saleVis = uint32(enums.EVisiblePrivate)
	}
	rentVis := req.GetRentVisibility()
	if rentVis == 0 {
		rentVis = uint32(enums.EVisiblePrivate)
	}
	vis := req.GetVisibility()
	if vis == 0 {
		vis = uint32(enums.EVisiblePrivate)
	}

	p := &entity.Product{
		Name:              req.GetName(),
		Code:              req.GetCode(),
		Description:       req.GetDescription(),
		GoogleMapLink:     req.GetGoogleMapLink(),
		Area:              req.GetArea(),
		OwnerID:           profileID,
		OwnerOf:           ownerOf,
		OwnerType:         ownerOf,
		MiningMode:        uint32(enums.EMiningModeAll),
		Version:           1,
		SaleStatus:        saleStatus,
		RentStatus:        rentStatus,
		Visibility:        vis,
		SaleVisibility:    saleVis,
		RentVisibility:    rentVis,
		TransactionType:   req.GetTransactionType(),
		SourceType:        req.GetSourceType(),
		SourceStatus:      sourceStatus,
		Note:              req.GetNote(),
		DocTypeNote:       req.GetCertificateHouseNote(),
		SourceContactNote: req.GetSourceContactNote(),
		CreatedBy:         &profileID,
	}
	if req.PropertyTypeId != nil {
		v := req.GetPropertyTypeId()
		p.PropertyTypeID = &v
	}
	if req.DocTypeId != nil {
		v := req.GetDocTypeId()
		p.DocTypeID = &v
	}
	if req.ProjectId != nil {
		v := req.GetProjectId()
		p.ProjectID = &v
	}
	if req.ProvinceId != nil {
		v := req.GetProvinceId()
		p.ProvinceID = &v
	}
	if req.WardId != nil {
		v := req.GetWardId()
		p.WardID = &v
	}
	if req.SourceContactId != nil {
		v := req.GetSourceContactId()
		p.SourceContactID = &v
	}
	if st := req.GetSourceTime(); st != "" {
		if t := _utils.ParseStringToTime(st); t != nil {
			p.SourceTime = t
		} else if ms, err := strconv.ParseInt(st, 10, 64); err == nil && ms > 0 {
			tt := time.UnixMilli(ms)
			p.SourceTime = &tt
		}
	}
	if items := req.GetMediaItems(); len(items) > 0 {
		p.MediaList = make([]entity.ProductMediaInput, 0, len(items))
		for _, m := range items {
			if m == nil || m.GetMediaUrl() == "" {
				continue
			}
			p.MediaList = append(p.MediaList, entity.ProductMediaInput{
				MediaURL:  m.GetMediaUrl(),
				MediaType: m.GetMediaType(),
				IsMain:    m.GetIsMain(),
				Order:     m.GetOrder(),
			})
		}
	}
	if hi := req.GetHouseInfo(); hi != nil {
		p.HouseInfo = mapHouseInfoFromPB(hi, profileID)
	}

	if price := req.GetPriceData(); price != nil {
		pp := &entity.ProductPrice{
			Currency:           price.GetCurrency(),
			SaleCommissionType: price.SaleCommissionType,
			RentCommissionType: price.RentCommissionType,
			RentPaymentCycle:   price.RentPaymentCycle,
			ChangedBy:          profileID,
			ChannelPrice:       price.GetChannelPrice(),
			ChangeNote:         price.GetChangeNote(),
		}
		if price.SalePrice != nil {
			v := price.GetSalePrice()
			pp.SalePrice = &v
		}
		if price.SaleCommission != nil {
			v := price.GetSaleCommission()
			pp.SaleCommission = &v
		}
		if price.RentPrice != nil {
			v := price.GetRentPrice()
			pp.RentPrice = &v
		}
		if price.RentCommission != nil {
			v := price.GetRentCommission()
			pp.RentCommission = &v
		}
		if price.Deposite != nil {
			v := price.GetDeposite()
			pp.Deposite = &v
		}
		p.Price = pp
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	if s.sync != nil {
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductMe, 0), now)
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyPropertyMe, 0), now)
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductDetail, p.ID), now)
	}

	return &bdspropb.Response{Id: p.ID, Message: "success"}, nil
}

func (s *LocalProductService) DeleteProduct(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.BadRequestException("id is required")
	}
	profileID := _utils.GetProfileIdWithContext(ctx)
	if err := s.repo.SoftDelete(ctx, id, profileID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return _errors.NotFoundException("product not found")
		}
		return err
	}
	now := time.Now().UnixMilli()
	if s.sync != nil {
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductMe, 0), now)
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyPropertyMe, 0), now)
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductDetail, id), now)
	}
	return nil
}

func mapProductToPB(p *entity.Product) *bdspropb.ProductResponse {
	if p == nil {
		return nil
	}
	ownerID := p.OwnerID
	out := &bdspropb.ProductResponse{
		Id:                   p.ID,
		Name:                 p.Name,
		Code:                 p.Code,
		Area:                 p.Area,
		Description:          p.Description,
		GoogleMapLink:        p.GoogleMapLink,
		OwnerId:              &ownerID,
		OwnerType:            p.OwnerOf,
		SaleStatus:           p.SaleStatus,
		RentStatus:           p.RentStatus,
		Visibility:           p.Visibility,
		SaleVisibility:       p.SaleVisibility,
		RentVisibility:       p.RentVisibility,
		TransactionType:      p.TransactionType,
		SourceType:           p.SourceType,
		SourceTypeName:       enums.SourceTypeMap[enums.ESourceType(p.SourceType)],
		Archived:             p.Archived,
		Note:                 p.Note,
		CertificateHouseNote: p.DocTypeNote,
		SourceContactNote:    p.SourceContactNote,
		Priority:             p.Priority,
		ShareCount:           p.ShareCount,
		AssetCount:           p.AssetCount,
		PostCount:            p.PostCount,
		DealCount:            p.DealCount,
		ContactCount:         p.ContactCount,
		AppointmentCount:     p.AppointmentCount,
		SourceStatus: &sharepb.EnumItem{
			Value: p.SourceStatus,
			Name:  enums.SourceStatusNames[enums.ESourceStatus(p.SourceStatus)],
		},
	}

	provinceName := p.ProvinceName
	wardName := p.WardName
	if p.Property != nil && p.Property.Location != nil {
		if provinceName == "" && p.Property.Location.Province != nil {
			provinceName = p.Property.Location.Province.Name
		}
		if wardName == "" && p.Property.Location.Ward != nil {
			wardName = p.Property.Location.Ward.Name
		}
	}
	out.ProvinceName = provinceName
	out.WardName = wardName

	addr := &sharepb.AddressV3Proto{
		Detail:       p.Address,
		ProvinceName: provinceName,
		WardName:     wardName,
	}
	if p.ProvinceID != nil {
		addr.ProvinceId = p.ProvinceID
		out.ProvinceId = p.ProvinceID
	}
	if p.WardID != nil {
		addr.WardId = p.WardID
		out.WardId = p.WardID
	}
	if p.Address != "" || p.ProvinceID != nil || p.WardID != nil || provinceName != "" || wardName != "" {
		out.Address = addr
	}
	if p.PropertyID != nil {
		out.PropertyId = p.PropertyID
	}
	if p.PropertyTypeID != nil {
		out.PropertyTypeId = p.PropertyTypeID
	}
	if p.Property != nil && p.Property.Info != nil && p.Property.Info.PropertyType != nil {
		pt := p.Property.Info.PropertyType
		out.PropertyType = &bdspropb.ItemResponse{Id: pt.ID, Name: pt.Name}
		name := pt.Name
		out.PropertyTypeName = &name
	}
	if p.DocTypeID != nil {
		out.DocTypeId = p.DocTypeID
	}
	if p.ProjectID != nil {
		out.ProjectId = p.ProjectID
	}
	if p.LastPriceID != nil {
		out.LastPriceId = p.LastPriceID
	}
	if p.SourceContactID != nil {
		out.SourceContactId = p.SourceContactID
	}
	if p.SourceTime != nil {
		out.SourceTime = _utils.FormatTimeToString(p.SourceTime)
	}
	if p.CreatedBy != nil {
		out.CreatedBy = p.CreatedBy
	}
	if p.CreatedAt != nil {
		out.CreatedAt = _utils.FormatTimeToString(p.CreatedAt)
	}
	if p.UpdatedAt != nil {
		out.UpdatedAt = _utils.FormatTimeToString(p.UpdatedAt)
	}
	if p.Price != nil {
		pp := &bdspropb.ProductPrice{
			Id:                 p.Price.ID,
			Currency:           p.Price.Currency,
			SalePrice:          p.Price.SalePrice,
			SaleCommission:     p.Price.SaleCommission,
			SaleCommissionType: p.Price.SaleCommissionType,
			RentPrice:          p.Price.RentPrice,
			RentCommission:     p.Price.RentCommission,
			RentCommissionType: p.Price.RentCommissionType,
			RentPaymentCycle:   p.Price.RentPaymentCycle,
			Deposite:           p.Price.Deposite,
			ChannelPrice:       p.Price.ChannelPrice,
			ChangeNote:         p.Price.ChangeNote,
		}
		if p.Price.ProductID != nil {
			pp.ProductId = *p.Price.ProductID
		}
		out.Price = pp
	}
	if p.HouseInfo != nil {
		out.HouseInfo = mapHouseInfoToPB(p.HouseInfo)
	}
	if p.Property != nil {
		out.Property = mapPropertyLineageToProductPB(p.Property)
		if p.Property.BuildingInfo != nil {
			out.BuildingInfo = mapBuildingInfoToPB(p.Property.BuildingInfo)
		}
		if len(p.Property.MediaList) > 0 {
			out.MediaList = make([]*bdspropb.MediaItem, 0, len(p.Property.MediaList))
			for _, m := range p.Property.MediaList {
				out.MediaList = append(out.MediaList, &bdspropb.MediaItem{
					Id:        m.ID,
					MediaUrl:  m.MediaURL,
					MediaType: m.MediaType,
				})
			}
			first := p.Property.MediaList[0]
			out.Avatar = &bdspropb.MediaItem{
				Id:        first.ID,
				MediaUrl:  first.MediaURL,
				MediaType: first.MediaType,
			}
		}
	}
	return out
}

func mapPropertyLineageToProductPB(p *entity.PropertyLineage) *bdspropb.PropertyResponse {
	if p == nil {
		return nil
	}
	out := &bdspropb.PropertyResponse{Id: p.ID}
	if p.UpdatedAt != nil {
		out.UpdatedAt = _utils.FormatTimeToString(p.UpdatedAt)
	}
	if p.Info != nil {
		out.Title = p.Info.Title
		if p.Info.PropertyTypeID != nil {
			out.PropertyTypeId = p.Info.PropertyTypeID
		}
		if p.Info.ProjectID != nil {
			out.ProjectId = p.Info.ProjectID
		}
		if p.Info.PropertyType != nil {
			out.PropertyType = &sharepb.TagItem{Id: p.Info.PropertyType.ID, Name: p.Info.PropertyType.Name}
		}
	}
	if p.LandInfo != nil && p.LandInfo.AreaTotal != nil {
		out.AreaTotal = p.LandInfo.AreaTotal
	}
	if p.Statistic != nil {
		pc := p.Statistic.ProductCount
		ac := p.Statistic.AssetCount
		lc := p.Statistic.ListingCount
		out.ProductCount = &pc
		out.AssetCount = &ac
		out.PostCount = &lc
	}
	if p.BuildingInfo != nil {
		out.BuildingInfo = mapBuildingInfoToPB(p.BuildingInfo)
	}
	if p.Location != nil {
		out.AddressDetail = p.Location.AddressDetail
		out.MapUrl = p.Location.MapURL
		if p.Location.Latitude != nil {
			out.Latitude = p.Location.Latitude
		}
		if p.Location.Longitude != nil {
			out.Longitude = p.Location.Longitude
		}
		if p.Location.AddressDetail != "" || p.Location.ProvinceID != nil || p.Location.WardID != nil {
			addr := &sharepb.AddressV3Proto{
				Detail:     p.Location.AddressDetail,
				ProvinceId: p.Location.ProvinceID,
				WardId:     p.Location.WardID,
			}
			if p.Location.Province != nil {
				addr.ProvinceName = p.Location.Province.Name
			}
			if p.Location.Ward != nil {
				addr.WardName = p.Location.Ward.Name
			}
			out.Address = addr
		}
	}
	if len(p.MediaList) > 0 {
		first := p.MediaList[0]
		out.Avatar = &bdspropb.PropertyMediaItem{
			Url:       first.MediaURL,
			ThumbUrl:  first.ThumbURL,
			MediaType: first.MediaType,
		}
		out.MediaList = make([]*bdspropb.MediaItem, 0, len(p.MediaList))
		for _, m := range p.MediaList {
			out.MediaList = append(out.MediaList, &bdspropb.MediaItem{
				Id:        m.ID,
				MediaUrl:  m.MediaURL,
				MediaType: m.MediaType,
			})
		}
	}
	return out
}

func mapHouseInfoFromPB(hi *bdspropb.HouseInfo, profileID uint64) *entity.HouseInfo {
	if hi == nil {
		return nil
	}
	out := &entity.HouseInfo{
		CreatedBy:        &profileID,
		OrientationHouse: hi.GetOrientationHouse(),
		CertificateHouse: hi.GetCertificateHouse(),
	}
	if hi.NumBedroom != nil {
		v := hi.GetNumBedroom()
		out.NumBedroom = &v
	}
	if hi.NumBathroom != nil {
		v := hi.GetNumBathroom()
		out.NumBathroom = &v
	}
	if hi.NumFloor != nil {
		v := hi.GetNumFloor()
		out.NumFloor = &v
	}
	if hi.NumFront != nil {
		v := hi.GetNumFront()
		out.NumFront = &v
	}
	if hi.NumCarPark != nil {
		v := hi.GetNumCarPark()
		out.NumCarPark = &v
	}
	if hi.NumToilet != nil {
		v := hi.GetNumToilet()
		out.NumToilet = &v
	}
	if hi.Furniture != nil {
		v := hi.GetFurniture()
		out.Furniture = &v
	}
	if hi.Orientation != nil {
		v := hi.GetOrientation()
		out.Orientation = &v
	}
	if hi.RoadWidth != nil {
		v := hi.GetRoadWidth()
		out.RoadWidth = &v
	}
	if hi.FrontWidth != nil {
		v := hi.GetFrontWidth()
		out.FrontWidth = &v
	}
	if hi.BackWidth != nil {
		v := hi.GetBackWidth()
		out.BackWidth = &v
	}
	if hi.Width != nil {
		v := hi.GetWidth()
		out.Width = &v
	}
	if hi.Height != nil {
		v := hi.GetHeight()
		out.Height = &v
	}
	return out
}

func mapHouseInfoToPB(hi *entity.HouseInfo) *bdspropb.HouseInfo {
	if hi == nil {
		return nil
	}
	out := &bdspropb.HouseInfo{
		NumBedroom:       hi.NumBedroom,
		NumBathroom:      hi.NumBathroom,
		NumFloor:         hi.NumFloor,
		NumFront:         hi.NumFront,
		NumCarPark:       hi.NumCarPark,
		NumToilet:        hi.NumToilet,
		Furniture:        hi.Furniture,
		Orientation:      hi.Orientation,
		RoadWidth:        hi.RoadWidth,
		FrontWidth:       hi.FrontWidth,
		BackWidth:        hi.BackWidth,
		Width:            hi.Width,
		Height:           hi.Height,
	}
	if hi.OrientationHouse > 0 {
		v := hi.OrientationHouse
		out.OrientationHouse = &v
		name := enums.EHouseOrientNames[enums.EHouseOrient(hi.OrientationHouse)]
		out.OrientationName = &name
	}
	if hi.CertificateHouse > 0 {
		v := hi.CertificateHouse
		out.CertificateHouse = &v
		name := enums.EHouseCertificateNames[enums.EHouseCertificate(hi.CertificateHouse)]
		out.CertificateName = &name
	}
	return out
}

func mapBuildingInfoToPB(bi *entity.PropertyBuildingInfo) *bdspropb.PropertyBuildingInfo {
	if bi == nil {
		return nil
	}
	return &bdspropb.PropertyBuildingInfo{
		BuildingType:     bi.BuildingType,
		AreaActual:       bi.AreaActual,
		AreaFloor:        bi.AreaFloor,
		AreaConstruction: bi.AreaConstruction,
		Note:             bi.Note,
		Floors:           bi.Floors,
		RoomNumber:       bi.RoomNumber,
		Bedrooms:         bi.Bedrooms,
		Bathrooms:        bi.Bathrooms,
		BuildStatus:      bi.BuildStatus,
		Direction:        bi.Direction,
		BalconyDirection: bi.BalconyDirection,
	}
}
