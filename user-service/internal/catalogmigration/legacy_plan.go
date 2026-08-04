package catalogmigration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"common/productcatalog"
)

type LegacyField string

const (
	FieldCreateProduct      LegacyField = "create_product"
	FieldLimitPost          LegacyField = "post_limit"
	FieldPublicProductLimit LegacyField = "public_product_limit"
	FieldCreateProductAI    LegacyField = "create_product_ai"
	FieldFriendLimit        LegacyField = "friend_limit"
	FieldCustomerLimit      LegacyField = "customer_limit"
	FieldJoinOrg            LegacyField = "join_org"
	FieldCreateOrganization LegacyField = "create_organization"
	FieldOrgMemberLimit     LegacyField = "org_member_limit"
	FieldOrgProductLimit    LegacyField = "org_product_limit"
)

var entitlementFields = []LegacyField{
	FieldCreateProduct,
	FieldLimitPost,
	FieldPublicProductLimit,
	FieldCreateProductAI,
	FieldFriendLimit,
	FieldCustomerLimit,
	FieldJoinOrg,
	FieldCreateOrganization,
	FieldOrgMemberLimit,
	FieldOrgProductLimit,
}

type ValueMode string

const (
	ValueModeFeatureFlag ValueMode = "feature_flag"
	ValueModeAmount      ValueMode = "amount"
)

type ZeroSemantics string

const (
	ZeroReject    ZeroSemantics = "reject"
	ZeroOmit      ZeroSemantics = "omit"
	ZeroUnlimited ZeroSemantics = "unlimited"
)

type LegacyPlanSnapshot struct {
	ID                 uint64
	PackageName        string
	DurationDays       int64
	Price              int64
	CreateProduct      int64
	LimitPost          int64
	PublicProductLimit int64
	CreateProductAI    int64
	FriendLimit        int64
	CustomerLimit      int64
	JoinOrg            int64
	CreateOrganization int64
	OrgMemberLimit     int64
	OrgProductLimit    int64
	Color              string
	Features           []string
	Visibility         bool
}

type EntitlementDecision struct {
	SourceField LegacyField
	ValueMode   ValueMode
	Zero        ZeroSemantics
	Code        string
	Kind        productcatalog.EntitlementKind
	FeatureCode string
	MeterCode   string
	Period      productcatalog.PeriodKind
}

type PriceDecision struct {
	Code        string
	Currency    string
	Multiplier  int64
	BillingUnit string
	Quantity    int64
}

type Decision struct {
	CatalogVersion            string
	ProductCode               string
	ProductDisplayName        string
	PlanCode                  string
	PlanVersion               string
	PlanDisplayName           string
	SubjectScope              productcatalog.SubjectScope
	DurationHandling          string
	MarketingFeaturesHandling string
	RecurringPrice            *PriceDecision
	PriceReviewReason         string
	Entitlements              []EntitlementDecision
	IgnoredFields             map[LegacyField]string
}

type Draft struct {
	LegacyPlanID         uint64
	SourceChecksum       string
	TermsChecksum        string
	SubscriptionTermDays int64
	LegacyColor          string
	LegacyFeatures       []string
	LegacyVisibility     bool
	Catalog              productcatalog.Catalog
}

