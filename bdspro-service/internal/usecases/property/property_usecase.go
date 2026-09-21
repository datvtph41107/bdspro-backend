package property_usecases

import (
	"bdspro/internal"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/modules"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	property_repo "bdspro/internal/repo/property"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_errors "common/errors"
	"common/logging"
	"common/pkg/fieldmask"
	"common/pkg/patch"
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"pb/clients"
	tqdpb "pb/types/tqd"
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
	SyncProvider             *_utils.SyncUtil
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
	ProvinceRepo             property_repo.ProvinceV2Repository
	WardRepo                 property_repo.WardV2Repository
	PropertyImpactRepo       property_repo.PropertyImpactRepository
	AuditLogRepo             property_repo.AuditLogRepository

	tqdGRPCClient *clients.TQDGrpcClient
	impactEngine  modules.ImpactDetectionEngine
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
	syncProvider *_utils.SyncUtil,
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
	provinceRepo property_repo.ProvinceV2Repository,
	wardRepo property_repo.WardV2Repository,
	propertyImpactRepo property_repo.PropertyImpactRepository,
	tqdGRPCClient *clients.TQDGrpcClient,
	impactEngine modules.ImpactDetectionEngine,
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
		ProvinceRepo:             provinceRepo,
		WardRepo:                 wardRepo,
		PropertyImpactRepo:       propertyImpactRepo,
		tqdGRPCClient:            tqdGRPCClient,
		impactEngine:             impactEngine,
	}
}

func (u *PropertyUsecase) GetPropertyAmenities(
	ctx context.Context,
	propertyID uint64,
	sourceType uint32,
) ([]*dto.PropertyAmenityDTO, error) {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("Unauthorized"), _errors.WithLegacyCode(401))
	}

	property, err := u.PropertyRepo.GetDetail(ctx, propertyID, sourceType)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, _errors.ReturnError(service.PropertyNotFound)
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
		return nil, _errors.ReturnError(service.PropertyNotFound, _errors.WithPublicMessage("Không tìm thấy BĐS"))
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
		return nil, _errors.ReturnError(service.PropertyAlreadyAdded)
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
		return _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("Unauthorized"), _errors.WithLegacyCode(401))
	}

	property, err := u.PropertyRepo.GetDetail(ctx, req.PropertyID, req.SourceType)
	if err != nil {
		return err
	}
	if property == nil {
		return _errors.ReturnError(service.PropertyNotFound)
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
		return nil, _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("Unauthorized"), _errors.WithLegacyCode(401))
	}

	property, err := u.PropertyRepo.GetDetail(ctx, propertyID, sourceType)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, _errors.ReturnError(service.PropertyNotFound)
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
		return nil, 0, _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("Unauthorized"), _errors.WithLegacyCode(401))
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
		return _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("Unauthorized"), _errors.WithLegacyCode(401))
	}

	property, err := u.PropertyRepo.GetDetail(ctx, req.PropertyID, req.SourceType)
	if err != nil {
		return err
	}
	if property == nil {
		return _errors.ReturnError(service.PropertyNotFound)
	}

	// validate tag ids (nếu có)
	if len(req.TagIDs) > 0 {
		tags, err := u.PropertyAmenityRepo.FindByIDs(ctx, req.TagIDs)
		if err != nil {
			return err
		}
		if len(tags) != len(req.TagIDs) {
			return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("Invalid tag ids"))
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
		return nil, 0, _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("origin_id required"), _errors.WithLegacyCode(401))
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
		return nil, _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("id is required"))
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
			return fmt.Errorf("create property info: %w", err)
		}

		// 2. Tạo Location
		if req.Location != nil {
			loc := req.Location
			loc.PropertyIdentifyID = identifyID
			if err := u.PropertyLocationRepo.Create(txCtx, loc); err != nil {
				return fmt.Errorf("create property location: %w", err)
			}
			locationID = &loc.ID
			infoID = &info.ID
		}

		// 3. Tạo LandInfo
		if req.LandInfo != nil {
			landInfo := req.LandInfo
			landInfo.PropertyIdentifyID = identifyID
			if err := u.PropertyLandInfoRepo.Create(txCtx, landInfo); err != nil {
				return fmt.Errorf("create property land info: %w", err)
			}
			landInfoID = &landInfo.ID
		}

		// 4. Tạo BuildingInfo
		if req.BuildingInfo != nil {
			buildingInfo := req.BuildingInfo
			buildingInfo.PropertyIdentifyID = identifyID
			if err := u.PropertyBuildingInfoRepo.Create(txCtx, buildingInfo); err != nil {
				return fmt.Errorf("create property building info: %w", err)
			}
			buildingID = &buildingInfo.ID
		}

		// 5. Tạo Edvidence
		if req.Edvidence != nil {
			edv := req.Edvidence
			edv.PropertyIdentifyID = identifyID
			if err := u.PropertyEdvidenceRepo.Create(txCtx, edv); err != nil {
				return fmt.Errorf("create property evidence: %w", err)
			}
			edvidenceID = &edv.ID
		}

		// 6. Tạo ExternalRef
		if req.ExternalRef != nil {
			extRef := req.ExternalRef
			extRef.PropertyIdentifyID = identifyID
			if err := u.PropertyExternalRefRepo.Create(txCtx, extRef); err != nil {
				return fmt.Errorf("create property external ref: %w", err)
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
			return fmt.Errorf("create property statistic: %w", err)
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
			return fmt.Errorf("create property lineage: %w", err)
		}

		// 9. Tạo PropertyUser
		if _, err := u.PropertyUserRepo.CreateOwner(txCtx, property.ID, originId); err != nil {
			return fmt.Errorf("create property user: %w", err)
		}

		// 7. Tạo MediaList
		if len(req.MediaList) > 0 {
			// for i := range req.MediaList {
			// 	req.MediaList[i].LineageID = &property.ID
			// }
			if err := u.PropertyMediaRepo.CreateBatch(txCtx, req.MediaList); err != nil {
				return fmt.Errorf("create property media: %w", err)
			}
			avatarID := req.MediaList[0].ID

			info.AvatarID = &avatarID
			if err := u.PropertyInfoRepo.Save(txCtx, info); err != nil {
				return fmt.Errorf("create property info: %w", err)
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
			u.SyncProvider.PutTimestamp(cloneCtx, u.SyncProvider.GetKey(cloneCtx, _utils.SyncKeyPropertyMe, 0), property.UpdatedAt.UnixMilli())
			u.SyncProvider.PutTimestamp(cloneCtx, u.SyncProvider.GetKey(cloneCtx, _utils.SyncKeyProductMe, 0), property.UpdatedAt.UnixMilli())
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
		return _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("Unauthorized"), _errors.WithLegacyCode(401))
	}

	if len(req.IDs) == 0 {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("No property IDs provided"))
	}

	for _, id := range req.IDs {
		property, err := u.PropertyRepo.GetOneByID(ctx, id)
		if err != nil {
			return fmt.Errorf("validate property id %d: %w", id, err)
		}
		if property == nil {
			return _errors.ReturnError(service.PropertyNotFound, _errors.WithPublicMessage(fmt.Sprintf("Property with ID %d not found", id)))
		}
	}

	return u.PropertyRepo.ArchiveBatchProperties(ctx, req.IDs, req.Archived, originId)
}

