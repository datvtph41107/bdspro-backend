package commercialprofile

import (
	"context"
	"errors"
	"testing"
	"time"

	plandomain "user/internal/domain/plan"
	subscriptiondomain "user/internal/domain/subscription"
	planadmin "user/internal/usecase/plan/admin"
)

type catalogStub struct{ page planadmin.Page }

func (s catalogStub) ListPlanVersions(context.Context, planadmin.Query) (planadmin.Page, error) {
	return s.page, nil
}

type subscriptionStub struct{}

func (subscriptionStub) GetUser(context.Context, uint64) (subscriptiondomain.AdminUserProjection, error) {
	return subscriptiondomain.AdminUserProjection{ProfileID: 42}, nil
}

func activePlan(id uint64, code string, rank int32, from, until *time.Time) plandomain.PlanVersionAggregate {
	return plandomain.PlanVersionAggregate{
		ID: id, ProductCode: "qhpro", ProductStatus: plandomain.StatusActive,
		PlanCode: code, PlanStatus: plandomain.StatusActive, TierRank: rank,
		Status: plandomain.StatusActive, SubjectScope: plandomain.SubjectScopeProfile,
		EffectiveFrom: from, EffectiveUntil: until,
	}
}

func TestListAvailablePlansReturnsOnlyEffectiveTermsInTierOrder(t *testing.T) {
	now := time.Date(2026, 8, 29, 0, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)
	service := NewService(catalogStub{page: planadmin.Page{PlanVersions: []plandomain.PlanVersionAggregate{
		activePlan(2, "pro", 20, &past, nil), activePlan(1, "basic", 10, &past, nil), activePlan(3, "future", 30, &future, nil),
	}}}, subscriptionStub{}, func() time.Time { return now })

	plans, err := service.ListAvailablePlans(context.Background(), "qhpro", "profile")
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 2 || plans[0].PlanCode != "basic" || plans[1].PlanCode != "pro" {
		t.Fatalf("plans = %+v", plans)
	}
}

func TestListAvailablePlansFailsClosedOnOverlappingVersions(t *testing.T) {
	now := time.Date(2026, 8, 29, 0, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	service := NewService(catalogStub{page: planadmin.Page{PlanVersions: []plandomain.PlanVersionAggregate{
		activePlan(1, "pro", 20, &past, nil), activePlan(2, "pro", 20, &past, nil),
	}}}, subscriptionStub{}, func() time.Time { return now })
	_, err := service.ListAvailablePlans(context.Background(), "qhpro", "profile")
	if !errors.Is(err, ErrAmbiguousActiveTerms) {
		t.Fatalf("error = %v", err)
	}
}