/** BuildDraft converts one reviewed legacy plan into a draft snapshot. */
func BuildDraft(snapshot LegacyPlanSnapshot, decision Decision) (Draft, error) {
	if snapshot.ID == 0 {
		return Draft{}, fmt.Errorf("legacy plan id is required")
	}
	if strings.TrimSpace(snapshot.PackageName) == "" {
		return Draft{}, fmt.Errorf("legacy package name is required")
	}
	if snapshot.DurationDays <= 0 {
		return Draft{}, fmt.Errorf("legacy duration must be positive")
	}
	if err := validateHandling(decision); err != nil {
		return Draft{}, err
	}

	entitlements, err := buildEntitlements(snapshot, decision)
	if err != nil {
		return Draft{}, err
	}
	prices, err := buildPrices(snapshot, decision)
	if err != nil {
		return Draft{}, err
	}

	plan := productcatalog.PlanVersion{
		ProductCode:  decision.ProductCode,
		PlanCode:     decision.PlanCode,
		Version:      decision.PlanVersion,
		DisplayName:  decision.PlanDisplayName,
		Status:       productcatalog.StatusDraft,
		SubjectScope: decision.SubjectScope,
		Entitlements: entitlements,
		Prices:       prices,
	}
	catalog := productcatalog.Catalog{
		Version: decision.CatalogVersion,
		Products: []productcatalog.Product{
			{
				Code:        decision.ProductCode,
				DisplayName: decision.ProductDisplayName,
				Plans:       []productcatalog.PlanVersion{plan},
			},
		},
	}
	if _, err := productcatalog.New(catalog); err != nil {
		return Draft{}, fmt.Errorf("invalid catalog draft: %w", err)
	}

	sourceChecksum, err := checksum(snapshot)
	if err != nil {
		return Draft{}, err
	}
	termsChecksum, err := productcatalog.TermsChecksum(plan)
	if err != nil {
		return Draft{}, err
	}

	return Draft{
		LegacyPlanID:         snapshot.ID,
		SourceChecksum:       sourceChecksum,
		TermsChecksum:        termsChecksum,
		SubscriptionTermDays: snapshot.DurationDays,
		LegacyColor:          snapshot.Color,
		LegacyFeatures:       append([]string(nil), snapshot.Features...),
		LegacyVisibility:     snapshot.Visibility,
		Catalog:              catalog,
	}, nil
}

/** SourceChecksum fingerprints the full legacy plan row. */
func SourceChecksum(snapshot LegacyPlanSnapshot) (string, error) {
	return checksum(snapshot)
}

func validateHandling(decision Decision) error {
	if decision.DurationHandling != "subscription_access_term" &&
		decision.DurationHandling != "manual_review" {
		return fmt.Errorf("duration handling requires an explicit decision")
	}
	if decision.MarketingFeaturesHandling != "catalog_display_only" &&
		decision.MarketingFeaturesHandling != "manual_review" {
		return fmt.Errorf("marketing features require an explicit decision")
	}
	if decision.RecurringPrice == nil && strings.TrimSpace(decision.PriceReviewReason) == "" {
		return fmt.Errorf("legacy price requires mapping or review reason")
	}
	if decision.RecurringPrice != nil && strings.TrimSpace(decision.PriceReviewReason) != "" {
		return fmt.Errorf("legacy price cannot be mapped and ignored together")
	}
	return nil
}

func buildEntitlements(
	snapshot LegacyPlanSnapshot,
	decision Decision,
) ([]productcatalog.Entitlement, error) {
	seen := make(map[LegacyField]struct{}, len(entitlementFields))
	entitlements := make([]productcatalog.Entitlement, 0, len(decision.Entitlements))

	for _, mapping := range decision.Entitlements {
		if !isEntitlementField(mapping.SourceField) {
			return nil, fmt.Errorf("unknown legacy field %q", mapping.SourceField)
		}
		if _, exists := seen[mapping.SourceField]; exists {
			return nil, fmt.Errorf("legacy field %q decided more than once", mapping.SourceField)
		}
		seen[mapping.SourceField] = struct{}{}

		value := snapshot.value(mapping.SourceField)
		entitlement, include, err := mapEntitlement(value, mapping)
		if err != nil {
			return nil, fmt.Errorf("legacy field %q: %w", mapping.SourceField, err)
		}
		if include {
			entitlements = append(entitlements, entitlement)
		}
	}

	ignoredFields := sortedIgnoredFields(decision.IgnoredFields)
	for _, field := range ignoredFields {
		reason := strings.TrimSpace(decision.IgnoredFields[field])
		if !isEntitlementField(field) {
			return nil, fmt.Errorf("unknown ignored legacy field %q", field)
		}
		if reason == "" {
			return nil, fmt.Errorf("ignored legacy field %q requires a reason", field)
		}
		if _, exists := seen[field]; exists {
			return nil, fmt.Errorf("legacy field %q mapped and ignored", field)
		}
		seen[field] = struct{}{}
	}

	for _, field := range entitlementFields {
		if _, exists := seen[field]; !exists {
			return nil, fmt.Errorf("legacy field %q has no migration decision", field)
		}
	}

	return entitlements, nil
}

