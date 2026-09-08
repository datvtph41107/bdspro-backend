package plan

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	stableCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*(?:\.[a-z0-9_]+)*$`)
	versionPattern    = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	currencyPattern   = regexp.MustCompile(`^[A-Z]{3}$`)
)

type Registry struct {
	version  string
	products map[string]string
	plans    map[string]map[string]PlanVersion
}

/** NewRegistry validates and freezes a catalog snapshot. */
func NewRegistry(catalog Snapshot) (*Registry, error) {
	if !versionPattern.MatchString(catalog.Version) {
		return nil, fmt.Errorf("invalid catalog version %q", catalog.Version)
	}

	registry := &Registry{
		version:  catalog.Version,
		products: make(map[string]string),
		plans:    make(map[string]map[string]PlanVersion),
	}

	productCodes := make(map[string]struct{}, len(catalog.Products))
	for _, product := range catalog.Products {
		if err := validateProduct(product); err != nil {
			return nil, err
		}
		if _, exists := productCodes[product.Code]; exists {
			return nil, fmt.Errorf("duplicate product code %q", product.Code)
		}
		productCodes[product.Code] = struct{}{}
		registry.products[product.Code] = product.DisplayName

		for _, plan := range product.Plans {
			if plan.ProductCode != product.Code {
				return nil, fmt.Errorf(
					"plan %q product %q does not match %q",
					plan.PlanCode,
					plan.ProductCode,
					product.Code,
				)
			}

			versions := registry.plans[plan.PlanCode]
			if versions == nil {
				versions = make(map[string]PlanVersion)
				registry.plans[plan.PlanCode] = versions
			}
			if _, exists := versions[plan.Version]; exists {
				return nil, fmt.Errorf(
					"duplicate plan version %s@%s",
					plan.PlanCode,
					plan.Version,
				)
			}
			versions[plan.Version] = copyPlan(plan)
		}
	}

	for planCode, versions := range registry.plans {
		if err := validateEffectiveWindows(planCode, versions); err != nil {
			return nil, err
		}
	}

	return registry, nil
}

/** Version returns the frozen catalog version. */
func (r *Registry) Version() string {
	if r == nil {
		return ""
	}
	return r.version
}

/** Plan returns one immutable plan version. */
func (r *Registry) Plan(planCode, version string) (PlanVersion, bool) {
	if r == nil {
		return PlanVersion{}, false
	}
	versions := r.plans[planCode]
	plan, ok := versions[version]
	if !ok {
		return PlanVersion{}, false
	}
	return copyPlan(plan), true
}

/** EffectivePlan resolves one active version at a point in time. */
func (r *Registry) EffectivePlan(planCode string, at time.Time) (PlanVersion, bool) {
	if r == nil || at.IsZero() {
		return PlanVersion{}, false
	}

	versions := r.plans[planCode]
	for _, plan := range versions {
		if plan.Status != StatusActive {
			continue
		}
		if at.Before(plan.EffectiveFrom) {
			continue
		}
		if plan.EffectiveUntil != nil && !at.Before(*plan.EffectiveUntil) {
			continue
		}
		return copyPlan(plan), true
	}
	return PlanVersion{}, false
}

/** Snapshot returns a deterministic defensive copy. */
func (r *Registry) Snapshot() Snapshot {
	if r == nil {
		return Snapshot{}
	}

	productPlans := make(map[string][]PlanVersion)
	for _, versions := range r.plans {
		for _, plan := range versions {
			productPlans[plan.ProductCode] = append(
				productPlans[plan.ProductCode],
				copyPlan(plan),
			)
		}
	}

	productCodes := make([]string, 0, len(productPlans))
	for productCode := range productPlans {
		productCodes = append(productCodes, productCode)
	}
	sort.Strings(productCodes)

	products := make([]Product, 0, len(productCodes))
	for _, productCode := range productCodes {
		plans := productPlans[productCode]
		sort.Slice(plans, func(i, j int) bool {
			if plans[i].PlanCode == plans[j].PlanCode {
				return plans[i].Version < plans[j].Version
			}
			return plans[i].PlanCode < plans[j].PlanCode
		})
		products = append(products, Product{
			Code:        productCode,
			DisplayName: r.products[productCode],
			Plans:       plans,
		})
	}

	return Snapshot{Version: r.version, Products: products}
}

func validateProduct(product Product) error {
	if !isStableCode(product.Code) {
		return fmt.Errorf("invalid product code %q", product.Code)
	}
	if strings.TrimSpace(product.DisplayName) == "" {
		return fmt.Errorf("product %q display name is required", product.Code)
	}
	if len(product.Plans) == 0 {
		return fmt.Errorf("product %q requires at least one plan", product.Code)
	}
	for _, plan := range product.Plans {
		if err := validatePlan(plan); err != nil {
			return err
		}
		if err := validateCommercialCompleteness(plan); err != nil {
			return err
		}
	}
	return nil
}

func validatePlan(plan PlanVersion) error {
	if !isStableCode(plan.ProductCode) {
		return fmt.Errorf("invalid product code %q", plan.ProductCode)
	}
	if !isStableCode(plan.PlanCode) {
		return fmt.Errorf("invalid plan code %q", plan.PlanCode)
	}
	if !versionPattern.MatchString(plan.Version) {
		return fmt.Errorf("invalid plan version %q", plan.Version)
	}
	if strings.TrimSpace(plan.DisplayName) == "" {
		return fmt.Errorf("plan %q display name is required", plan.PlanCode)
	}
	if !validStatus(plan.Status) {
		return fmt.Errorf("plan %q has invalid status %q", plan.PlanCode, plan.Status)
	}
	if !validSubjectScope(plan.SubjectScope) {
		return fmt.Errorf(
			"plan %q has invalid subject scope %q",
			plan.PlanCode,
			plan.SubjectScope,
		)
	}
	if plan.SubscriptionTermDays < 0 {
		return fmt.Errorf("plan %q subscription term must not be negative", plan.PlanCode)
	}
	if plan.Status != StatusDraft && plan.EffectiveFrom.IsZero() {
		return fmt.Errorf("plan %q requires effective_from", plan.PlanCode)
	}
	if plan.EffectiveUntil != nil && !plan.EffectiveUntil.After(plan.EffectiveFrom) {
		return fmt.Errorf("plan %q effective window is invalid", plan.PlanCode)
	}
	if err := validateEntitlements(plan); err != nil {
		return err
	}
	if err := validateOperationBindings(plan); err != nil {
		return err
	}
	if err := validatePrices(plan); err != nil {
		return err
	}
	return nil
}

// validateCommercialCompleteness is intentionally stricter than checksum
// validation so legacy published rows can still verify their historical hash
// while every newly active contract must carry explicit commercial terms.
func validateCommercialCompleteness(plan PlanVersion) error {
	if plan.Status != StatusDraft && plan.SubscriptionTermDays <= 0 {
		return fmt.Errorf("plan %q requires subscription_term_days", plan.PlanCode)
	}
	return nil
}

func validateOperationBindings(plan PlanVersion) error {
	features := make(map[string]struct{}, len(plan.Entitlements))
	usageMeters := make(map[string]struct{}, len(plan.Entitlements))
	for _, entitlement := range plan.Entitlements {
		switch entitlement.Kind {
		case EntitlementFeatureAccess:
			features[entitlement.FeatureCode] = struct{}{}
		case EntitlementUsageAllowance:
			usageMeters[entitlement.MeterCode] = struct{}{}
		}
	}

	codes := make(map[string]struct{}, len(plan.Operations))
	for _, policy := range plan.Operations {
		if !isStableCode(policy.Code) {
			return fmt.Errorf("invalid operation code %q", policy.Code)
		}
		if _, exists := codes[policy.Code]; exists {
			return fmt.Errorf("duplicate operation policy %q", policy.Code)
		}
		codes[policy.Code] = struct{}{}
		if !isStableCode(policy.FeatureCode) {
			return fmt.Errorf("operation %q requires feature_code", policy.Code)
		}
		if _, exists := features[policy.FeatureCode]; !exists {
			return fmt.Errorf(
				"operation %q feature %q is not granted by this plan",
				policy.Code,
				policy.FeatureCode,
			)
		}
		if policy.MeterCode == "" {
			if policy.UnitsPerAction != 0 {
				return fmt.Errorf(
					"unmetered operation %q must have zero units_per_action",
					policy.Code,
				)
			}
			continue
		}
		if !isStableCode(policy.MeterCode) {
			return fmt.Errorf("operation %q has invalid meter_code", policy.Code)
		}
		if policy.UnitsPerAction <= 0 {
			return fmt.Errorf("operation %q units_per_action must be positive", policy.Code)
		}
		if _, exists := usageMeters[policy.MeterCode]; !exists {
			return fmt.Errorf(
				"operation %q meter %q has no usage allowance",
				policy.Code,
				policy.MeterCode,
			)
		}
	}
	return nil
}

func validateEntitlements(plan PlanVersion) error {
	codes := make(map[string]struct{}, len(plan.Entitlements))
	targets := make(map[string]struct{}, len(plan.Entitlements))
	for _, entitlement := range plan.Entitlements {
		if !isStableCode(entitlement.Code) {
			return fmt.Errorf("invalid entitlement code %q", entitlement.Code)
		}
		if _, exists := codes[entitlement.Code]; exists {
			return fmt.Errorf("duplicate entitlement code %q", entitlement.Code)
		}
		codes[entitlement.Code] = struct{}{}

		target, err := validateEntitlement(entitlement)
		if err != nil {
			return err
		}
		if _, exists := targets[target]; exists {
			return fmt.Errorf("duplicate entitlement target %q", target)
		}
		targets[target] = struct{}{}
	}
	return nil
}

func validateEntitlement(entitlement Entitlement) (string, error) {
	switch entitlement.Kind {
	case EntitlementFeatureAccess:
		if !isStableCode(entitlement.FeatureCode) {
			return "", fmt.Errorf(
				"feature entitlement %q requires feature_code",
				entitlement.Code,
			)
		}
		if entitlement.MeterCode != "" || entitlement.Amount != 0 ||
			entitlement.Unlimited || entitlement.Period != PeriodNone {
			return "", fmt.Errorf(
				"feature entitlement %q contains quota fields",
				entitlement.Code,
			)
		}
		return "feature:" + entitlement.FeatureCode, nil

	case EntitlementUsageAllowance:
		if !isStableCode(entitlement.MeterCode) {
			return "", fmt.Errorf(
				"usage entitlement %q requires meter_code",
				entitlement.Code,
			)
		}
		if entitlement.FeatureCode != "" {
			return "", fmt.Errorf(
				"usage entitlement %q contains feature_code",
				entitlement.Code,
			)
		}
		if err := validateAmount(entitlement); err != nil {
			return "", err
		}
		if entitlement.Period != PeriodDay &&
			entitlement.Period != PeriodCalendarMonth &&
			entitlement.Period != PeriodSubscriptionCycle &&
			entitlement.Period != PeriodLifetime {
			return "", fmt.Errorf(
				"usage entitlement %q has invalid period %q",
				entitlement.Code,
				entitlement.Period,
			)
		}
		return "usage:" + entitlement.MeterCode + ":" + string(entitlement.Period), nil

	case EntitlementCapacityLimit:
		if !isStableCode(entitlement.MeterCode) {
			return "", fmt.Errorf(
				"capacity entitlement %q requires meter_code",
				entitlement.Code,
			)
		}
		if entitlement.FeatureCode != "" {
			return "", fmt.Errorf(
				"capacity entitlement %q contains feature_code",
				entitlement.Code,
			)
		}
		if err := validateAmount(entitlement); err != nil {
			return "", err
		}
		if entitlement.Period != PeriodNone {
			return "", fmt.Errorf(
				"capacity entitlement %q must use period none",
				entitlement.Code,
			)
		}
		return "capacity:" + entitlement.MeterCode, nil

	default:
		return "", fmt.Errorf(
			"entitlement %q has invalid kind %q",
			entitlement.Code,
			entitlement.Kind,
		)
	}
}

func validateAmount(entitlement Entitlement) error {
	if entitlement.Unlimited {
		if entitlement.Amount != 0 {
			return fmt.Errorf(
				"entitlement %q unlimited amount must be zero",
				entitlement.Code,
			)
		}
		return nil
	}
	if entitlement.Amount <= 0 {
		return fmt.Errorf(
			"entitlement %q amount must be positive",
			entitlement.Code,
		)
	}
	return nil
}

func validatePrices(plan PlanVersion) error {
	codes := make(map[string]struct{}, len(plan.Prices))
	for _, price := range plan.Prices {
		if !isStableCode(price.Code) {
			return fmt.Errorf("invalid price item code %q", price.Code)
		}
		if _, exists := codes[price.Code]; exists {
			return fmt.Errorf("duplicate price item code %q", price.Code)
		}
		codes[price.Code] = struct{}{}
		if !currencyPattern.MatchString(price.Currency) {
			return fmt.Errorf("price %q has invalid currency", price.Code)
		}
		if price.AmountMinor < 0 {
			return fmt.Errorf("price %q amount cannot be negative", price.Code)
		}
		if !isStableCode(price.BillingUnit) {
			return fmt.Errorf("price %q requires billing_unit", price.Code)
		}
		if price.Quantity <= 0 {
			return fmt.Errorf("price %q quantity must be positive", price.Code)
		}

		switch price.Kind {
		case PriceRecurring:
			if price.MeterCode != "" {
				return fmt.Errorf("recurring price %q cannot target a meter", price.Code)
			}
		case PriceAddOn, PriceOverage:
			if !isStableCode(price.MeterCode) {
				return fmt.Errorf("price %q requires meter_code", price.Code)
			}
		default:
			return fmt.Errorf("price %q has invalid kind %q", price.Code, price.Kind)
		}
	}
	return nil
}

func validateEffectiveWindows(planCode string, versions map[string]PlanVersion) error {
	active := make([]PlanVersion, 0, len(versions))
	for _, plan := range versions {
		if plan.Status == StatusActive {
			active = append(active, plan)
		}
	}
	sort.Slice(active, func(i, j int) bool {
		return active[i].EffectiveFrom.Before(active[j].EffectiveFrom)
	})

	for i := 1; i < len(active); i++ {
		previous := active[i-1]
		current := active[i]
		if previous.EffectiveUntil == nil || current.EffectiveFrom.Before(*previous.EffectiveUntil) {
			return fmt.Errorf("plan %q has overlapping active versions", planCode)
		}
	}
	return nil
}

func copyPlan(plan PlanVersion) PlanVersion {
	copyValue := plan
	copyValue.Entitlements = append([]Entitlement(nil), plan.Entitlements...)
	copyValue.Operations = append([]OperationBinding(nil), plan.Operations...)
	copyValue.Prices = append([]PriceItem(nil), plan.Prices...)
	if plan.EffectiveUntil != nil {
		until := *plan.EffectiveUntil
		copyValue.EffectiveUntil = &until
	}
	return copyValue
}

func isStableCode(value string) bool {
	return stableCodePattern.MatchString(value)
}

func validStatus(status Status) bool {
	return status == StatusDraft || status == StatusActive || status == StatusRetired
}

func validSubjectScope(scope SubjectScope) bool {
	return scope == SubjectScopeProfile ||
		scope == SubjectScopeOrganization ||
		scope == SubjectScopeAny
}
