package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	property_repo "bdspro/internal/repo/property"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_errors "common/errors"
	_provider "common/provider"
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type PropertyUsecase struct {
	PropertyRepo             repo.PropertyRepo
	PropertyUserRepo         repo.PropertyUserRepo
	PropertyInfoRepo         repo.PropertyInfoRepo
	PropertyStatisticRepo    repo.PropertyStatisticRepo
	PropertyIdentifyRepo     repo.PropertyIdentifyRepo
	PropertyLocationRepo     repo.PropertyLocationRepo
	PropertyExternalRefRepo  repo.PropertyExternalRefRepo
	PropertyEdvidenceRepo    repo.PropertyEdvidenceRepo
	PropertyRelationRepo     repo.PropertyRelationRepo
	PropertyProductRepo      repo.PropertyProductRepo
	PropertyLandInfoRepo     repo.PropertyLandInfoRepo
	PropertyBuildingInfoRepo repo.PropertyBuildingInfoRepo
	PropertyMediaRepo        repo.PropertyMediaRepo
	PropertyAmenityRepo      repo.PropertyAmenityRepo
	AssetLegalRepo           repo.AssetLegalRepo
	ProductRepo              repo.ProductRepo
	ProductUserRepo          repo.ProductUserRepo
	ProductPriceRepo         repo.ProductPriceRepo
	ProductMediaRepo         repo.ProductMediaRepo
	AssetRepo                repo.AssetRepo
	Transaction              provider.TransactionProvider
	NotifyProvider           provider.NotificationProvider
	HubProvider              provider.HubProvider
	SyncProvider             *_provider.SyncProvider
	LineageRepo              property_repo.PropertyLineageRepository
	InfoRepo                 property_repo.PropertyInfoRepository
	LocationRepo             property_repo.PropertyLocationRepository
	LandInfoRepo             property_repo.PropertyLandInfoRepository
	BuildingInfoRepo         property_repo.PropertyBuildingInfoRepository
	EvidenceRepo             property_repo.PropertyEvidenceRepository
	MediaRepo                property_repo.PropertyMediaRepository
	AmenityRepo              property_repo.PropertyAmenityRepository
	ExternalRefRepo          property_repo.PropertyExternalRefRepository
}

func NewPropertyUsecase(
	propertyRepo repo.PropertyRepo,
	propertyUserRepo repo.PropertyUserRepo,
	propertyInfoRepo repo.PropertyInfoRepo,
	propertyStatisticRepo repo.PropertyStatisticRepo,
	propertyIdentifyRepo repo.PropertyIdentifyRepo,
	propertyLocationRepo repo.PropertyLocationRepo,
	propertyExternalRefRepo repo.PropertyExternalRefRepo,
	propertyEdvidenceRepo repo.PropertyEdvidenceRepo,
	propertyRelationRepo repo.PropertyRelationRepo,
	propertyProductRepo repo.PropertyProductRepo,
	propertyLandInfoRepo repo.PropertyLandInfoRepo,
	propertyBuildingInfoRepo repo.PropertyBuildingInfoRepo,
	propertyMediaRepo repo.PropertyMediaRepo,
	propertyAmenityRepo repo.PropertyAmenityRepo,
	assetLegalRepo repo.AssetLegalRepo,
	productRepo repo.ProductRepo,
	productUserRepo repo.ProductUserRepo,
	productPriceRepo repo.ProductPriceRepo,
	productMediaRepo repo.ProductMediaRepo,
	assetRepo repo.AssetRepo,
	transaction provider.TransactionProvider,
	notifyProvider provider.NotificationProvider,
	hubProvider provider.HubProvider,
	syncProvider *_provider.SyncProvider,
	lineageRepo property_repo.PropertyLineageRepository,
	inforRepo property_repo.PropertyInfoRepository,
	locationRepo property_repo.PropertyLocationRepository,
	landInforRepo property_repo.PropertyLandInfoRepository,
	buildingInforRepo property_repo.PropertyBuildingInfoRepository,
	evidenceRepo property_repo.PropertyEvidenceRepository,
	mediaRepo property_repo.PropertyMediaRepository,
	amenityRepo property_repo.PropertyAmenityRepository,
	externalRefRepo property_repo.PropertyExternalRefRepository,
) *PropertyUsecase {
	return &PropertyUsecase{
		PropertyRepo:             propertyRepo,
		PropertyUserRepo:         propertyUserRepo,
		PropertyInfoRepo:         propertyInfoRepo,
		PropertyStatisticRepo:    propertyStatisticRepo,
		PropertyIdentifyRepo:     propertyIdentifyRepo,
		PropertyLocationRepo:     propertyLocationRepo,
		PropertyExternalRefRepo:  propertyExternalRefRepo,
		PropertyEdvidenceRepo:    propertyEdvidenceRepo,
		PropertyRelationRepo:     propertyRelationRepo,
		PropertyProductRepo:      propertyProductRepo,
		PropertyLandInfoRepo:     propertyLandInfoRepo,
		PropertyBuildingInfoRepo: propertyBuildingInfoRepo,
		PropertyMediaRepo:        propertyMediaRepo,
		PropertyAmenityRepo:      propertyAmenityRepo,
		AssetLegalRepo:           assetLegalRepo,
		ProductRepo:              productRepo,
		ProductUserRepo:          productUserRepo,
		ProductPriceRepo:         productPriceRepo,
		ProductMediaRepo:         productMediaRepo,
		AssetRepo:                assetRepo,
		Transaction:              transaction,
		NotifyProvider:           notifyProvider,
		HubProvider:              hubProvider,
		SyncProvider:             syncProvider,
		LineageRepo:              lineageRepo,
		InfoRepo:                 inforRepo,
		LocationRepo:             locationRepo,
		LandInfoRepo:             landInforRepo,
		BuildingInfoRepo:         buildingInforRepo,
		EvidenceRepo:             evidenceRepo,
		MediaRepo:                mediaRepo,
		ExternalRefRepo:          externalRefRepo,
	}
}

func (u *PropertyUsecase) GetPropertyAmenities(
	ctx context.Context,
	propertyID uint64,
	sourceType uint32,
) ([]*dto.PropertyAmenityDTO, error) {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, _errors.UnauthorizedException()
	}

	property, err := u.PropertyRepo.GetDetail(ctx, propertyID, sourceType)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, _errors.NotFoundException("Property not found")
	}

	amenities, err := u.PropertyAmenityRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	selectedIDs, err := u.PropertyAmenityRepo.GetAmenityIDsByProperty(ctx, propertyID)
	if err != nil {
		return nil, err
	}

	selectedMap := map[uint64]bool{}
	for _, id := range selectedIDs {
		selectedMap[id] = true
	}

	result := make([]*dto.PropertyAmenityDTO, 0, len(amenities))
	for _, a := range amenities {
		result = append(result, &dto.PropertyAmenityDTO{
			ID:      a.ID,
			Name:    a.Name,
			Active:  true,
			Checked: selectedMap[a.ID],
		})
	}

	return result, nil
}

func (u *PropertyUsecase) CloneProperty(
	ctx context.Context,
	propertyID, userID uint64,
) (*domain.PropertyLineage, error) {
	var result *domain.PropertyLineage
	err := u.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		prop, err := u.PropertyRepo.GetOneByID(ctx, propertyID)
		if err != nil {
			return err
		}
		// if prop.SourceType != enums.PropertySourceCanonical {
		// 	return _errors.InternalServerException("Invalid property is not system")
		// }
		newProp := u.cloneBase(prop, userID)
		if err := u.PropertyRepo.Create(ctx, newProp); err != nil {
			return err
		}
		if _, err := u.PropertyUserRepo.CreateOwner(ctx, newProp.ID, userID); err != nil {
			return err
		}
		result = newProp
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *PropertyUsecase) cloneBase(prop *domain.PropertyLineage, userID uint64) *domain.PropertyLineage {
	newProp := *prop
	newProp.ID = 0

	// newProp.SourceType = enums.PropertySourceBase
	// newProp.CreatedBy = userID

	// newProp.OriginProfileId = &prop.ID
	// newProp.CreatedAt = time.Now()
	// newProp.UpdatedAt = nil
	// newProp.DeletedAt = gorm.DeletedAt{}

	// newProp.Province = nil
	// newProp.District = nil
	// newProp.Ward = nil
	// newProp.Project = nil
	// // newProp.MediaList = nil
	// newProp.Amenities = nil
	return &newProp
}

func (u *PropertyUsecase) UpdatePropertyAmenities(
	ctx context.Context,
	req *dto.UpdatePropertyAmenitiesRequest,
) error {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return _errors.UnauthorizedException()
	}

	property, err := u.PropertyRepo.GetDetail(ctx, req.PropertyID, req.SourceType)
	if err != nil {
		return err
	}
	if property == nil {
		return _errors.NotFoundException("Property not found")
	}

	return u.PropertyAmenityRepo.UpdatePropertyAmenityIds(
		ctx,
		req.PropertyID,
		req.AmenityIDs,
	)
}

