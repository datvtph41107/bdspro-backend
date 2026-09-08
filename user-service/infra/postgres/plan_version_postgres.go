package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	_db "common/db"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"user/internal/domain/plan"
	"user/internal/models"
	"user/internal/usecase/plan/admin"
	"user/internal/usecase/plan/publish"
)

// PlanVersionRepo persists serving Product/PlanVersion state. It intentionally does
// not implement transition/shadow contracts.
type PlanVersionRepo struct {
	transactionRepo *_db.TransactionRepo
}

/** CreatePlanVersionDraft persists one complete draft aggregate atomically. */
func (r *PlanVersionRepo) CreatePlanVersionDraft(ctx context.Context, input admin.DraftRecord) (uint64, error) {
	if r == nil || r.transactionRepo == nil {
		return 0, fmt.Errorf("plan repository is not configured")
	}
	var versionID uint64
	err := r.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		db := r.transactionRepo.GetDB(txCtx)
		product, err := r.findOrCreateProduct(db, input)
		if err != nil {
			return err
		}
		planRow, err := r.findOrCreatePlan(db, product.ID, input)
		if err != nil {
			return err
		}
		version := models.CatalogPlanVersion{
			PlanID:               planRow.ID,
			Version:              input.Terms.Version,
			DisplayName:          input.Terms.DisplayName,
			Status:               plan.StatusDraft,
			SubjectScope:         input.Terms.SubjectScope,
			SubscriptionTermDays: int32Pointer(input.Terms.SubscriptionTermDays),
			Source:               models.CatalogSourceNative,
			TermsChecksum:        input.TermsChecksum,
			CreatedBy:            uint64Pointer(input.ActorID),
			UpdatedBy:            uint64Pointer(input.ActorID),
		}
		if err := db.Create(&version).Error; err != nil {
			return catalogMutationError(err)
		}
		if err := createDraftTerms(db, version.ID, input); err != nil {
			return err
		}
		versionID = version.ID
		return nil
	})
	return versionID, err
}

/** UpdatePlanVersionDraft replaces all editable terms while holding the version lock. */
func (r *PlanVersionRepo) UpdatePlanVersionDraft(ctx context.Context, planVersionID uint64, input admin.DraftRecord) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("plan repository is not configured")
	}
	return r.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		db := r.transactionRepo.GetDB(txCtx)
		var version models.CatalogPlanVersion
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&version, planVersionID).Error; err != nil {
			return catalogNotFoundError(err)
		}
		if version.Status != plan.StatusDraft || version.PublishedAt != nil {
			return admin.ErrPlanVersionNotDraft
		}
		var planRow models.CatalogPlan
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&planRow, version.PlanID).Error; err != nil {
			return catalogNotFoundError(err)
		}
		var product models.CatalogProduct
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, planRow.ProductID).Error; err != nil {
			return catalogNotFoundError(err)
		}
		if product.Code != input.Terms.ProductCode || planRow.Code != input.Terms.PlanCode || version.Version != input.Terms.Version {
			return fmt.Errorf("product code, plan code and version are immutable")
		}
		if product.DisplayName != input.ProductDisplayName || int32Value(planRow.TierRank) != input.TierRank {
			return admin.ErrStablePlanFieldsImmutable
		}
		if err := db.Model(&version).Updates(map[string]any{
			"display_name":           input.Terms.DisplayName,
			"subject_scope":          input.Terms.SubjectScope,
			"subscription_term_days": input.Terms.SubscriptionTermDays,
			"terms_checksum":         input.TermsChecksum,
			"updated_by":             input.ActorID,
			"updated_at":             gorm.Expr("NOW()"),
		}).Error; err != nil {
			return err
		}
		for _, model := range []any{&models.CatalogPlanEntitlement{}, &models.CatalogPlanOperationPolicy{}, &models.CatalogPriceItem{}} {
			if err := db.Where("plan_version_id = ?", planVersionID).Delete(model).Error; err != nil {
				return err
			}
		}
		return createDraftTerms(db, planVersionID, input)
	})
}

