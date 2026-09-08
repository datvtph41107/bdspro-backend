package models

import "testing"

func TestCatalogV2TableNames(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"product", CatalogProduct{}.TableName(), "catalog_products"},
		{"plan", CatalogPlan{}.TableName(), "catalog_plans"},
		{"plan version", CatalogPlanVersion{}.TableName(), "catalog_plan_versions"},
		{"entitlement", CatalogPlanEntitlement{}.TableName(), "catalog_plan_entitlements"},
		{"price", CatalogPriceItem{}.TableName(), "catalog_price_items"},
		{"subscription", CatalogSubscription{}.TableName(), "catalog_subscriptions"},
		{"subscription event", CatalogSubscriptionEvent{}.TableName(), "catalog_subscription_events"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("TableName() = %q, want %q", test.got, test.want)
			}
		})
	}
}

func TestSubscriptionStatusValues(t *testing.T) {
	statuses := []SubscriptionStatus{
		SubscriptionPending,
		SubscriptionActive,
		SubscriptionPastDue,
		SubscriptionCanceled,
		SubscriptionExpired,
	}
	want := []string{"pending", "active", "past_due", "canceled", "expired"}
	for i, status := range statuses {
		if string(status) != want[i] {
			t.Fatalf("status[%d] = %q, want %q", i, status, want[i])
		}
	}
}

func TestCatalogV2EnumValues(t *testing.T) {
	if CatalogSourceLegacyBackfill != "legacy_backfill" {
		t.Fatalf("legacy source = %q", CatalogSourceLegacyBackfill)
	}
	if SubscriptionSubjectProfile != "profile" ||
		SubscriptionSubjectOrganization != "organization" {
		t.Fatal("subscription subject enum changed")
	}
}