func (u *PropertyUsecase) GetPropertyTags(
	ctx context.Context,
	propertyID uint64,
	sourceType uint32,
) ([]*dto.PropertyTagDTO, error) {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, _errors.UnauthorizedException()
	}

	property, err := u.PropertyRepo.GetDetail(ctx, propertyID, sourceType)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, _errors.NotFoundException("Property not found")
	}

	tagIDs, err := u.PropertyAmenityRepo.GetTagIDsByProperty(ctx, propertyID)
	if err != nil {
		return nil, err
	}

	if len(tagIDs) == 0 {
		return []*dto.PropertyTagDTO{}, nil
	}

	tags, err := u.PropertyAmenityRepo.FindByIDs(ctx, tagIDs)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.PropertyTagDTO, 0, len(tags))
	for _, tag := range tags {
		result = append(result, &dto.PropertyTagDTO{
			ID:   tag.ID,
			Name: tag.Name,
			Type: tag.Type,
			Icon: tag.Icon,
		})
	}

	return result, nil
}

func (u *PropertyUsecase) SearchTags(
	ctx context.Context,
	search *dto.TagSearchDTO,
) ([]*dto.TagDTO, int64, error) {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, 0, _errors.UnauthorizedException()
	}

	tags, total, err := u.PropertyAmenityRepo.Search(ctx, search)
	if err != nil {
		return nil, 0, err
	}

	items := make([]*dto.TagDTO, 0, len(tags))
	for _, tag := range tags {
		items = append(items, &dto.TagDTO{
			ID:          tag.ID,
			Name:        tag.Name,
			Type:        tag.Type,
			Description: tag.Description,
			Icon:        tag.Icon,
		})
	}

	return items, total, nil
}

func (u *PropertyUsecase) UpdatePropertyTags(
	ctx context.Context,
	req *dto.UpdatePropertyTagsRequest,
) error {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return _errors.UnauthorizedException()
	}

	property, err := u.PropertyRepo.GetDetail(ctx, req.PropertyID, req.SourceType)
	if err != nil {
		return err
	}
	if property == nil {
		return _errors.NotFoundException("Property not found")
	}

	// validate tag ids (nếu có)
	if len(req.TagIDs) > 0 {
		tags, err := u.PropertyAmenityRepo.FindByIDs(ctx, req.TagIDs)
		if err != nil {
			return err
		}
		if len(tags) != len(req.TagIDs) {
			return _errors.BadRequestException("Invalid tag ids")
		}
	}

	if err := u.PropertyAmenityRepo.ReplacePropertyTags(
		ctx,
		req.PropertyID,
		req.TagIDs,
	); err != nil {
		return err
	}

	return nil
}

// LinkProductToProperty link product với property của product cha (nếu có)
func (u *PropertyUsecase) LinkProductToProperty(ctx context.Context, productID uint64, parentProductID uint64) error {
	// Tìm property của product cha
	parentRelation, err := u.PropertyRelationRepo.GetByProductID(ctx, parentProductID)
	if err != nil {
		return err
	}

	// Nếu không có property relation cho parent, không làm gì
	if parentRelation == nil {
		return nil
	}

	// Tạo link product với property
	propertyProduct := &domain.PropertyProduct{
		PropertyID: parentRelation.PropertyID,
		ProductID:  productID,
	}

	return u.PropertyProductRepo.Create(ctx, propertyProduct)
}

// LinkAssetToProperty link asset với property của product (nếu có)
func (u *PropertyUsecase) LinkAssetToProperty(ctx context.Context, assetID uint64, productID uint64) error {
	// Tìm property của product
	productRelation, err := u.PropertyRelationRepo.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}

	// Nếu không có property relation cho product, không làm gì
	if productRelation == nil {
		return nil
	}

	// Tạo relation mới với type Asset hoặc update existing relation nếu chưa có asset relation
	if productRelation.RelationType == enums.ERelationTypeProduct {
		// Nếu đã có product relation, tạo thêm asset relation mới
		assetRelation := &domain.PropertyRelation{
			PropertyID:   productRelation.PropertyID,
			RelationID:   &assetID,
			RelationType: enums.ERelationTypeAsset,
		}
		return u.PropertyRelationRepo.Create(ctx, assetRelation)
	}
	// Nếu chưa có relation, tạo mới
	productRelation.RelationID = &assetID
	productRelation.RelationType = enums.ERelationTypeAsset
	return u.PropertyRelationRepo.Update(ctx, productRelation)
}

// GetList lấy danh sách Property với phân trang
func (u *PropertyUsecase) GetList(ctx context.Context, searchDTO *dto.PropertySearchDTO) ([]*domain.PropertyLineage, int64, error) {
	if searchDTO == nil {
		searchDTO = &dto.PropertySearchDTO{}
	}
	var userId uint64
	if searchDTO.Owned {
		userId = _utils.GetOriginIdFromContext(ctx)
	}
	list, total, err := u.PropertyRepo.Search(ctx, searchDTO, userId)
	if err != nil {
		return nil, 0, err
	}
	regionIDsMap := make(map[uint64]struct{})
	// for _, p := range list {
	// 	if p.RegionID != nil {
	// 		regionIDsMap[*p.RegionID] = struct{}{}
	// 	}
	// }

	var regionIDs []uint64
	for id := range regionIDsMap {
		regionIDs = append(regionIDs, id)
	}

	// regionsMap, err := u.HubProvider.GetRegionsByIdsMap(ctx, regionIDs)
	// if err != nil {
	// 	return nil, 0, err
	// }

	// for _, p := range list {
	// 	if p.RegionID == nil {
	// 		continue
	// 	}

	// 	region, ok := regionsMap[*p.RegionID]
	// 	if !ok || region == nil {
	// 		continue
	// 	}

	// 	p.RegionName = region.Name
	// 	p.RegionCode = region.Code
	// }

	return list, total, nil
}

// ListMe lấy danh sách Property của user (raw SQL), chỉ trả title, address, project, buildingInfo, avatar
func (u *PropertyUsecase) ListMe(ctx context.Context, page, size uint32) ([]*domain.PropertyLineage, int64, error) {
	originId := _utils.GetOriginIdFromContext(ctx)
	if originId == 0 {
		return nil, 0, _errors.UnauthorizedException("origin_id required")
	}
	return u.PropertyRepo.ListMe(ctx, originId, _dto.Pagable{Page: page, Size: size})
}

// Detail lấy chi tiết Property theo ID
func (u *PropertyUsecase) Detail(
	ctx context.Context,
	id uint64,
	sourceType uint32,
) (*domain.PropertyLineage, error) {
	if id == 0 {
		return nil, _errors.BadRequestException("id is required")
	}

	property, err := u.PropertyRepo.GetDetail2(ctx, id, sourceType)
	if err != nil {
		return nil, err
	}

	return property, nil
}

// GetRelation lấy danh sách PropertyRelation với đầy đủ thông tin Asset/Product/Post
func (u *PropertyUsecase) GetRelation(ctx context.Context, propertyID uint64) ([]domain.PropertyRelation, error) {
	// Lấy tất cả PropertyRelations với propertyID này
	relations, err := u.PropertyRelationRepo.GetAllByPropertyID(ctx, propertyID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.PropertyRelation, 0, len(relations))

	// Load thêm thông tin Asset/Product/Post cho mỗi relation
	for _, relation := range relations {
		item := domain.PropertyRelation{
			PropertyID:   relation.PropertyID,
			RelationID:   relation.RelationID,
			RelationType: relation.RelationType,
		}

		if relation.RelationID == nil {
			result = append(result, item)
			continue
		}

		switch relation.RelationType {
		case enums.ERelationTypeProduct:
			// Load Product nếu có
			if product, err := u.ProductRepo.GetProductByID(ctx, relation.RelationID); err == nil && product != nil {
				item.Product = product
			}
		case enums.ERelationTypeAsset:
			// Load Asset nếu có
			if asset, err := u.AssetRepo.GetByID(*relation.RelationID); err == nil && asset != nil {
				item.Asset = asset
			}
		}

		result = append(result, item)
	}

	return result, nil
}

