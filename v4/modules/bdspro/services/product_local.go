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

const productMeMaxRecords uint32 = 50

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
	if req.Size > 0 && req.Size < size {
		size = req.Size
	}

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
		Id:                p.ID,
		Name:              p.Name,
		Code:              p.Code,
		Area:              p.Area,
		Description:       p.Description,
		GoogleMapLink:     p.GoogleMapLink,
		OwnerId:           &ownerID,
		OwnerType:         p.OwnerOf,
		SaleStatus:        p.SaleStatus,
		RentStatus:        p.RentStatus,
		Visibility:        p.Visibility,
		SaleVisibility:    p.SaleVisibility,
		RentVisibility:    p.RentVisibility,
		TransactionType:   p.TransactionType,
		SourceType:        p.SourceType,
		Archived:          p.Archived,
		Note:              p.Note,
		CertificateHouseNote: p.DocTypeNote,
		SourceContactNote: p.SourceContactNote,
		SourceStatus: &sharepb.EnumItem{
			Value: p.SourceStatus,
			Name:  enums.SourceStatusNames[enums.ESourceStatus(p.SourceStatus)],
		},
	}
	if p.Address != "" {
		out.Address = &sharepb.AddressV3Proto{Detail: p.Address}
	}
	if p.PropertyID != nil {
		out.PropertyId = p.PropertyID
	}
	if p.PropertyTypeID != nil {
		out.PropertyTypeId = p.PropertyTypeID
	}
	if p.DocTypeID != nil {
		out.DocTypeId = p.DocTypeID
	}
	if p.ProjectID != nil {
		out.ProjectId = p.ProjectID
	}
	if p.ProvinceID != nil {
		out.ProvinceId = p.ProvinceID
	}
	if p.WardID != nil {
		out.WardId = p.WardID
	}
	if p.LastPriceID != nil {
		out.LastPriceId = p.LastPriceID
	}
	if p.SourceContactID != nil {
		out.SourceContactId = p.SourceContactID
	}
	if p.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(p.CreatedAt.UnixMilli(), 10)
	}
	if p.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(p.UpdatedAt.UnixMilli(), 10)
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
	return out
}