func (u *PropertyUsecase) HideProperties(ctx context.Context, ids []uint64, hidden bool) error {
	originId := _utils.GetOriginIdFromContext(ctx)
	if originId == 0 {
		return _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("Unauthorized"), _errors.WithLegacyCode(401))
	}

	if len(ids) == 0 {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("No property IDs provided"))
	}

	for _, id := range ids {
		property, err := u.PropertyRepo.GetOneByID(ctx, id)
		if err != nil {
			return fmt.Errorf("validate property id %d: %w", id, err)
		}
		if property == nil {
			return _errors.ReturnError(service.PropertyNotFound, _errors.WithPublicMessage(fmt.Sprintf("Property with ID %d not found", id)))
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
//
// ─────────────────────────────────────────────────────────────────────────────
func (u *PropertyUsecase) UpdateProperty(
	ctx context.Context,
	c *dto.UpdatePropertyDTO,
	userID uint64,
) (*dto.UpdatePropertyResult, error) {
	var result *dto.UpdatePropertyResult
	var impactEval *modules.ImpactEvaluation

	err := u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// Step 1: Load current lineage
		lineage, err := u.loadLineage(txCtx, c.LineageID)
		if err != nil {
			return err
		}

		// Step 2: Check permission and get actor role
		actorRole, err := u.checkPermission(txCtx, lineage, userID)
		if err != nil {
			return err
		}

		// Step 3: Evaluate impact before commit
		evaluation, err := u.evaluateImpact(txCtx, lineage, c, userID, actorRole)
		if err != nil {
			return err
		}
		impactEval = evaluation

		// Step 4: Apply action policy based on evaluation
		if err := u.applyActionPolicy(evaluation, c.ImpactAssessment, actorRole); err != nil {
			return err
		}

		// Step 5: Validate location update
		if err := u.validateLocationUpdate(txCtx, lineage, c); err != nil {
			return fmt.Errorf("location validation failed: %w", err)
		}

		// Step 6: Execute the update
		result, err = u.executeUpdate(txCtx, lineage, c, userID)
		if err != nil {
			return err
		}

		// Step 7: Post-update tasks (async within transaction)
		u.postUpdateTasks(txCtx, lineage.ID, userID, impactEval, c.ImpactAssessment)

		return nil
	})

	// Handle require confirm error (return evaluation for UI)
	if impactErr, ok := err.(*ImpactRequiredError); ok {
		return &dto.UpdatePropertyResult{
			LineageID:        c.LineageID,
			RequiresConfirm:  true,
			ImpactEvaluation: u.convertToResponse(impactErr.Evaluation),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	// Async sync events after transaction commit
	go u.triggerSyncEvents(context.Background(), result.LineageID)

	return result, nil
}

// SubmitReport gửi báo cáo/đề xuất thay đổi
func (u *PropertyUsecase) SubmitReport(
	ctx context.Context,
	cmd *dto.SubmitReportDTO,
	reporterOriginID uint64,
) (*dto.SubmitReportResult, error) {
	var result *dto.SubmitReportResult

	err := u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// Load original lineage
		lineage, err := u.loadLineage(txCtx, cmd.LineageID)
		if err != nil {
			return err
		}

		// Check for duplicate report
		if err := u.checkDuplicateReport(txCtx, lineage.ID, reporterOriginID); err != nil {
			return err
		}

		// Create propose lineage if there's contribute data
		proposeLineageID, err := u.createProposeIfNeeded(txCtx, lineage, cmd.ContributeData)
		if err != nil {
			return err
		}

		// Create report record
		report, err := u.createReportRecord(txCtx, lineage, reporterOriginID, cmd, proposeLineageID)
		if err != nil {
			return err
		}

		result = u.buildReportResult(report, lineage.ID, proposeLineageID, cmd)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Async sync after transaction
	go u.triggerReportSync(context.Background(), result.LineageID)

	return result, nil
}

// =============================================================================
// PRIVATE METHODS - LOAD & VALIDATION
// =============================================================================

func (u *PropertyUsecase) loadLineage(ctx context.Context, lineageID uint64) (*domain.PropertyLineage, error) {
	lineage, err := u.LineageRepo.GetByID(ctx, lineageID)
	if err != nil {
		return nil, err
	}
	if lineage == nil {
		return nil, fmt.Errorf("property lineage id=%d not found", lineageID)
	}
	return lineage, nil
}

func (u *PropertyUsecase) checkPermission(ctx context.Context, lineage *domain.PropertyLineage, userID uint64) (enums.ActorRole, error) {
	propertyUser, err := u.PropertyUserRepo.GetByLineageAndOwner(ctx, lineage.ID, userID)
	if err != nil {
		return enums.ActorRoleCollaborator, err
	}
	if propertyUser == nil {
		return enums.ActorRoleCollaborator, _errors.ReturnError(service.PropertyUpdateDenied)
	}
	return u.determineActorRole(propertyUser, userID), nil
}

func (u *PropertyUsecase) determineActorRole(propertyUser *domain.PropertyUser, userID uint64) enums.ActorRole {
	if propertyUser != nil && propertyUser.OwnerOriginID == userID {
		return enums.ActorRoleOwner
	}
	// TODO: Check admin role from permission service
	return enums.ActorRoleCollaborator
}

func (u *PropertyUsecase) evaluateImpact(
	ctx context.Context,
	lineage *domain.PropertyLineage,
	c *dto.UpdatePropertyDTO,
	userID uint64,
	actorRole enums.ActorRole,
) (*modules.ImpactEvaluation, error) {
	evalReq := &modules.EvaluateRequest{
		PropertyID:      lineage.ID,
		CurrentVersion:  lineage,
		ProposedChanges: c,
		ActorID:         userID,
		ActorRole:       actorRole,
		PreviewMode:     c.ImpactAssessment != nil && c.ImpactAssessment.PreviewOnly,
		Confirmed:       c.ImpactAssessment != nil && c.ImpactAssessment.Confirmed,
	}

	evaluation, err := u.impactEngine.Evaluate(ctx, evalReq)
	if err != nil {
		return nil, fmt.Errorf("impact evaluation failed: %w", err)
	}
	return evaluation, nil
}

func (u *PropertyUsecase) applyActionPolicy(
	evaluation *modules.ImpactEvaluation,
	assessment *dto.ImpactAssessmentDTO,
	actorRole enums.ActorRole,
) error {
	switch evaluation.ActionPolicy {
	case enums.ActionPolicyBlock:
		return fmt.Errorf("%s", evaluation.BlockReason)

	case enums.ActionPolicyRequireConfirm:
		if assessment == nil || assessment.PreviewOnly {
			return &ImpactRequiredError{Evaluation: evaluation}
		}
		if !assessment.Confirmed {
			return &ImpactRequiredError{Evaluation: evaluation}
		}

	case enums.ActionPolicyRequireReview:
		if !actorRole.HasFullAccess() {
			return fmt.Errorf("Thay đổi này cần được admin xem xét trước khi áp dụng")
		}
		slog.Info("admin override for critical property change")
	}

	return nil
}

func (u *PropertyUsecase) validateLocationUpdate(
	ctx context.Context,
	lineage *domain.PropertyLineage,
	cmd *dto.UpdatePropertyDTO,
) error {
	logger := logging.FromContext(ctx)

	if cmd.Location == nil {
		return nil
	}

	// Validate ward belongs to province
	if lineage.LocationID != nil {
		existing, err := u.LocationRepo.GetByID(ctx, *lineage.LocationID)
		if err == nil && existing != nil && existing.ProvinceID != nil {
			if cmd.Location.WardID != nil {
				ward, err := u.WardRepo.GetByID(ctx, *cmd.Location.WardID)
				if err == nil && ward != nil {
					if ward.ProvinceID != *existing.ProvinceID {
						return fmt.Errorf("ward %d does not belong to current province %d",
							*cmd.Location.WardID, *existing.ProvinceID)
					}
				}
			}
		}
	}

	// Warning if province_id provided but not in mask
	if cmd.Location.ProvinceID != nil && !cmd.Mask.Allows("location.province_id") {
		logger.Warn(
			"province ID ignored because field mask excludes it",
		)
	}

	return nil
}

// =============================================================================
// PRIVATE METHODS - UPDATE EXECUTION
// =============================================================================

func (u *PropertyUsecase) executeUpdate(
	ctx context.Context,
	lineage *domain.PropertyLineage,
	cmd *dto.UpdatePropertyDTO,
	userID uint64,
) (*dto.UpdatePropertyResult, error) {
	logger := logging.FromContext(ctx)

	fkPatch := make(FKPatch)

	// Handle side blocks (media, amenities, etc.)
	if err := u.handleSideBlocks(ctx, lineage.ID, cmd, fkPatch); err != nil {
		return nil, err
	}

	// Auto-set avatar from media if not provided
	if err := u.setAvatarFromMedia(ctx, lineage.ID, cmd); err != nil {
		logger.Warn(
			"set property avatar from media failed",
			slog.Any("error", err),
		)
	}

	// Handle location resolution from TQD (async with timeout)
	u.handleLocationResolution(ctx, lineage, cmd, userID)

	// Handle FK blocks (info, location, land_info, etc.)
	if err := u.handleFKBlocks(ctx, lineage, cmd, fkPatch); err != nil {
		return nil, err
	}

	// Update lineage with new FK values
	if len(fkPatch) > 0 {
		if err := u.LineageRepo.UpdateFields(ctx, lineage.ID, fkPatch); err != nil {
			return nil, err
		}
	}

	return &dto.UpdatePropertyResult{
		LineageID: lineage.ID,
		UpdatedAt: _utils.TimeNowRFC3339(),
	}, nil
}

func (u *PropertyUsecase) handleFKBlocks(
	ctx context.Context,
	lineage *domain.PropertyLineage,
	cmd *dto.UpdatePropertyDTO,
	fkPatch FKPatch,
) error {
	handlers := u.buildHandlers(lineage, cmd)
	for _, entry := range handlers {
		newFK, err := RunDirect(ctx, entry.currentFK, entry.h)
		if err != nil {
			return err
		}
		fkPatch.Set(entry.h.FKColumn(), newFK)
	}
	return nil
}

func (u *PropertyUsecase) buildHandlers(lineage *domain.PropertyLineage, cmd *dto.UpdatePropertyDTO) []struct {
	currentFK *uint64
	h         Handler
} {
	return []struct {
		currentFK *uint64
		h         Handler
	}{
		{lineage.PropertyInfoID, &InfoHandler{Patch: cmd.Info, Mask: cmd.Mask, Repo: u.InfoRepo}},
		{lineage.LocationID, &LocationHandler{Patch: cmd.Location, Mask: cmd.Mask, Repo: u.LocationRepo}},
		{lineage.LandInfoID, &LandInfoHandler{Patch: cmd.LandInfo, Mask: cmd.Mask, Repo: u.LandInfoRepo}},
		{lineage.BuildingInfoID, &BuildingHandler{Patch: cmd.BuildingInfo, Mask: cmd.Mask, Repo: u.BuildingInfoRepo}},
		{lineage.EdvidenceID, &EvidenceHandler{Patch: cmd.Evidence, Mask: cmd.Mask, Repo: u.EvidenceRepo}},
	}
}

func (u *PropertyUsecase) handleLocationResolution(
	ctx context.Context,
	lineage *domain.PropertyLineage,
	cmd *dto.UpdatePropertyDTO,
	userID uint64,
) {
	logger := logging.FromContext(ctx)

	if !u.needResolveLocation(cmd) {
		return
	}

	var currentLocation *domain.PropertyLocation
	if lineage.LocationID != nil {
		currentLocation, _ = u.LocationRepo.GetByID(ctx, *lineage.LocationID)
	}

	ch := make(chan *ResolveLocationResult, 1)
	resolve := &ResolveLocationFromTQD{
		Latitude:  *cmd.Location.Latitude,
		Longitude: *cmd.Location.Longitude,
		Mask:      cmd.Mask,
		Result:    ch,
		userID:    userID,
	}

	go u.resolveLocationFromTQD(ctx, resolve)

	select {
	case result := <-ch:
		if result.Error == nil {
			u.applyResolveResult(cmd, result, currentLocation)
		}
	case <-time.After(3 * time.Second):
		logger.Warn(
			"property location resolution timed out",
		)
	}
}

func (u *PropertyUsecase) needResolveLocation(cmd *dto.UpdatePropertyDTO) bool {
	if cmd.Location == nil {
		return false
	}
	if cmd.Location.Latitude == nil || cmd.Location.Longitude == nil {
		return false
	}
	if !cmd.Mask.Allows("location.latitude") && !cmd.Mask.Allows("location.longitude") {
		return false
	}
	if cmd.Location.ProvinceID != nil {
		return false
	}
	return true
}

func (u *PropertyUsecase) resolveLocationFromTQD(
	ctx context.Context,
	resolve *ResolveLocationFromTQD,
) {
	logger := logging.FromContext(ctx)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := &ResolveLocationResult{}

	defer func() {
		if r := recover(); r != nil {
			logger.Error(
				"panic while resolving property location",
				slog.Any("panic", r),
			)
			result.Error = fmt.Errorf("internal panic: %v", r)
		}
		select {
		case resolve.Result <- result:
		default:
			logger.Warn(
				"property location result channel full or closed",
			)
		}
		close(resolve.Result)
	}()

	resp, err := u.tqdGRPCClient.GetNearestLocation(ctx, resolve.Latitude, resolve.Longitude)
	if err != nil {
		result.Error = fmt.Errorf("tqd error: %w", err)
		return
	}

	if resp == nil || resp.Province == nil {
		logger.Warn(
			"TQD returned no province for property location",
			slog.Float64("latitude", resolve.Latitude),
			slog.Float64("longitude", resolve.Longitude),
		)
		return
	}

	province, err := u.findOrCreateProvince(ctx, resp.Province, resolve.userID)
	if err != nil {
		result.Error = fmt.Errorf("province error: %w", err)
		return
	}

	if resolve.Mask.Allows("location.province_id") {
		result.ProvinceID = &province.ID
	}

	if resp.Ward != nil && resolve.Mask.Allows("location.ward_id") {
		ward, err := u.findOrCreateWard(ctx, resp.Ward, province.ID, resolve.userID)
		if err == nil && ward != nil {
			result.WardID = &ward.ID
		} else if err != nil {
			logger.Warn(
				"resolve property ward failed",
				slog.Any("error", err),
			)
		}
	}
}

func (u *PropertyUsecase) applyResolveResult(
	cmd *dto.UpdatePropertyDTO,
	result *ResolveLocationResult,
	currentLocation *domain.PropertyLocation,
) {
	if result.ProvinceID != nil && cmd.Mask.Allows("location.province_id") {
		if currentLocation == nil ||
			currentLocation.ProvinceID == nil ||
			*result.ProvinceID != *currentLocation.ProvinceID {
			cmd.Location.ProvinceID = result.ProvinceID
			cmd.Location.ResolvedFromTQD = true
			slog.Info(
				"property province resolved",
				slog.Uint64("province_id", *result.ProvinceID),
			)
		}
	}

	if result.WardID != nil && cmd.Mask.Allows("location.ward_id") {
		if currentLocation == nil ||
			currentLocation.WardID == nil ||
			*result.WardID != *currentLocation.WardID {
			cmd.Location.WardID = result.WardID
			slog.Info(
				"property ward resolved",
				slog.Uint64("ward_id", *result.WardID),
			)
		}
	}
}

func (u *PropertyUsecase) findOrCreateProvince(ctx context.Context, p *tqdpb.Province, userID uint64) (*domain.ProvinceV2, error) {
	logger := logging.FromContext(ctx)

	logger.Info(
		"creating province from TQD location",
		slog.String("tqd_id", p.Id),
		slog.String("code", p.Code),
		slog.String("full_name", p.FullName),
		slog.String("short_name", p.ShortName),
		slog.String("division_type", p.DivisionType),
		slog.String("phone_code", p.PhoneCode),
		slog.Float64("latitude", p.Latitude),
		slog.Float64("longitude", p.Longitude),
	)
	logger.Info(
		"province creation requested by user",
		slog.Uint64("user_id", userID),
	)
	province, err := u.ProvinceRepo.GetByTQDID(ctx, p.Id)
	if err == nil && province != nil {
		return province, nil
	}

	newProvince := &domain.ProvinceV2{
		// BaseEntity: _models.BaseEntity{
		// 	AuditBase: _models.AuditBase{
		// 		CreatedBy: &userID,
		// 		UpdatedBy: &userID,
		// 	},
		// },
		TQDID:        p.Id,
		Code:         p.Code,
		Name:         p.FullName,
		Codename:     p.ShortName,
		DivisionType: getFirstWord(p.DivisionType),
		PhoneCode:    p.PhoneCode,
		Lat:          p.Latitude,
		Lng:          p.Longitude,
	}

	logger.Info(
		"attempting to create province",
		slog.Any("tqd_id", newProvince.TQDID),
		slog.Any("code", newProvince.Code),
		slog.Any("name", newProvince.Name),
		slog.Any("codename", newProvince.Codename),
		slog.Any("division_type", newProvince.DivisionType),
		slog.Any("phone_code", newProvince.PhoneCode),
		slog.Any("lat", newProvince.Lat),
		slog.Any("lng", newProvince.Lng),
	)
	err = u.ProvinceRepo.Create(ctx, newProvince)
	if err == nil {
		logger.Info(
			"created province",
			slog.Any("province_id", newProvince.ID),
		)
		return newProvince, nil
	}

	if isDuplicateKeyError(err) {
		logger.Warn(
			"duplicate province detected; retrying lookup",
		)
		province, findErr := u.ProvinceRepo.GetByTQDID(ctx, p.Id)
		if findErr == nil && province != nil {
			return province, nil
		}
	}

	return nil, err
}

func (u *PropertyUsecase) findOrCreateWard(ctx context.Context, w *tqdpb.Ward, provinceID uint64, userID uint64) (*domain.WardV2, error) {
	logger := logging.FromContext(ctx)

	ward, err := u.WardRepo.GetByTQDID(ctx, w.Id)
	if err == nil && ward != nil {
		if ward.ProvinceID == provinceID {
			return ward, nil
		}
		return nil, fmt.Errorf("ward province mismatch")
	}

	newWard := &domain.WardV2{
		// BaseEntity: _models.BaseEntity{
		// 	AuditBase: _models.AuditBase{
		// 		CreatedBy: &userID,
		// 		UpdatedBy: &userID,
		// 	},
		// },
		TQDID:         w.Id,
		Code:          w.Code,
		Name:          w.FullName,
		Codename:      w.ShortName,
		DivisionType:  getFirstWord(w.DivisionType),
		ShortCodename: w.ShortName,
		Lat:           w.Latitude,
		Lng:           w.Longitude,
		ProvinceID:    provinceID,
	}

	err = u.WardRepo.Create(ctx, newWard)
	if err == nil {
		logger.Info(
			"created ward",
			slog.Any("ward_id", newWard.ID),
		)
		return newWard, nil
	}

	if isDuplicateKeyError(err) {
		ward, findErr := u.WardRepo.GetByTQDID(ctx, w.Id)
		if findErr == nil && ward != nil && ward.ProvinceID == provinceID {
			return ward, nil
		}
	}

	return nil, err
}

// =============================================================================
// PRIVATE METHODS - SIDE BLOCKS
// =============================================================================

func (u *PropertyUsecase) handleSideBlocks(
	ctx context.Context,
	lineageID uint64,
	cmd *dto.UpdatePropertyDTO,
	fkPatch FKPatch,
) error {
	originID := _utils.GetOriginIdFromContext(ctx)
	if originID == 0 {
		return _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("origin_id required"), _errors.WithLegacyCode(401))
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
		u.applyLineageMeta(cmd.LineageMeta, cmd.Mask, fkPatch)
	}
	if cmd.Personalization != nil {
		if err := u.updatePersonalization(ctx, cmd.Mask, lineageID, originID, cmd.Personalization); err != nil {
			return fmt.Errorf("personalization: %w", err)
		}
	}
	return nil
}

func (u *PropertyUsecase) handleMedia(ctx context.Context, lineageID uint64, data *dto.UpdateMediaDTO) error {
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
			fields := map[string]any{
				"media_type": item.MediaType,
				"media_url":  item.MediaURL,
				"thumb_url":  item.ThumbURL,
				"sort_order": orderIdx + 1,
			}
			if err := u.MediaRepo.UpdateFields(ctx, item.ID, fields); err != nil {
				return fmt.Errorf("update media id=%d: %w", item.ID, err)
			}
		} else {
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

func (u *PropertyUsecase) handleAmenities(ctx context.Context, lineageID uint64, data *dto.UpdateAmenityDTO) error {
	return u.AmenityRepo.ReplaceAmenities(ctx, lineageID, data.AmenityIDs)
}

func (u *PropertyUsecase) handleAreaRegions(ctx context.Context, lineageID uint64, data *dto.UpdateAreaRegionsDTO) error {
	return u.AreaRegionRepo.ReplaceAreaRegions(ctx, lineageID, data.AreaRegionIDs)
}

func (u *PropertyUsecase) handleExternalRefs(ctx context.Context, data *dto.UpdateExternalRefDTO) error {
	for _, item := range data.Items {
		if item.ID > 0 {
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
		DateString("personalization.archived_at", "archived_at", data.ArchivedAt, nil).
		DateString("personalization.hidden_at", "hidden_at", data.HiddenAt, nil).
		Build()

	return u.PropertyUserRepo.UpdateFieldsByLineageAndOwner(ctx, lineageID, originID, fields)
}

func (u *PropertyUsecase) setAvatarFromMedia(ctx context.Context, lineageID uint64, cmd *dto.UpdatePropertyDTO) error {
	mediaList, err := u.MediaRepo.GetByLineageID(ctx, lineageID)
	if err != nil || len(mediaList) == 0 {
		return nil
	}

	var avatarID uint64
	for _, m := range mediaList {
		if m.MediaType == "image" {
			avatarID = m.ID
			break
		}
	}
	if avatarID == 0 {
		avatarID = mediaList[0].ID
	}

	if cmd.Info == nil {
		cmd.Info = &dto.UpdateInfoDTO{}
	}
	cmd.Info.AvatarID = &avatarID
	return nil
}

// =============================================================================
// PRIVATE METHODS - POST UPDATE
// =============================================================================

func (u *PropertyUsecase) postUpdateTasks(
	ctx context.Context,
	propertyID uint64,
	userID uint64,
	evaluation *modules.ImpactEvaluation,
	assessment *dto.ImpactAssessmentDTO,
) {
	logger := logging.FromContext(ctx)

	// Update entity states based on projected states
	if err := u.updateEntityStates(ctx, evaluation); err != nil {
		logger.Warn(
			"update property entity states failed",
			slog.Any("error", err),
		)
	}

	// Create new snapshots for affected entities
	if err := u.createAffectedSnapshots(ctx, propertyID, evaluation); err != nil {
		logger.Warn(
			"create property impact snapshots failed",
			slog.Any("error", err),
		)
	}

	// Create audit log
	if err := u.createAuditLog(ctx, propertyID, userID, evaluation, assessment); err != nil {
		logger.Warn(
			"create property audit log failed",
			slog.Any("error", err),
		)
	}
}

func (u *PropertyUsecase) updateEntityStates(ctx context.Context, evaluation *modules.ImpactEvaluation) error {
	logger := logging.FromContext(ctx)

	if evaluation == nil || len(evaluation.SnapshotComparisons) == 0 {
		return nil
	}

	for _, comp := range evaluation.SnapshotComparisons {
		switch comp.EntityType {
		case "product":
			if err := u.updateProductState(ctx, comp.EntityID, comp.ProjectedState); err != nil {
				logger.Warn(
					"update product state failed",
					slog.Any("product_id", comp.EntityID),
					slog.Any("error", err),
				)
			}
		case "listing":
			if err := u.updateListingState(ctx, comp.EntityID, comp.ProjectedState); err != nil {
				logger.Warn(
					"update listing state failed",
					slog.Any("listing_id", comp.EntityID),
					slog.Any("error", err),
				)
			}
		}
	}
	return nil
}

func (u *PropertyUsecase) updateProductState(ctx context.Context, productID uint64, state enums.EffectiveState) error {
	logger := logging.FromContext(ctx)

	// TODO: Call product service to update state
	logger.Info(
		"updating product state",
		slog.Uint64("product_id", productID),
		slog.String("state", state.String()),
	)
	return nil
}

func (u *PropertyUsecase) updateListingState(ctx context.Context, listingID uint64, state enums.EffectiveState) error {
	logger := logging.FromContext(ctx)

	// TODO: Call listing service to update state
	logger.Info(
		"updating listing state",
		slog.Uint64("listing_id", listingID),
		slog.String("state", state.String()),
	)
	return nil
}

func (u *PropertyUsecase) createAffectedSnapshots(ctx context.Context, propertyID uint64, evaluation *modules.ImpactEvaluation) error {
	logger := logging.FromContext(ctx)

	if evaluation == nil || len(evaluation.SnapshotComparisons) == 0 {
		return nil
	}

	lineage, err := u.LineageRepo.GetByID(ctx, propertyID)
	if err != nil {
		return err
	}

	for _, comp := range evaluation.SnapshotComparisons {
		snapshotData, err := u.buildSnapshotData(lineage, comp.EntityType)
		if err != nil {
			logger.Warn(
				"build property impact snapshot data failed",
				slog.String("entity_type", comp.EntityType),
				slog.Any("entity_id", comp.EntityID),
				slog.Any("error", err),
			)
			continue
		}

		snapshot := &domain.PropertyImpactVersion{
			EntityType:   comp.EntityType,
			EntityID:     comp.EntityID,
			PropertyID:   propertyID,
			SnapshotData: snapshotData,
			Version:      u.getNextVersion(ctx, comp.EntityType, comp.EntityID),
			ValidFrom:    time.Now(),
			CreatedBy:    u.getCurrentUserID(ctx),
		}

		if err := u.PropertyImpactRepo.Create(ctx, snapshot); err != nil {
			logger.Warn(
				"create property impact snapshot failed",
				slog.String("entity_type", comp.EntityType),
				slog.Any("entity_id", comp.EntityID),
				slog.Any("error", err),
			)
		}
	}

	return nil
}

func (u *PropertyUsecase) buildSnapshotData(lineage *domain.PropertyLineage, entityType string) (domain.JSONB, error) {
	data := make(map[string]interface{})

	if lineage.PropertyInfo != nil {
		data["title"] = lineage.PropertyInfo.Title
		data["property_type_id"] = lineage.PropertyInfo.PropertyTypeID
		data["legal_status"] = lineage.PropertyInfo.LegalStatus
	}

	if lineage.Location != nil {
		data["latitude"] = lineage.Location.Latitude
		data["longitude"] = lineage.Location.Longitude
		data["province_id"] = lineage.Location.ProvinceID
		data["ward_id"] = lineage.Location.WardID
	}

	if lineage.LandInfo != nil {
		data["area_total"] = lineage.LandInfo.AreaTotal
		data["plot"] = lineage.LandInfo.Plot
	}

	if lineage.BuildingInfo != nil {
		data["area_actual"] = lineage.BuildingInfo.AreaActual
		data["floors"] = lineage.BuildingInfo.Floors
	}

	return domain.MarshalJSONB(data)
}

func (u *PropertyUsecase) getNextVersion(ctx context.Context, entityType string, entityID uint64) uint32 {
	snapshot, err := u.PropertyImpactRepo.GetCurrent(ctx, entityType, entityID)
	if err != nil || snapshot == nil {
		return 1
	}
	return snapshot.Version + 1
}

func (u *PropertyUsecase) createAuditLog(
	ctx context.Context,
	propertyID, userID uint64,
	evaluation *modules.ImpactEvaluation,
	assessment *dto.ImpactAssessmentDTO,
) error {
	// Tạo map changes
	changesMap := make(domain.AuditChangesJSON)
	changesMap["changed_fields"] = len(evaluation.ChangeDetection.ChangedFields)
	changesMap["impact_level"] = evaluation.ImpactLevel.String()
	changesMap["action_policy"] = evaluation.ActionPolicy.String()

	// Thêm thông tin chi tiết các field thay đổi nếu có
	if len(evaluation.ChangeDetection.ChangedFields) > 0 {
		fieldChanges := make([]map[string]interface{}, 0, len(evaluation.ChangeDetection.ChangedFields))
		for _, change := range evaluation.ChangeDetection.ChangedFields {
			fieldChanges = append(fieldChanges, map[string]interface{}{
				"field":       change.FieldPath,
				"old_value":   change.OldValue,
				"new_value":   change.NewValue,
				"change_type": change.ChangeType,
			})
		}
		changesMap["field_changes"] = fieldChanges
	}

	// Thêm thông tin affected entities nếu có
	if evaluation.AffectedSummary != nil {
		changesMap["affected_summary"] = map[string]interface{}{
			"products":  evaluation.AffectedSummary.Products,
			"listings":  evaluation.AffectedSummary.Listings,
			"assets":    evaluation.AffectedSummary.Assets,
			"deals":     evaluation.AffectedSummary.Deals,
			"crm_notes": evaluation.AffectedSummary.CrmNotes,
			"total":     evaluation.AffectedSummary.Total,
		}
	}

	// Thêm thông tin projected states nếu có
	if evaluation.ProjectedStates != nil {
		changesMap["projected_states"] = map[string]interface{}{
			"active":     evaluation.ProjectedStates.Active,
			"restricted": evaluation.ProjectedStates.Restricted,
			"frozen":     evaluation.ProjectedStates.Frozen,
			"archived":   evaluation.ProjectedStates.Archived,
		}
	}

	auditLog := &domain.PropertyAuditLog{
		PropertyID:       propertyID,
		UserID:           userID,
		Action:           "UPDATE",
		ImpactLevel:      int32(evaluation.ImpactLevel),
		AffectedProducts: evaluation.AffectedSummary.Products,
		AffectedListings: evaluation.AffectedSummary.Listings,
		AffectedAssets:   evaluation.AffectedSummary.Assets,
		Changes:          changesMap, // Đây là domain.AuditChangesJSON, không phải []byte
		CreatedAt:        time.Now(),
	}

	if assessment != nil {
		if assessment.Confirmed {
			now := time.Now()
			auditLog.ConfirmedAt = &now
		}
		auditLog.SessionID = assessment.SessionID
	}

	if evaluation.ActionPolicy == enums.ActionPolicyBlock {
		auditLog.BlockedReason = evaluation.BlockReason
	}

	return u.AuditLogRepo.Create(ctx, auditLog)
}

func (u *PropertyUsecase) getCurrentUserID(ctx context.Context) uint64 {
	return _utils.GetProfileIdWithContext(ctx)
}

// =============================================================================
// PRIVATE METHODS - PROPOSE & REPORT
// =============================================================================

func (u *PropertyUsecase) proposeChange(
	ctx context.Context,
	origin *domain.PropertyLineage,
	cmd *dto.UpdatePropertyDTO,
) (*dto.UpdatePropertyResult, error) {
	// Create propose lineage with copied FKs from origin
	propose := &domain.PropertyLineage{
		PropertyInfoID: origin.PropertyInfoID,
		LocationID:     origin.LocationID,
		LandInfoID:     origin.LandInfoID,
		BuildingInfoID: origin.BuildingInfoID,
		EdvidenceID:    origin.EdvidenceID,
	}
	if err := u.LineageRepo.Create(ctx, propose); err != nil {
		return nil, fmt.Errorf("create propose lineage: %w", err)
	}

	fkPatch := make(FKPatch)

	// Always insert new records for propose (no update)
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

func (u *PropertyUsecase) checkDuplicateReport(ctx context.Context, lineageID, reporterID uint64) error {
	existing, err := u.ReportRepo.GetByLineageAndReporter(ctx, lineageID, reporterID)
	if err != nil {
		return err
	}
	if existing != nil {
		return _errors.ReturnError(service.PropertyReportAlreadyPending)
	}
	return nil
}

func (u *PropertyUsecase) createProposeIfNeeded(
	ctx context.Context,
	lineage *domain.PropertyLineage,
	contributeData *dto.UpdatePropertyDTO,
) (*uint64, error) {
	if contributeData == nil {
		return nil, nil
	}
	proposeResult, err := u.proposeChange(ctx, lineage, contributeData)
	if err != nil {
		return nil, fmt.Errorf("propose change: %w", err)
	}
	return &proposeResult.LineageID, nil
}

func (u *PropertyUsecase) createReportRecord(
	ctx context.Context,
	lineage *domain.PropertyLineage,
	reporterOriginID uint64,
	cmd *dto.SubmitReportDTO,
	proposeLineageID *uint64,
) (*domain.PropertyReport, error) {
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
	if err := u.ReportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("create report: %w", err)
	}
	return report, nil
}

func (u *PropertyUsecase) buildReportResult(
	report *domain.PropertyReport,
	lineageID uint64,
	proposeLineageID *uint64,
	cmd *dto.SubmitReportDTO,
) *dto.SubmitReportResult {
	return &dto.SubmitReportResult{
		ReportID:         report.ID,
		LineageID:        lineageID,
		ProposeLineageID: proposeLineageID,
		ReportType:       cmd.ReportType,
		Note:             derefStr(cmd.Note),
		CreatedAt:        _utils.TimeNowRFC3339(),
	}
}

func (u *PropertyUsecase) triggerReportSync(ctx context.Context, lineageID uint64) {
	cloneCtx := _utils.CloneContext(ctx)
	now := time.Now().UnixMilli()
	u.SyncProvider.PutTimestamp(cloneCtx,
		u.SyncProvider.GetKey(cloneCtx, _utils.SyncKeyPropertyDetail, lineageID), now)
}

// =============================================================================
// PRIVATE METHODS - SYNC & UTILITIES
// =============================================================================

func (u *PropertyUsecase) triggerSyncEvents(ctx context.Context, lineageID uint64) {
	cloneCtx := _utils.CloneContext(ctx)
	profileID := _utils.GetProfileIdWithContext(cloneCtx)
	now := time.Now().UnixMilli()

	u.SyncProvider.PutTimestamp(cloneCtx,
		u.SyncProvider.GetKey(cloneCtx, _utils.SyncKeyPropertyMe, profileID), now)
	u.SyncProvider.PutTimestamp(cloneCtx,
		u.SyncProvider.GetKey(cloneCtx, _utils.SyncKeyPropertyDetail, lineageID), now)
}

func (u *PropertyUsecase) convertToResponse(evaluation *modules.ImpactEvaluation) *dto.ImpactEvaluationResponse {
	if evaluation == nil {
		return nil
	}
	return &dto.ImpactEvaluationResponse{
		ImpactLevel:     evaluation.ImpactLevel.String(),
		ActionPolicy:    evaluation.ActionPolicy.String(),
		RequiresConfirm: evaluation.RequiresConfirm,
		BlockReason:     evaluation.BlockReason,
		Warnings:        evaluation.Warnings,
		Suggestions:     evaluation.Suggestions,
		AffectedSummary: &dto.AffectedSummaryResponse{
			Products: evaluation.AffectedSummary.Products,
			Listings: evaluation.AffectedSummary.Listings,
			Assets:   evaluation.AffectedSummary.Assets,
			Deals:    evaluation.AffectedSummary.Deals,
			CrmNotes: evaluation.AffectedSummary.CrmNotes,
			Total:    evaluation.AffectedSummary.Total,
		},
		ProjectedStates: &dto.ProjectedStatesResponse{
			Active:     evaluation.ProjectedStates.Active,
			Restricted: evaluation.ProjectedStates.Restricted,
			Frozen:     evaluation.ProjectedStates.Frozen,
			Archived:   evaluation.ProjectedStates.Archived,
		},
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func getFirstWord(input string) string {
	parts := strings.SplitN(input, " ", 2)
	return parts[0]
}

// ResolveLocationResult kết quả resolve location
type ResolveLocationResult struct {
	ProvinceID *uint64
	WardID     *uint64
	Error      error
}

// ResolveLocationFromTQD request resolve location
type ResolveLocationFromTQD struct {
	Latitude  float64
	Longitude float64
	Mask      fieldmask.FieldMask
	Result    chan<- *ResolveLocationResult
	userID    uint64
}
