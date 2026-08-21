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
	"user/internal/models"
	"user/internal/modules/catalog/admin"
	catalogdomain "user/internal/modules/catalog/domain"
	"user/internal/modules/catalog/publication"
	catalogshadow "user/internal/modules/catalog/shadow"
	subscriptiondomain "user/internal/modules/subscription/domain"
	subscriptionshadow "user/internal/modules/subscription/shadow"
)

/** CatalogRepo persists catalog v2 aggregates. */
// @bind: user/internal/modules/catalog/shadow.Repository
// @bind: user/internal/modules/subscription/shadow.Repository
type CatalogRepo struct {
	transactionRepo *_db.TransactionRepo
}

var (
	_ publication.Repository        = (*CatalogRepo)(nil)
	_ admin.Repository              = (*CatalogRepo)(nil)
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
	if errors.Is(err, publication.ErrPlanVersionNotFound) {
		return catalogdomain.PlanVersionAggregate{}, catalogshadow.ErrNotFound
	}
	return aggregate, err
}

/** ListPlanVersions returns a bounded page for the Admin control plane. */
func (r *CatalogRepo) ListPlanVersions(
	ctx context.Context,
	query admin.Query,
) (admin.Page, error) {
	if r == nil || r.transactionRepo == nil {
		return admin.Page{}, fmt.Errorf("catalog repository is not configured")
	}
	db := r.transactionRepo.GetDB(ctx)
	base := db.Model(&models.CatalogPlanVersion{}).
		Joins("JOIN catalog_plans ON catalog_plans.id = catalog_plan_versions.plan_id").
		Joins("JOIN catalog_products ON catalog_products.id = catalog_plans.product_id")
	if query.ProductCode != "" {
		base = base.Where("catalog_products.code = ?", query.ProductCode)
	}
	if query.PlanCode != "" {
		base = base.Where("catalog_plans.code = ?", query.PlanCode)
	}
	if query.Status != "" {
		base = base.Where("catalog_plan_versions.status = ?", query.Status)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return admin.Page{}, err
	}
	offset := int((query.Page - 1) * query.PageSize)
	var versions []models.CatalogPlanVersion
	if err := base.
		Order("catalog_plan_versions.created_at DESC, catalog_plan_versions.id DESC").
		Offset(offset).
		Limit(int(query.PageSize)).
		Find(&versions).Error; err != nil {
		return admin.Page{}, err
	}
	page := admin.Page{
		PlanVersions: make([]catalogdomain.PlanVersionAggregate, 0, len(versions)),
		Total:        uint64(total),
	}
	for _, version := range versions {
		aggregate, err := r.loadPlanVersion(ctx, version.ID, false)
		if err != nil {
			return admin.Page{}, err
		}
		page.PlanVersions = append(page.PlanVersions, aggregate)
	}
	return page, nil
}

/** MarkPlanVersionPublished performs a guarded draft-to-active update. */
func (r *CatalogRepo) MarkPlanVersionPublished(
	ctx context.Context,
	input publication.PublishRecord,
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
		return publication.ErrConcurrentPublish
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
		return publication.ErrConcurrentPublish
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
		return publication.ErrConcurrentPublish
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
				publication.ErrPlanVersionNotFound
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
	var operations []models.CatalogPlanOperationPolicy
	if err := db.Where("plan_version_id = ?", version.ID).
		Order("operation_code ASC").
		Find(&operations).Error; err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}
	var prices []models.CatalogPriceItem
	if err := db.Where("plan_version_id = ?", version.ID).
		Order("code ASC").
		Find(&prices).Error; err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}

	return toPlanVersionAggregate(product, plan, version, entitlements, operations, prices), nil
}

func toPlanVersionAggregate(
	product models.CatalogProduct,
	plan models.CatalogPlan,
	version models.CatalogPlanVersion,
	entitlements []models.CatalogPlanEntitlement,
	operations []models.CatalogPlanOperationPolicy,
	prices []models.CatalogPriceItem,
) catalogdomain.PlanVersionAggregate {
	aggregate := catalogdomain.PlanVersionAggregate{
		ID:                   version.ID,
		ProductID:            product.ID,
		PlanID:               plan.ID,
		ProductCode:          product.Code,
		ProductDisplayName:   product.DisplayName,
		ProductStatus:        product.Status,
		PlanCode:             plan.Code,
		TierRank:             int32Value(plan.TierRank),
		PlanStatus:           plan.Status,
		Version:              version.Version,
		DisplayName:          version.DisplayName,
		Status:               version.Status,
		SubjectScope:         version.SubjectScope,
		SubscriptionTermDays: int32Value(version.SubscriptionTermDays),
		EffectiveFrom:        copyTime(version.EffectiveFrom),
		EffectiveUntil:       copyTime(version.EffectiveUntil),
		PublishedAt:          copyTime(version.PublishedAt),
		TermsChecksum:        version.TermsChecksum,
		Entitlements:         make([]productcatalog.Entitlement, 0, len(entitlements)),
		Operations:           make([]productcatalog.OperationPolicy, 0, len(operations)),
		Prices:               make([]productcatalog.PriceItem, 0, len(prices)),
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
	for _, entity := range operations {
		aggregate.Operations = append(
			aggregate.Operations,
			productcatalog.OperationPolicy{
				Code:           entity.OperationCode,
				FeatureCode:    entity.FeatureCode,
				MeterCode:      valueOrEmpty(entity.MeterCode),
				UnitsPerAction: entity.UnitsPerAction,
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

func int32Value(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
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
