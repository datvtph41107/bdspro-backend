package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"common/productcatalog"
	"user/internal/modules/catalog/domain"
	"user/internal/modules/catalog/publication"
	"user/internal/modules/catalog/shadow"
)

var ErrPlanVersionNotFound = errors.New("catalog plan version not found")

/** Query scopes a bounded catalog control-plane read. */
type Query struct {
	Page        uint32
	PageSize    uint32
	ProductCode string
	PlanCode    string
	Status      string
}

/** Page is one immutable control-plane list result. */
type Page struct {
	PlanVersions []domain.PlanVersionAggregate
	Total        uint64
}

/** Repository reads catalog aggregates without changing terms. */
type Repository interface {
	GetPlanVersion(context.Context, uint64) (domain.PlanVersionAggregate, error)
	ListPlanVersions(context.Context, Query) (Page, error)
}

/** Publisher owns the guarded transactional draft-to-active transition. */
type Publisher interface {
	Publish(context.Context, publication.Command) (publication.Result, error)
}

/** Service exposes the small read/publish catalog control-plane slice. */
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
		return Page{}, errors.New("catalog admin service is not configured")
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
) (domain.PlanVersionAggregate, error) {
	if s == nil || s.repository == nil {
		return domain.PlanVersionAggregate{}, errors.New("catalog admin service is not configured")
	}
	if planVersionID == 0 {
		return domain.PlanVersionAggregate{}, errors.New("plan version id is required")
	}
	aggregate, err := s.repository.GetPlanVersion(ctx, planVersionID)
	if errors.Is(err, publication.ErrPlanVersionNotFound) ||
		errors.Is(err, shadow.ErrNotFound) {
		return domain.PlanVersionAggregate{}, ErrPlanVersionNotFound
	}
	if err != nil {
		return domain.PlanVersionAggregate{}, err
	}
	return aggregate.Copy(), nil
}

/** Publish promotes exactly one reviewed draft while retaining actor evidence. */
func (s *Service) Publish(
	ctx context.Context,
	command publication.Command,
) (publication.Result, error) {
	if s == nil || s.publisher == nil {
		return publication.Result{}, errors.New("catalog publication is not configured")
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
	if query.Status != "" && !validStatus(productcatalog.Status(query.Status)) {
		return Query{}, fmt.Errorf("invalid catalog status %q", query.Status)
	}
	return query, nil
}

func validStatus(status productcatalog.Status) bool {
	return status == productcatalog.StatusDraft ||
		status == productcatalog.StatusActive ||
		status == productcatalog.StatusRetired
}