func mapEntitlement(
	value int64,
	decision EntitlementDecision,
) (productcatalog.Entitlement, bool, error) {
	if value < 0 {
		return productcatalog.Entitlement{}, false, fmt.Errorf("negative values are invalid")
	}

	entitlement := productcatalog.Entitlement{
		Code:        decision.Code,
		Kind:        decision.Kind,
		FeatureCode: decision.FeatureCode,
		MeterCode:   decision.MeterCode,
		Period:      decision.Period,
	}

	switch decision.ValueMode {
	case ValueModeFeatureFlag:
		if decision.Kind != productcatalog.EntitlementFeatureAccess {
			return productcatalog.Entitlement{}, false, fmt.Errorf(
				"feature flag requires feature_access",
			)
		}
		if value == 0 {
			return productcatalog.Entitlement{}, false, nil
		}
		if value != 1 {
			return productcatalog.Entitlement{}, false, fmt.Errorf(
				"feature flag requires legacy value 0 or 1",
			)
		}
		return entitlement, true, nil

	case ValueModeAmount:
		if decision.Kind != productcatalog.EntitlementUsageAllowance &&
			decision.Kind != productcatalog.EntitlementCapacityLimit {
			return productcatalog.Entitlement{}, false, fmt.Errorf(
				"amount mode requires usage or capacity kind",
			)
		}
		if value > 0 {
			entitlement.Amount = value
			return entitlement, true, nil
		}
		switch decision.Zero {
		case ZeroOmit:
			return productcatalog.Entitlement{}, false, nil
		case ZeroUnlimited:
			entitlement.Unlimited = true
			return entitlement, true, nil
		default:
			return productcatalog.Entitlement{}, false, fmt.Errorf(
				"zero value semantics are ambiguous",
			)
		}

	default:
		return productcatalog.Entitlement{}, false, fmt.Errorf(
			"invalid value mode %q",
			decision.ValueMode,
		)
	}
}

func buildPrices(
	snapshot LegacyPlanSnapshot,
	decision Decision,
) ([]productcatalog.PriceItem, error) {
	if decision.RecurringPrice == nil {
		return nil, nil
	}
	price := *decision.RecurringPrice
	if price.Multiplier <= 0 {
		return nil, fmt.Errorf("price multiplier must be positive")
	}
	if snapshot.Price < 0 {
		return nil, fmt.Errorf("legacy price cannot be negative")
	}
	if snapshot.Price > math.MaxInt64/price.Multiplier {
		return nil, fmt.Errorf("legacy price conversion overflows int64")
	}
	return []productcatalog.PriceItem{
		{
			Code:        price.Code,
			Kind:        productcatalog.PriceRecurring,
			Currency:    price.Currency,
			AmountMinor: snapshot.Price * price.Multiplier,
			BillingUnit: price.BillingUnit,
			Quantity:    price.Quantity,
		},
	}, nil
}

func (snapshot LegacyPlanSnapshot) value(field LegacyField) int64 {
	switch field {
	case FieldCreateProduct:
		return snapshot.CreateProduct
	case FieldLimitPost:
		return snapshot.LimitPost
	case FieldPublicProductLimit:
		return snapshot.PublicProductLimit
	case FieldCreateProductAI:
		return snapshot.CreateProductAI
	case FieldFriendLimit:
		return snapshot.FriendLimit
	case FieldCustomerLimit:
		return snapshot.CustomerLimit
	case FieldJoinOrg:
		return snapshot.JoinOrg
	case FieldCreateOrganization:
		return snapshot.CreateOrganization
	case FieldOrgMemberLimit:
		return snapshot.OrgMemberLimit
	case FieldOrgProductLimit:
		return snapshot.OrgProductLimit
	default:
		return 0
	}
}

func isEntitlementField(field LegacyField) bool {
	for _, candidate := range entitlementFields {
		if field == candidate {
			return true
		}
	}
	return false
}

func sortedIgnoredFields(values map[LegacyField]string) []LegacyField {
	fields := make([]LegacyField, 0, len(values))
	for field := range values {
		fields = append(fields, field)
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i] < fields[j]
	})
	return fields
}

func checksum(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal checksum payload: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
