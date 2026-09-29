package usecase

import (
	"context"

	_errors "common/errors"
	tqdpb "pb/types/tqd"
	"tqd/internal"
	"tqd/internal/domain"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"

	"gorm.io/datatypes"
)

type SubscriptionUsecase interface {
	// Parcel
	ListParcelSubscriptions(ctx context.Context, userID uint64, status *uint32, page, limit int) ([]domain.UserSubscription, int64, error)
	GetParcelStatus(ctx context.Context, userID uint64, parcelID uint64) (bool, *domain.UserSubscription, error)
	CreateParcelSubscription(ctx context.Context, userID uint64, parcelID uint64, scope *tqdpb.SubscriptionScopeMessage, triggerConfig string) (*domain.UserSubscription, error)
	DeleteParcelSubscription(ctx context.Context, userID uint64, parcelID uint64) error

	// Region
	CreateRegionSubscription(ctx context.Context, userID uint64, regionID uint64, scope *tqdpb.SubscriptionScopeMessage, triggerConfig string) (*domain.UserSubscription, error)
	DeleteRegionSubscription(ctx context.Context, userID uint64, subscriptionID uint64) error
	ListUserRegionSubscriptions(ctx context.Context, userID uint64, page, limit int) ([]domain.UserSubscription, []qh_domain.QHRegion, int64, error)

	// General
	ListUserSubscriptions(ctx context.Context, userID uint64, targetType *string, status *string, page, limit int) ([]domain.UserSubscription, int64, error)

	// Admin
	AdminList(ctx context.Context, userID *uint64, targetType *string, status *uint32, page, limit int) ([]domain.UserSubscription, int64, error)
	AdminDelete(ctx context.Context, subscriptionID uint64) error
}

type subscriptionUsecase struct {
	subRepo    repo.SubscriptionRepository
	regionRepo repo.RegionRepository
	parcelRepo repo.IParcelRepo
}

func NewSubscriptionUsecase(
	subRepo repo.SubscriptionRepository,
	regionRepo repo.RegionRepository,
	parcelRepo repo.IParcelRepo,
) SubscriptionUsecase {
	return &subscriptionUsecase{
		subRepo:    subRepo,
		regionRepo: regionRepo,
		parcelRepo: parcelRepo,
	}
}

// ==================== Parcel ====================
func (u *subscriptionUsecase) ListParcelSubscriptions(ctx context.Context, userID uint64, status *uint32, page, limit int) ([]domain.UserSubscription, int64, error) {
	return u.subRepo.ListByUserWithStatus(ctx, userID, status, page, limit)
}

func (u *subscriptionUsecase) GetParcelStatus(ctx context.Context, userID uint64, parcelID uint64) (bool, *domain.UserSubscription, error) {
	sub, err := u.subRepo.GetByUserAndTarget(ctx, userID, "parcel", parcelID)
	if err != nil || sub == nil {
		return false, nil, err
	}
	return true, sub, nil
}

