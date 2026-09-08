package seeders

import (
	"context"
	"fmt"
	"regexp"

	"gorm.io/gorm"
)

var identifierPattern = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

type adminAuthRow struct {
	ID     uint64
	UserID uint64
}

func seedAdmin(ctx context.Context, db *gorm.DB, spec adminSeedSpec) (Result, error) {
	var result Result
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sequence := range [][2]string{
			{"user_profile", "profile_id"},
			{"roles", "id"},
			{"role_profiles", "id"},
			{"auth_method", "id"},
			{"admin_profiles", "id"},
		} {
			if err := repairSequence(tx, sequence[0], sequence[1]); err != nil {
				return err
			}
		}

		if err := tx.Raw(`
INSERT INTO roles (
    role_name, role_description, key, is_default, domain_type,
    organization_id, allow_assign, role_key, created_at, updated_at
) VALUES (?, ?, ?, false, 'SYSTEM', 0, false, ?, NOW(), NOW())
ON CONFLICT (organization_id, key) WHERE deleted_at IS NULL
DO UPDATE SET
    role_name = EXCLUDED.role_name,
    role_description = EXCLUDED.role_description,
    domain_type = 'SYSTEM',
    allow_assign = false,
    role_key = EXCLUDED.role_key,
    updated_at = NOW()
RETURNING id
`, spec.RoleName, spec.RoleDescription, spec.RoleKey, spec.RoleCode).Scan(&result.RoleID).Error; err != nil {
			return fmt.Errorf("upsert seed admin role: %w", err)
		}

		permissionCount, err := replaceRolePermissions(tx, result.RoleID, spec)
		if err != nil {
			return err
		}
		result.PermissionCount = permissionCount

		var existingAuth adminAuthRow
		if err := tx.Raw(`
SELECT id, COALESCE(user_id, 0) AS user_id
  FROM auth_method
 WHERE provider = 'ADMIN'
   AND lower(auth_name) = lower(?)
   AND deleted_at IS NULL
 ORDER BY id
 LIMIT 1
`, spec.Username).Scan(&existingAuth).Error; err != nil {
			return fmt.Errorf("find seed admin credential: %w", err)
		}

		if existingAuth.UserID != 0 {
			result.ProfileID = existingAuth.UserID
			if err := tx.Exec(`
UPDATE user_profile
   SET phone = ?, email = ?, full_name = ?, role_type = ?, role_key = ?, role_id = ?,
       status = 10, updated_at = NOW(), deleted_at = NULL
 WHERE profile_id = ?
`, nullableText(spec.Phone), nullableText(spec.Email), spec.FullName, spec.RoleCode, spec.RoleCode, result.RoleID, result.ProfileID).Error; err != nil {
				return fmt.Errorf("update seed admin profile: %w", err)
			}
		} else {
			if err := tx.Raw(`
INSERT INTO user_profile (
    phone, email, full_name, role_type, role_key, role_id, status,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, 10, NOW(), NOW())
RETURNING profile_id
`, nullableText(spec.Phone), nullableText(spec.Email), spec.FullName, spec.RoleCode, spec.RoleCode, result.RoleID).
				Scan(&result.ProfileID).Error; err != nil {
				return fmt.Errorf("create seed admin profile: %w", err)
			}
		}

		if existingAuth.ID != 0 {
			result.AuthID = existingAuth.ID
			if err := tx.Exec(`
UPDATE auth_method
   SET password = ?, full_name = ?, email = ?, phone = ?, role_key = ?, status = 1,
       user_id = ?, updated_at = NOW(), deleted_at = NULL
 WHERE id = ?
`, spec.PasswordHash, spec.FullName, nullableText(spec.Email), nullableText(spec.Phone), spec.RoleCode,
				result.ProfileID, result.AuthID).Error; err != nil {
				return fmt.Errorf("update seed admin credential: %w", err)
			}
		} else {
			if err := tx.Raw(`
INSERT INTO auth_method (
    provider, auth_name, password, full_name, email, phone,
    role_key, status, user_id, created_at, updated_at
) VALUES ('ADMIN', ?, ?, ?, ?, ?, ?, 1, ?, NOW(), NOW())
RETURNING id
`, spec.Username, spec.PasswordHash, spec.FullName, nullableText(spec.Email), nullableText(spec.Phone),
				spec.RoleCode, result.ProfileID).Scan(&result.AuthID).Error; err != nil {
				return fmt.Errorf("create seed admin credential: %w", err)
			}
		}

		if err := tx.Exec(`
INSERT INTO user_info(profile_id, status, role_id, created_at, updated_at)
VALUES (?, 10, ?, NOW(), NOW())
ON CONFLICT (profile_id) DO UPDATE
   SET status = 10, role_id = EXCLUDED.role_id, updated_at = NOW(), deleted_at = NULL
`, result.ProfileID, result.RoleID).Error; err != nil {
			return fmt.Errorf("upsert seed admin user state: %w", err)
		}

		var assignmentExists bool
		if err := tx.Raw(`
SELECT EXISTS (
    SELECT 1 FROM role_profiles
     WHERE profile_id = ? AND role_id = ? AND deleted_at IS NULL
)`, result.ProfileID, result.RoleID).Scan(&assignmentExists).Error; err != nil {
			return fmt.Errorf("find seed admin role assignment: %w", err)
		}
		if !assignmentExists {
			if err := tx.Exec(`
INSERT INTO role_profiles(profile_id, role_id, created_at, updated_at)
VALUES (?, ?, NOW(), NOW())
`, result.ProfileID, result.RoleID).Error; err != nil {
				return fmt.Errorf("assign seed admin role: %w", err)
			}
		}

		var adminProfileID uint64
		if err := tx.Raw(`
SELECT id FROM admin_profiles
 WHERE profile_id = ? AND deleted_at IS NULL
 ORDER BY id
 LIMIT 1
`, result.ProfileID).Scan(&adminProfileID).Error; err != nil {
			return fmt.Errorf("find seed admin projection: %w", err)
		}
		if adminProfileID != 0 {
			if err := tx.Exec(`
UPDATE admin_profiles
   SET auth_id = ?, full_name = ?, email = ?, phone = ?, role_id = ?, role_key = ?,
       role_type = ?, role_description = ?, status = 10, updated_at = NOW(), deleted_at = NULL
 WHERE id = ?
`, result.AuthID, spec.FullName, nullableText(spec.Email), nullableText(spec.Phone), result.RoleID,
				spec.RoleKey, spec.RoleType, spec.AdminRoleDescription, adminProfileID).Error; err != nil {
				return fmt.Errorf("update seed admin projection: %w", err)
			}
		} else {
			if err := tx.Exec(`
INSERT INTO admin_profiles (
    auth_id, profile_id, full_name, email, phone, role_id,
    role_key, role_type, role_description, status, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 10, NOW(), NOW())
`, result.AuthID, result.ProfileID, spec.FullName, nullableText(spec.Email), nullableText(spec.Phone),
				result.RoleID, spec.RoleKey, spec.RoleType, spec.AdminRoleDescription).Error; err != nil {
				return fmt.Errorf("create seed admin projection: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return Result{}, err
	}
	return result, nil
}

func replaceRolePermissions(tx *gorm.DB, roleID uint64, spec adminSeedSpec) (int64, error) {
	if err := tx.Exec(`DELETE FROM role_permissions WHERE role_id = ?`, roleID).Error; err != nil {
		return 0, fmt.Errorf("clear seed admin permissions: %w", err)
	}

	if spec.GrantAllPermissions {
		if err := tx.Exec(`
INSERT INTO role_permissions(role_id, permission_id)
SELECT ?, permission.id
  FROM permissions permission
 WHERE permission.deleted_at IS NULL
   AND permission.key IS NOT NULL
ON CONFLICT (role_id, permission_id) DO NOTHING
`, roleID).Error; err != nil {
			return 0, fmt.Errorf("grant development admin permissions: %w", err)
		}
	} else {
		var available int64
		if err := tx.Raw(`
SELECT COUNT(*)
  FROM permissions permission
 WHERE permission.deleted_at IS NULL
   AND permission.key IN ?
`, spec.PermissionKeys).Scan(&available).Error; err != nil {
			return 0, fmt.Errorf("count seed admin permissions: %w", err)
		}
		if available != int64(len(spec.PermissionKeys)) {
			return 0, fmt.Errorf("seed admin permission catalog incomplete: have %d want %d", available, len(spec.PermissionKeys))
		}
		if err := tx.Exec(`
INSERT INTO role_permissions(role_id, permission_id)
SELECT ?, permission.id
  FROM permissions permission
 WHERE permission.deleted_at IS NULL
   AND permission.key IN ?
ON CONFLICT (role_id, permission_id) DO NOTHING
`, roleID, spec.PermissionKeys).Error; err != nil {
			return 0, fmt.Errorf("grant seed admin permissions: %w", err)
		}
	}

	var granted int64
	if err := tx.Raw(`SELECT COUNT(*) FROM role_permissions WHERE role_id = ?`, roleID).Scan(&granted).Error; err != nil {
		return 0, fmt.Errorf("count granted seed admin permissions: %w", err)
	}
	return granted, nil
}

func repairSequence(tx *gorm.DB, table, column string) error {
	if !identifierPattern.MatchString(table) || !identifierPattern.MatchString(column) {
		return fmt.Errorf("invalid sequence identifier %q.%q", table, column)
	}
	statement := fmt.Sprintf(`
SELECT setval(
    pg_get_serial_sequence('%s', '%s'),
    GREATEST(COALESCE((SELECT MAX(%s) FROM %s), 0) + 1, 1),
    false
)
`, table, column, column, table)
	if err := tx.Exec(statement).Error; err != nil {
		return fmt.Errorf("repair %s.%s sequence: %w", table, column, err)
	}
	return nil
}

func nullableText(value string) any {
	if value == "" {
		return nil
	}
	return value
}
