package services

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"auth/internal"
	"auth/internal/permissioncatalog"

	_errors "common/errors"
	"common/identity"
	_utils "common/utils"
	"pb/clients"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

/**
 * MatchMode makes multi-permission semantics explicit.
 */
type MatchMode string

const (
	MatchAll MatchMode = "ALL"
	MatchAny MatchMode = "ANY"
)

/**
 * SnapshotStatus exposes authorization freshness for readiness and audit.
 */
type SnapshotStatus struct {
	Version    string
	LoadedAt   time.Time
	RoleCount  int
	CodeCount  int
	Ready      bool
	Stale      bool
	LastError  string
	MaxStaleAt time.Time
}

type permissionSnapshot struct {
	catalog  permissioncatalog.Snapshot
	loadedAt time.Time
}

type roleCacheEntry struct {
	roleIDs   []uint64
	expiresAt time.Time
}

type PermissionConfig struct {
	RefreshInterval time.Duration
	MaxStaleness    time.Duration
	RoleCacheTTL    time.Duration
	RequestTimeout  time.Duration
}

/**
 * PermissionService is the authorization decision owner.
 */
type PermissionService struct {
	mu       sync.RWMutex
	snapshot permissionSnapshot
	lastErr  error

	roleMu    sync.Mutex
	roleCache map[uint64]roleCacheEntry

	userClient      *clients.UserGrpcClient
	refreshInterval time.Duration
	maxStaleness    time.Duration
	roleCacheTTL    time.Duration
	requestTimeout  time.Duration
}

/**
 * NewPermissionService constructs permission authority state without starting
 * process work. The Auth process root owns Run(ctx).
 */
func NewPermissionService(
	userClient *clients.UserGrpcClient,
	cfg PermissionConfig,
) *PermissionService {
	return &PermissionService{
		userClient:      userClient,
		roleCache:       make(map[uint64]roleCacheEntry),
		refreshInterval: cfg.RefreshInterval,
		maxStaleness:    cfg.MaxStaleness,
		roleCacheTTL:    cfg.RoleCacheTTL,
		requestTimeout:  cfg.RequestTimeout,
	}
}

