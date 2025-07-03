package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	pb_bdspro "pb/types/bdspro"
)

type ProductMapper struct {
	PostMapper *PostMapper
}

func NewProductMapper(postMapper *PostMapper) *ProductMapper {
	return &ProductMapper{
		PostMapper: postMapper,
	}
}

func (m *ProductMapper) MapProductPb(product *domain.Product) *pb_bdspro.ProductResponse {
	result := &pb_bdspro.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Area:        product.Area,
		Address:     product.Address,
		SaleStatus:  uint32(product.SaleStatus),
		RentStatus:  uint32(product.RentStatus),
	}

	if product.District != nil {
		result.DistrictName = product.District.Name
	}
	if product.Province != nil {
		result.ProvinceName = product.Province.Name
	}
	if product.Ward != nil {
		result.WardName = product.Ward.Name
	}

	if len(product.MediaList) > 0 {
		result.MediaList = make([]*pb_bdspro.MediaItem, len(product.MediaList))
		for i, media := range product.MediaList {
			result.MediaList[i] = &pb_bdspro.MediaItem{
				Id:        media.ID,
				MediaUrl:  media.MediaURL,
				MediaType: media.MediaType,
			}
		}
	}

	if product.Price != nil {
		result.Price = &pb_bdspro.ProductPrice{
			// PriceOwner:     product.Price.PriceOwner,
			Id:             product.Price.ID,
			Currency:       product.Price.Currency,
			SalePrice:      product.Price.SalePrice,
			SaleCommission: product.Price.SaleCommission,
			RentPrice:      product.Price.RentPrice,
			RentCommission: product.Price.RentCommission,
			Deposite:       product.Price.Deposite,
		}

	}

	return result
}

func (m *ProductMapper) MapProductPbList(products []domain.Product) []*pb_bdspro.ProductResponse {
	productList := make([]*pb_bdspro.ProductResponse, len(products))
	for i, product := range products {
		productList[i] = m.MapProductPb(&product)
	}
	return productList
}

func MapProductAttachmentPbList(products []domain.Product) []*pb_bdspro.ProductAttachment {
	productList := make([]*pb_bdspro.ProductAttachment, len(products))
	for i, product := range products {
		productList[i] = &pb_bdspro.ProductAttachment{
			Id:             product.ID,
			Name:           product.Name,
			Code:           product.Code,
			Area:           product.Area,
			AssetId:        product.AssetId,
			DepositeId:     product.DepositeID,
			RentStatus:     uint32(product.RentStatus),
			SaleStatus:     uint32(product.SaleStatus),
			SourceType:     int32(product.SourceType),
			SourceTypeName: enums.SourceTypeMap[product.SourceType],
		}
		if product.Price != nil {
			productList[i].Price = *product.Price.SalePrice
		}
		if product.MediaList != nil && len(product.MediaList) > 0 {
			productList[i].Image = product.MediaList[0].MediaURL
		}
	}
	return productList
}

func (m *ProductMapper) SaveProductToEntity(req *pb_bdspro.SaveChildRequest) *dto.ProductSaveRequest {
	result := &dto.ProductSaveRequest{
		ParentID:        &req.ParentId,
		Name:            req.Name,
		Code:            req.Code,
		Area:            req.Area,
		Description:     req.Description,
		Note:            req.Note,
		PropertyTypeId:  req.PropertyTypeId,
		DocTypeId:       req.DocTypeId,
		AmenityIds:      req.AmenityIds,
		ProjectId:       req.ProjectId,
		ProvinceID:      req.ProvinceId,
		DistrictID:      req.DistrictId,
		WardID:          req.WardId,
		TransactionType: int(req.TransactionType),
		GoogleMapLink:   req.GoogleMapLink,
		SourceType:      uint(req.SourceType),
		SaleStatus:      enums.EProductStatus(req.SaleStatus),
		RentStatus:      enums.EProductStatus(req.RentStatus),
		SaleVisibility:  enums.EVisibility(req.SaleVisibility),
		RentVisibility:  enums.EVisibility(req.RentVisibility),
		// OwnerType:      enums.EOwnerType(req.OwnerType),
	}

	// PostCreate * PostSaveWithProduct
	// AssetCreate * user_dto.AssetCreateWithProduct

	if req.PriceData != nil {
		result.PriceData = &domain.ProductPrice{
			PriceOwner:     req.PriceData.PriceOwner,
			Currency:       req.PriceData.Currency,
			SalePrice:      req.PriceData.SalePrice,
			SaleCommission: req.PriceData.SaleCommission,
			RentPrice:      req.PriceData.RentPrice,
			RentCommission: req.PriceData.RentCommission,
			Deposite:       req.PriceData.Deposite,
		}
	}

	if req.MediaItems != nil {
		result.MediaItems = make([]dto.MediaItem, len(req.MediaItems))
		for i, media := range req.MediaItems {
			result.MediaItems[i] = dto.MediaItem{
				ID:        media.Id,
				MediaURL:  media.MediaUrl,
				MediaType: media.MediaType,
				IsMain:    media.IsMain,
				Order:     int(media.Order),
			}
		}
	}

	if req.HouseInfo != nil {
		result.HouseInfo = &dto.HouseInfoUpdate{
			NumBedroom:  int(req.HouseInfo.NumBedroom),
			NumBathroom: int(req.HouseInfo.NumBathroom),
			NumFloor:    int(req.HouseInfo.NumFloor),
			NumFront:    int(req.HouseInfo.NumFront),
			NumCarPark:  int(req.HouseInfo.NumCarPark),
			Furniture:   req.HouseInfo.Furniture,
			Orientation: req.HouseInfo.Orientation,
		}
	}

	if req.ProductPrivate != nil {
		result.PrivateData = &domain.ProductPrivate{
			PrivateDocs:   req.ProductPrivate.PrivateDocs,
			ImportPrice:   &req.ProductPrivate.ImportPrice,
			OperatingCost: &req.ProductPrivate.OperatingCost,
			InternalNote:  req.ProductPrivate.InternalNote,
			TargetProfit:  &req.ProductPrivate.TargetProfit,
		}
	}

	// todo: map nôt asset và post
	// if req.Post != nil {
	// 	result.PostCreate = &dto.PostSaveWithProduct{
	// 		ProductID:  req.Post.ProductId,
	// 		Title:      req.Post.Title,
	// 		Content:    req.Post.Content,
	// 		Price:      req.Post.PostPrice,
	// 		Visibility: enums.EVisibility(req.Post.Visibility),
	// 	}
	// }

	// if req.Asset != nil {
	// 	result.AssetCreate = &user_dto.AssetCreateWithProduct{
	// 		PurchasePrice: req.Asset.PurchasePrice,
	// 		PurchaseDate:  req.Asset.PurchaseDate,
	// 		LegalStatus:   req.Asset.LegalStatus,
	// 		RentStatus:    req.Asset.RentStatus,
	// 		Description:   req.Asset.Description,
	// 		LegalItems:    req.Asset.LegalItems,
	// 	}
	// }

	return result
}

func (m *ProductMapper) MapDevideChildRequest(req *pb_bdspro.DevideChildRequest) *dto.DevideChildRequest {
	childs := make([]dto.ProductSaveRequest, len(req.Childs))
	for i, child := range req.Childs {
		r := m.SaveProductToEntity(child)
		if r != nil {
			childs[i] = *r
		}
	}
	return &dto.DevideChildRequest{
		ParentID: &req.ParentId,
		Childs:   childs,
	}
}