// UserCreate tạo Property mới cùng với Product, Asset, Post và các thông tin liên quan
// Thứ tự tạo: BasicInfo (Identify) -> Location, LandInfo, BuildingInfo, ExternalRef, Edvidence, MediaList -> Lineage (cuối) -> PropertyUser
// Lưu ý: Product, Asset, Post sẽ được tạo ở handler layer trước khi gọi usecase này
func (u *PropertyUsecase) UserCreate(ctx context.Context, req *dto.CreatePropertyProductRequest) (*domain.PropertyLineage, error) {
	var (
		property    *domain.PropertyLineage
		productID   *uint64
		assetID     *uint64
		identifyID  *uint64
		locationID  *uint64
		infoID      *uint64
		landInfoID  *uint64
		buildingID  *uint64
		edvidenceID *uint64
		// externalRefID *uint64
	)

	originId := _utils.GetOriginIdFromContext(ctx)
	_ = _utils.GetProfileIdWithContext(ctx)

	err := u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {

		// 1. Tạo PropertyIdentify (BasicInfo) đầu tiên - bắt buộc để có ID cho các bảng con
		var info *domain.PropertyInfo
		if req.Info != nil {
			info = req.Info
		} else {
			info = &domain.PropertyInfo{
				SourceType: enums.EPropertySourceType(enums.SourceTypeOther),
				// Visibility:      enums.EVisiblePrivate,
				RecordStatus: enums.EPropertyStatusActive,
				LegalStatus:  enums.EHouseCertificateNone,
				// PrivacyLevel:    enums.EVisiblePublic,
				Note:            "Note",
				UnitCode:        "UnitCode",
				Identifier:      "Identifier",
				Level:           "Level",
				OriginProfileId: &originId,
			}
		}
		if err := u.PropertyInfoRepo.Save(txCtx, info); err != nil {
			return _errors.InternalServerException("Create PropertyInfo Error: ", err.Error())
		}

		// 2. Tạo Location
		if req.Location != nil {
			loc := req.Location
			loc.PropertyIdentifyID = identifyID
			if err := u.PropertyLocationRepo.Create(txCtx, loc); err != nil {
				return _errors.InternalServerException("Create PropertyLocation Error: ", err.Error())
			}
			locationID = &loc.ID
			infoID = &info.ID
		}

		// 3. Tạo LandInfo
		if req.LandInfo != nil {
			landInfo := req.LandInfo
			landInfo.PropertyIdentifyID = identifyID
			if err := u.PropertyLandInfoRepo.Create(txCtx, landInfo); err != nil {
				return _errors.InternalServerException("Create PropertyLandInfo Error: ", err.Error())
			}
			landInfoID = &landInfo.ID
		}

		// 4. Tạo BuildingInfo
		if req.BuildingInfo != nil {
			buildingInfo := req.BuildingInfo
			buildingInfo.PropertyIdentifyID = identifyID
			if err := u.PropertyBuildingInfoRepo.Create(txCtx, buildingInfo); err != nil {
				return _errors.InternalServerException("Create PropertyBuildingInfo Error: ", err.Error())
			}
			buildingID = &buildingInfo.ID
		}

		// 5. Tạo Edvidence
		if req.Edvidence != nil {
			edv := req.Edvidence
			edv.PropertyIdentifyID = identifyID
			if err := u.PropertyEdvidenceRepo.Create(txCtx, edv); err != nil {
				return _errors.InternalServerException("Create PropertyEdvidence Error: ", err.Error())
			}
			edvidenceID = &edv.ID
		}

		// 6. Tạo ExternalRef
		if req.ExternalRef != nil {
			extRef := req.ExternalRef
			extRef.PropertyIdentifyID = identifyID
			if err := u.PropertyExternalRefRepo.Create(txCtx, extRef); err != nil {
				return _errors.InternalServerException("Create PropertyExternalRef Error: ", err.Error())
			}
			// externalRefID = &extRef.ID
		}

		productCount := int64(0)
		assetCount := int64(0)
		listCount := int64(0)
		if req.Product != nil {
			productCount = 1
		}
		if req.Asset != nil {
			assetCount = 1
		}
		if req.Post != nil {
			listCount = 1
		}
		// 7. Tạo PropertyStatistic
		statistic := &domain.PropertyStatistic{
			ProductCount: productCount,
			AssetCount:   assetCount,
			ListingCount: listCount,
		}
		if err := u.PropertyStatisticRepo.Save(txCtx, statistic); err != nil {
			return _errors.InternalServerException("Create PropertyStatistic Error: ", err.Error())
		}
		statisticID := &statistic.ID

		// 8. Tạo PropertyLineage cuối cùng - chứa ID của các bảng trên
		lineageData := &domain.PropertyLineage{}
		property = &domain.PropertyLineage{
			PropertyIdentifyID: identifyID,
			PropertyInfoID:     infoID,
			LocationID:         locationID,
			LandInfoID:         landInfoID,
			BuildingInfoID:     buildingID,
			EdvidenceID:        edvidenceID,
			// ExternalRefID:      externalRefID,
			NationalID:         lineageData.NationalID,
			StatisticID:        statisticID,
			VerifiedNationalAt: lineageData.VerifiedNationalAt,
		}

		if err := u.PropertyRepo.Create(txCtx, property); err != nil {
			return _errors.InternalServerException("Create PropertyLineage Error: ", err.Error())
		}

		// 9. Tạo PropertyUser
		if _, err := u.PropertyUserRepo.CreateOwner(txCtx, property.ID, originId); err != nil {
			return _errors.InternalServerException("Create PropertyUser Error: ", err.Error())
		}

		// 7. Tạo MediaList
		if len(req.MediaList) > 0 {
			for i := range req.MediaList {
				req.MediaList[i].LineageID = &property.ID
			}
			if err := u.PropertyMediaRepo.CreateBatch(txCtx, req.MediaList); err != nil {
				return _errors.InternalServerException("Create PropertyMedia Error: ", err.Error())
			}
			avatarID := req.MediaList[0].ID

			info.AvatarID = &avatarID
			if err := u.PropertyInfoRepo.Save(txCtx, info); err != nil {
				return _errors.InternalServerException("Create PropertyInfo Error: ", err.Error())
			}
		}

		// 10. Tạo Product
		var err error
		if req.Product != nil {
			req.Product.PropertyID = &property.ID
			productID, err = u.CreateProductFromProperty(txCtx, req.Product)
			if err != nil {
				return err
			}
		}

		// 11. Tạo Asset
		if req.Asset != nil {
			assetID, err = u.CreateAssetFromProperty(txCtx, req.Asset)
			if err != nil {
				return err
			}
		}

		// 12. PropertyRelation
		if productID != nil {
			if err := u.PropertyRelationRepo.Create(txCtx, &domain.PropertyRelation{
				PropertyID:   property.ID,
				RelationID:   productID,
				RelationType: enums.ERelationTypeProduct,
			}); err != nil {
				return err
			}
		}
		if assetID != nil {
			if err := u.PropertyRelationRepo.Create(txCtx, &domain.PropertyRelation{
				PropertyID:   property.ID,
				RelationID:   assetID,
				RelationType: enums.ERelationTypeAsset,
			}); err != nil {
				return err
			}
		}

		// 13. Amenity
		if len(req.AmenityIds) > 0 {
			if err := u.UpdatePropertyAmenity(txCtx, property.ID, req.AmenityIds); err != nil {
				return err
			}
		}

		// 14. AssetLegal
		if req.LegalInfo != nil && assetID != nil {
			legalInfo := req.LegalInfo
			legalInfo.AssetID = *assetID
			if err := u.AssetLegalRepo.Create(txCtx, legalInfo); err != nil {
				return err
			}
		}

		go func() {
			cloneCtx := _utils.CloneContext(ctx)
			u.SyncProvider.PutTimestamp(cloneCtx, u.SyncProvider.GetKey(cloneCtx, _provider.SyncKeyPropertyMe, 0), property.UpdatedAt.UnixMilli())
			u.SyncProvider.PutTimestamp(cloneCtx, u.SyncProvider.GetKey(cloneCtx, _provider.SyncKeyProductMe, 0), property.UpdatedAt.UnixMilli())
		}()

		return nil
	})

	if err != nil {
		return nil, err
	}

	if property != nil && u.NotifyProvider != nil {
		u.NotifyProvider.RegistedEventProperty(ctx,
			&dto.CreatePropertyActivityDTO{
				SubjectID:   property.ID,
				SubjectType: _enum.SUBJECT_TYPE_PROPERTY,
				Action:      _enum.ACTION_PROPERTY_CREATE,
				ActorID:     _utils.GetProfileIdWithContext(ctx),
				Description: "Sản phẩm được tạo ở đây",
			})
	}

	return property, nil
}

func (u *PropertyUsecase) buildProperty(req *domain.PropertyLineage, userId uint64) (*domain.PropertyLineage, error) {
	// if req.SourceType == 0 {
	// 	return nil, errors.New("source_type is required")
	// }

	// if req.CreatedAt.IsZero() {
	// 	req.CreatedAt = time.Now()
	// }

	// if req.RecordStatus == 0 {
	// 	req.RecordStatus = enums.EPropertyStatusActive
	// }

	property := &domain.PropertyLineage{
		// Title:              req.Title,
		// PropertyTypeID:     req.PropertyTypeID,
		// AddressDetail:      req.AddressDetail,
		// Latitude:           req.Latitude,
		// Longitude:          req.Longitude,
		// MapURL:             req.MapURL,
		// ProjectID:          req.ProjectID,
		// Level:              req.Level,
		// UnitCode:           req.UnitCode,
		// Identifier:         req.Identifier,
		NationalID: req.NationalID,
		// NationalIDVerified: req.NationalIDVerified,
		// AreaTotal:          req.AreaTotal,
		// AreaLand:           req.AreaLand,
		// AreaResidential:    req.AreaResidential,
		// LegalNote:      req.LegalNote,
		// LegalStatus:    req.LegalStatus,
		// ProvinceID:     req.ProvinceID,
		// DistrictID:     req.DistrictID,
		// WardID:         req.WardID,
		// AvatarID:       req.AvatarID,
		BuildingInfoID: req.BuildingInfoID,
		// SourceType:     req.SourceType,
		// RecordStatus:   enums.EPropertyStatusActive,
		// Scope:          req.Scope,
		// CreatedBy:      userId,
	}

	return property, nil
}