func (s *PermissionService) Run(ctx context.Context) {
	retryDelay := 500 * time.Millisecond
	for {
		refreshCtx, cancel := context.WithTimeout(ctx, s.requestTimeout)
		err := s.refresh(refreshCtx)
		cancel()
		if err == nil {
			break
		}
		s.setLastError(err)
		slog.Warn(
			"auth permission initial refresh failed",
			slog.Any("error", err),
			slog.Duration("retry_delay", retryDelay),
		)

		timer := time.NewTimer(retryDelay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		}
		if retryDelay < 10*time.Second {
			retryDelay *= 2
		}
	}

	ticker := time.NewTicker(s.refreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			refreshCtx, cancel := context.WithTimeout(ctx, s.requestTimeout)
			err := s.refresh(refreshCtx)
			cancel()
			if err != nil {
				s.setLastError(err)
				slog.Warn(
					"auth permission refresh failed",
					slog.Any("error", err),
				)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *PermissionService) refresh(ctx context.Context) error {
	if s == nil || s.userClient == nil || s.userClient.InternalClient == nil || s.userClient.PermissionClient == nil {
		return fmt.Errorf("user IAM clients are not ready")
	}

	roles, err := s.userClient.GetRolePermissions(ctx, &sharepb.IdRequest{})
	if err != nil {
		return fmt.Errorf("load role permissions: %w", err)
	}
	permissions, err := s.userClient.GetAllPermissions(ctx)
	if err != nil {
		return fmt.Errorf("load permission codes: %w", err)
	}

	next, err := buildPermissionSnapshot(roles.GetData(), permissions)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.snapshot = next
	s.lastErr = nil
	s.mu.Unlock()
	slog.Info(
		"auth permission snapshot refreshed",
		slog.String("version", next.catalog.Version()),
		slog.Int("roles", next.catalog.RoleCount()),
		slog.Int("codes", next.catalog.CodeCount()),
	)
	return nil
}

func buildPermissionSnapshot(
	items []*userpb.RolePermissionItem,
	permissions []*sharepb.Permission,
) (permissionSnapshot, error) {
	catalogPermissions := make([]permissioncatalog.Permission, 0, len(permissions))
	for _, permission := range permissions {
		if permission == nil || permission.Id == 0 {
			continue
		}
		catalogPermissions = append(catalogPermissions, permissioncatalog.Permission{
			ID:   permission.Id,
			Code: permission.Key,
		})
	}

	rolePermissions := make([]permissioncatalog.RolePermissions, 0, len(items))
	for _, item := range items {
		if item == nil || item.Role == nil || item.Role.Id == 0 {
			continue
		}
		rolePermissions = append(rolePermissions, permissioncatalog.RolePermissions{
			RoleID:        item.Role.Id,
			PermissionIDs: append([]uint64(nil), item.PermissionIds...),
		})
	}

	catalog, err := permissioncatalog.Build(rolePermissions, catalogPermissions)
	if err != nil {
		return permissionSnapshot{}, fmt.Errorf("build permission catalog: %w", err)
	}
	return permissionSnapshot{catalog: catalog, loadedAt: time.Now().UTC()}, nil
}

func (s *PermissionService) setLastError(err error) {
	s.mu.Lock()
	s.lastErr = err
	s.mu.Unlock()
}

/**
 * Status returns freshness without exposing mutable maps.
 */
func (s *PermissionService) Status(now time.Time) SnapshotStatus {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	s.mu.RLock()
	snapshot := s.snapshot
	lastErr := s.lastErr
	s.mu.RUnlock()
	status := SnapshotStatus{
		Version:    snapshot.catalog.Version(),
		LoadedAt:   snapshot.loadedAt,
		RoleCount:  snapshot.catalog.RoleCount(),
		CodeCount:  snapshot.catalog.CodeCount(),
		Ready:      !snapshot.loadedAt.IsZero(),
		MaxStaleAt: snapshot.loadedAt.Add(s.maxStaleness),
	}
	status.Stale = !status.Ready || !now.Before(status.MaxStaleAt)
	if lastErr != nil {
		status.LastError = lastErr.Error()
	}
	return status
}

func (s *PermissionService) resolvePermissionIDs(
	permissionIDs []uint32,
	permissionCodes []string,
) ([]uint32, error) {
	s.mu.RLock()
	snapshot := s.snapshot
	s.mu.RUnlock()
	if snapshot.loadedAt.IsZero() || time.Since(snapshot.loadedAt) >= s.maxStaleness {
		return nil, _errors.ReturnError(service.PermissionSnapshotUnavailable)
	}

	set := make(map[uint32]struct{}, len(permissionIDs)+len(permissionCodes))
	for _, permissionID := range permissionIDs {
		if permissionID == 0 || !snapshot.catalog.ContainsID(permissionID) {
			return nil, _errors.ReturnError(service.PermissionIDNotFound)
		}
		set[permissionID] = struct{}{}
	}
	for _, rawCode := range permissionCodes {
		code := strings.TrimSpace(rawCode)
		if code == "" {
			return nil, _errors.ReturnError(service.PermissionCodeInvalid)
		}
		permissionID, ok := snapshot.catalog.ResolveCode(code)
		if !ok {
			return nil, _errors.ReturnError(service.PermissionCodeDenied)
		}
		set[permissionID] = struct{}{}
	}

	result := make([]uint32, 0, len(set))
	for permissionID := range set {
		result = append(result, permissionID)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result, nil
}

func profileIDFromContext(ctx context.Context) uint64 {
	if actor, ok := identity.ActorFromContext(ctx); ok && actor.ProfileID != 0 {
		return actor.ProfileID
	}
	return _utils.GetProfileIdWithContext(ctx)
}

func (s *PermissionService) resolveRoleIDs(ctx context.Context, profileID uint64) ([]uint64, error) {
	now := time.Now()
	s.roleMu.Lock()
	cached, ok := s.roleCache[profileID]
	if ok && now.Before(cached.expiresAt) {
		roleIDs := append([]uint64(nil), cached.roleIDs...)
		s.roleMu.Unlock()
		return roleIDs, nil
	}
	s.roleMu.Unlock()

	callCtx, cancel := context.WithTimeout(ctx, s.requestTimeout)
	defer cancel()
	roleIDs, err := s.userClient.GetRoleIdsByProfileId(callCtx, profileID)
	if err != nil {
		return nil, _errors.ReturnError(service.RoleAssignmentUnavailable)
	}
	roleIDs = normalizeRoleIDs(roleIDs)
	s.roleMu.Lock()
	s.roleCache[profileID] = roleCacheEntry{
		roleIDs:   append([]uint64(nil), roleIDs...),
		expiresAt: now.Add(s.roleCacheTTL),
	}
	s.roleMu.Unlock()
	return roleIDs, nil
}

func normalizeRoleIDs(roleIDs []uint64) []uint64 {
	set := make(map[uint64]struct{}, len(roleIDs))
	for _, roleID := range roleIDs {
		if roleID != 0 {
			set[roleID] = struct{}{}
		}
	}
	result := make([]uint64, 0, len(set))
	for roleID := range set {
		result = append(result, roleID)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

/**
 * RequiredPermissions evaluates current assignments with explicit semantics.
 */
func (s *PermissionService) RequiredPermissions(
	ctx context.Context,
	permissionIDs []uint32,
	permissionCodes []string,
	mode MatchMode,
) error {
	requested, err := s.resolvePermissionIDs(permissionIDs, permissionCodes)
	if err != nil {
		return err
	}
	if len(requested) == 0 {
		return nil
	}
	if mode != MatchAll && mode != MatchAny {
		return _errors.ReturnError(service.PermissionMatchModeInvalid)
	}

	profileID := profileIDFromContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(service.Unauthenticated)
	}
	roleIDs, err := s.resolveRoleIDs(ctx, profileID)
	if err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return _errors.ReturnError(service.PermissionDenied)
	}

	allowed := false
	if mode == MatchAll {
		allowed = s.hasAllPermissions(roleIDs, requested)
	} else {
		allowed = s.hasAnyPermission(roleIDs, requested)
	}
	if !allowed {
		return _errors.ReturnError(service.PermissionDenied)
	}
	return nil
}

func (s *PermissionService) hasAllPermissions(roleIDs []uint64, permissionIDs []uint32) bool {
	for _, permissionID := range permissionIDs {
		found := false
		for _, roleID := range roleIDs {
			if s.roleHasPermission(roleID, permissionID) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (s *PermissionService) hasAnyPermission(roleIDs []uint64, permissionIDs []uint32) bool {
	for _, permissionID := range permissionIDs {
		for _, roleID := range roleIDs {
			if s.roleHasPermission(roleID, permissionID) {
				return true
			}
		}
	}
	return false
}

func (s *PermissionService) roleHasPermission(roleID uint64, permissionID uint32) bool {
	s.mu.RLock()
	snapshot := s.snapshot
	s.mu.RUnlock()
	return snapshot.catalog.RoleHas(roleID, permissionID)
}

/**
 * InvalidateProfileRoles limits revocation staleness after IAM changes.
 */
func (s *PermissionService) InvalidateProfileRoles(profileID uint64) {
	s.roleMu.Lock()
	delete(s.roleCache, profileID)
	s.roleMu.Unlock()
}

/**
 * ParseProfileID accepts administrative invalidation input.
 */
func ParseProfileID(value string) (uint64, error) {
	profileID, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil || profileID == 0 {
		return 0, fmt.Errorf("invalid profile ID")
	}
	return profileID, nil
}
