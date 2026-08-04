package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	_db "common/db"
	"common/productcatalog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"user/internal/catalogdomain"
	"user/internal/catalogpublication"
	"user/internal/catalogshadow"
	"user/internal/models"
	"user/internal/subscriptiondomain"
	"user/internal/subscriptionshadow"
)

/** CatalogRepo persists catalog v2 aggregates. */
type CatalogRepo struct {
	transactionRepo *_db.TransactionRepo
}

var (
	_ catalogpublication.Repository = (*CatalogRepo)(nil)
	_ catalogshadow.Repository      = (*CatalogRepo)(nil)
	_ subscriptionshadow.Repository = (*CatalogRepo)(nil)
)

/** NewCatalogRepo creates the catalog v2 repository. */
func NewCatalogRepo(transactionRepo *_db.TransactionRepo) *CatalogRepo {
	return &CatalogRepo{transactionRepo: transactionRepo}
}

/** LoadPlanVersionForUpdate locks one draft aggregate. */
func (r *CatalogRepo) LoadPlanVersionForUpdate(
	ctx context.Context,
	planVersionID uint64,
) (catalogdomain.PlanVersionAggregate, error) {
	if r == nil || r.transactionRepo == nil {
		return catalogdomain.PlanVersionAggregate{}, fmt.Errorf(
			"catalog repository is not configured",
		)
	}
	return r.loadPlanVersion(ctx, planVersionID, true)
}

/** GetPlanVersion loads one catalog version for shadow reads. */
func (r *CatalogRepo) GetPlanVersion(
	ctx context.Context,
	planVersionID uint64,
) (catalogdomain.PlanVersionAggregate, error) {
	if r == nil || r.transactionRepo == nil {
		return catalogdomain.PlanVersionAggregate{}, fmt.Errorf(
			"catalog repository is not configured",
		)
	}
	aggregate, err := r.loadPlanVersion(ctx, planVersionID, false)
	if errors.Is(err, catalogpublication.ErrPlanVersionNotFound) {
		return catalogdomain.PlanVersionAggregate{}, catalogshadow.ErrNotFound
	}
	return aggregate, err
}

/** MarkPlanVersionPublished performs a guarded draft-to-active update. */
func (r *CatalogRepo) MarkPlanVersionPublished(
	ctx context.Context,
	input catalogpublication.PublishRecord,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("catalog repository is not configured")
	}
	parentUpdates := map[string]any{
		"status":     productcatalog.StatusActive,
		"updated_by": input.PublishedBy,
		"updated_at": gorm.Expr("NOW()"),
	}
	db := r.transactionRepo.GetDB(ctx)
	productResult := db.Model(&models.CatalogProduct{}).
		Where("id = ? AND status IN ?", input.ProductID, []productcatalog.Status{
			productcatalog.StatusDraft,
			productcatalog.StatusActive,
		}).
		Updates(parentUpdates)
	if productResult.Error != nil {
		return productResult.Error
	}
	if productResult.RowsAffected != 1 {
		return catalogpublication.ErrConcurrentPublish
	}
	planResult := db.Model(&models.CatalogPlan{}).
		Where("id = ? AND product_id = ? AND status IN ?", input.PlanID, input.ProductID, []productcatalog.Status{
			productcatalog.StatusDraft,
			productcatalog.StatusActive,
		}).
		Updates(parentUpdates)
	if planResult.Error != nil {
		return planResult.Error
	}
	if planResult.RowsAffected != 1 {
		return catalogpublication.ErrConcurrentPublish
	}

	updates := map[string]any{
		"status":          productcatalog.StatusActive,
		"effective_from":  input.EffectiveFrom,
		"effective_until": input.EffectiveUntil,
		"published_at":    input.PublishedAt,
		"published_by":    input.PublishedBy,
		"updated_by":      input.PublishedBy,
		"terms_checksum":  input.TermsChecksum,
		"updated_at":      gorm.Expr("NOW()"),
	}
	result := db.
		Model(&models.CatalogPlanVersion{}).
		Where(
			"id = ? AND status = ? AND published_at IS NULL",
			input.PlanVersionID,
			productcatalog.StatusDraft,
		).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return catalogpublication.ErrConcurrentPublish
	}
	return nil
}

