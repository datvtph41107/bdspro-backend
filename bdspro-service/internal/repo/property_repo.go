package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud3"
	_dto "common/domain/dto"
	"context"
)

type PropertyRepo interface {
	crud3.CrudRepo[domain.Property]
	Search(ctx context.Context, dto *dto.PropertySearchDTO, userId uint64) ([]*domain.Property, int64, error)
	GetDetail(ctx context.Context, id uint64, sourceType uint32) (*domain.Property, error)
	UpdateOption(ctx context.Context, id uint64, property *domain.Property) error
	GetOneByID(c context.Context, id uint64) (*domain.Property, error)
	DeleteBatchProperties(ctx context.Context, ids []uint64) error
	ArchiveBatchProperties(ctx context.Context, ids []uint64, archived bool) error
	GetRefreshSyncIds(c context.Context, req *dto.CheckVersionSyncRequest) ([]uint64, error)
	GetUpdatedAt(c context.Context, id uint64) (int64, error)

	// Property User
}

type PropertyRelationRepo interface {
	Create(ctx context.Context, relation *domain.PropertyRelation) error
	Update(ctx context.Context, relation *domain.PropertyRelation) error
	GetByPropertyID(ctx context.Context, propertyID uint64) (*domain.PropertyRelation, error)
	GetAllByPropertyID(ctx context.Context, propertyID uint64) ([]domain.PropertyRelation, error)
	GetByAssetID(ctx context.Context, assetID uint64) (*domain.PropertyRelation, error)
	GetByProductID(ctx context.Context, productID uint64) (*domain.PropertyRelation, error)
	GetCountsByPropertyID(ctx context.Context, propertyID uint64) (productCount uint32, assetCount uint32, postCount uint32, err error)
}

type PropertyProductRepo interface {
	Create(ctx context.Context, propertyProduct *domain.PropertyProduct) error
	GetByProductID(ctx context.Context, productID uint64) ([]domain.PropertyProduct, error)
	GetByPropertyID(ctx context.Context, propertyID uint64) ([]domain.PropertyProduct, error)
}

type PropertyLandInfoRepo interface {
	Create(ctx context.Context, landInfo *domain.PropertyLandInfo) error
	Update(ctx context.Context, landInfo *domain.PropertyLandInfo) error
	GetByPropertyID(ctx context.Context, propertyID uint64) (*domain.PropertyLandInfo, error)
}

type PropertyBuildingInfoRepo interface {
	Create(ctx context.Context, buildingInfo *domain.PropertyBuildingInfo) error
	Update(ctx context.Context, buildingInfo *domain.PropertyBuildingInfo) error
	GetByPropertyID(ctx context.Context, propertyID uint64) (*domain.PropertyBuildingInfo, error)
}

type PropertyMediaRepo interface {
	Create(ctx context.Context, media *domain.PropertyMedia) error
	CreateBatch(ctx context.Context, mediaList []domain.PropertyMedia) error
	GetByPropertyID(ctx context.Context, propertyID uint64) ([]domain.PropertyMedia, error)
	DeleteByPropertyID(ctx context.Context, propertyID uint64) error
}

type PropertyAmenityRepo interface {
	GetAll(ctx context.Context, search *dto.AmenitySearchDTO) ([]domain.AmenityItem, error)
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) ([]_dto.ItemDTO, error)

	GetAmenityIDsByProperty(ctx context.Context, propertyID uint64) ([]uint64, error)

	UpdatePropertyAmenityIds(ctx context.Context, propertyID uint64, amenityIDs []uint64) error
	GetTagIDsByProperty(ctx context.Context, propertyID uint64) ([]uint64, error)
	ReplacePropertyTags(ctx context.Context, propertyID uint64, tagIDs []uint64) error
	Search(ctx context.Context, search *dto.TagSearchDTO) ([]*domain.Tag, int64, error)
	FindByIDs(ctx context.Context, ids []uint64) ([]*domain.Tag, error)
}
