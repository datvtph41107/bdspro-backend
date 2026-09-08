package subscription

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	commonidentity "common/identity"
	subscriptiondomain "user/internal/domain/subscription"
)

const PermissionView = "COMMERCIAL_SUBSCRIPTION_VIEW"

var (
	ErrAdminActorRequired             = errors.New("trusted subscription admin actor identity is required")
	ErrPermissionDenied               = errors.New("subscription permission denied")
	ErrPermissionAuthorityUnavailable = errors.New("subscription permission authority unavailable")
)

// PermissionAuthorizer is the subscription consumer's narrow IAM boundary.
// Transport provenance and human authority remain separate checks.
type PermissionAuthorizer interface {
	HasPermission(ctx context.Context, actorID uint64, permissionCode string) (bool, error)
}

func RequireViewPermission(ctx context.Context, authorizer PermissionAuthorizer) (uint64, error) {
	if _, ok := commonidentity.ServiceCallerFromContext(ctx); !ok {
		return 0, ErrAdminActorRequired
	}
	actor, ok := commonidentity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(actor.TokenType, "ACCESS") {
		return 0, ErrAdminActorRequired
	}
	if authorizer == nil {
		return 0, ErrPermissionAuthorityUnavailable
	}
	allowed, err := authorizer.HasPermission(ctx, actor.ProfileID, PermissionView)
	if err != nil {
		return 0, errors.Join(ErrPermissionAuthorityUnavailable, err)
	}
	if !allowed {
		return 0, ErrPermissionDenied
	}
	return actor.ProfileID, nil
}

type AdminQuery struct {
	Page        uint32
	PageSize    uint32
	ProfileID   uint64
	Status      string
	ProductCode string
	PlanCode    string
}

type AdminPage struct {
	Subscriptions []subscriptiondomain.AdminProjection
	Total         uint64
	Page          uint32
	PageSize      uint32
}

// AdminRepository is owned by the subscription usecase consumer.
type AdminRepository interface {
	List(context.Context, AdminQuery) (AdminPage, error)
	ListEffectiveByProfiles(context.Context, []uint64) ([]subscriptiondomain.AdminUserProjection, error)
}

type AdminService struct{ repository AdminRepository }

func NewAdminService(repository AdminRepository) *AdminService {
	return &AdminService{repository: repository}
}

func (s *AdminService) List(ctx context.Context, query AdminQuery) (AdminPage, error) {
	if s == nil || s.repository == nil {
		return AdminPage{}, errors.New("subscription admin projection is not configured")
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		return AdminPage{}, fmt.Errorf("page size must be between 1 and 100")
	}
	query.Status = strings.TrimSpace(query.Status)
	query.ProductCode = strings.TrimSpace(query.ProductCode)
	query.PlanCode = strings.TrimSpace(query.PlanCode)
	return s.repository.List(ctx, query)
}

func (s *AdminService) GetUser(ctx context.Context, profileID uint64) (subscriptiondomain.AdminUserProjection, error) {
	if profileID == 0 {
		return subscriptiondomain.AdminUserProjection{}, errors.New("profile id is required")
	}
	items, err := s.ListUsers(ctx, []uint64{profileID})
	if err != nil {
		return subscriptiondomain.AdminUserProjection{}, err
	}
	return items[0], nil
}

func (s *AdminService) ListUsers(ctx context.Context, profileIDs []uint64) ([]subscriptiondomain.AdminUserProjection, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("subscription admin projection is not configured")
	}
	if len(profileIDs) == 0 || len(profileIDs) > 100 {
		return nil, errors.New("between 1 and 100 profile ids are required")
	}
	seen := make(map[uint64]struct{}, len(profileIDs))
	normalized := make([]uint64, 0, len(profileIDs))
	for _, id := range profileIDs {
		if id == 0 {
			return nil, errors.New("profile id is required")
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		normalized = append(normalized, id)
	}
	rows, err := s.repository.ListEffectiveByProfiles(ctx, normalized)
	if err != nil {
		return nil, err
	}
	byProfile := make(map[uint64]subscriptiondomain.AdminUserProjection, len(rows))
	for _, row := range rows {
		byProfile[row.ProfileID] = row
	}
	result := make([]subscriptiondomain.AdminUserProjection, 0, len(normalized))
	for _, id := range normalized {
		if row, ok := byProfile[id]; ok {
			result = append(result, row)
		} else {
			result = append(result, subscriptiondomain.AdminUserProjection{ProfileID: id})
		}
	}
	return result, nil
}

func ParseProfileID(subjectID string) (uint64, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(subjectID), 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid profile subscription subject")
	}
	return id, nil
}
