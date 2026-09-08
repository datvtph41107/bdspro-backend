package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	catalogdomain "user/internal/domain/plan"
	"user/internal/usecase/plan/publish"
)

var ErrPlanVersionNotFound = errors.New("plan version not found")

var (
	ErrPlanVersionNotDraft       = errors.New("only a draft plan version can be changed")
	ErrPlanVersionNotActive      = errors.New("only an active plan version can be retired")
	ErrStablePlanFieldsImmutable = errors.New("product identity and plan tier are immutable after the first version")
	ErrPlanVersionConflict       = errors.New("plan version already exists")
)

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
	CreatePlanVersionDraft(context.Context, DraftRecord) (uint64, error)
	UpdatePlanVersionDraft(context.Context, uint64, DraftRecord) error
	DeletePlanVersionDraft(context.Context, uint64) error
	RetirePlanVersion(context.Context, uint64, uint64, time.Time) error
}

// DraftCommand is the complete candidate contract entered by an administrator.
// Replacing the entire draft prevents partially updated entitlement/price data
// from becoming an implicit commercial contract.
type DraftCommand struct {
	ActorID            uint64
	ProductDisplayName string
	TierRank           int32
	Terms              catalogdomain.PlanVersion
}

// DraftRecord is the validated persistence command.
type DraftRecord struct {
	ActorID            uint64
	ProductDisplayName string
	TierRank           int32
	Terms              catalogdomain.PlanVersion
	TermsChecksum      string
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

// CreateDraft creates a new version candidate without affecting the serving catalog.
func (s *Service) CreateDraft(ctx context.Context, command DraftCommand) (catalogdomain.PlanVersionAggregate, error) {
	if s == nil || s.repository == nil {
		return catalogdomain.PlanVersionAggregate{}, errors.New("plan admin service is not configured")
	}
	record, err := validateDraftCommand(command)
	if err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}
	id, err := s.repository.CreatePlanVersionDraft(ctx, record)
	if err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}
	return s.GetPlanVersion(ctx, id)
}

// UpdateDraft atomically replaces one draft's complete terms.
func (s *Service) UpdateDraft(ctx context.Context, planVersionID uint64, command DraftCommand) (catalogdomain.PlanVersionAggregate, error) {
	if s == nil || s.repository == nil {
		return catalogdomain.PlanVersionAggregate{}, errors.New("plan admin service is not configured")
	}
	if planVersionID == 0 {
		return catalogdomain.PlanVersionAggregate{}, errors.New("plan version id is required")
	}
	record, err := validateDraftCommand(command)
	if err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}
	if err := s.repository.UpdatePlanVersionDraft(ctx, planVersionID, record); err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}
	return s.GetPlanVersion(ctx, planVersionID)
}

// ValidateDraft applies the same domain rules used by publication and returns
// the deterministic terms checksum shown to operators before publishing.
func (s *Service) ValidateDraft(ctx context.Context, planVersionID uint64) (string, error) {
	aggregate, err := s.GetPlanVersion(ctx, planVersionID)
	if err != nil {
		return "", err
	}
	if aggregate.Status != catalogdomain.StatusDraft {
		return "", ErrPlanVersionNotDraft
	}
	if aggregate.SubscriptionTermDays <= 0 {
		return "", fmt.Errorf("plan %q requires subscription_term_days", aggregate.PlanCode)
	}
	return catalogdomain.TermsChecksum(aggregate.CurrentContract())
}

// DeleteDraft removes only a never-published candidate. Active or retired
// versions remain durable evidence for orders, subscriptions and usage.
func (s *Service) DeleteDraft(ctx context.Context, planVersionID, actorID uint64) error {
	if s == nil || s.repository == nil {
		return errors.New("plan admin service is not configured")
	}
	if planVersionID == 0 {
		return errors.New("plan version id is required")
	}
	if actorID == 0 {
		return errors.New("actor id is required")
	}
	return s.repository.DeletePlanVersionDraft(ctx, planVersionID)
}

// Retire removes one published version from future catalog discovery while
// retaining the version and its terms as durable commercial evidence.
func (s *Service) Retire(ctx context.Context, planVersionID, actorID uint64) (catalogdomain.PlanVersionAggregate, error) {
	if s == nil || s.repository == nil {
		return catalogdomain.PlanVersionAggregate{}, errors.New("plan admin service is not configured")
	}
	if planVersionID == 0 {
		return catalogdomain.PlanVersionAggregate{}, errors.New("plan version id is required")
	}
	if actorID == 0 {
		return catalogdomain.PlanVersionAggregate{}, errors.New("actor id is required")
	}
	if err := s.repository.RetirePlanVersion(ctx, planVersionID, actorID, time.Now().UTC()); err != nil {
		return catalogdomain.PlanVersionAggregate{}, err
	}
	return s.GetPlanVersion(ctx, planVersionID)
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

func validateDraftCommand(command DraftCommand) (DraftRecord, error) {
	if command.ActorID == 0 {
		return DraftRecord{}, errors.New("actor id is required")
	}
	command.ProductDisplayName = strings.TrimSpace(command.ProductDisplayName)
	if command.ProductDisplayName == "" {
		return DraftRecord{}, errors.New("product display name is required")
	}
	if command.TierRank <= 0 {
		return DraftRecord{}, errors.New("tier rank must be positive")
	}
	command.Terms.ProductCode = strings.TrimSpace(command.Terms.ProductCode)
	command.Terms.PlanCode = strings.TrimSpace(command.Terms.PlanCode)
	command.Terms.Version = strings.TrimSpace(command.Terms.Version)
	command.Terms.DisplayName = strings.TrimSpace(command.Terms.DisplayName)
	command.Terms.Status = catalogdomain.StatusDraft
	command.Terms.EffectiveFrom = time.Time{}
	command.Terms.EffectiveUntil = nil
	checksum, err := catalogdomain.TermsChecksum(command.Terms)
	if err != nil {
		return DraftRecord{}, err
	}
	return DraftRecord{
		ActorID:            command.ActorID,
		ProductDisplayName: command.ProductDisplayName,
		TierRank:           command.TierRank,
		Terms:              command.Terms,
		TermsChecksum:      checksum,
	}, nil
}
