package seeders

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type clientSeedSpec struct {
	Username     string
	PasswordHash string
	FullName     string
	Email        string
	Phone        string
	RoleCode     int
}

type clientAuthRow struct {
	ID     uint64
	UserID uint64
}

func seedClient(ctx context.Context, db *gorm.DB, spec clientSeedSpec) (Result, error) {
	var result Result
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sequence := range [][2]string{
			{"user_profile", "profile_id"},
			{"auth_method", "id"},
		} {
			if err := repairSequence(tx, sequence[0], sequence[1]); err != nil {
				return err
			}
		}

		var existingAuth clientAuthRow
		if err := tx.Raw(`
SELECT id, COALESCE(user_id, 0) AS user_id
  FROM auth_method
 WHERE provider IN ('PASSWORD', 'ADMIN')
   AND auth_name = ?
   AND deleted_at IS NULL
 ORDER BY CASE provider WHEN 'PASSWORD' THEN 0 ELSE 1 END, id
 LIMIT 1
`, spec.Username).Scan(&existingAuth).Error; err != nil {
			return fmt.Errorf("find acceptance client credential: %w", err)
		}

		if existingAuth.UserID != 0 {
			result.ProfileID = existingAuth.UserID
		} else {
			if err := tx.Raw(`
SELECT profile_id
  FROM user_profile
 WHERE phone = ? AND deleted_at IS NULL
 ORDER BY profile_id
 LIMIT 1
`, spec.Phone).Scan(&result.ProfileID).Error; err != nil {
				return fmt.Errorf("find acceptance client profile: %w", err)
			}
		}

		if result.ProfileID == 0 {
			if err := tx.Raw(`
INSERT INTO user_profile (
    phone, email, full_name, role_type, role_key, status,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, 10, NOW(), NOW())
RETURNING profile_id
`, spec.Phone, spec.Email, spec.FullName, spec.RoleCode, spec.RoleCode).
				Scan(&result.ProfileID).Error; err != nil {
				return fmt.Errorf("create acceptance client profile: %w", err)
			}
		} else {
			if err := tx.Exec(`
UPDATE user_profile
   SET phone = ?, email = ?, full_name = ?, role_type = ?, role_key = ?, status = 10,
       updated_at = NOW(), deleted_at = NULL
 WHERE profile_id = ?
`, spec.Phone, spec.Email, spec.FullName, spec.RoleCode, spec.RoleCode, result.ProfileID).Error; err != nil {
				return fmt.Errorf("update acceptance client profile: %w", err)
			}
		}

		if existingAuth.ID == 0 {
			if err := tx.Raw(`
INSERT INTO auth_method (
    provider, auth_name, password, full_name, email, phone,
    role_key, status, user_id, created_at, updated_at
) VALUES ('PASSWORD', ?, ?, ?, ?, ?, ?, 1, ?, NOW(), NOW())
RETURNING id
`, spec.Username, spec.PasswordHash, spec.FullName, spec.Email, spec.Phone, spec.RoleCode, result.ProfileID).
				Scan(&result.AuthID).Error; err != nil {
				return fmt.Errorf("create acceptance client credential: %w", err)
			}
		} else {
			result.AuthID = existingAuth.ID
			if err := tx.Exec(`
UPDATE auth_method
   SET provider = 'PASSWORD', password = ?, full_name = ?, email = ?, phone = ?,
       role_key = ?, status = 1, user_id = ?, updated_at = NOW(), deleted_at = NULL
 WHERE id = ?
`, spec.PasswordHash, spec.FullName, spec.Email, spec.Phone, spec.RoleCode,
				result.ProfileID, result.AuthID).Error; err != nil {
				return fmt.Errorf("update acceptance client credential: %w", err)
			}
		}

		if err := tx.Exec(`
INSERT INTO user_info(profile_id, status, role_id, created_at, updated_at)
VALUES (?, 10, NULL, NOW(), NOW())
ON CONFLICT (profile_id) DO UPDATE
   SET status = 10, role_id = NULL, updated_at = NOW(), deleted_at = NULL
`, result.ProfileID).Error; err != nil {
			return fmt.Errorf("upsert acceptance client user state: %w", err)
		}

		if err := tx.Exec(`
INSERT INTO user_status(auth_id, active, verified, updated_at)
VALUES (?, true, true, NOW())
ON CONFLICT (auth_id) DO UPDATE
   SET active = true, verified = true, locked_until = NULL, updated_at = NOW()
`, result.AuthID).Error; err != nil {
			return fmt.Errorf("upsert acceptance client status: %w", err)
		}

		return nil
	})
	if err != nil {
		return Result{}, err
	}
	return result, nil
}
