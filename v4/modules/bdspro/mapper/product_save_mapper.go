package mapper

import (
	"strconv"
	"time"

	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"
)

// ProductSaveToEntity maps ProductSaveRequest to a new Product row (parity product/new).
func ProductSaveToEntity(req *bdspropb.ProductSaveRequest, profileID uint64) *entity.Product {
	if req == nil {
		return nil
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
		SourceType:        enums.ESourceType(req.GetSourceType()),
		SourceStatus:      enums.ESourceStatus(sourceStatus),
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
		p.HouseInfo = houseInfoFromPB(hi, profileID)
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
	return p
}

func houseInfoFromPB(hi *bdspropb.HouseInfo, profileID uint64) *entity.HouseInfo {
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
	if hi.RoadWidth != nil {
		v := hi.GetRoadWidth()
		out.RoadWidth = &v
	}
	if hi.FrontWidth != nil {
		v := hi.GetFrontWidth()
		out.FrontWidth = &v
	}
	return out
}
