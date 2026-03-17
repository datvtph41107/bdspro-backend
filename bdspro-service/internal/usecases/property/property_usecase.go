package property_usecases

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
	"common/pkg/fieldmask"
	"common/pkg/patch"
	_provider "common/provider"
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
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
	AreaRegionRepo           property_repo.AreaRegionRepository
	ExternalRefRepo          property_repo.PropertyExternalRefRepository
	ReportRepo               property_repo.PropertyReportRepository
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
	areaRegionRepo property_repo.AreaRegionRepository,
	externalRefRepo property_repo.PropertyExternalRefRepository,
	reportRepo property_repo.PropertyReportRepository,
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
		AreaRegionRepo:           areaRegionRepo,
		AmenityRepo:              amenityRepo,
		ReportRepo:               reportRepo,
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

func (u *PropertyUsecase) AddPropertyUser(
	ctx context.Context,
	propertyIdentifyID, userID uint64,
) (*domain.PropertyLineage, error) {
	lineage, err := u.LineageRepo.GetByPropertyIdentifyID(ctx, propertyIdentifyID)
	if err != nil {
		return nil, err
	}
	if lineage == nil {
		return nil, _errors.NotFoundException("Không tìm thấy BĐS")
	}

	existed, err := u.PropertyUserRepo.GetByLineageAndOwnerUnscoped(ctx, lineage.ID, userID)
	if err != nil {
		return nil, err
	}
	if existed != nil {
		if existed.DeletedAt != nil {
			if err := u.PropertyUserRepo.RestoreByID(ctx, existed.ID); err != nil {
				return nil, err
			}
			return lineage, nil
		}
		return nil, _errors.ConflictException("Đã thêm BĐS này")
	}

	if _, err := u.PropertyUserRepo.CreateOwner(ctx, lineage.ID, userID); err != nil {
		return nil, err
	}

	return lineage, nil
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
		// statisticID *uint64
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
				LegalStatus: enums.EHouseCertificateNone,
				// PrivacyLevel:    enums.EVisiblePublic,
				Note:       "Note",
				UnitCode:   "UnitCode",
				Identifier: "Identifier",
				Level:      "Level",
				// OriginProfileId: &originId,
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
			// for i := range req.MediaList {
			// 	req.MediaList[i].LineageID = &property.ID
			// }
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
			return _errors.NotFoundException(fmt.Sprintf("Property with ID %d not found", id))
		}
	}

	return u.PropertyRepo.ArchiveBatchProperties(ctx, req.IDs, req.Archived, originId)
}

func (u *PropertyUsecase) HideProperties(ctx context.Context, ids []uint64, hidden bool) error {
	originId := _utils.GetOriginIdFromContext(ctx)
	if originId == 0 {
		return _errors.UnauthorizedException()
	}

	if len(ids) == 0 {
		return _errors.BadRequestException("No property IDs provided")
	}

	for _, id := range ids {
		property, err := u.PropertyRepo.GetOneByID(ctx, id)
		if err != nil {
			return _errors.InternalServerException("Failed to validate property ID: " + err.Error())
		}
		if property == nil {
			return _errors.NotFoundException(fmt.Sprintf("Property with ID %d not found", id))
		}
	}

	return u.PropertyRepo.HideBatchProperties(ctx, ids, hidden, originId)
}

// // UPDATE V2 PROPERTY
// ─────────────────────────────────────────────────────────────────────────────
// UpdateProperty — entry point
//
// Chỉ làm 3 việc:
//   1. Load lineage (duy nhất 1 SELECT bắt buộc)
//   2. Routing: hệ thống + non-admin → propose; còn lại → direct
//   3. Delegate xuống proposeChange hoặc directUpdate
// ─────────────────────────────────────────────────────────────────────────────

