package checkoutpostgres

import (
	"testing"
	"time"

	plan "user/internal/domain/plan"
	"user/internal/usecase/subscription/checkout"
)

func TestCheckoutCommandRowRoundTripPreservesFrozenTerms(t *testing.T) {
	createdAt := time.Date(2026, 8, 19, 11, 0, 0, 123, time.UTC)
	command := checkout.CheckoutCommand{
		Subject:     checkout.Subject{Kind: checkout.SubjectProfile, ID: "42"},
		PlanCode:    "qhpro.pro",
		RequestHash: checkout.RequestHashFor(checkout.Subject{Kind: checkout.SubjectProfile, ID: "42"}, "qhpro.pro"),
		CommandKey:  "checkout-1",
		CreatedAt:   createdAt,
		Terms: checkout.PlanTerms{
			ProductID:            1,
			ProductCode:          "qhpro",
			PlanID:               2,
			PlanCode:             "qhpro.pro",
			PlanVersionID:        3,
			PlanVersion:          "2.0.0",
			TierRank:             20,
			SubscriptionTermDays: 30,
			TermsChecksum:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			SubjectScope:         plan.SubjectScopeProfile,
			Price: checkout.Money{
				Currency:    "VND",
				AmountMinor: 19900000,
			},
		},
	}

	row := checkoutSnapshotToRow(command)
	got := checkoutSnapshotFromRow(row)
	if got != command {
		t.Fatalf("round trip changed checkout command:\n got  %+v\n want %+v", got, command)
	}
}

func TestCheckoutCommandFromRowRejectsBrokenSnapshotAtApplicationBoundary(t *testing.T) {
	row := checkoutSnapshotRow{
		SubjectKind:          checkout.SubjectProfile,
		SubjectID:            "42",
		CommandKey:           "checkout-1",
		RequestHash:          checkout.RequestHashFor(checkout.Subject{Kind: checkout.SubjectProfile, ID: "42"}, "qhpro.pro"),
		ProductID:            1,
		ProductCode:          "qhpro",
		PlanID:               2,
		PlanCode:             "qhpro.pro",
		PlanVersionID:        3,
		PlanVersion:          "2.0.0",
		TierRank:             20,
		SubscriptionTermDays: 0,
		TermsChecksum:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		SubjectScope:         plan.SubjectScopeProfile,
		Currency:             "VND",
		AmountMinor:          19900000,
		CreatedAt:            time.Date(2026, 8, 19, 11, 0, 0, 0, time.UTC),
	}
	if command := checkoutSnapshotFromRow(row); command.IsValid() {
		t.Fatalf("broken durable snapshot became valid: %+v", command)
	}
}