// CreateProductFromProperty tạo Product nếu chưa có ID, trả về productID
func (u *PropertyUsecase) CreateProductFromProperty(txCtx context.Context, product *domain.Product) (*uint64, error) {
	if product == nil {
		return nil, nil
	}

	// Nếu Product đã có ID thì không cần tạo
	if product.ID != 0 {
		return &product.ID, nil
	}

	// Set OwnerID và OwnerOf nếu chưa có
	if product.OwnerID == 0 {
		ownerID := _utils.GetProfileIdWithContext(txCtx)
		product.OwnerID = ownerID
	}
	if product.OwnerOf == 0 {
		product.OwnerOf = enums.EOwnerOfMember
	}

	// Tạo Price nếu có
	if product.Price != nil {
		if err := u.ProductPriceRepo.Create(txCtx, product.Price); err != nil {
			return nil, err
		}
		product.LastPriceID = &product.Price.ID
	}

	// Tạo Product
	if err := u.ProductRepo.Create(txCtx, product); err != nil {
		return nil, err
	}

	// Update ProductID vào Price sau khi tạo Product
	if product.Price != nil {
		if err := u.ProductPriceRepo.UpdateProductID(txCtx, &product.ID, &product.Price.ID); err != nil {
			return nil, err
		}
	}

	// Tạo ProductUser để link Product với User
	// Chỉ tạo nếu OwnerOf là Member (không phải Organization)
	if product.OwnerOf == enums.EOwnerOfMember {
		profileID := _utils.GetProfileIdWithContext(txCtx)
		_, err := u.ProductUserRepo.CreateOwner(txCtx, product.ID, profileID)
		if err != nil {
			return nil, err
		}
	}

	return &product.ID, nil
}

// CreateAssetFromProperty tạo Asset nếu chưa có ID, trả về assetID
func (u *PropertyUsecase) CreateAssetFromProperty(txCtx context.Context, asset *domain.Asset) (*uint64, error) {
	if asset == nil {
		return nil, nil
	}

	// Nếu Asset đã có ID thì không cần tạo
	if asset.ID != 0 {
		return &asset.ID, nil
	}

	// Set OwnerID và OwnerOf nếu chưa có
	if asset.OwnerID == nil || *asset.OwnerID == 0 {
		ownerID := _utils.GetProfileIdWithContext(txCtx)
		asset.OwnerID = &ownerID
	}
	if asset.OwnerOf == 0 {
		asset.OwnerOf = enums.EOwnerOfMember
	}

	// Set default values nếu chưa có
	if asset.Archived == false && asset.ID == 0 {
		asset.Archived = false
	}

	// Tạo Asset
	if err := u.AssetRepo.Create(txCtx, asset); err != nil {
		return nil, err
	}

	return &asset.ID, nil
}

// CreateWithProduct tạo Property từ Product
// Nếu Product đã có AssetId thì không tạo Property mới (vì đã có PropertyRelation)
func (u *PropertyUsecase) CreateWithProduct(ctx context.Context, product *domain.Product) (*domain.PropertyLineage, error) {
	// Nếu Product đã có AssetId, kiểm tra xem đã có PropertyRelation chưa
	if product.AssetId != nil {
		relation, err := u.PropertyRelationRepo.GetByAssetID(ctx, *product.AssetId)
		if err != nil {
			return nil, err
		}
		// Nếu đã có PropertyRelation, không tạo Property mới
		if relation != nil {
			// Lấy Property hiện có
			property, err := u.PropertyRepo.GetByID(relation.PropertyID)
			if err != nil {
				return nil, err
			}
			return property, nil
		}
	}

	// Tạo Property từ thông tin Product
	property := &domain.PropertyLineage{
		// PropertyTypeID: product.PropertyTypeId,
		// Title:          product.Name,
		// ProjectID:      product.ProjectID,
		// // AreaTotal:      &product.Area,
		// AddressDetail: product.Address,
		// MapURL:        product.GoogleMapLink,
		// RecordStatus:  enums.EPropertyStatusActive,
	}

	// Tạo Property
	if err := u.PropertyRepo.Create(ctx, property); err != nil {
		return nil, err
	}

	// Tạo PropertyRelation với type Product
	relation := &domain.PropertyRelation{
		PropertyID:   property.ID,
		RelationID:   &product.ID,
		RelationType: enums.ERelationTypeProduct,
	}

	if err := u.PropertyRelationRepo.Create(ctx, relation); err != nil {
		return nil, err
	}

	// Nếu product có AssetId, tạo thêm relation với Asset
	if product.AssetId != nil {
		assetRelation := &domain.PropertyRelation{
			PropertyID:   property.ID,
			RelationID:   product.AssetId,
			RelationType: enums.ERelationTypeAsset,
		}
		if err := u.PropertyRelationRepo.Create(ctx, assetRelation); err != nil {
			return nil, err
		}
	}

	// Copy amenityIds từ Product sang Property nếu có
	amenityIds, err := u.ProductRepo.GetAmenityIDsByProductID(ctx, product.ID)
	if err == nil && len(amenityIds) > 0 {
		if err := u.UpdatePropertyAmenity(ctx, property.ID, amenityIds); err != nil {
			// Log error nhưng không fail transaction
			// return nil, err
		}
	}

	return property, nil
}

// UpdatePropertyAmenity cập nhật amenityIds cho Property
func (u *PropertyUsecase) UpdatePropertyAmenity(ctx context.Context, propertyID uint64, amenityIDs []uint64) error {
	if len(amenityIDs) > 0 {
		return u.PropertyAmenityRepo.UpdatePropertyAmenityIds(ctx, propertyID, amenityIDs)
	}
	return nil
}

// GetOrCreatePropertyFromProductID lấy PropertyID từ ProductID, nếu chưa có thì tạo mới
// Trả về PropertyID và error
func (u *PropertyUsecase) GetOrCreatePropertyFromProductID(ctx context.Context, productID uint64) (uint64, error) {
	// Kiểm tra xem đã có PropertyRelation chưa
	relation, err := u.PropertyRelationRepo.GetByProductID(ctx, productID)
	if err != nil {
		return 0, err
	}

	// Nếu đã có PropertyRelation, trả về PropertyID
	if relation != nil {
		return relation.PropertyID, nil
	}

	// Nếu chưa có, lấy Product để tạo Property
	product, err := u.ProductRepo.GetByIDContext(ctx, productID)
	if err != nil {
		return 0, err
	}
	if product == nil {
		return 0, errors.New("product not found")
	}

	// Tạo Property từ thông tin Product
	property := &domain.PropertyLineage{
		// PropertyTypeID: product.PropertyTypeId,
		// Title:          product.Name,
		// ProjectID:      product.ProjectID,
		// AreaTotal:      &product.Area,
		// AddressDetail: product.Address,
		// MapURL:        product.GoogleMapLink,
		// RecordStatus:  enums.EPropertyStatusActive,
	}

	// Tạo Property
	if err := u.PropertyRepo.Create(ctx, property); err != nil {
		return 0, err
	}

	// Tạo PropertyRelation với type Product
	relation = &domain.PropertyRelation{
		PropertyID:   property.ID,
		RelationID:   &productID,
		RelationType: enums.ERelationTypeProduct,
	}

	if err := u.PropertyRelationRepo.Create(ctx, relation); err != nil {
		return 0, err
	}

	// Nếu product có AssetId, tạo thêm relation với Asset
	if product.AssetId != nil {
		assetRelation := &domain.PropertyRelation{
			PropertyID:   property.ID,
			RelationID:   product.AssetId,
			RelationType: enums.ERelationTypeAsset,
		}
		if err := u.PropertyRelationRepo.Create(ctx, assetRelation); err != nil {
			return 0, err
		}
	}

	// Copy amenityIds từ Product sang Property nếu có
	amenityIds, err := u.ProductRepo.GetAmenityIDsByProductID(ctx, productID)
	if err == nil && len(amenityIds) > 0 {
		if err := u.UpdatePropertyAmenity(ctx, property.ID, amenityIds); err != nil {
			// Log error nhưng không fail transaction
		}
	}

	return property.ID, nil
}