/** GetLegacyMembershipMapping loads one membership migration mapping. */
func (r *CatalogRepo) GetLegacyMembershipMapping(
	ctx context.Context,
	legacyMembershipID uint64,
) (subscriptionshadow.Mapping, error) {
	if r == nil || r.transactionRepo == nil {
		return subscriptionshadow.Mapping{}, fmt.Errorf(
			"catalog repository is not configured",
		)
	}
	var entity models.LegacyMembershipSubscriptionMapping
	err := r.transactionRepo.GetDB(ctx).
		Where("legacy_membership_id = ?", legacyMembershipID).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return subscriptionshadow.Mapping{}, subscriptionshadow.ErrNotFound
	}
	if err != nil {
		return subscriptionshadow.Mapping{}, err
	}
	return subscriptionshadow.Mapping{
		LegacyMembershipID: entity.LegacyMembershipID,
		SubscriptionID:     copyUint64(entity.SubscriptionID),
		SourceChecksum:     entity.SourceChecksum,
		Status:             entity.Status,
	}, nil
}

/** GetLegacyPlanMappingForSubscription loads the approved plan mapping. */
func (r *CatalogRepo) GetLegacyPlanMappingForSubscription(
	ctx context.Context,
	legacyPlanID uint64,
) (subscriptionshadow.PlanMapping, error) {
	mapping, err := r.GetLegacyPlanMapping(ctx, legacyPlanID)
	if errors.Is(err, catalogshadow.ErrNotFound) {
		return subscriptionshadow.PlanMapping{}, subscriptionshadow.ErrNotFound
	}
	if err != nil {
		return subscriptionshadow.PlanMapping{}, err
	}
	return subscriptionshadow.PlanMapping{
		LegacyPlanID:  mapping.LegacyPlanID,
		PlanVersionID: copyUint64(mapping.PlanVersionID),
		Status:        mapping.Status,
	}, nil
}

/** GetSubscription loads one durable subscription for shadow reads. */
func (r *CatalogRepo) GetSubscription(
	ctx context.Context,
	subscriptionID uint64,
) (subscriptiondomain.Aggregate, error) {
	if r == nil || r.transactionRepo == nil {
		return subscriptiondomain.Aggregate{}, fmt.Errorf(
			"catalog repository is not configured",
		)
	}
	db := r.transactionRepo.GetDB(ctx)
	var subscription models.CatalogSubscription
	if err := db.First(&subscription, subscriptionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return subscriptiondomain.Aggregate{}, subscriptionshadow.ErrNotFound
		}
		return subscriptiondomain.Aggregate{}, err
	}
	var version models.CatalogPlanVersion
	if err := db.First(&version, subscription.PlanVersionID).Error; err != nil {
		return subscriptiondomain.Aggregate{}, err
	}
	var plan models.CatalogPlan
	if err := db.First(&plan, version.PlanID).Error; err != nil {
		return subscriptiondomain.Aggregate{}, err
	}
	if plan.ProductID != subscription.ProductID {
		return subscriptiondomain.Aggregate{}, fmt.Errorf(
			"subscription product does not match plan version",
		)
	}
	var product models.CatalogProduct
	if err := db.First(&product, subscription.ProductID).Error; err != nil {
		return subscriptiondomain.Aggregate{}, err
	}
	return toSubscriptionAggregate(subscription, product, plan, version), nil
}

/** GetLegacyPlanMapping loads one migration mapping. */
func (r *CatalogRepo) GetLegacyPlanMapping(
	ctx context.Context,
	legacyPlanID uint64,
) (catalogshadow.Mapping, error) {
	if r == nil || r.transactionRepo == nil {
		return catalogshadow.Mapping{}, fmt.Errorf(
			"catalog repository is not configured",
		)
	}
	var entity models.LegacyPlanCatalogMapping
	err := r.transactionRepo.GetDB(ctx).
		Where("legacy_plan_id = ?", legacyPlanID).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return catalogshadow.Mapping{}, catalogshadow.ErrNotFound
	}
	if err != nil {
		return catalogshadow.Mapping{}, err
	}
	return catalogshadow.Mapping{
		LegacyPlanID:   entity.LegacyPlanID,
		PlanVersionID:  copyUint64(entity.PlanVersionID),
		SourceChecksum: entity.SourceChecksum,
		Status:         entity.Status,
	}, nil
}