// ─────────────────────────────────────────────────────────────────────────────
// UpdateProperty — entry point của usecase
//
// Trách nhiệm:
//  1. Transaction boundary
//  2. Load lineage + auth check
//  3. Dispatch đúng flow: directUpdate hoặc proposeChange
//  4. Side effect async: sync timestamp (không block response)
// ─────────────────────────────────────────────────────────────────────────────

func (u *PropertyUsecase) UpdateProperty(
	ctx context.Context,
	c *dto.UpdatePropertyDTO,
	userID uint64,
) (*dto.UpdatePropertyResult, error) {
	var result *dto.UpdatePropertyResult

	err := u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		lineage, err := u.LineageRepo.GetByID(txCtx, c.LineageID)
		if err != nil {
			return err
		}
		if lineage == nil {
			return fmt.Errorf("property lineage id=%d not found", c.LineageID)
		}

		propertyUser, err := u.PropertyUserRepo.GetByLineageAndOwner(txCtx, lineage.ID, userID)
		if err != nil {
			return err
		}
		if propertyUser == nil {
			return _errors.ForbiddenException("You dont have permission action")
		}

		// Dispatch flow — mở rộng ở đây khi có propose/report flow
		result, err = u.directUpdate(txCtx, lineage, c)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Async side effect — không block response, không affect transaction
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		profileID := _utils.GetProfileIdWithContext(cloneCtx)
		now := time.Now().UnixMilli()
		u.SyncProvider.PutTimestamp(cloneCtx,
			u.SyncProvider.GetKey(cloneCtx, _provider.SyncKeyPropertyMe, profileID), now)
		u.SyncProvider.PutTimestamp(cloneCtx,
			u.SyncProvider.GetKey(cloneCtx, _provider.SyncKeyPropertyDetail, result.LineageID), now)
	}()

	return result, nil
}

// ═════════════════════════════════════════════════════════════════════════════
// FLOW 1 — DirectUpdate
//
// Dùng cho: owner update bất động sản của mình.
//
// Logic:
//   Mỗi sub-domain chạy qua RunDirect:
//     FK nil → Insert (tạo record mới) → thêm FK vào patch
//     FK có  → Update (partial, chỉ field được gửi)
//     absent → skip hoàn toàn
//
//   Side blocks (media, amenities, refs, meta) xử lý riêng.
//   Gom tất cả FK mới → 1 UPDATE lineage duy nhất cuối cùng.
// ═════════════════════════════════════════════════════════════════════════════

func (u *PropertyUsecase) directUpdate(
	ctx context.Context,
	lineage *domain.PropertyLineage,
	cmd *dto.UpdatePropertyDTO,
) (*dto.UpdatePropertyResult, error) {
	fkPatch := make(FKPatch)

	// Side blocks: media, amenities, refs, area_regions, lineage_meta
	if err := u.handleSideBlocks(ctx, lineage.ID, cmd, fkPatch); err != nil {
		return nil, err
	}

	// Lấy ảnh đầu tiên trong media làm avatar, gán ID vào info
	if list, err := u.MediaRepo.GetByLineageID(ctx, lineage.ID); err == nil && len(list) > 0 {
		var avatarID uint64
		for i := range list {
			if list[i].MediaType == "image" {
				avatarID = list[i].ID
				break
			}
		}
		if avatarID == 0 {
			avatarID = list[0].ID // không có image thì lấy media đầu tiên
		}
		if cmd.Info == nil {
			cmd.Info = &dto.UpdateInfoDTO{}
		}
		cmd.Info.AvatarID = &avatarID
	}

	// Khai báo handlers — thứ tự không quan trọng, độc lập nhau
	// Mask được truyền vào từng handler để patch.Builder guard đúng
	handlers := []struct {
		currentFK *uint64
		h         Handler
	}{
		{
			lineage.PropertyInfoID,
			&InfoHandler{Patch: cmd.Info, Mask: cmd.Mask, Repo: u.InfoRepo},
		},
		{
			lineage.LocationID,
			&LocationHandler{Patch: cmd.Location, Mask: cmd.Mask, Repo: u.LocationRepo},
		},
		{
			lineage.LandInfoID,
			&LandInfoHandler{Patch: cmd.LandInfo, Mask: cmd.Mask, Repo: u.LandInfoRepo},
		},
		{
			lineage.BuildingInfoID,
			&BuildingHandler{Patch: cmd.BuildingInfo, Mask: cmd.Mask, Repo: u.BuildingInfoRepo},
		},
		{
			lineage.EdvidenceID,
			&EvidenceHandler{Patch: cmd.Evidence, Mask: cmd.Mask, Repo: u.EvidenceRepo},
		},
		{
			lineage.EdvidenceID,
			&EvidenceHandler{Patch: cmd.Evidence, Mask: cmd.Mask, Repo: u.EvidenceRepo},
		},
	}

	for _, entry := range handlers {
		newFK, err := RunDirect(ctx, entry.currentFK, entry.h)
		if err != nil {
			return nil, err
		}
		fkPatch.Set(entry.h.FKColumn(), newFK)
	}

	// Flush tất cả FK + lineage_meta fields vào lineage — 1 UPDATE duy nhất
	if len(fkPatch) > 0 {
		if err := u.LineageRepo.UpdateFields(ctx, lineage.ID, fkPatch); err != nil {
			return nil, fmt.Errorf("lineage update: %w", err)
		}
	}

	return &dto.UpdatePropertyResult{
		LineageID: lineage.ID,
		UpdatedAt: _utils.TimeNowRFC3339(),
	}, nil
}

