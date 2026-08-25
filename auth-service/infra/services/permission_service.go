package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_algorithm "common/algorithm"
	_enum "common/domain/enum"
	_errors "common/errors"
	_identitycontext "common/identitycontext"
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
	rolePermissions  map[uint64]*_algorithm.BitCache
	permissionByCode map[string]uint32
	version          string
	loadedAt         time.Time
}

type roleCacheEntry struct {
	roleIDs   []uint64
	expiresAt time.Time
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

	// nhớ là khi trừ thì trừ trực tiếp trong bitcache luôn sau đó debounce khoảng 20s thì mới submit vào db.
	// toàn bộ lưu db là do user quản lý không liên quan gì trong đây cả
	// gói đó được phép gọi những hàm nào
	packPermissions map[uint32]_algorithm.BitCache
	// hàm đó được gọi bao nhiêu lần
	permissionQuota map[uint16]uint8
	// số lần gọi hàm còn lại của user theo permission
	userRest map[uint64]map[uint16]uint8
	// quotaCache map[uint16]

	userClient      *clients.UserGrpcClient
	refreshInterval time.Duration
	maxStaleness    time.Duration
	roleCacheTTL    time.Duration
	requestTimeout  time.Duration
	stopCh          chan struct{}
}

/**
 * NewPermissionService starts a bounded refresh lifecycle.
 */
func NewPermissionService(envRuntime string) *PermissionService {
	svc := &PermissionService{
		userClient:      clients.MakeUserGrpcClient(envRuntime),
		roleCache:       make(map[uint64]roleCacheEntry),
		refreshInterval: durationFromEnv("QHPRO_PERMISSION_REFRESH_INTERVAL", 30*time.Second),
		maxStaleness:    durationFromEnv("QHPRO_PERMISSION_MAX_STALENESS", 2*time.Minute),
		roleCacheTTL:    durationFromEnv("QHPRO_PERMISSION_ROLE_CACHE_TTL", 15*time.Second),
		requestTimeout:  durationFromEnv("QHPRO_PERMISSION_REQUEST_TIMEOUT", 2*time.Second),
		stopCh:          make(chan struct{}),
	}
	go svc.runRefreshLoop()
	return svc
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		log.Printf("[PermissionService] invalid %s=%q, using %s", key, value, fallback)
		return fallback
	}
	return parsed
}

func (s *PermissionService) runRefreshLoop() {
	retryDelay := 500 * time.Millisecond
	for {
		ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
		err := s.refresh(ctx)
		cancel()
		if err == nil {
			break
		}
		s.setLastError(err)
		log.Printf("[PermissionService] initial refresh failed: %v", err)
		select {
		case <-time.After(retryDelay):
		case <-s.stopCh:
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
			ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
			err := s.refresh(ctx)
			cancel()
			if err != nil {
				s.setLastError(err)
				log.Printf("[PermissionService] refresh failed: %v", err)
			}
		case <-s.stopCh:
			return
		}
	}
}

/**
 * Close stops the refresh worker.
 */