func (r *CatalogRepo) loadPlanVersion(
	ctx context.Context,
	planVersionID uint64,
	forUpdate bool,
) (catalogdomain.PlanVersionAggregate, error) {
	db := r.transactionRepo.GetDB(ctx)
	query := db.Model(&models.CatalogPlanVersion{})
	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var version models.CatalogPlanVersion
	if err := query.First(&version, planVersionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return catalogdomain.PlanVersionAggregate{},
				catalogpublication.ErrPlanVersionNotFound
		}
		return catalogdomain.PlanVersionAggregate{}, err
	}

	var plan models.CatalogPlan
	planQuery := db.Model(&models.CatalogPlan{})
	if forUpdate {
		planQuery = planQuery.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := planQuery.First(&plan, version.PlanID).Error; err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}
	var product models.CatalogProduct
	productQuery := db.Model(&models.CatalogProduct{})
	if forUpdate {
		productQuery = productQuery.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := productQuery.First(&product, plan.ProductID).Error; err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}

	var entitlements []models.CatalogPlanEntitlement
	if err := db.Where("plan_version_id = ?", version.ID).
		Order("code ASC").
		Find(&entitlements).Error; err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}
	var prices []models.CatalogPriceItem
	if err := db.Where("plan_version_id = ?", version.ID).
		Order("code ASC").
		Find(&prices).Error; err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}

	return toPlanVersionAggregate(product, plan, version, entitlements, prices), nil
}

func toPlanVersionAggregate(
	product models.CatalogProduct,
	plan models.CatalogPlan,
	version models.CatalogPlanVersion,
	entitlements []models.CatalogPlanEntitlement,
	prices []models.CatalogPriceItem,
) catalogdomain.PlanVersionAggregate {
	aggregate := catalogdomain.PlanVersionAggregate{
		ID:                 version.ID,
		ProductID:          product.ID,
		PlanID:             plan.ID,
		ProductCode:        product.Code,
		ProductDisplayName: product.DisplayName,
		ProductStatus:      product.Status,
		PlanCode:           plan.Code,
		PlanStatus:         plan.Status,
		Version:            version.Version,
		DisplayName:        version.DisplayName,
		Status:             version.Status,
		SubjectScope:       version.SubjectScope,
		EffectiveFrom:      copyTime(version.EffectiveFrom),
		EffectiveUntil:     copyTime(version.EffectiveUntil),
		PublishedAt:        copyTime(version.PublishedAt),
		TermsChecksum:      version.TermsChecksum,
		Entitlements:       make([]productcatalog.Entitlement, 0, len(entitlements)),
		Prices:             make([]productcatalog.PriceItem, 0, len(prices)),
	}
	for _, entity := range entitlements {
		aggregate.Entitlements = append(
			aggregate.Entitlements,
			productcatalog.Entitlement{
				Code:        entity.Code,
				Kind:        entity.Kind,
				FeatureCode: valueOrEmpty(entity.FeatureCode),
				MeterCode:   valueOrEmpty(entity.MeterCode),
				Amount:      entity.Amount,
				Unlimited:   entity.Unlimited,
				Period:      entity.Period,
			},
		)
	}
	for _, entity := range prices {
		aggregate.Prices = append(
			aggregate.Prices,
			productcatalog.PriceItem{
				Code:        entity.Code,
				Kind:        entity.Kind,
				Currency:    entity.Currency,
				AmountMinor: entity.AmountMinor,
				BillingUnit: entity.BillingUnit,
				MeterCode:   valueOrEmpty(entity.MeterCode),
				Quantity:    entity.Quantity,
			},
		)
	}
	return aggregate
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func copyUint64(value *uint64) *uint64 {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}

func toSubscriptionAggregate(
	subscription models.CatalogSubscription,
	product models.CatalogProduct,
	plan models.CatalogPlan,
	version models.CatalogPlanVersion,
) subscriptiondomain.Aggregate {
	return subscriptiondomain.Aggregate{
		ID:                 subscription.ID,
		SubscriptionKey:    subscription.SubscriptionKey,
		SubjectKind:        subscription.SubjectKind,
		SubjectID:          subscription.SubjectID,
		ProductID:          subscription.ProductID,
		ProductCode:        product.Code,
		PlanVersionID:      subscription.PlanVersionID,
		PlanCode:           plan.Code,
		PlanVersion:        version.Version,
		PlanTermsChecksum:  version.TermsChecksum,
		PlanStatus:         version.Status,
		Status:             subscription.Status,
		StartedAt:          subscription.StartedAt,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		AccessUntil:        copyTime(subscription.AccessUntil),
		AutoRenew:          subscription.AutoRenew,
		OrderReference:     subscription.OrderReference,
		CanceledAt:         copyTime(subscription.CanceledAt),
		EndedAt:            copyTime(subscription.EndedAt),
		CreatedAt:          subscription.CreatedAt,
		UpdatedAt:          subscription.UpdatedAt,
	}
}