// RetirePlanVersion closes the selling window without deleting the immutable
// terms referenced by historical orders and subscriptions.
func (r *PlanVersionRepo) RetirePlanVersion(ctx context.Context, planVersionID, actorID uint64, retiredAt time.Time) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("plan repository is not configured")
	}
	if retiredAt.IsZero() {
		return fmt.Errorf("retired at is required")
	}
	return r.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		db := r.transactionRepo.GetDB(txCtx)
		result := db.Model(&models.CatalogPlanVersion{}).
			Where("id = ? AND status = ? AND published_at IS NOT NULL", planVersionID, plan.StatusActive).
			Updates(map[string]any{
				"status":          plan.StatusRetired,
				"effective_until": retiredAt.UTC(),
				"updated_by":      actorID,
				"updated_at":      gorm.Expr("NOW()"),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 1 {
			return nil
		}
		var count int64
		if err := db.Model(&models.CatalogPlanVersion{}).Where("id = ?", planVersionID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return admin.ErrPlanVersionNotFound
		}
		return admin.ErrPlanVersionNotActive
	})
}

/** DeletePlanVersionDraft removes only a never-published candidate and its terms. */
func (r *PlanVersionRepo) DeletePlanVersionDraft(ctx context.Context, planVersionID uint64) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("plan repository is not configured")
	}
	return r.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		db := r.transactionRepo.GetDB(txCtx)
		var version models.CatalogPlanVersion
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&version, planVersionID).Error; err != nil {
			return catalogNotFoundError(err)
		}
		if version.Status != plan.StatusDraft || version.PublishedAt != nil {
			return admin.ErrPlanVersionNotDraft
		}
		var subscriptions int64
		if err := db.Model(&models.CatalogSubscription{}).Where("plan_version_id = ? OR pending_plan_version_id = ?", planVersionID, planVersionID).Count(&subscriptions).Error; err != nil {
			return err
		}
		if subscriptions != 0 {
			return admin.ErrPlanVersionNotDraft
		}
		for _, model := range []any{&models.CatalogPlanEntitlement{}, &models.CatalogPlanOperationPolicy{}, &models.CatalogPriceItem{}} {
			if err := db.Where("plan_version_id = ?", planVersionID).Delete(model).Error; err != nil {
				return err
			}
		}
		return db.Delete(&models.CatalogPlanVersion{}, planVersionID).Error
	})
}

func (r *PlanVersionRepo) findOrCreateProduct(db *gorm.DB, input admin.DraftRecord) (models.CatalogProduct, error) {
	var product models.CatalogProduct
	err := db.Where("code = ?", input.Terms.ProductCode).First(&product).Error
	if err == nil {
		return product, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.CatalogProduct{}, err
	}
	product = models.CatalogProduct{Code: input.Terms.ProductCode, DisplayName: input.ProductDisplayName, Status: plan.StatusDraft, CreatedBy: uint64Pointer(input.ActorID), UpdatedBy: uint64Pointer(input.ActorID)}
	if err := db.Create(&product).Error; err != nil {
		return models.CatalogProduct{}, catalogMutationError(err)
	}
	return product, nil
}

func (r *PlanVersionRepo) findOrCreatePlan(db *gorm.DB, productID uint64, input admin.DraftRecord) (models.CatalogPlan, error) {
	var planRow models.CatalogPlan
	err := db.Where("code = ?", input.Terms.PlanCode).First(&planRow).Error
	if err == nil {
		if planRow.ProductID != productID {
			return models.CatalogPlan{}, admin.ErrPlanVersionConflict
		}
		return planRow, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.CatalogPlan{}, err
	}
	planRow = models.CatalogPlan{ProductID: productID, Code: input.Terms.PlanCode, TierRank: int32Pointer(input.TierRank), Status: plan.StatusDraft, CreatedBy: uint64Pointer(input.ActorID), UpdatedBy: uint64Pointer(input.ActorID)}
	if err := db.Create(&planRow).Error; err != nil {
		return models.CatalogPlan{}, catalogMutationError(err)
	}
	return planRow, nil
}

