// Package permissioncatalog builds the immutable authorization snapshot used by
// Auth. Durable permission IDs/codes and role assignments come from User IAM;
// this package must not encode the current database sequence as a source-code
// ceiling.
package permissioncatalog

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

const MaxReasonablePermissionID uint64 = 1_000_000

var (
	ErrInvalidPermission       = errors.New("invalid permission catalog entry")
	ErrDuplicatePermissionCode = errors.New("permission code maps to multiple IDs")
	ErrDuplicatePermissionID   = errors.New("permission ID maps to multiple codes")
	ErrUnknownRolePermission   = errors.New("role references permission outside durable catalog")
)

type Permission struct {
	ID   uint64
	Code string
}

type RolePermissions struct {
	RoleID        uint64
	PermissionIDs []uint64
}

// Snapshot is immutable after Build returns. Maps are intentionally private so
// authorization callers can query but cannot mutate the cached authority.
type Snapshot struct {
	rolePermissions  map[uint64]map[uint32]struct{}
	permissionByCode map[string]uint32
	permissionByID   map[uint32]string
	version          string
}

func Build(roles []RolePermissions, permissions []Permission) (Snapshot, error) {
	byCode := make(map[string]uint32, len(permissions))
	byID := make(map[uint32]string, len(permissions))
	lines := make([]string, 0, len(roles)+len(permissions))

	for _, permission := range permissions {
		code := strings.TrimSpace(permission.Code)
		if permission.ID == 0 || permission.ID > MaxReasonablePermissionID || code == "" {
			return Snapshot{}, fmt.Errorf("%w: id=%d code=%q", ErrInvalidPermission, permission.ID, permission.Code)
		}
		id := uint32(permission.ID)
		if previous, found := byCode[code]; found && previous != id {
			return Snapshot{}, fmt.Errorf("%w: %s", ErrDuplicatePermissionCode, code)
		}
		if previous, found := byID[id]; found && previous != code {
			return Snapshot{}, fmt.Errorf("%w: %d", ErrDuplicatePermissionID, id)
		}
		byCode[code] = id
		byID[id] = code
		lines = append(lines, fmt.Sprintf("permission:%s:%d", code, id))
	}

	roleSets := make(map[uint64]map[uint32]struct{}, len(roles))
	for _, role := range roles {
		if role.RoleID == 0 {
			continue
		}
		ids := append([]uint64(nil), role.PermissionIDs...)
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		set := make(map[uint32]struct{}, len(ids))
		for _, rawID := range ids {
			if rawID == 0 || rawID > MaxReasonablePermissionID {
				return Snapshot{}, fmt.Errorf("%w: role=%d permission=%d", ErrUnknownRolePermission, role.RoleID, rawID)
			}
			id := uint32(rawID)
			if _, known := byID[id]; !known {
				return Snapshot{}, fmt.Errorf("%w: role=%d permission=%d", ErrUnknownRolePermission, role.RoleID, id)
			}
			set[id] = struct{}{}
		}
		roleSets[role.RoleID] = set
		lines = append(lines, fmt.Sprintf("role:%d:%v", role.RoleID, ids))
	}

	sort.Strings(lines)
	hash := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return Snapshot{
		rolePermissions:  roleSets,
		permissionByCode: byCode,
		permissionByID:   byID,
		version:          hex.EncodeToString(hash[:]),
	}, nil
}

func (s Snapshot) Version() string { return s.version }
func (s Snapshot) RoleCount() int  { return len(s.rolePermissions) }
func (s Snapshot) CodeCount() int  { return len(s.permissionByCode) }

func (s Snapshot) ResolveCode(code string) (uint32, bool) {
	id, ok := s.permissionByCode[strings.TrimSpace(code)]
	return id, ok
}

func (s Snapshot) ContainsID(id uint32) bool {
	_, ok := s.permissionByID[id]
	return ok
}

func (s Snapshot) RoleHas(roleID uint64, permissionID uint32) bool {
	set := s.rolePermissions[roleID]
	if set == nil {
		return false
	}
	_, ok := set[permissionID]
	return ok
}
