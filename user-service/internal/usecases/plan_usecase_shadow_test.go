package usecases

import (
	"context"
	"errors"
	"testing"

	_dto "common/domain/dto"
	"user/internal/models"
)

type planRepoStub struct {
	plans []models.Plan
	err   error
}

func (r *planRepoStub) Create(context.Context, *models.Plan) error            { return nil }
func (r *planRepoStub) CreateBatch(context.Context, []*models.Plan) error     { return nil }
func (r *planRepoStub) Update(context.Context, uint64, *models.Plan) error    { return nil }
func (r *planRepoStub) Delete(context.Context, uint64) error                  { return nil }
func (r *planRepoStub) GetAll(context.Context) ([]models.Plan, error)         { return r.plans, r.err }
func (r *planRepoStub) GetByID(context.Context, uint64) (*models.Plan, error) { return nil, r.err }
func (r *planRepoStub) GetList(context.Context, _dto.IPagable) ([]models.Plan, int64, error) {
	return r.plans, int64(len(r.plans)), r.err
}
func (r *planRepoStub) GetListByIDs(context.Context, []uint64) ([]models.Plan, error) {
	return r.plans, r.err
}
func (r *planRepoStub) GetDetail(context.Context, uint64) (*models.Plan, error) { return nil, r.err }
func (r *planRepoStub) BeforeSave(context.Context, *uint64, *models.Plan) error { return nil }
func (r *planRepoStub) AfterSave(context.Context, *uint64, *models.Plan) error  { return nil }
func (r *planRepoStub) GetVisiblePlans(context.Context) ([]models.Plan, error) {
	return append([]models.Plan(nil), r.plans...), r.err
}

type shadowObserverStub struct {
	calls int
	plans []models.Plan
}

func (o *shadowObserverStub) ObserveVisiblePlans(_ context.Context, plans []models.Plan) {
	o.calls++
	o.plans = append([]models.Plan(nil), plans...)
}

func TestGetVisiblePlansKeepsLegacyResponseAndObserves(t *testing.T) {
	plan := models.Plan{PackageName: "Legacy Pro"}
	plan.ID = 7
	repository := &planRepoStub{plans: []models.Plan{plan}}
	observer := &shadowObserverStub{}
	usecase := NewPlanUsecase(repository, observer)

	plans, errDTO := usecase.GetVisiblePlans(context.Background())
	if errDTO != nil {
		t.Fatalf("GetVisiblePlans() error = %v", errDTO)
	}
	if len(plans) != 1 || plans[0].PackageName != "Legacy Pro" {
		t.Fatalf("legacy response changed: %+v", plans)
	}
	if observer.calls != 1 || len(observer.plans) != 1 {
		t.Fatal("shadow observer was not called")
	}
	observer.plans[0].PackageName = "changed"
	if plans[0].PackageName != "Legacy Pro" {
		t.Fatal("observer mutated the legacy response")
	}
}

func TestGetVisiblePlansDoesNotObserveLegacyFailure(t *testing.T) {
	observer := &shadowObserverStub{}
	usecase := NewPlanUsecase(
		&planRepoStub{err: errors.New("legacy database unavailable")},
		observer,
	)

	plans, errDTO := usecase.GetVisiblePlans(context.Background())
	if errDTO == nil || plans != nil {
		t.Fatal("expected legacy read failure")
	}
	if observer.calls != 0 {
		t.Fatal("shadow read ran after legacy failure")
	}
}