func createDraftTerms(db *gorm.DB, planVersionID uint64, input admin.DraftRecord) error {
	for _, value := range input.Terms.Entitlements {
		row := models.CatalogPlanEntitlement{PlanVersionID: planVersionID, Code: value.Code, Kind: value.Kind, FeatureCode: stringPointer(value.FeatureCode), MeterCode: stringPointer(value.MeterCode), Amount: value.Amount, Unlimited: value.Unlimited, Period: value.Period, CreatedBy: uint64Pointer(input.ActorID)}
		if err := db.Create(&row).Error; err != nil {
			return catalogMutationError(err)
		}
	}
	for _, value := range input.Terms.Operations {
		row := models.CatalogPlanOperationPolicy{PlanVersionID: planVersionID, OperationCode: value.Code, FeatureCode: value.FeatureCode, MeterCode: stringPointer(value.MeterCode), UnitsPerAction: value.UnitsPerAction, CreatedBy: uint64Pointer(input.ActorID)}
		if err := db.Create(&row).Error; err != nil {
			return catalogMutationError(err)
		}
	}
	for _, value := range input.Terms.Prices {
		row := models.CatalogPriceItem{PlanVersionID: planVersionID, Code: value.Code, Kind: value.Kind, Currency: value.Currency, AmountMinor: value.AmountMinor, BillingUnit: value.BillingUnit, MeterCode: stringPointer(value.MeterCode), Quantity: value.Quantity, CreatedBy: uint64Pointer(input.ActorID)}
		if err := db.Create(&row).Error; err != nil {
			return catalogMutationError(err)
		}
	}
	return nil
}

func catalogNotFoundError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return admin.ErrPlanVersionNotFound
	}
	return err
}

func catalogMutationError(err error) error {
	var pqError *pq.Error
	if errors.As(err, &pqError) && pqError.Code == "23505" {
		return admin.ErrPlanVersionConflict
	}
	if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
		return admin.ErrPlanVersionConflict
	}
	return err
}

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}
	copy := value
	return &copy
}

func uint64Pointer(value uint64) *uint64 { return &value }

func int32Pointer(value int32) *int32 { return &value }

var (
	_ publish.Repository = (*PlanVersionRepo)(nil)
	_ admin.Repository   = (*PlanVersionRepo)(nil)
)

// NewPlanVersionRepo creates the serving Plan persistence adapter.
func NewPlanVersionRepo(transactionRepo *_db.TransactionRepo) *PlanVersionRepo {
	return &PlanVersionRepo{transactionRepo: transactionRepo}
}

/** LoadPlanVersionForUpdate locks one draft aggregate. */
func (r *PlanVersionRepo) LoadPlanVersionForUpdate(
	ctx context.Context,
	planVersionID uint64,
) (plan.PlanVersionAggregate, error) {
	if r == nil || r.transactionRepo == nil {
		return plan.PlanVersionAggregate{}, fmt.Errorf(
			"plan repository is not configured",
		)
	}
	return r.loadPlanVersion(ctx, planVersionID, true)
}

/** GetPlanVersion loads one serving PlanVersion aggregate. */
func (r *PlanVersionRepo) GetPlanVersion(
	ctx context.Context,
	planVersionID uint64,
) (plan.PlanVersionAggregate, error) {
	if r == nil || r.transactionRepo == nil {
		return plan.PlanVersionAggregate{}, fmt.Errorf(
			"plan repository is not configured",
		)
	}
	return r.loadPlanVersion(ctx, planVersionID, false)
}

/** ListPlanVersions returns a bounded page for the Admin control plane. */
func (r *PlanVersionRepo) ListPlanVersions(
	ctx context.Context,
	query admin.Query,
) (admin.Page, error) {
	if r == nil || r.transactionRepo == nil {
		return admin.Page{}, fmt.Errorf("plan repository is not configured")
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
		PlanVersions: make([]plan.PlanVersionAggregate, 0, len(versions)),
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
func (r *PlanVersionRepo) MarkPlanVersionPublished(
	ctx context.Context,
	input publish.PublishRecord,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("plan repository is not configured")
	}
	parentUpdates := map[string]any{
		"status":     plan.StatusActive,
		"updated_by": input.PublishedBy,
		"updated_at": gorm.Expr("NOW()"),
	}
	db := r.transactionRepo.GetDB(ctx)
	productResult := db.Model(&models.CatalogProduct{}).
		Where("id = ? AND status IN ?", input.ProductID, []plan.Status{
			plan.StatusDraft,
			plan.StatusActive,
		}).
		Updates(parentUpdates)
	if productResult.Error != nil {
		return productResult.Error
	}
	if productResult.RowsAffected != 1 {
		return publish.ErrConcurrentPublish
	}
	planResult := db.Model(&models.CatalogPlan{}).
		Where("id = ? AND product_id = ? AND status IN ?", input.PlanID, input.ProductID, []plan.Status{
			plan.StatusDraft,
			plan.StatusActive,
		}).
		Updates(parentUpdates)
	if planResult.Error != nil {
		return planResult.Error
	}
	if planResult.RowsAffected != 1 {
		return publish.ErrConcurrentPublish
	}

	updates := map[string]any{
		"status":          plan.StatusActive,
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
			plan.StatusDraft,
		).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return publish.ErrConcurrentPublish
	}
	return nil
}

