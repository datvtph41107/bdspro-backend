package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	catalogdomain "user/internal/domain/plan"
	"user/internal/usecase/plan/publish"
)

var ErrPlanVersionNotFound = errors.New("plan version not found")

/** Query scopes a bounded PlanVersion read. */
type Query struct {
	Page        uint32
	PageSize    uint32
	ProductCode string
	PlanCode    string
	Status      string
}

/** Page is one immutable control-plane list result. */
type Page struct {
	PlanVersions []catalogdomain.PlanVersionAggregate
	Total        uint64
}

/** Repository reads PlanVersion aggregates without changing terms. */
type Repository interface {
	GetPlanVersion(context.Context, uint64) (catalogdomain.PlanVersionAggregate, error)
	ListPlanVersions(context.Context, Query) (Page, error)
}

/** Publisher owns the guarded transactional draft-to-active transition. */
type Publisher interface {
	Publish(context.Context, publish.Command) (publish.Result, error)
}

/** Service exposes bounded PlanVersion read/publish operations. */
type Service struct {
	repository Repository
	publisher  Publisher
}

/** NewService creates a control-plane service from explicit dependencies. */
func NewService(repository Repository, publisher Publisher) *Service {
	return &Service{repository: repository, publisher: publisher}
}

/** ListPlanVersions returns a validated, bounded page of immutable terms. */
func (s *Service) ListPlanVersions(ctx context.Context, query Query) (Page, error) {
	if s == nil || s.repository == nil {
		return Page{}, errors.New("plan admin service is not configured")
	}
	query, err := normalizeQuery(query)
	if err != nil {
		return Page{}, err
	}
	return s.repository.ListPlanVersions(ctx, query)
}

/** GetPlanVersion returns one defensive aggregate snapshot. */
func (s *Service) GetPlanVersion(
	ctx context.Context,
	planVersionID uint64,
) (catalogdomain.PlanVersionAggregate, error) {
	if s == nil || s.repository == nil {
		return catalogdomain.PlanVersionAggregate{}, errors.New("plan admin service is not configured")
	}
	if planVersionID == 0 {
		return catalogdomain.PlanVersionAggregate{}, errors.New("plan version id is required")
	}
	aggregate, err := s.repository.GetPlanVersion(ctx, planVersionID)
	if errors.Is(err, publish.ErrPlanVersionNotFound) {
		return catalogdomain.PlanVersionAggregate{}, ErrPlanVersionNotFound
	}
	if err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}
	return aggregate.Copy(), nil
}

/** Publish promotes exactly one reviewed draft while retaining actor evidence. */
func (s *Service) Publish(
	ctx context.Context,
	command publish.Command,
) (publish.Result, error) {
	if s == nil || s.publisher == nil {
		return publish.Result{}, errors.New("plan publication is not configured")
	}
	return s.publisher.Publish(ctx, command)
}

func normalizeQuery(query Query) (Query, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		return Query{}, fmt.Errorf("page size must be between 1 and 100")
	}
	query.ProductCode = strings.TrimSpace(query.ProductCode)
	query.PlanCode = strings.TrimSpace(query.PlanCode)
	query.Status = strings.TrimSpace(query.Status)
	if query.Status != "" && !validStatus(catalogdomain.Status(query.Status)) {
		return Query{}, fmt.Errorf("invalid plan status %q", query.Status)
	}
	return query, nil
}

func validStatus(status catalogdomain.Status) bool {
	return status == catalogdomain.StatusDraft ||
		status == catalogdomain.StatusActive ||
		status == catalogdomain.StatusRetired
}