// UpdatePropertyFromProduct cập nhật Property, BuildingInfo và LandInfo của Property từ ProductID
func (u *PropertyUsecase) UpdatePropertyFromProduct(ctx context.Context, productID uint64, houseInfo *dto.UpdateProductRequest) error {
	if houseInfo == nil {
		return nil
	}

	// Lấy Property từ Product
	propertyRelation, err := u.PropertyRelationRepo.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}

	// Nếu không có Property, không làm gì
	if propertyRelation == nil {
		return nil
	}

	// Update Property với các field: AreaLand, AreaTotal, ProjectID, PropertyTypeID
	property, err := u.PropertyRepo.GetDetail(ctx, propertyRelation.PropertyID, houseInfo.SourceType)
	if err != nil {
		return err
	}
	if property != nil {
		// Tạo property update chỉ với các field cần update
		updateProperty := &domain.PropertyLineage{
			// AreaLand:       &houseInfo.AreaLand,
			// AreaTotal:      &houseInfo.AreaTotal,
			// ProjectID:      houseInfo.ProjectID,
			// PropertyTypeID: houseInfo.PropertyTypeID,
		}

		// Update Property nếu có thay đổi
		if err := u.PropertyRepo.UpdateOption(ctx, propertyRelation.PropertyID, updateProperty); err != nil {
			return err
		}
	}

	// Convert từ UpdateProductRequest sang PropertyBuildingInfo
	buildingInfo := &domain.PropertyBuildingInfo{
		// PropertyID: propertyRelation.PropertyID,
	}

	// Map các field từ UpdateProductRequest sang PropertyBuildingInfo
	if houseInfo.Floors != nil {
		buildingInfo.Floors = houseInfo.Floors
	}
	if houseInfo.Bedrooms != nil {
		buildingInfo.Bedrooms = houseInfo.Bedrooms
	}
	if houseInfo.Bathrooms != nil {
		buildingInfo.Bathrooms = houseInfo.Bathrooms
	}
	if houseInfo.Direction != 0 {
		buildingInfo.Direction = houseInfo.Direction
	}
	if houseInfo.BalconyDirection != 0 {
		buildingInfo.BalconyDirection = houseInfo.BalconyDirection
	}
	// Lấy existing building info để update hoặc create
	existingBuildingInfo, err := u.PropertyBuildingInfoRepo.GetByPropertyID(ctx, propertyRelation.PropertyID)
	if err != nil {
		return err
	}

	if existingBuildingInfo != nil {
		// Update existing
		buildingInfo.ID = existingBuildingInfo.ID
		if err := u.PropertyBuildingInfoRepo.Update(ctx, buildingInfo); err != nil {
			return err
		}
	} else {
		// Create new
		if err := u.PropertyBuildingInfoRepo.Create(ctx, buildingInfo); err != nil {
			return err
		}
		// Link BuildingInfoID vào Property
		property, err := u.PropertyRepo.GetOneByID(ctx, propertyRelation.PropertyID)
		if err != nil {
			return err
		}
		if property != nil {
			property.BuildingInfoID = &buildingInfo.ID
			if err := u.PropertyRepo.Update(ctx, propertyRelation.PropertyID, property); err != nil {
				return err
			}
		}
	}

	// Update LandInfo nếu có RoadWidth (RoadWidth lưu vào LandInfo, không phải BuildingInfo)
	if houseInfo.RoadWidth != nil {
		landInfo := &domain.PropertyLandInfo{
			// PropertyID: propertyRelation.PropertyID,
			// RoadWidth:  houseInfo.RoadWidth,
		}

		// Lấy existing land info để update hoặc create
		existingLandInfo, err := u.PropertyLandInfoRepo.GetByPropertyID(ctx, propertyRelation.PropertyID)
		if err != nil {
			return err
		}

		if existingLandInfo != nil {
			// Update existing
			landInfo.ID = existingLandInfo.ID
			// Giữ nguyên các field khác nếu có
			// if existingLandInfo.Frontage != nil {
			// 	landInfo.Frontage = existingLandInfo.Frontage
			// }
			if existingLandInfo.Depth != nil {
				landInfo.Depth = existingLandInfo.Depth
			}
			if existingLandInfo.LandNote != "" {
				landInfo.LandNote = existingLandInfo.LandNote
			}
			if err := u.PropertyLandInfoRepo.Update(ctx, landInfo); err != nil {
				return err
			}
		} else {
			// Create new
			if err := u.PropertyLandInfoRepo.Create(ctx, landInfo); err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateMediaItems cập nhật danh sách media items cho Property
// Dựa trên pattern từ đoạn code 411-423: xóa media cũ, tạo media mới, set ảnh đầu tiên làm avatar
func (u *PropertyUsecase) UpdateMediaItems(ctx context.Context, propertyID uint64, mediaItems []dto.MediaItem) ([]uint64, error) {
	// Xóa các media items cũ của property
	if err := u.PropertyMediaRepo.DeleteByPropertyID(ctx, propertyID); err != nil {
		return nil, err
	}

	// Nếu danh sách media items rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(mediaItems) == 0 {
		return nil, nil
	}

	// Convert từ dto.MediaItem sang domain.PropertyMedia và gán PropertyID
	var propertyMediaList []domain.PropertyMedia
	for i, mediaItem := range mediaItems {
		propertyMedia := domain.PropertyMedia{
			// PropertyID: propertyID,
			MediaType: mediaItem.MediaType,
			MediaURL:  mediaItem.MediaURL,
			SortOrder: int32(i + 1), // Sử dụng Order từ MediaItem hoặc index + 1
			IsCover:   mediaItem.IsMain,
		}
		// Sử dụng Order từ MediaItem nếu có
		if mediaItem.Order > 0 {
			propertyMedia.SortOrder = int32(mediaItem.Order)
		}
		propertyMediaList = append(propertyMediaList, propertyMedia)
	}

	// Tạo batch media items mới
	if err := u.PropertyMediaRepo.CreateBatch(ctx, propertyMediaList); err != nil {
		return nil, err
	}

	// Lấy danh sách ID đã tạo
	var mediaIDs []uint64
	var avatarID *uint64
	for _, media := range propertyMediaList {
		mediaIDs = append(mediaIDs, media.ID)
		// Tìm ảnh có IsCover = true làm avatar, hoặc lấy ảnh đầu tiên
		if avatarID == nil {
			if media.IsCover {
				avatarID = &media.ID
			}
		}
	}

	// Nếu không có ảnh IsCover, lấy ảnh đầu tiên
	if avatarID == nil && len(mediaIDs) > 0 {
		avatarID = &mediaIDs[0]
	}

	// Update AvatarID của Property
	if avatarID != nil {
		property, err := u.PropertyRepo.GetOneByID(ctx, propertyID)
		if err != nil {
			return nil, err
		}
		if property != nil {
			// property.AvatarID = avatarID
			if err := u.PropertyRepo.Update(ctx, propertyID, property); err != nil {
				return nil, err
			}
		}
	} else {
		// Nếu không có media items, xóa AvatarID
		property, err := u.PropertyRepo.GetOneByID(ctx, propertyID)
		if err != nil {
			return nil, err
		}
		if property != nil {
			// property.AvatarID = nil
			if err := u.PropertyRepo.Update(ctx, propertyID, property); err != nil {
				return nil, err
			}
		}
	}

	return mediaIDs, nil
}

// CreateWithAsset tạo Property từ Asset
// Nếu Asset đã có ProductID thì không tạo Property mới (vì đã có PropertyRelation)
func (u *PropertyUsecase) CreateWithAsset(ctx context.Context, asset *domain.Asset) (*domain.PropertyLineage, error) {
	// Nếu Asset đã có ProductID, kiểm tra xem đã có PropertyRelation chưa
	if asset.ProductID != nil {
		relation, err := u.PropertyRelationRepo.GetByProductID(ctx, *asset.ProductID)
		if err != nil {
			return nil, err
		}
		// Nếu đã có PropertyRelation, không tạo Property mới
		if relation != nil {
			// Lấy Property hiện có
			property, err := u.PropertyRepo.GetOneByID(ctx, relation.PropertyID)
			if err != nil {
				return nil, err
			}
			return property, nil
		}
	}

	// Tạo Property từ thông tin Asset
	property := &domain.PropertyLineage{
		// PropertyTypeID: asset.PropertyTypeId,
		// Title:          asset.Name,
		// AreaTotal:      &asset.Area,
		// AddressDetail: asset.Address,
		// RecordStatus:  enums.EPropertyStatusActive, // Default status
	}

	// Map legal status từ Asset - giữ nguyên dạng số hoặc convert nếu cần
	// LegalStatus trong Property là string, có thể để trống hoặc convert sau
	// property.LegalStatus = "" // Có thể map sau nếu cần

	// Tạo Property
	if err := u.PropertyRepo.Create(ctx, property); err != nil {
		return nil, err
	}

	// Tạo PropertyRelation với type Asset
	assetRelation := &domain.PropertyRelation{
		PropertyID:   property.ID,
		RelationID:   &asset.ID,
		RelationType: enums.ERelationTypeAsset,
	}

	if err := u.PropertyRelationRepo.Create(ctx, assetRelation); err != nil {
		return nil, err
	}

	// Nếu asset có ProductID, tạo thêm relation với Product
	if asset.ProductID != nil {
		productRelation := &domain.PropertyRelation{
			PropertyID:   property.ID,
			RelationID:   asset.ProductID,
			RelationType: enums.ERelationTypeProduct,
		}
		if err := u.PropertyRelationRepo.Create(ctx, productRelation); err != nil {
			return nil, err
		}
	}

	return property, nil
}

// Update cập nhật Property và các thông tin liên quan (LandInfo, BuildingInfo)
func (u *PropertyUsecase) DeleteProperties(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	return u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := u.PropertyRepo.DeleteBatchProperties(txCtx, ids); err != nil {
			return err
		}

		return nil
	})
}

func (u *PropertyUsecase) ArchiveProperties(ctx context.Context, req *dto.ArchivePropertiesDTO) error {
	originId := _utils.GetOriginIdFromContext(ctx)

	if originId == 0 {
		return _errors.UnauthorizedException()
	}

	if len(req.IDs) == 0 {
		return _errors.BadRequestException("No property IDs provided")
	}

	for _, id := range req.IDs {
		property, err := u.PropertyRepo.GetOneByID(ctx, id)
		if err != nil {
			return _errors.InternalServerException("Failed to validate property ID: " + err.Error())
		}
		if property == nil {
			return _errors.NotFoundException("Property with ID " + fmt.Sprintf("%d", id) + " not found")
		}
	}
	return u.PropertyRepo.ArchiveBatchProperties(ctx, req.IDs, req.Archived)
}

// // UPDATE V2 PROPERTY
// UpdateProperty — entry point
//
// Flow:
//  1. Load existing PropertyLineage (bao gồm FK IDs của các sub-domain)
//  2. Với mỗi sub-domain có trong DTO: build update map → diff với DB → chỉ
//     write nếu thực sự có thay đổi

