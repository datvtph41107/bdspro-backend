package bootstrapadmin

import (
	"context"
	"errors"
	"fmt"

	"user/internal/usecase/bootstrapadmin"

	"gorm.io/gorm"
)

const (
	rootRoleKey  = "QHPRO_SYSTEM_ROOT"
	rootRoleName = "Quản trị hệ thống gốc"
	rootRoleType = "ROLE_SUPER_ADMIN"
	rootRoleCode = 30
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateRoot(
	ctx context.Context,
	input bootstrapadmin.Input,
	passwordHash string,
) (bootstrapadmin.Result, error) {
	if s == nil || s.db == nil {
		return bootstrapadmin.Result{}, errors.New("bootstrap database is not configured")
	}

	var created bootstrapadmin.Result
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Serialize the one-time transition even when two operators invoke the
		// command concurrently. The lock lives only for this transaction.
		if err := tx.Exec(`SELECT pg_advisory_xact_lock(hashtext('qhpro:user:root-operator'))`).Error; err != nil {
			return fmt.Errorf("lock root operator bootstrap: %w", err)
		}

		var rootExists bool
		if err := tx.Raw(`
SELECT EXISTS (
    SELECT 1
      FROM role_profiles assignment
      JOIN roles role ON role.id = assignment.role_id
     WHERE role.organization_id = 0
       AND role.key = ?
       AND role.deleted_at IS NULL
       AND assignment.deleted_at IS NULL
)`, rootRoleKey).Scan(&rootExists).Error; err != nil {
			return fmt.Errorf("check existing root operator: %w", err)
		}
		if rootExists {
			return bootstrapadmin.ErrAlreadyBootstrapped
		}

		var usernameExists bool
		if err := tx.Raw(`
SELECT EXISTS (
    SELECT 1 FROM auth_method
     WHERE provider = 'ADMIN' AND lower(auth_name) = lower(?) AND deleted_at IS NULL
)`, input.Username).Scan(&usernameExists).Error; err != nil {
			return fmt.Errorf("check administrator username: %w", err)
		}
		if usernameExists {
			return bootstrapadmin.ErrUsernameInUse
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
`, rootRoleName, "Vai trò one-time bootstrap; không gán qua Admin API", rootRoleKey, rootRoleCode).
			Scan(&created.RoleID).Error; err != nil {
			return fmt.Errorf("create root role: %w", err)
		}

		permissionInsert := tx.Exec(`
INSERT INTO role_permissions(role_id, permission_id)
SELECT ?, permission.id
  FROM permissions permission
 WHERE permission.deleted_at IS NULL
   AND permission.key IS NOT NULL
ON CONFLICT (role_id, permission_id) DO NOTHING
`, created.RoleID)
		if permissionInsert.Error != nil {
			return fmt.Errorf("grant root permissions: %w", permissionInsert.Error)
		}
		created.PermissionCount = permissionInsert.RowsAffected

		if err := tx.Raw(`
INSERT INTO user_profile (
    phone, email, full_name, role_type, role_key, role_id, status,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, 10, NOW(), NOW())
RETURNING profile_id
`, input.Phone, input.Email, input.FullName, rootRoleCode, rootRoleCode, created.RoleID).
			Scan(&created.ProfileID).Error; err != nil {
			return fmt.Errorf("create root profile: %w", err)
		}

		if err := tx.Raw(`
INSERT INTO auth_method (
    provider, auth_name, password, full_name, email, phone,
    role_key, status, user_id, created_at, updated_at
) VALUES ('ADMIN', ?, ?, ?, ?, ?, ?, 1, ?, NOW(), NOW())
RETURNING id
`, input.Username, passwordHash, input.FullName, input.Email, input.Phone, rootRoleCode, created.ProfileID).
			Scan(&created.AuthID).Error; err != nil {
			return fmt.Errorf("create root credential: %w", err)
		}

		if err := tx.Exec(`
INSERT INTO user_info(profile_id, status, role_id, created_at, updated_at)
VALUES (?, 10, ?, NOW(), NOW())
`, created.ProfileID, created.RoleID).Error; err != nil {
			return fmt.Errorf("create root user state: %w", err)
		}

		if err := tx.Exec(`
INSERT INTO role_profiles(profile_id, role_id, created_at, updated_at)
VALUES (?, ?, NOW(), NOW())
`, created.ProfileID, created.RoleID).Error; err != nil {
			return fmt.Errorf("assign root role: %w", err)
		}

		if err := tx.Exec(`
INSERT INTO admin_profiles (
    auth_id, profile_id, full_name, email, phone, role_id,
    role_key, role_type, role_description, status, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 10, NOW(), NOW())
`, created.AuthID, created.ProfileID, input.FullName, input.Email, input.Phone,
			created.RoleID, rootRoleKey, rootRoleType, "Quản trị viên gốc được cấp bởi operator").Error; err != nil {
			return fmt.Errorf("create root administrator profile: %w", err)
		}

		return nil
	})
	if err != nil {
		return bootstrapadmin.Result{}, err
	}
	return created, nil
}