func (u *subscriptionUsecase) CreateParcelSubscription(ctx context.Context, userID uint64, parcelID uint64, scope *tqdpb.SubscriptionScopeMessage, triggerConfig string) (*domain.UserSubscription, error) {
	// check parcel exists
	_, err := u.parcelRepo.GetByID(ctx, parcelID)
	if err != nil {
		return nil, _errors.ReturnError(service.ParcelNotFound)
	}
	// check existing subscription
	existing, _ := u.subRepo.GetByUserAndTarget(ctx, userID, "parcel", parcelID)
	if existing != nil {
		return nil, _errors.ReturnError(service.SubscriptionAlreadyExists)
	}
	// convert scope proto to datatypes.JSON
	scopeJSON, err := dto.ProtoToSubscriptionScopeJSON(scope)
	if err != nil {
		return nil, _errors.ReturnError(service.SubscriptionScopeInvalid)
	}
	// default scope if empty
	if scope == nil || len(scopeJSON) == 0 {
		defaultScope := &tqdpb.SubscriptionScopeMessage{
			Layers: &tqdpb.LayerScopeMessage{
				WatchAllLayers:           true,
				IncludeNewLayerAvailable: true,
			},
			Legal: &tqdpb.LegalScopeMessage{
				WatchStatusChange:       true,
				PriorityEffectiveChange: true,
			},
			Version: &tqdpb.VersionScopeMessage{
				WatchVersionChange: true,
			},
		}
		scopeJSON, _ = dto.ProtoToSubscriptionScopeJSON(defaultScope)
	}
	// convert trigger config
	var triggerJSON datatypes.JSON
	if triggerConfig != "" {
		triggerJSON = datatypes.JSON(triggerConfig)
	}
	sub := &domain.UserSubscription{
		UserID:            userID,
		TargetType:        "parcel",
		TargetID:          parcelID,
		SubscriptionScope: scopeJSON,
		TriggerConfig:     triggerJSON,
		Status:            enums.SubStatusActive,
	}
	if err := u.subRepo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (u *subscriptionUsecase) DeleteParcelSubscription(ctx context.Context, userID uint64, parcelID uint64) error {
	sub, err := u.subRepo.GetByUserAndTarget(ctx, userID, "parcel", parcelID)
	if err != nil || sub == nil {
		return _errors.ReturnError(service.SubscriptionNotFound)
	}
	return u.subRepo.Delete(ctx, sub.ID)
}

// ==================== Region ====================

func (u *subscriptionUsecase) CreateRegionSubscription(ctx context.Context, userID uint64, regionID uint64, scope *tqdpb.SubscriptionScopeMessage, triggerConfig string) (*domain.UserSubscription, error) {
	// check region exists
	_, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return nil, _errors.ReturnError(service.RegionNotFound)
	}
	// check existing subscription
	existing, _ := u.subRepo.GetByUserAndTarget(ctx, userID, "region", regionID)
	if existing != nil {
		return nil, _errors.ReturnError(service.SubscriptionAlreadyExists)
	}
	// convert scope proto to datatypes.JSON
	scopeJSON, err := dto.ProtoToSubscriptionScopeJSON(scope)
	if err != nil {
		return nil, _errors.ReturnError(service.SubscriptionScopeInvalid)
	}
	// default scope if empty
	if scope == nil || len(scopeJSON) == 0 {
		defaultScope := &tqdpb.SubscriptionScopeMessage{
			Layers: &tqdpb.LayerScopeMessage{
				WatchAllLayers:           true,
				IncludeNewLayerAvailable: true,
			},
			Legal: &tqdpb.LegalScopeMessage{
				WatchStatusChange:       true,
				PriorityEffectiveChange: true,
			},
		}
		scopeJSON, _ = dto.ProtoToSubscriptionScopeJSON(defaultScope)
	}
	// convert trigger config
	var triggerJSON datatypes.JSON
	if triggerConfig != "" {
		triggerJSON = datatypes.JSON(triggerConfig)
	}
	sub := &domain.UserSubscription{
		UserID:            userID,
		TargetType:        "region",
		TargetID:          regionID,
		SubscriptionScope: scopeJSON,
		TriggerConfig:     triggerJSON,
		Status:            enums.SubStatusActive,
	}
	if err := u.subRepo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (u *subscriptionUsecase) DeleteRegionSubscription(ctx context.Context, userID uint64, subscriptionID uint64) error {
	sub, err := u.subRepo.GetByID(ctx, subscriptionID)
	if err != nil || sub == nil {
		return _errors.ReturnError(service.SubscriptionNotFound)
	}
	if sub.UserID != userID {
		return _errors.ReturnError(service.SubscriptionPermissionDenied)
	}
	return u.subRepo.Delete(ctx, subscriptionID)
}

func (u *subscriptionUsecase) ListUserRegionSubscriptions(ctx context.Context, userID uint64, page, limit int) ([]domain.UserSubscription, []qh_domain.QHRegion, int64, error) {
	targetType := "region"
	subs, total, err := u.subRepo.ListByUser(ctx, userID, &targetType, nil, page, limit)
	if err != nil {
		return nil, nil, 0, err
	}
	regionIDs := make([]uint64, 0, len(subs))
	for _, sub := range subs {
		id := sub.TargetID
		regionIDs = append(regionIDs, id)
	}
	regions, err := u.regionRepo.GetByIDs(ctx, regionIDs)
	if err != nil {
		return subs, nil, total, err
	}
	regionMap := make(map[uint64]qh_domain.QHRegion)
	for _, r := range regions {
		regionMap[r.ID] = r
	}
	orderedRegions := make([]qh_domain.QHRegion, 0, len(subs))
	for _, sub := range subs {
		id := sub.TargetID
		if r, ok := regionMap[id]; ok {
			orderedRegions = append(orderedRegions, r)
		}
	}
	return subs, orderedRegions, total, nil
}

func (u *subscriptionUsecase) ListUserSubscriptions(ctx context.Context, userID uint64, targetType *string, status *string, page, limit int) ([]domain.UserSubscription, int64, error) {
	return u.subRepo.ListByUser(ctx, userID, targetType, status, page, limit)
}

// ==================== Admin ====================

func (u *subscriptionUsecase) AdminList(ctx context.Context, userID *uint64, targetType *string, status *uint32, page, limit int) ([]domain.UserSubscription, int64, error) {
	return u.subRepo.AdminList(ctx, userID, targetType, status, page, limit)
}

func (u *subscriptionUsecase) AdminDelete(ctx context.Context, subscriptionID uint64) error {
	return u.subRepo.AdminDelete(ctx, subscriptionID)
}
