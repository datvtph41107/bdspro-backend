package iampermission

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// Checker evaluates durable User IAM role assignments by permission code.
// It intentionally bypasses role-name shortcuts and hardcoded numeric permission IDs.
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
      JOIN role_permissions rperm
        ON rperm.role_id = r.id
      JOIN permissions p
        ON p.id = rperm.permission_id
       AND p.deleted_at IS NULL
     WHERE rp.profile_id = ?
       AND rp.deleted_at IS NULL
       AND p.key = ?
)
`, actorID, permissionCode).Scan(&allowed)
	if result.Error != nil {
		return false, result.Error
	}
	return allowed, nil
}
