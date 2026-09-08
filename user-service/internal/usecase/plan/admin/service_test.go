package admin

import (
	"context"
	"errors"
	"testing"
	"time"

	catalogdomain "user/internal/domain/plan"
	"user/internal/usecase/plan/publish"
)

type fakeRepository struct {
	aggregate catalogdomain.PlanVersionAggregate
	page      Page
	err       error
	query     Query
	created   *DraftRecord
	updated   *DraftRecord
	deletedID uint64
	retiredID uint64
}

func (f *fakeRepository) GetPlanVersion(
	context.Context,
	uint64,
) (catalogdomain.PlanVersionAggregate, error) {
	if f.err != nil {
		return catalogdomain.PlanVersionAggregate{}, f.err
	}
	return f.aggregate.Copy(), nil
}

func (f *fakeRepository) ListPlanVersions(_ context.Context, query Query) (Page, error) {
	f.query = query
	if f.err != nil {
		return Page{}, f.err
	}
	return f.page, nil
}

func (f *fakeRepository) CreatePlanVersionDraft(_ context.Context, record DraftRecord) (uint64, error) {
	f.created = &record
	if f.err != nil {
		return 0, f.err
	}
	f.aggregate = aggregateFromDraft(17, record)
	return 17, nil
}

func (f *fakeRepository) UpdatePlanVersionDraft(_ context.Context, id uint64, record DraftRecord) error {
	f.updated = &record
	if f.err != nil {
		return f.err
	}
	f.aggregate = aggregateFromDraft(id, record)
	return nil
}

func (f *fakeRepository) DeletePlanVersionDraft(_ context.Context, id uint64) error {
	f.deletedID = id
	return f.err
}

func (f *fakeRepository) RetirePlanVersion(_ context.Context, id, _ uint64, retiredAt time.Time) error {
	f.retiredID = id
	if f.err != nil {
		return f.err
	}
	f.aggregate.Status = catalogdomain.StatusRetired
	f.aggregate.EffectiveUntil = &retiredAt
	return nil
}

type fakePublisher struct {
	command publish.Command
	err     error
}

func (f *fakePublisher) Publish(
	_ context.Context,
	command publish.Command,
) (publish.Result, error) {
	f.command = command
	if f.err != nil {
		return publish.Result{}, f.err
	}
	return publish.Result{PlanVersionID: command.PlanVersionID}, nil
}

