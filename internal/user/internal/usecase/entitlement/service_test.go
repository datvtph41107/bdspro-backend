package evaluate

import (
	"common/operation"
	"context"
	"errors"
	"testing"
	"time"
	"user/internal/domain/entitlement"
)

type fakeSubscriptionStore struct {
	subscription Subscription
	found        bool
	err          error
}

func (f fakeSubscriptionStore) FindSubscriptionForOperation(
	ctx context.Context,
	subjectType string,
	subjectID string,
	operationCode operation.Code,
	now time.Time,
) (Subscription, bool, error) {
	return f.subscription, f.found, f.err
}

type fakePlanAccessStore struct {
	access PlanAccess
	found  bool
	err    error
}

func (f fakePlanAccessStore) FindPlanAccess(
	ctx context.Context,
	planVersionID uint64,
	operationCode operation.Code,
) (PlanAccess, bool, error) {
	return f.access, f.found, f.err
}

func TestGetAccessReturnsQuotaLimit(t *testing.T) {
	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}

	service := NewService(
		fakeSubscriptionStore{
			found: true,
			subscription: Subscription{
				ID:            77,
				Subject:       subject,
				PlanVersionID: 8,
				PlanCode:      "qhpro.pro",
				PlanVersion:   "1.0.0",
				Status:        "active",
				StartedAt:     periodStart,
				PeriodStart:   periodStart,
				PeriodEnd:     periodEnd,
			},
		},
		fakePlanAccessStore{
			found: true,
			access: PlanAccess{
				Allowed:        true,
				FeatureCode:    "workspace.report.generate",
				MeterCode:      "workspace.report_generation.accepted",
				UnitsPerAction: 1,
				Limit:          100,
				Period:         access.PeriodSubscriptionCycle,
			},
		},
	)

	result, err := service.GetAccess(
		context.Background(),
		subject,
		operation.Code("workspace.report.generate"),
		now,
	)
	if err != nil {
		t.Fatalf("GetAccess() error = %v", err)
	}
	if !result.Allowed || result.Limit != 100 || result.SubscriptionID != 77 {
		t.Fatalf("unexpected access result: %+v", result)
	}
	if !result.UsesQuota() {
		t.Fatal("expected access to use quota")
	}
	if result.Metering.MeterCode != "workspace.report_generation.accepted" || result.Metering.UnitsPerAction != 1 {
		t.Fatalf("missing metering evidence: %+v", result.Metering)
	}
}

func TestGetAccessDeniesWhenSubscriptionMissing(t *testing.T) {
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	service := NewService(
		fakeSubscriptionStore{found: false},
		fakePlanAccessStore{},
	)

	result, err := service.GetAccess(
		context.Background(),
		subject,
		operation.Code("workspace.report.generate"),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("GetAccess() error = %v", err)
	}
	if result.Allowed {
		t.Fatal("expected missing subscription to be denied")
	}
}

func TestGetAccessFailsWhenUsagePeriodBoundaryIsUnknown(t *testing.T) {
	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	service := NewService(
		fakeSubscriptionStore{
			found: true,
			subscription: Subscription{
				ID:            77,
				Subject:       subject,
				PlanVersionID: 8,
				Status:        "active",
				StartedAt:     now.Add(-time.Hour),
				PeriodStart:   time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:     time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		fakePlanAccessStore{
			found: true,
			access: PlanAccess{
				Allowed:        true,
				FeatureCode:    "workspace.report.generate",
				MeterCode:      "workspace.report_generation.accepted",
				UnitsPerAction: 1,
				Limit:          10,
				Period:         access.PeriodDay,
			},
		},
	)

	_, err := service.GetAccess(context.Background(), subject, operation.Code("workspace.report.generate"), now)
	if !errors.Is(err, ErrUsagePeriod) {
		t.Fatalf("GetAccess() error = %v, want ErrUsagePeriod", err)
	}
}