func (u *PropertyUsecase) UpdateProperty(
	ctx context.Context,
	c *dto.UpdatePropertyDTO,
) (*dto.UpdatePropertyResult, error) {
	var result *dto.UpdatePropertyResult

	err := u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		existing, err := u.LineageRepo.GetByID(txCtx, c.LineageID)
		if err != nil {
			return err
		}
		if existing == nil {
			return errors.New("property lineage not found")
		}

		result = &dto.UpdatePropertyResult{
			LineageID: existing.ID,
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		}

		if existing.PropertyIdentifyID != nil {
			// tạo đề xuất thay đổi với bản chính định danh
			proposeChangePropertyIdentify()
		}

		if c.Info != nil {
			if err := u.handleInfo(txCtx, c.LineageID, existing.PropertyInfoID, c.Info, c.FieldMask); err != nil {
				return err
			}
			result.Info, _ = u.InfoRepo.GetByLineageID(txCtx, c.LineageID)
		}

		if c.Location != nil {
			if err := u.handleLocation(txCtx, c.LineageID, existing.LocationID, c.Location, c.FieldMask); err != nil {
				return err
			}
			result.Location, _ = u.LocationRepo.GetByLineageID(txCtx, c.LineageID)
		}

		if c.LandInfo != nil {
			if err := u.handleLandInfo(txCtx, c.LineageID, existing.LandInfoID, c.LandInfo, c.FieldMask); err != nil {
				return err
			}
			result.LandInfo, _ = u.LandInfoRepo.GetByLineageID(txCtx, c.LineageID)
		}

		if c.BuildingInfo != nil {
			newBuildingID, err := u.handleBuildingInfo(txCtx, c.LineageID, existing.BuildingInfoID, c.BuildingInfo, c.FieldMask)
			if err != nil {
				return err
			}
			if newBuildingID != nil {
				if err := u.LineageRepo.UpdateFields(txCtx, c.LineageID, map[string]any{
					"building_info_id": *newBuildingID,
				}); err != nil {
					return err
				}
			}
			result.BuildingInfo, _ = u.BuildingInfoRepo.GetByLineageID(txCtx, c.LineageID)
		}

		if c.Evidence != nil {
			if err := u.handleEvidence(txCtx, c.LineageID, existing.EdvidenceID, c.Evidence, c.FieldMask); err != nil {
				return err
			}
			result.Evidence, _ = u.EvidenceRepo.GetByLineageID(txCtx, c.LineageID)
		}

		if c.Media != nil {
			if err := u.handleMedia(txCtx, c.LineageID, c.Media); err != nil {
				return err
			}
			result.Media, _ = u.MediaRepo.GetByLineageID(txCtx, c.LineageID)
		}

		if c.Amenities != nil {
			if err := u.AmenityRepo.ReplaceAmenities(txCtx, c.LineageID, c.Amenities.AmenityIDs); err != nil {
				return err
			}
			ids, _ := u.AmenityRepo.GetAmenityIDsByLineageID(txCtx, c.LineageID)
			items := make([]domain.AmenityItem, 0, len(ids))
			for _, id := range ids {
				items = append(items, domain.AmenityItem{ID: id})
			}
			result.Amenities = items
		}

		if c.ExternalRefs != nil {
			if err := u.handleExternalRefs(txCtx, c.LineageID, c.ExternalRefs); err != nil {
				return err
			}
			result.ExternalRefs, _ = u.ExternalRefRepo.GetByLineageID(txCtx, c.LineageID)
		}

		if c.LineageMeta != nil {
			if err := u.handleLineageMeta(txCtx, c.LineageID, existing, c.LineageMeta, c.FieldMask); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		profileId := _utils.GetProfileIdWithContext(cloneCtx)
		u.SyncProvider.PutTimestamp(cloneCtx, u.SyncProvider.GetKey(cloneCtx, _provider.SyncKeyPropertyMe, profileId), time.Now().UnixMilli())
		u.SyncProvider.PutTimestamp(cloneCtx, u.SyncProvider.GetKey(cloneCtx, _provider.SyncKeyPropertyDetail, result.LineageID), time.Now().UnixMilli())
	}()
	return result, nil
}

func proposeChangePropertyIdentify() {

}

// ─────────────────────────────────────────────────────────────────────────────
// Sub-domain handlers
// Pattern: build update map → diff với existing record → write nếu có thay đổi
// ─────────────────────────────────────────────────────────────────────────────

func (u *PropertyUsecase) handleInfo(
	ctx context.Context,
	lineageID uint64,
	infoID *uint64,
	c *dto.UpdateInfoDTO,
	mask []string,
) error {
	if infoID == nil {
		// Record chưa tồn tại → insert, không cần diff
		entity := mapInfoDTOToEntity(lineageID, c)
		if err := u.InfoRepo.Create(ctx, entity); err != nil {
			return err
		}
		return u.LineageRepo.UpdateFields(ctx, lineageID, map[string]any{
			"property_info_id": entity.ID,
		})
	}

	// Load existing để diff
	existing, err := u.InfoRepo.GetByID(ctx, *infoID)
	if err != nil {
		return err
	}

	updates := buildInfoUpdates(c, mask, existing)
	if len(updates) == 0 {
		return nil // Không có gì thay đổi
	}
	return u.InfoRepo.UpdateFields(ctx, *infoID, updates)
}

func (u *PropertyUsecase) handleLocation(
	ctx context.Context,
	lineageID uint64,
	locationID *uint64,
	c *dto.UpdateLocationDTO,
	mask []string,
) error {
	if locationID == nil {
		entity := mapLocationDTOToEntity(lineageID, c)
		if err := u.LocationRepo.Create(ctx, entity); err != nil {
			return err
		}
		return u.LineageRepo.UpdateFields(ctx, lineageID, map[string]any{
			"location_id": entity.ID,
		})
	}

	existing, err := u.LocationRepo.GetByID(ctx, *locationID)
	if err != nil {
		return err
	}

	updates := buildLocationUpdates(c, mask, existing)
	if len(updates) == 0 {
		return nil
	}
	return u.LocationRepo.UpdateFields(ctx, *locationID, updates)
}

func (u *PropertyUsecase) handleLandInfo(
	ctx context.Context,
	lineageID uint64,
	landInfoID *uint64,
	c *dto.UpdateLandInfoDTO,
	mask []string,
) error {
	if landInfoID == nil {
		entity := mapLandInfoDTOToEntity(lineageID, c)
		if err := u.LandInfoRepo.Create(ctx, entity); err != nil {
			return err
		}
		return u.LineageRepo.UpdateFields(ctx, lineageID, map[string]any{
			"land_info_id": entity.ID,
		})
	}

	existing, err := u.LandInfoRepo.GetByID(ctx, *landInfoID)
	if err != nil {
		return err
	}

	updates := buildLandInfoUpdates(c, mask, existing)
	if len(updates) == 0 {
		return nil
	}
	return u.LandInfoRepo.UpdateFields(ctx, *landInfoID, updates)
}

// handleBuildingInfo trả về newID nếu record được insert mới.
func (u *PropertyUsecase) handleBuildingInfo(
	ctx context.Context,
	lineageID uint64,
	buildingID *uint64,
	c *dto.UpdateBuildingDTO,
	mask []string,
) (*uint64, error) {
	if buildingID == nil {
		entity := mapBuildingDTOToEntity(lineageID, c)
		if err := u.BuildingInfoRepo.Create(ctx, entity); err != nil {
			return nil, err
		}
		return &entity.ID, nil
	}

	existing, err := u.BuildingInfoRepo.GetByID(ctx, *buildingID)
	if err != nil {
		return nil, err
	}

	updates := buildBuildingUpdates(c, mask, existing)
	if len(updates) == 0 {
		return nil, nil
	}
	return nil, u.BuildingInfoRepo.UpdateFields(ctx, *buildingID, updates)
}

func (u *PropertyUsecase) handleEvidence(
	ctx context.Context,
	lineageID uint64,
	evidenceID *uint64,
	c *dto.UpdateEvidenceDTO,
	mask []string,
) error {
	if evidenceID == nil {
		entity := &domain.PropertyEdvidence{
			Title:       derefStr(c.Title),
			Description: derefStr(c.Description),
			FileID:      c.FileID,
		}
		if err := u.EvidenceRepo.Create(ctx, entity); err != nil {
			return err
		}
		return u.LineageRepo.UpdateFields(ctx, lineageID, map[string]any{
			"edvidence_id": entity.ID,
		})
	}

	existing, err := u.EvidenceRepo.GetByID(ctx, *evidenceID)
	if err != nil {
		return err
	}

	updates := buildEvidenceUpdates(c, mask, existing)
	if len(updates) == 0 {
		return nil
	}
	return u.EvidenceRepo.UpdateFields(ctx, *evidenceID, updates)
}

