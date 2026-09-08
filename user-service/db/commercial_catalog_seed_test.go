package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"user/internal/domain/plan"
)

func TestCommercialCatalogSeedChecksumsMatchDomainContract(t *testing.T) {
	t.Parallel()

	contents, err := os.ReadFile(filepath.Join("..", "database", "migrations", "000003_seed_qhpro_commercial_catalog.up.sql"))
	require.NoError(t, err)
	migration := string(contents)

	for _, terms := range commercialSeedTerms() {
		checksum, checksumErr := plan.TermsChecksum(terms)
		require.NoError(t, checksumErr)
		require.Contains(t, migration, checksum, "seed checksum must be produced by plan.TermsChecksum")
		require.Contains(t, migration, terms.PlanCode)
	}

	require.Equal(t, 1, strings.Count(migration, "'recurring'"), "seed must have one canonical recurring-price insert")
	require.Contains(t, migration, "99000::BIGINT")
	require.Contains(t, migration, "199000::BIGINT")
}

func commercialSeedTerms() []plan.PlanVersion {
	return []plan.PlanVersion{
		commercialSeedPlan("qhpro.basic", "Cơ bản", 3, 99000),
		commercialSeedPlan("qhpro.pro", "Chuyên nghiệp", 100, 199000),
	}
}

func commercialSeedPlan(code, displayName string, allowance, amountMinor int64) plan.PlanVersion {
	return plan.PlanVersion{
		ProductCode:          "qhpro",
		PlanCode:             code,
		Version:              "1.0.0",
		DisplayName:          displayName,
		Status:               plan.StatusDraft,
		SubjectScope:         plan.SubjectScopeAny,
		SubscriptionTermDays: 30,
		Entitlements: []plan.Entitlement{
			{Code: "workspace.report.feature", Kind: plan.EntitlementFeatureAccess, FeatureCode: "workspace.report.generate", Period: plan.PeriodNone},
			{Code: "workspace.report.allowance", Kind: plan.EntitlementUsageAllowance, MeterCode: "workspace.report_generation.accepted", Amount: allowance, Period: plan.PeriodSubscriptionCycle},
		},
		Operations: []plan.OperationBinding{
			{Code: "workspace.report.generate", FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1},
		},
		Prices: []plan.PriceItem{
			{Code: code + ".monthly", Kind: plan.PriceRecurring, Currency: "VND", AmountMinor: amountMinor, BillingUnit: "subscription.cycle", Quantity: 1},
		},
	}
}