func TestListPlanVersionsNormalizesAndBoundsQuery(t *testing.T) {
	repository := &fakeRepository{page: Page{Total: 1}}
	service := NewService(repository, &fakePublisher{})
	_, err := service.ListPlanVersions(context.Background(), Query{
		ProductCode: " qhpro ",
		PlanCode:    " qhpro.pro ",
		Status:      "draft",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repository.query.Page != 1 || repository.query.PageSize != 20 ||
		repository.query.ProductCode != "qhpro" || repository.query.PlanCode != "qhpro.pro" {
		t.Fatalf("normalized query = %#v", repository.query)
	}

	_, err = service.ListPlanVersions(context.Background(), Query{PageSize: 101})
	if err == nil {
		t.Fatal("expected page size validation error")
	}
}

func TestGetPlanVersionReturnsCopyAndMapsNotFound(t *testing.T) {
	from := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	repository := &fakeRepository{aggregate: catalogdomain.PlanVersionAggregate{
		ID:            9,
		Status:        catalogdomain.StatusDraft,
		EffectiveFrom: &from,
	}}
	service := NewService(repository, &fakePublisher{})
	aggregate, err := service.GetPlanVersion(context.Background(), 9)
	if err != nil {
		t.Fatal(err)
	}
	*aggregate.EffectiveFrom = aggregate.EffectiveFrom.Add(time.Hour)
	if !repository.aggregate.EffectiveFrom.Equal(from) {
		t.Fatal("service leaked mutable aggregate time")
	}

	repository.err = publish.ErrPlanVersionNotFound
	_, err = service.GetPlanVersion(context.Background(), 9)
	if !errors.Is(err, ErrPlanVersionNotFound) {
		t.Fatalf("error = %v", err)
	}
}

func TestPublishDelegatesOnlyToPublisher(t *testing.T) {
	publisher := &fakePublisher{}
	service := NewService(&fakeRepository{}, publisher)
	command := publish.Command{PlanVersionID: 9, ActorID: 42}
	if _, err := service.Publish(context.Background(), command); err != nil {
		t.Fatal(err)
	}
	if publisher.command != command {
		t.Fatalf("command = %#v, want %#v", publisher.command, command)
	}
}

func TestDraftLifecycleValidatesAndDelegates(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository, &fakePublisher{})
	command := validDraftCommand()

	created, err := service.CreateDraft(context.Background(), command)
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != 17 || repository.created == nil || len(repository.created.TermsChecksum) != 64 {
		t.Fatalf("create result=%#v record=%#v", created, repository.created)
	}

	command.Terms.DisplayName = "QHPro Basic mới"
	updated, err := service.UpdateDraft(context.Background(), 17, command)
	if err != nil {
		t.Fatal(err)
	}
	if updated.DisplayName != command.Terms.DisplayName || repository.updated == nil {
		t.Fatalf("update result=%#v", updated)
	}
	checksum, err := service.ValidateDraft(context.Background(), 17)
	if err != nil || len(checksum) != 64 {
		t.Fatalf("validate checksum=%q err=%v", checksum, err)
	}
	if err := service.DeleteDraft(context.Background(), 17, 42); err != nil || repository.deletedID != 17 {
		t.Fatalf("delete id=%d err=%v", repository.deletedID, err)
	}
}

func TestDraftRejectsInvalidCommercialTerms(t *testing.T) {
	service := NewService(&fakeRepository{}, &fakePublisher{})
	command := validDraftCommand()
	command.Terms.Operations[0].MeterCode = "missing.meter"
	if _, err := service.CreateDraft(context.Background(), command); err == nil {
		t.Fatal("CreateDraft accepted an operation without a matching usage allowance")
	}
}

func TestRetirePreservesAggregateAndRemovesItFromServingState(t *testing.T) {
	repository := &fakeRepository{aggregate: catalogdomain.PlanVersionAggregate{ID: 17, Status: catalogdomain.StatusActive}}
	service := NewService(repository, &fakePublisher{})
	retired, err := service.Retire(context.Background(), 17, 42)
	if err != nil {
		t.Fatal(err)
	}
	if repository.retiredID != 17 || retired.Status != catalogdomain.StatusRetired || retired.EffectiveUntil == nil {
		t.Fatalf("retired aggregate=%#v repository=%#v", retired, repository)
	}
}

func validDraftCommand() DraftCommand {
	return DraftCommand{
		ActorID: 42, ProductDisplayName: "QHPro", TierRank: 10,
		Terms: catalogdomain.PlanVersion{
			ProductCode: "qhpro", PlanCode: "qhpro.basic", Version: "2.0.0",
			DisplayName: "QHPro Basic", Status: catalogdomain.StatusDraft,
			SubjectScope: catalogdomain.SubjectScopeProfile, SubscriptionTermDays: 30,
			Entitlements: []catalogdomain.Entitlement{
				{Code: "workspace.report.feature", Kind: catalogdomain.EntitlementFeatureAccess, FeatureCode: "workspace.report.generate", Period: catalogdomain.PeriodNone},
				{Code: "workspace.report.allowance", Kind: catalogdomain.EntitlementUsageAllowance, MeterCode: "workspace.report_generation.accepted", Amount: 10, Period: catalogdomain.PeriodSubscriptionCycle},
			},
			Operations: []catalogdomain.OperationBinding{{Code: "workspace.report.generate", FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1}},
			Prices:     []catalogdomain.PriceItem{{Code: "qhpro.basic.monthly", Kind: catalogdomain.PriceRecurring, Currency: "VND", AmountMinor: 99000, BillingUnit: "subscription_cycle", Quantity: 1}},
		},
	}
}

func aggregateFromDraft(id uint64, record DraftRecord) catalogdomain.PlanVersionAggregate {
	return catalogdomain.PlanVersionAggregate{
		ID: id, ProductCode: record.Terms.ProductCode, ProductDisplayName: record.ProductDisplayName,
		PlanCode: record.Terms.PlanCode, Version: record.Terms.Version, DisplayName: record.Terms.DisplayName,
		Status: catalogdomain.StatusDraft, SubjectScope: record.Terms.SubjectScope,
		SubscriptionTermDays: record.Terms.SubscriptionTermDays, TierRank: record.TierRank,
		TermsChecksum: record.TermsChecksum, Entitlements: record.Terms.Entitlements,
		Operations: record.Terms.Operations, Prices: record.Terms.Prices,
	}
}
