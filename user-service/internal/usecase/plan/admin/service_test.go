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