// ═════════════════════════════════════════════════════════════════════════════
// FLOW 2 — ProposeChange
//
// Dùng cho: user đề xuất thay đổi bất động sản của người khác (report/contribute).
//
// Logic:
//   Tạo lineage mới (propose), copy tất cả FK từ origin.
//   Các sub-domain có dữ liệu → Insert record mới → update FK propose.
//   Sub-domain absent → giữ FK từ origin (shared read, không tốn storage).
// ═════════════════════════════════════════════════════════════════════════════

func (u *PropertyUsecase) proposeChange(
	ctx context.Context,
	origin *domain.PropertyLineage,
	cmd *dto.UpdatePropertyDTO,
) (*dto.UpdatePropertyResult, error) {
	// Tạo propose lineage — copy FK từ origin làm base
	propose := &domain.PropertyLineage{
		PropertyInfoID: origin.PropertyInfoID,
		LocationID:     origin.LocationID,
		LandInfoID:     origin.LandInfoID,
		BuildingInfoID: origin.BuildingInfoID,
		EdvidenceID:    origin.EdvidenceID,
		// Status, proposed_at, proposed_by, ref_origin_id... set ở đây
	}
	if err := u.LineageRepo.Create(ctx, propose); err != nil {
		return nil, fmt.Errorf("create propose lineage: %w", err)
	}

	fkPatch := make(FKPatch)

	// Propose: luôn Insert mới nếu có dữ liệu (không Update origin)
	proposeHandlers := []Handler{
		&InfoHandler{Patch: cmd.Info, Mask: cmd.Mask, Repo: u.InfoRepo},
		&LocationHandler{Patch: cmd.Location, Mask: cmd.Mask, Repo: u.LocationRepo},
		&LandInfoHandler{Patch: cmd.LandInfo, Mask: cmd.Mask, Repo: u.LandInfoRepo},
		&BuildingHandler{Patch: cmd.BuildingInfo, Mask: cmd.Mask, Repo: u.BuildingInfoRepo},
		&EvidenceHandler{Patch: cmd.Evidence, Mask: cmd.Mask, Repo: u.EvidenceRepo},
	}

	for _, h := range proposeHandlers {
		newID, err := RunPropose(ctx, h)
		if err != nil {
			return nil, err
		}
		if newID != 0 {
			fkPatch[h.FKColumn()] = newID
		}
	}

	if err := u.handleSideBlocks(ctx, propose.ID, cmd, fkPatch); err != nil {
		return nil, err
	}

	if len(fkPatch) > 0 {
		if err := u.LineageRepo.UpdateFields(ctx, propose.ID, fkPatch); err != nil {
			return nil, fmt.Errorf("propose fk update: %w", err)
		}
	}

	return &dto.UpdatePropertyResult{
		LineageID: propose.ID,
		UpdatedAt: _utils.TimeNowRFC3339(),
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// handleSideBlocks
//
// Side blocks có semantic khác với FK blocks:
//   media/external_refs : upsert per item theo lineage_id
//   amenities/area_regions: REPLACE toàn bộ list
//   lineage_meta        : update trực tiếp vào fkPatch (cùng lineage record)
//
// fkPatch được truyền vào để lineage_meta merge vào cùng 1 UPDATE.
// ─────────────────────────────────────────────────────────────────────────────

func (u *PropertyUsecase) handleSideBlocks(
	ctx context.Context,
	lineageID uint64,
	cmd *dto.UpdatePropertyDTO,
	fkPatch FKPatch,
) error {
	originID := _utils.GetOriginIdFromContext(ctx)
	if originID == 0 {
		return _errors.UnauthorizedException("origin_id required")
	}
	if cmd.Media != nil {
		if err := u.handleMedia(ctx, lineageID, cmd.Media); err != nil {
			return fmt.Errorf("media: %w", err)
		}
	}
	if cmd.Amenities != nil {
		if err := u.handleAmenities(ctx, lineageID, cmd.Amenities); err != nil {
			return fmt.Errorf("amenities: %w", err)
		}
	}
	if cmd.AreaRegions != nil {
		if err := u.handleAreaRegions(ctx, lineageID, cmd.AreaRegions); err != nil {
			return fmt.Errorf("area_regions: %w", err)
		}
	}
	if cmd.ExternalRefs != nil {
		if err := u.handleExternalRefs(ctx, cmd.ExternalRefs); err != nil {
			return fmt.Errorf("external_refs: %w", err)
		}
	}
	if cmd.LineageMeta != nil {
		// Merge vào fkPatch — cùng UPDATE lineage record cuối cùng
		u.applyLineageMeta(cmd.LineageMeta, cmd.Mask, fkPatch)
	}
	if cmd.Personalization != nil {
		if err := u.updatePersonalization(ctx, cmd.Mask, lineageID, originID, cmd.Personalization); err != nil {
			return fmt.Errorf("personalization: %w", err)
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// handleMedia — xóa media không còn trong list, rồi upsert per item
//
// id=0 → INSERT mới
// id>0 → UPDATE record đó
// Media có lineage_id nhưng không nằm trong data.Items → DELETE
// ─────────────────────────────────────────────────────────────────────────────

func (u *PropertyUsecase) handleMedia(
	ctx context.Context,
	lineageID uint64,
	data *dto.UpdateMediaDTO,
) error {
	keepIDs := make([]uint64, 0, len(data.Items))
	for _, item := range data.Items {
		if item.ID > 0 {
			keepIDs = append(keepIDs, item.ID)
		}
	}
	if err := u.MediaRepo.DeleteByLineageIDExcept(ctx, lineageID, keepIDs); err != nil {
		return fmt.Errorf("delete media: %w", err)
	}
	for orderIdx, item := range data.Items {
		if item.ID > 0 {
			// UPDATE — partial update, chỉ field non-zero
			// Media dùng scalar (không *T) → update tất cả field có trong item
			fields := map[string]any{
				"media_type": item.MediaType,
				"media_url":  item.MediaURL,
				"thumb_url":  item.ThumbURL,
				"sort_order": orderIdx + 1,
				// "is_cover":   item.IsCover,
			}
			if err := u.MediaRepo.UpdateFields(ctx, item.ID, fields); err != nil {
				return fmt.Errorf("update media id=%d: %w", item.ID, err)
			}
		} else {
			// INSERT mới
			entity := &domain.PropertyMedia{
				LineageID: lineageID,
				MediaType: item.MediaType,
				MediaURL:  item.MediaURL,
				ThumbURL:  item.ThumbURL,
				SortOrder: int32(orderIdx + 1),
				IsCover:   item.IsCover,
			}
			if err := u.MediaRepo.Create(ctx, entity); err != nil {
				return fmt.Errorf("create media: %w", err)
			}
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// handleAmenities — REPLACE toàn bộ list
// ─────────────────────────────────────────────────────────────────────────────

func (u *PropertyUsecase) handleAmenities(
	ctx context.Context,
	lineageID uint64,
	data *dto.UpdateAmenityDTO,
) error {
	// AmenityIDs nil = sub-domain không gửi → đã skip ở handleSideBlocks
	// AmenityIDs [] = gửi list rỗng = xóa hết
	return u.AmenityRepo.ReplaceAmenities(ctx, lineageID, data.AmenityIDs)
}

// ─────────────────────────────────────────────────────────────────────────────
// handleAreaRegions — REPLACE toàn bộ list
// ─────────────────────────────────────────────────────────────────────────────

func (u *PropertyUsecase) handleAreaRegions(
	ctx context.Context,
	lineageID uint64,
	data *dto.UpdateAreaRegionsDTO,
) error {
	return u.AreaRegionRepo.ReplaceAreaRegions(ctx, lineageID, data.AreaRegionIDs)
}

// ─────────────────────────────────────────────────────────────────────────────
// handleExternalRefs — upsert per item
// ─────────────────────────────────────────────────────────────────────────────

func (u *PropertyUsecase) handleExternalRefs(
	ctx context.Context,
	data *dto.UpdateExternalRefDTO,
) error {
	for _, item := range data.Items {
		if item.ID > 0 {
			// UPDATE — scalar fields, update tất cả
			fields := map[string]any{
				"external_ref_id": item.ExternalRefID,
				"source_system":   item.SourceSystem,
				"source_code":     item.SourceCode,
				"confidence":      item.Confidence,
				"sync_status":     item.SyncStatus,
				"note":            item.Note,
			}
			if err := u.ExternalRefRepo.UpdateFields(ctx, item.ID, fields); err != nil {
				return fmt.Errorf("update external_ref id=%d: %w", item.ID, err)
			}
		} else {
			// INSERT mới
			entity := &domain.PropertyExternalRef{
				ExternalRefID: item.ExternalRefID,
				SourceSystem:  enums.ERefSourceSys(item.SourceSystem),
				SourceCode:    item.SourceCode,
				Confidence:    item.Confidence,
				SyncStatus:    enums.ESyncStatus(item.SyncStatus),
				Note:          item.Note,
			}
			if err := u.ExternalRefRepo.Create(ctx, entity); err != nil {
				return fmt.Errorf("create external_ref: %w", err)
			}
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// applyLineageMeta — merge lineage_meta fields vào fkPatch
//
// Dùng patch.Builder để đảm bảo cùng quy chuẩn presence + diff + mask.
// Merge vào fkPatch → cùng 1 UPDATE lineage record cuối cùng.
// ─────────────────────────────────────────────────────────────────────────────

func (u *PropertyUsecase) applyLineageMeta(
	meta *dto.UpdateLineageMetaDTO,
	mask fieldmask.FieldMask,
	fkPatch FKPatch,
) {
	fields := patch.WithMask(mask).
		String("lineage_meta.national_id", "national_id", meta.NationalID, "").
		Int32("lineage_meta.lineage_status", "lineage_status", meta.LineageStatus, 10).
		DateString("lineage_meta.verified_national_at", "verified_national_at", meta.VerifiedNationalAt, nil).
		Build()

	for k, v := range fields {
		fkPatch[k] = v
	}
}
func (u *PropertyUsecase) updatePersonalization(
	ctx context.Context,
	mask fieldmask.FieldMask,
	lineageID uint64,
	originID uint64,
	data *dto.UpdatePersonalizationDTO,
) error {
	fields := patch.WithMask(mask).
		DateString("lineage_meta.archived_at", "archived_at", data.ArchivedAt, nil).
		DateString("lineage_meta.hidden_at", "hidden_at", data.HiddenAt, nil).
		Build()

	return u.PropertyUserRepo.UpdateFieldsByLineageAndOwner(ctx, lineageID, originID, fields)
}

// ═════════════════════════════════════════════════════════════════════════════
// SubmitReport — API báo cáo sai / gắn cờ / đóng góp BĐS
//
// Flow:
//  1. Load lineage gốc, verify tồn tại.
//  2. Check duplicate report (cùng reporter + lineage + status pending/reviewing).
//  3. Nếu có propose data → chạy proposeChange (tái sử dụng flow hiện tại)
//     tạo propose lineage mới copy từ origin.
//  4. Tạo bản ghi property_report liên kết lineage gốc + propose lineage.
//  5. Async: bump sync timestamp.
//
// NOTE:
//   - Reporter KHÔNG cần là owner của lineage (khác với UpdateProperty).
//   - proposeChange được tái sử dụng nguyên vẹn — không sửa.
//   - Usecase này chỉ thêm orchestration layer (report record) bên ngoài.
//
// ═════════════════════════════════════════════════════════════════════════════
func (u *PropertyUsecase) SubmitReport(
	ctx context.Context,
	cmd *dto.SubmitReportDTO,
	reporterOriginID uint64,
) (*dto.SubmitReportResult, error) {
	var result *dto.SubmitReportResult

	err := u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		lineage, err := u.LineageRepo.GetByID(txCtx, cmd.LineageID)
		if err != nil {
			return err
		}
		if lineage == nil {
			return fmt.Errorf("property lineage id=%d not found", cmd.LineageID)
		}

		existing, err := u.ReportRepo.GetByLineageAndReporter(txCtx, lineage.ID, reporterOriginID)
		if err != nil {
			return err
		}
		if existing != nil {
			return _errors.ConflictException("you already have a pending report for this property")
		}

		var proposeLineageID *uint64
		if cmd.ContributeData != nil {
			proposeResult, err := u.proposeChange(txCtx, lineage, cmd.ContributeData)
			if err != nil {
				return fmt.Errorf("propose change: %w", err)
			}
			proposeLineageID = &proposeResult.LineageID
		}

		report := &domain.PropertyReport{
			LineageID:          lineage.ID,
			ReporterOriginID:   reporterOriginID,
			ReportType:         enums.EReportType(cmd.ReportType),
			SubIssues:          marshalJSON(cmd.SubIssues),
			Status:             enums.EReportStatusPending,
			Description:        derefStr(cmd.Note),
			RelatedLineageID:   cmd.RelatedLineageID,
			AttachmentMediaIDs: marshalJSON(cmd.AttachmentMediaIDs),
			AttachmentFileIDs:  marshalJSON(cmd.AttachmentFileIDs),
		}
		if err := u.ReportRepo.Create(txCtx, report); err != nil {
			return fmt.Errorf("create report: %w", err)
		}

		result = &dto.SubmitReportResult{
			ReportID:         report.ID,
			LineageID:        lineage.ID,
			ProposeLineageID: proposeLineageID,
			ReportType:       cmd.ReportType,
			Note:             derefStr(cmd.Note),
			CreatedAt:        _utils.TimeNowRFC3339(),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		now := time.Now().UnixMilli()
		u.SyncProvider.PutTimestamp(cloneCtx,
			u.SyncProvider.GetKey(cloneCtx, _provider.SyncKeyPropertyDetail, result.LineageID), now)
	}()

	return result, nil
}