func (r *PlanVersionRepo) loadPlanVersion(
	ctx context.Context,
	planVersionID uint64,
	forUpdate bool,
) (plan.PlanVersionAggregate, error) {
	db := r.transactionRepo.GetDB(ctx)
	query := db.Model(&models.CatalogPlanVersion{})
	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var version models.CatalogPlanVersion
	if err := query.First(&version, planVersionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return plan.PlanVersionAggregate{},
				publish.ErrPlanVersionNotFound
		}
		return plan.PlanVersionAggregate{}, err
	}

	var planRow models.CatalogPlan
	planQuery := db.Model(&models.CatalogPlan{})
	if forUpdate {
		planQuery = planQuery.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := planQuery.First(&planRow, version.PlanID).Error; err != nil {
		return plan.PlanVersionAggregate{}, err
	}
	var product models.CatalogProduct
	productQuery := db.Model(&models.CatalogProduct{})
	if forUpdate {
		productQuery = productQuery.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := productQuery.First(&product, planRow.ProductID).Error; err != nil {
		return plan.PlanVersionAggregate{}, err
	}

	var entitlements []models.CatalogPlanEntitlement
	if err := db.Where("plan_version_id = ?", version.ID).
		Order("code ASC").
		Find(&entitlements).Error; err != nil {
		return plan.PlanVersionAggregate{}, err
	}
	var operations []models.CatalogPlanOperationPolicy
	if err := db.Where("plan_version_id = ?", version.ID).
		Order("operation_code ASC").
		Find(&operations).Error; err != nil {
		return plan.PlanVersionAggregate{}, err
	}
	var prices []models.CatalogPriceItem
	if err := db.Where("plan_version_id = ?", version.ID).
		Order("code ASC").
		Find(&prices).Error; err != nil {
		return plan.PlanVersionAggregate{}, err
	}

	return toPlanVersionAggregate(product, planRow, version, entitlements, operations, prices), nil
}

func toPlanVersionAggregate(
	product models.CatalogProduct,
	planRow models.CatalogPlan,
	version models.CatalogPlanVersion,
	entitlements []models.CatalogPlanEntitlement,
	operations []models.CatalogPlanOperationPolicy,
	prices []models.CatalogPriceItem,
) plan.PlanVersionAggregate {
	aggregate := plan.PlanVersionAggregate{
		ID:                   version.ID,
		ProductID:            product.ID,
		PlanID:               planRow.ID,
		ProductCode:          product.Code,
		ProductDisplayName:   product.DisplayName,
		ProductStatus:        product.Status,
		PlanCode:             planRow.Code,
		TierRank:             int32Value(planRow.TierRank),
		PlanStatus:           planRow.Status,
		Version:              version.Version,
		DisplayName:          version.DisplayName,
		Status:               version.Status,
		SubjectScope:         version.SubjectScope,
		SubscriptionTermDays: int32Value(version.SubscriptionTermDays),
		EffectiveFrom:        copyTime(version.EffectiveFrom),
		EffectiveUntil:       copyTime(version.EffectiveUntil),
		PublishedAt:          copyTime(version.PublishedAt),
		TermsChecksum:        version.TermsChecksum,
		Entitlements:         make([]plan.Entitlement, 0, len(entitlements)),
		Operations:           make([]plan.OperationBinding, 0, len(operations)),
		Prices:               make([]plan.PriceItem, 0, len(prices)),
	}
	for _, entity := range entitlements {
		aggregate.Entitlements = append(
			aggregate.Entitlements,
			plan.Entitlement{
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
			plan.OperationBinding{
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
			plan.PriceItem{
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

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}