func (u *PropertyUsecase) handleMedia(
	ctx context.Context,
	lineageID uint64,
	c *dto.UpdateMediaDTO,
) error {
	for _, item := range c.Items {
		if item.ID == 0 {
			entity := &domain.PropertyMedia{
				LineageID: &lineageID,
				MediaType: item.MediaType,
				MediaURL:  item.MediaURL,
				ThumbURL:  item.ThumbURL,
				SortOrder: item.SortOrder,
				IsCover:   item.IsCover,
			}
			if err := u.MediaRepo.Create(ctx, entity); err != nil {
				return err
			}
		} else {
			// Load existing để diff
			existing, err := u.MediaRepo.GetByID(ctx, item.ID)
			if err != nil {
				return err
			}
			updates := map[string]any{}
			if item.MediaURL != "" && item.MediaURL != existing.MediaURL {
				updates["media_url"] = item.MediaURL
			}
			if item.ThumbURL != "" && item.ThumbURL != existing.ThumbURL {
				updates["thumb_url"] = item.ThumbURL
			}
			if item.SortOrder != existing.SortOrder {
				updates["sort_order"] = item.SortOrder
			}
			if item.IsCover != existing.IsCover {
				updates["is_cover"] = item.IsCover
			}
			if len(updates) == 0 {
				continue
			}
			if err := u.MediaRepo.UpdateFields(ctx, item.ID, updates); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *PropertyUsecase) handleExternalRefs(
	ctx context.Context,
	lineageID uint64,
	c *dto.UpdateExternalRefDTO,
) error {
	for _, item := range c.Items {
		if item.ID == 0 {
			entity := &domain.PropertyExternalRef{
				ExternalRefID: item.ExternalRefID,
				SourceSystem:  enums.ERefSourceSys(item.SourceSystem),
				SourceCode:    item.SourceCode,
				Confidence:    item.Confidence,
				SyncStatus:    enums.ESyncStatus(item.SyncStatus),
				Note:          item.Note,
			}
			if err := u.ExternalRefRepo.Create(ctx, entity); err != nil {
				return err
			}
		} else {
			existing, err := u.ExternalRefRepo.GetByID(ctx, item.ID)
			if err != nil {
				return err
			}
			updates := map[string]any{}
			if item.SourceCode != existing.SourceCode {
				updates["source_code"] = item.SourceCode
			}
			if item.Confidence != existing.Confidence {
				updates["confidence"] = item.Confidence
			}
			if item.SyncStatus != int32(existing.SyncStatus) {
				updates["sync_status"] = item.SyncStatus
			}
			if item.Note != existing.Note {
				updates["note"] = item.Note
			}
			if len(updates) == 0 {
				continue
			}
			if err := u.ExternalRefRepo.UpdateFields(ctx, item.ID, updates); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *PropertyUsecase) handleLineageMeta(
	ctx context.Context,
	lineageID uint64,
	existing *domain.PropertyLineage,
	c *dto.UpdateLineageMetaDTO,
	mask []string,
) error {
	updates := map[string]any{}

	if inMask(mask, "lineage_meta.national_id") && c.NationalID != nil {
		if *c.NationalID != existing.NationalID {
			updates["national_id"] = *c.NationalID
		}
	}

	if inMask(mask, "lineage_meta.lineage_status") && c.LineageStatus != nil && *c.LineageStatus != 0 {
		updates["lineage_status"] = *c.LineageStatus
	}

	if inMask(mask, "lineage_meta.verified_national_at") && c.VerifiedNationalAt != nil {
		if *c.VerifiedNationalAt == "" {
			// Clear
			if existing.VerifiedNationalAt != nil {
				updates["verified_national_at"] = nil
			}
		} else {
			t, err := time.Parse(time.RFC3339, *c.VerifiedNationalAt)
			if err != nil {
				return errors.New("invalid verified_national_at format, expected RFC3339")
			}
			if existing.VerifiedNationalAt == nil || !t.Equal(*existing.VerifiedNationalAt) {
				updates["verified_national_at"] = t
			}
		}
	}

	if len(updates) == 0 {
		return nil
	}
	return u.LineageRepo.UpdateFields(ctx, lineageID, updates)
}

// ─────────────────────────────────────────────────────────────────────────────
// Field update builders — diff-aware
//
// Nhận DTO + mask + existing record.
// Chỉ thêm vào map nếu:
//   1. Path nằm trong mask (hoặc mask rỗng)
//   2. Giá trị thực sự khác với DB (diff check)
// ─────────────────────────────────────────────────────────────────────────────

func buildInfoUpdates(c *dto.UpdateInfoDTO, mask []string, ex *domain.PropertyInfo) map[string]any {
	u := make(map[string]any)

	diffStr(u, mask, "info.title", "title", c.Title, &ex.Title)
	diffStr(u, mask, "info.note", "note", c.Note, &ex.Note)
	diffStr(u, mask, "info.unit_code", "unit_code", c.UnitCode, &ex.UnitCode)
	diffStr(u, mask, "info.identifier", "identifier", c.Identifier, &ex.Identifier)
	diffStr(u, mask, "info.level", "level", c.Level, &ex.Level)
	diffInt32(u, mask, "info.scope", "scope", c.Scope, int32(ex.Scope))
	diffInt32(u, mask, "info.source_type", "source_type", c.SourceType, int32(ex.SourceType))
	diffInt32(u, mask, "info.legal_status", "legal_status", c.LegalStatus, int32(ex.LegalStatus))
	diffInt32(u, mask, "info.record_status", "record_status", c.RecordStatus, int32(ex.RecordStatus))

	// Nullable FK: nil = bỏ qua, *0 = clear
	diffNullableUint64(u, mask, "info.avatar_id", "avatar_id", c.AvatarID, ex.AvatarID)
	diffNullableUint64(u, mask, "info.property_type_id", "property_type_id", c.PropertyTypeID, ex.PropertyTypeID)
	diffNullableUint64(u, mask, "info.project_id", "project_id", c.ProjectID, ex.ProjectID)
	diffNullableUint64(u, mask, "info.origin_profile_id", "origin_profile_id", c.OriginProfileID, ex.OriginProfileId)

	return u
}

func buildLocationUpdates(c *dto.UpdateLocationDTO, mask []string, ex *domain.PropertyLocation) map[string]any {
	u := make(map[string]any)

	diffStr(u, mask, "location.address_detail", "address_detail", c.AddressDetail, &ex.AddressDetail)
	diffStr(u, mask, "location.map_url", "map_url", c.MapURL, &ex.MapURL)
	diffNullableFloat64(u, mask, "location.latitude", "latitude", c.Latitude, ex.Latitude)
	diffNullableFloat64(u, mask, "location.longitude", "longitude", c.Longitude, ex.Longitude)
	diffNullableUint64(u, mask, "location.region_id", "region_id", c.RegionID, ex.RegionID)
	diffNullableUint64(u, mask, "location.province_id", "province_id", c.ProvinceID, ex.ProvinceID)
	diffNullableUint64(u, mask, "location.district_id", "district_id", c.DistrictID, ex.DistrictID)
	diffNullableUint64(u, mask, "location.ward_id", "ward_id", c.WardID, ex.WardID)

	return u
}

func buildLandInfoUpdates(c *dto.UpdateLandInfoDTO, mask []string, ex *domain.PropertyLandInfo) map[string]any {
	u := make(map[string]any)

	diffStr(u, mask, "land_info.note", "note", c.Note, &ex.Note)
	diffStr(u, mask, "land_info.land_note", "land_note", c.LandNote, &ex.LandNote)
	diffInt32(u, mask, "land_info.purpose_used", "purpose_used", c.PurposeUsed, int32(ex.PurposeUsed))

	if c.Basic != nil {
		diffStr(u, mask, "land_info.basic.document_no", "document_no", c.Basic.DocumentNo, &ex.DocumentNo)
		diffStr(u, mask, "land_info.basic.issuring_auth", "issuring_auth", c.Basic.IssuringAuth, &ex.IssuringAuth)
		docType := ex.DocumentType.Value()
		diffUint32(u, mask, "land_info.basic.document_type", "document_type", c.Basic.DocumentType, &docType)
		diffUint32(u, mask, "land_info.basic.plot", "plot", c.Basic.Plot, &ex.Plot)
		diffUint32(u, mask, "land_info.basic.sheet", "sheet", c.Basic.Sheet, &ex.Sheet)
	}

	if c.Area != nil {
		diffNullableFloat64(u, mask, "land_info.area.area_total", "area_total", c.Area.AreaTotal, ex.AreaTotal)
		diffNullableFloat64(u, mask, "land_info.area.area_land", "area_land", c.Area.AreaLand, ex.AreaLand)
		diffNullableFloat64(u, mask, "land_info.area.area_plant", "area_plant", c.Area.AreaPlant, ex.AreaPlant)
	}

	if c.Dimension != nil {
		diffNullableFloat64(u, mask, "land_info.dimension.front_width", "front_width", c.Dimension.FrontWidth, ex.FrontWidth)
		diffNullableFloat64(u, mask, "land_info.dimension.depth", "depth", c.Dimension.Depth, ex.Depth)
		diffNullableFloat64(u, mask, "land_info.dimension.street_width", "street_width", c.Dimension.StreetWidth, ex.StreetWidth)
	}

	if c.Expiry != nil {
		diffDateStr(u, mask, "land_info.expiry.expired_land", "expired_land", c.Expiry.ExpiredLand, ex.ExpiredLand)
		diffDateStr(u, mask, "land_info.expiry.expired_plant", "expired_plant", c.Expiry.ExpiredPlant, ex.ExpiredPlant)
	}

	return u
}

func buildBuildingUpdates(c *dto.UpdateBuildingDTO, mask []string, ex *domain.PropertyBuildingInfo) map[string]any {
	u := make(map[string]any)

	diffStr(u, mask, "building_info.note", "note", c.Note, &ex.Note)
	diffInt32(u, mask, "building_info.build_status", "build_status", c.BuildStatus, int32(ex.BuildStatus))
	diffInt32(u, mask, "building_info.building_type", "building_type", c.BuildingType, int32(ex.BuildingType))
	diffInt32(u, mask, "building_info.direction", "direction", c.Direction, int32(ex.Direction))
	diffInt32(u, mask, "building_info.balcony_direction", "balcony_direction", c.BalconyDir, int32(ex.BalconyDirection))

	diffNullableFloat64(u, mask, "building_info.area_actual", "area_actual", c.AreaActual, ex.AreaActual)
	diffNullableFloat64(u, mask, "building_info.area_floor", "area_floor", c.AreaFloor, ex.AreaFloor)
	diffNullableFloat64(u, mask, "building_info.area_construction", "area_construction", c.AreaConstruction, ex.AreaConstruction)
	diffNullableUint32(u, mask, "building_info.floors", "floors", c.Floors, ex.Floors)
	diffNullableUint32(u, mask, "building_info.room_number", "room_number", c.RoomNumber, ex.RoomNumber)
	diffNullableUint32(u, mask, "building_info.bedrooms", "bedrooms", c.Bedrooms, ex.Bedrooms)
	diffNullableUint32(u, mask, "building_info.bathrooms", "bathrooms", c.Bathrooms, ex.Bathrooms)

	return u
}

func buildEvidenceUpdates(c *dto.UpdateEvidenceDTO, mask []string, ex *domain.PropertyEdvidence) map[string]any {
	u := make(map[string]any)

	diffStr(u, mask, "evidence.title", "title", c.Title, &ex.Title)
	diffStr(u, mask, "evidence.description", "description", c.Description, &ex.Description)
	diffNullableUint64(u, mask, "evidence.file_id", "file_id", c.FileID, ex.FileID)

	return u
}

// ─────────────────────────────────────────────────────────────────────────────
// Domain mappers (DTO → domain entity khi INSERT mới)
// ─────────────────────────────────────────────────────────────────────────────

func mapInfoDTOToEntity(lineageID uint64, c *dto.UpdateInfoDTO) *domain.PropertyInfo {
	e := &domain.PropertyInfo{
		Title:           derefStr(c.Title),
		Note:            derefStr(c.Note),
		UnitCode:        derefStr(c.UnitCode),
		Identifier:      derefStr(c.Identifier),
		Level:           derefStr(c.Level),
		AvatarID:        c.AvatarID,
		PropertyTypeID:  c.PropertyTypeID,
		ProjectID:       c.ProjectID,
		OriginProfileId: c.OriginProfileID,
	}
	if c.SourceType != nil {
		e.SourceType = enums.EPropertySourceType(*c.SourceType)
	}
	if c.LegalStatus != nil {
		e.LegalStatus = enums.EHouseCertificate(*c.LegalStatus)
	}
	if c.RecordStatus != nil {
		e.RecordStatus = enums.EPropertyStatus(*c.RecordStatus)
	}
	if c.Scope != nil {
		e.Scope = enums.EPropertyScope(*c.Scope)
	}
	return e
}

func mapLocationDTOToEntity(lineageID uint64, c *dto.UpdateLocationDTO) *domain.PropertyLocation {
	return &domain.PropertyLocation{
		AddressDetail: derefStr(c.AddressDetail),
		MapURL:        derefStr(c.MapURL),
		Latitude:      c.Latitude,
		Longitude:     c.Longitude,
		RegionID:      c.RegionID,
		ProvinceID:    c.ProvinceID,
		DistrictID:    c.DistrictID,
		WardID:        c.WardID,
	}
}

func mapLandInfoDTOToEntity(lineageID uint64, c *dto.UpdateLandInfoDTO) *domain.PropertyLandInfo {
	e := &domain.PropertyLandInfo{
		Note:     derefStr(c.Note),
		LandNote: derefStr(c.LandNote),
	}
	if c.PurposeUsed != nil {
		e.PurposeUsed = enums.EPurposeUsed(*c.PurposeUsed)
	}
	if c.Basic != nil {
		e.DocumentNo = derefStr(c.Basic.DocumentNo)
		e.IssuringAuth = derefStr(c.Basic.IssuringAuth)
		e.DocumentType = enums.EDocUnknown.Parse(c.Basic.DocumentType)
		e.Plot = derefUint32(c.Basic.Plot)
		e.Sheet = derefUint32(c.Basic.Sheet)
	}
	if c.Area != nil {
		e.AreaTotal = c.Area.AreaTotal
		e.AreaLand = c.Area.AreaLand
		e.AreaPlant = c.Area.AreaPlant
	}
	if c.Dimension != nil {
		e.FrontWidth = c.Dimension.FrontWidth
		e.Depth = c.Dimension.Depth
		e.StreetWidth = c.Dimension.StreetWidth
	}
	if c.Expiry != nil {
		if c.Expiry.ExpiredLand != nil && *c.Expiry.ExpiredLand != "" {
			if t, err := time.Parse(time.RFC3339, *c.Expiry.ExpiredLand); err == nil {
				e.ExpiredLand = &t
			}
		}
		if c.Expiry.ExpiredPlant != nil && *c.Expiry.ExpiredPlant != "" {
			if t, err := time.Parse(time.RFC3339, *c.Expiry.ExpiredPlant); err == nil {
				e.ExpiredPlant = &t
			}
		}
	}
	return e
}

func mapBuildingDTOToEntity(lineageID uint64, c *dto.UpdateBuildingDTO) *domain.PropertyBuildingInfo {
	e := &domain.PropertyBuildingInfo{
		Note:             derefStr(c.Note),
		AreaActual:       c.AreaActual,
		AreaFloor:        c.AreaFloor,
		AreaConstruction: c.AreaConstruction,
		Floors:           c.Floors,
		RoomNumber:       c.RoomNumber,
		Bedrooms:         c.Bedrooms,
		Bathrooms:        c.Bathrooms,
	}
	if c.BuildStatus != nil {
		e.BuildStatus = enums.EBuildStatus(*c.BuildStatus)
	}
	if c.BuildingType != nil {
		e.BuildingType = enums.EBuildingType(*c.BuildingType)
	}
	if c.Direction != nil {
		e.Direction = enums.EHouseOrient(*c.Direction)
	}
	if c.BalconyDir != nil {
		e.BalconyDirection = enums.EHouseOrient(*c.BalconyDir)
	}
	return e
}

// ─────────────────────────────────────────────────────────────────────────────
// Diff helpers
//
// Convention:
//   diffXxx(updates, mask, maskPath, col, incoming, existing)
//   - Chỉ thêm vào updates nếu maskPath khớp VÀ giá trị thực sự thay đổi
// ─────────────────────────────────────────────────────────────────────────────

// diffStr: incoming *string nil = bỏ qua; non-nil = so sánh với existing
func diffStr(u map[string]any, mask []string, path, col string, incoming *string, existing *string) {
	if !inMask(mask, path) || incoming == nil {
		return
	}
	exVal := ""
	if existing != nil {
		exVal = *existing
	}
	if *incoming != exVal {
		u[col] = *incoming
	}
}

func diffInt32(u map[string]any, mask []string, path, col string, incoming *int32, existing int32) {
	if !inMask(mask, path) || incoming == nil {
		return
	}
	if *incoming != existing {
		u[col] = *incoming
	}
}

func diffUint32(u map[string]any, mask []string, path, col string, incoming *uint32, existing *uint32) {
	if !inMask(mask, path) || incoming == nil {
		return
	}
	exVal := uint32(0)
	if existing != nil {
		exVal = *existing
	}
	if *incoming != exVal {
		u[col] = *incoming
	}
}

// diffNullableUint64: nil incoming = bỏ qua; *0 incoming = set NULL; *non0 = set value
func diffNullableUint64(u map[string]any, mask []string, path, col string, incoming *uint64, existing *uint64) {
	if !inMask(mask, path) || incoming == nil {
		return
	}
	// incoming not nil → so sánh với existing
	if existing == nil && *incoming == 0 {
		return // both null, no change
	}
	if existing != nil && *incoming == *existing {
		return // same value
	}
	if *incoming == 0 {
		u[col] = nil // clear FK
	} else {
		u[col] = *incoming
	}
}

func diffNullableFloat64(u map[string]any, mask []string, path, col string, incoming *float64, existing *float64) {
	if !inMask(mask, path) || incoming == nil {
		return
	}
	if existing == nil {
		u[col] = *incoming
		return
	}
	if *incoming != *existing {
		u[col] = *incoming
	}
}

func diffNullableUint32(u map[string]any, mask []string, path, col string, incoming *uint32, existing *uint32) {
	if !inMask(mask, path) || incoming == nil {
		return
	}
	if existing == nil || *incoming != *existing {
		u[col] = *incoming
	}
}

// diffDateStr: incoming nil = bỏ qua; "" = clear (set NULL); ISO8601 = so sánh & set
func diffDateStr(u map[string]any, mask []string, path, col string, incoming *string, existing *time.Time) {
	if !inMask(mask, path) || incoming == nil {
		return
	}
	if *incoming == "" {
		if existing != nil {
			u[col] = nil // clear
		}
		return
	}
	t, err := time.Parse(time.RFC3339, *incoming)
	if err != nil {
		return // invalid format → bỏ qua, không lỗi hard (caller có thể validate trước)
	}
	if existing == nil || !t.Equal(*existing) {
		u[col] = t
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Utilities
// ─────────────────────────────────────────────────────────────────────────────

// inMask: mask rỗng = cho phép tất cả (backward compat).
func inMask(mask []string, path string) bool {
	if len(mask) == 0 {
		return true
	}
	for _, p := range mask {
		if p == path || strings.HasPrefix(p, path+".") {
			return true
		}
	}
	return false
}

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func derefUint32(p *uint32) uint32 {
	if p == nil {
		return 0
	}
	return *p
}