func (s *PermissionService) Close() {
	if s == nil {
		return
	}
	select {
	case <-s.stopCh:
		return
	default:
		close(s.stopCh)
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
	log.Printf(
		"[PermissionService] snapshot=%s roles=%d codes=%d",
		next.version,
		len(next.rolePermissions),
		len(next.permissionByCode),
	)
	return nil
}

func buildPermissionSnapshot(
	items []*userpb.RolePermissionItem,
	permissions []*sharepb.Permission,
) (permissionSnapshot, error) {
	rolePermissions := make(map[uint64]*_algorithm.BitCache, len(items))
	lines := make([]string, 0, len(items)+len(permissions))
	for _, item := range items {
		if item == nil || item.Role == nil || item.Role.Id == 0 {
			continue
		}
		ids := append([]uint64(nil), item.PermissionIds...)
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		bitmap := _algorithm.NewBitCache(_enum.PermissionIDMax)
		for _, permissionID := range ids {
			if permissionID == 0 || permissionID > uint64(_enum.PermissionIDMax) {
				return permissionSnapshot{}, fmt.Errorf("permission %d exceeds bitmap contract", permissionID)
			}
			bitmap.Set(int(permissionID), true)
		}
		rolePermissions[item.Role.Id] = bitmap
		lines = append(lines, fmt.Sprintf("role:%d:%v", item.Role.Id, ids))
	}

	permissionByCode := make(map[string]uint32, len(permissions))
	for _, permission := range permissions {
		if permission == nil || permission.Id == 0 {
			continue
		}
		code := strings.TrimSpace(permission.Key)
		if code == "" {
			return permissionSnapshot{}, fmt.Errorf("permission %d has empty code", permission.Id)
		}
		if permission.Id > uint64(_enum.PermissionIDMax) {
			return permissionSnapshot{}, fmt.Errorf("permission %d exceeds bitmap contract", permission.Id)
		}
		if previous, exists := permissionByCode[code]; exists && previous != uint32(permission.Id) {
			return permissionSnapshot{}, fmt.Errorf("permission code %s maps to multiple IDs", code)
		}
		permissionByCode[code] = uint32(permission.Id)
		lines = append(lines, fmt.Sprintf("permission:%s:%d", code, permission.Id))
	}

	sort.Strings(lines)
	hash := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return permissionSnapshot{
		rolePermissions:  rolePermissions,
		permissionByCode: permissionByCode,
		version:          hex.EncodeToString(hash[:]),
		loadedAt:         time.Now().UTC(),
	}, nil
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
		Version:    snapshot.version,
		LoadedAt:   snapshot.loadedAt,
		RoleCount:  len(snapshot.rolePermissions),
		CodeCount:  len(snapshot.permissionByCode),
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
		return nil, _errors.ReturnError(503, "permission snapshot chưa sẵn sàng hoặc đã quá hạn")
	}

	set := make(map[uint32]struct{}, len(permissionIDs)+len(permissionCodes))
	for _, permissionID := range permissionIDs {
		if permissionID == 0 || permissionID > uint32(_enum.PermissionIDMax) {
			return nil, _errors.ReturnError(400, "permissionId không hợp lệ")
		}
		set[permissionID] = struct{}{}
	}
	for _, rawCode := range permissionCodes {
		code := strings.TrimSpace(rawCode)
		if code == "" {
			return nil, _errors.ReturnError(400, "permission code không hợp lệ")
		}
		permissionID, ok := snapshot.permissionByCode[code]
		if !ok {
			return nil, _errors.ReturnError(403, "permission code không tồn tại trong snapshot")
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
	if actor, ok := _identitycontext.FromContext(ctx); ok && actor.ProfileID != 0 {
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
		return nil, _errors.ReturnError(503, "không thể xác minh role assignment hiện tại")
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
		return _errors.ReturnError(400, "permission match mode không hợp lệ")
	}

	profileID := profileIDFromContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(401, "Unauthorized")
	}
	roleIDs, err := s.resolveRoleIDs(ctx, profileID)
	if err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return _errors.ReturnError(403, "Bạn không có quyền truy cập")
	}

	allowed := false
	if mode == MatchAll {
		allowed = s.hasAllPermissions(roleIDs, requested)
	} else {
		allowed = s.hasAnyPermission(roleIDs, requested)
	}
	if !allowed {
		return _errors.ReturnError(403, "Bạn không có quyền truy cập")
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
	bitmap := s.snapshot.rolePermissions[roleID]
	s.mu.RUnlock()
	return bitmap != nil && bitmap.Get(int(permissionID))
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

func (s *PermissionService) QuotaCheck(permissionID uint16) uint8 {
	// s.mu.RLock()
	// quota := s.permissionQuota[permissionID]
	// s.mu.RUnlock()
	return 0
}

func (s *PermissionService) QuotaCommit(permissionID uint16) uint8 {
	// s.mu.RLock()
	// quota := s.permissionQuota[permissionID]
	// s.mu.RUnlock()
	return 0
}
