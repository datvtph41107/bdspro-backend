package iampermission

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// Checker evaluates durable User IAM role assignments by permission code.
// It intentionally bypasses role-name shortcuts and hardcoded numeric permission IDs.
//
// @bind: user/internal/usecase/useradmin.PermissionAuthorizer
type Checker struct {
	db *gorm.DB
}

func NewChecker(db *gorm.DB) *Checker {
	return &Checker{db: db}
}

func (c *Checker) HasPermission(
	ctx context.Context,
	actorID uint64,
	permissionCode string,
) (bool, error) {
	permissionCode = strings.TrimSpace(permissionCode)
	if c == nil || c.db == nil {
		return false, errors.New("IAM permission database is not configured")
	}
	if actorID == 0 || permissionCode == "" {
		return false, nil
	}

	var allowed bool
	result := c.db.WithContext(ctx).Raw(`
SELECT EXISTS (
    SELECT 1
      FROM role_profiles rp
      JOIN roles r
        ON r.id = rp.role_id
       AND r.deleted_at IS NULL
     WHERE rp.profile_id = ?
       AND rp.deleted_at IS NULL
       AND (
           (
               r.organization_id = 0
               AND r.key = 'QHPRO_SYSTEM_ROOT'
               AND EXISTS (
                   SELECT 1 FROM permissions root_permission
                    WHERE root_permission.key = ?
                      AND root_permission.deleted_at IS NULL
               )
           )
           OR EXISTS (
               SELECT 1
                 FROM role_permissions rperm
                 JOIN permissions permission
                   ON permission.id = rperm.permission_id
                  AND permission.deleted_at IS NULL
                WHERE rperm.role_id = r.id
                  AND permission.key = ?
           )
       )
)
`, actorID, permissionCode, permissionCode).Scan(&allowed)
	if result.Error != nil {
		return false, result.Error
	}
	return allowed, nil
}
