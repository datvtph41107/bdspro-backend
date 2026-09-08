package commercialprofile

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	plandomain "user/internal/domain/plan"
	subscriptiondomain "user/internal/domain/subscription"
	planadmin "user/internal/usecase/plan/admin"
	subscriptionusecase "user/internal/usecase/subscription"
)

var (
	ErrUnavailable          = errors.New("commercial profile is unavailable")
	ErrAmbiguousActiveTerms = errors.New("multiple active plan versions overlap")
)

type CatalogReader interface {
	ListPlanVersions(context.Context, planadmin.Query) (planadmin.Page, error)
}

type SubscriptionReader interface {
	GetUser(context.Context, uint64) (subscriptiondomain.AdminUserProjection, error)
}

// Service returns only User-owned commercial projections. It never reads
// Payment evidence, TQD usage, or Redis quota state.
type Service struct {
	catalog       CatalogReader
	subscriptions SubscriptionReader
	now           func() time.Time
}

func NewService(catalog CatalogReader, subscriptions SubscriptionReader, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{catalog: catalog, subscriptions: subscriptions, now: now}
}

func (s *Service) ListAvailablePlans(ctx context.Context, productCode, subjectKind string) ([]plandomain.PlanVersionAggregate, error) {
	if s == nil || s.catalog == nil || s.now == nil {
		return nil, ErrUnavailable
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	productCode = strings.TrimSpace(productCode)
	subjectKind = strings.TrimSpace(subjectKind)
	if subjectKind == "" {
		subjectKind = string(plandomain.SubjectScopeProfile)
	}
	if subjectKind != string(plandomain.SubjectScopeProfile) && subjectKind != string(plandomain.SubjectScopeOrganization) {
		return nil, errors.New("invalid subject kind")
	}
	page, err := s.catalog.ListPlanVersions(ctx, planadmin.Query{
		Page: 1, PageSize: 100, ProductCode: productCode, Status: string(plandomain.StatusActive),
	})
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	byPlan := make(map[string]plandomain.PlanVersionAggregate)
	for _, candidate := range page.PlanVersions {
		if !available(candidate, subjectKind, now) {
			continue
		}
		key := candidate.ProductCode + "\x00" + candidate.PlanCode
		if _, exists := byPlan[key]; exists {
			return nil, ErrAmbiguousActiveTerms
		}
		byPlan[key] = candidate.Copy()
	}
	result := make([]plandomain.PlanVersionAggregate, 0, len(byPlan))
	for _, item := range byPlan {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].ProductCode != result[j].ProductCode {
			return result[i].ProductCode < result[j].ProductCode
		}
		if result[i].TierRank != result[j].TierRank {
			return result[i].TierRank < result[j].TierRank
		}
		return result[i].PlanCode < result[j].PlanCode
	})
	return result, nil
}

func (s *Service) GetProfile(ctx context.Context, profileID uint64) (subscriptiondomain.AdminUserProjection, error) {
	if s == nil || s.subscriptions == nil {
		return subscriptiondomain.AdminUserProjection{}, ErrUnavailable
	}
	if profileID == 0 {
		return subscriptiondomain.AdminUserProjection{}, errors.New("profile id is required")
	}
	return s.subscriptions.GetUser(ctx, profileID)
}

func available(item plandomain.PlanVersionAggregate, subjectKind string, now time.Time) bool {
	if item.Status != plandomain.StatusActive || item.ProductStatus != plandomain.StatusActive || item.PlanStatus != plandomain.StatusActive {
		return false
	}
	if item.SubjectScope != plandomain.SubjectScopeAny && string(item.SubjectScope) != subjectKind {
		return false
	}
	if item.EffectiveFrom == nil || item.EffectiveFrom.After(now) {
		return false
	}
	return item.EffectiveUntil == nil || item.EffectiveUntil.After(now)
}

var _ CatalogReader = (*planadmin.Service)(nil)
var _ SubscriptionReader = (*subscriptionusecase.AdminService)(nil)
