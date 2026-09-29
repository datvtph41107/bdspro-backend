package accesspostgres

import (
	"common/operation"
	"context"
	"errors"
	"user/internal/domain/entitlement"
	catalogdomain "user/internal/domain/plan"
	"user/internal/models"
	useraccess "user/internal/usecase/entitlement"

	"gorm.io/gorm"
)

/**
 * PlanAccessStore đọc feature access và usage allowance của một operation.
 */
type PlanAccessStore struct {
	db *gorm.DB
}

func NewPlanAccessStore(db *gorm.DB) *PlanAccessStore {
	return &PlanAccessStore{db: db}
}

func (s *PlanAccessStore) FindPlanAccess(
	ctx context.Context,
	planVersionID uint64,
	operationCode operation.Code,
) (useraccess.PlanAccess, bool, error) {
	if s == nil || s.db == nil {
		return useraccess.PlanAccess{}, false, errors.New("plan access database is not configured")
	}

	var policy models.CatalogPlanOperationPolicy
	err := s.db.WithContext(ctx).
		Where("plan_version_id = ?", planVersionID).
		Where("operation_code = ?", string(operationCode)).
		First(&policy).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return useraccess.PlanAccess{}, false, nil
	}
	if err != nil {
		return useraccess.PlanAccess{}, false, err
	}
	var feature models.CatalogPlanEntitlement
	err = s.db.WithContext(ctx).
		Where("plan_version_id = ?", planVersionID).
		Where("kind = ?", catalogdomain.EntitlementFeatureAccess).
		Where("feature_code = ?", policy.FeatureCode).
		First(&feature).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return useraccess.PlanAccess{}, false, nil
	}
	if err != nil {
		return useraccess.PlanAccess{}, false, err
	}

	result := useraccess.PlanAccess{
		Allowed:     true,
		FeatureCode: policy.FeatureCode,
		Period:      access.PeriodNone,
	}

	if policy.MeterCode == nil || *policy.MeterCode == "" {
		if policy.UnitsPerAction != 0 {
			return useraccess.PlanAccess{}, false, errors.New("unmetered operation policy must have zero units per action")
		}
		return result, true, nil
	}
	if policy.UnitsPerAction <= 0 {
		return useraccess.PlanAccess{}, false, errors.New("metered operation policy must have positive units per action")
	}

	var usage models.CatalogPlanEntitlement
	err = s.db.WithContext(ctx).
		Where("plan_version_id = ?", planVersionID).
		Where("kind = ?", catalogdomain.EntitlementUsageAllowance).
		Where("meter_code = ?", *policy.MeterCode).
		First(&usage).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return useraccess.PlanAccess{}, false, errors.New("usage limit is missing for this operation")
	}
	if err != nil {
		return useraccess.PlanAccess{}, false, err
	}
	if usage.Unlimited {
		if usage.Amount != 0 {
			return useraccess.PlanAccess{}, false, errors.New("unlimited usage entitlement must have zero amount")
		}
	} else if usage.Amount <= 0 {
		return useraccess.PlanAccess{}, false, errors.New("usage entitlement amount must be positive")
	}
	if usage.Period == catalogdomain.PeriodNone {
		return useraccess.PlanAccess{}, false, errors.New("usage entitlement period is missing")
	}

	result.Unlimited = usage.Unlimited
	result.MeterCode = *policy.MeterCode
	result.UnitsPerAction = policy.UnitsPerAction
	result.Limit = usage.Amount
	result.Period = access.Period(usage.Period)

	return result, true, nil
}
