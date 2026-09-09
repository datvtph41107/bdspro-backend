\set ON_ERROR_STOP on

-- Canonical development identities. This file is deliberately outside every
-- service migration directory: production schema evolution must never create
-- credentials with known passwords. The root Makefile executes it only when
-- QHPRO_ENVIRONMENT=development.
--
-- Admin: admin / admin123
-- Client: 0900000000 (or client) / client123

BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $development_identities$
DECLARE
    v_admin_profile_id bigint;
    v_admin_role_id bigint;
    v_admin_auth_id bigint;
    v_client_profile_id bigint;
    v_client_phone_auth_id bigint;
    v_client_password_auth_id bigint;
BEGIN
    IF current_database() <> 'user_service' THEN
        RAISE EXCEPTION 'development identities target user_service, got %', current_database();
    END IF;

    -- Explicit identifiers in historical snapshots may leave BIGSERIAL
    -- sequences behind the durable data. Repair only the insert paths owned by
    -- this development bootstrap before requesting new IDs.
    PERFORM setval(
        pg_get_serial_sequence('user_profile', 'profile_id'),
        GREATEST(COALESCE((SELECT MAX(profile_id) FROM user_profile), 0) + 1, 1),
        false
    );
    PERFORM setval(
        pg_get_serial_sequence('roles', 'id'),
        GREATEST(COALESCE((SELECT MAX(id) FROM roles), 0) + 1, 1),
        false
    );
    PERFORM setval(
        pg_get_serial_sequence('role_profiles', 'id'),
        GREATEST(COALESCE((SELECT MAX(id) FROM role_profiles), 0) + 1, 1),
        false
    );
    PERFORM setval(
        pg_get_serial_sequence('auth_method', 'id'),
        GREATEST(COALESCE((SELECT MAX(id) FROM auth_method), 0) + 1, 1),
        false
    );
    PERFORM setval(
        pg_get_serial_sequence('admin_profiles', 'id'),
        GREATEST(COALESCE((SELECT MAX(id) FROM admin_profiles), 0) + 1, 1),
        false
    );

    -- Administrator profile and full development IAM role.
    SELECT user_id
      INTO v_admin_profile_id
      FROM auth_method
     WHERE provider = 'ADMIN'
       AND lower(auth_name) = 'admin'
       AND deleted_at IS NULL
     ORDER BY id
     LIMIT 1;

    IF v_admin_profile_id IS NULL OR v_admin_profile_id = 0 THEN
        SELECT profile_id
          INTO v_admin_profile_id
          FROM user_profile
         WHERE phone = '0900000001' AND deleted_at IS NULL
         ORDER BY profile_id
         LIMIT 1;
    END IF;

    IF v_admin_profile_id IS NULL THEN
        INSERT INTO user_profile (
            phone, email, full_name, role_type, role_key, status,
            created_at, updated_at
        ) VALUES (
            '0900000001', 'admin@development.invalid',
            'BDSPro Development Admin', 30, 30, 10, NOW(), NOW()
        ) RETURNING profile_id INTO v_admin_profile_id;
    ELSE
        UPDATE user_profile
           SET phone = '0900000001',
               email = 'admin@development.invalid',
               full_name = 'BDSPro Development Admin',
               role_type = 30,
               role_key = 30,
               status = 10,
               deleted_at = NULL,
               updated_at = NOW()
         WHERE profile_id = v_admin_profile_id;
    END IF;

    SELECT id
      INTO v_admin_role_id
      FROM roles
     WHERE organization_id = 0
       AND key = 'QHPRO_DEVELOPMENT_ADMIN'
       AND deleted_at IS NULL
     ORDER BY id
     LIMIT 1;

    IF v_admin_role_id IS NULL THEN
        INSERT INTO roles (
            role_name, role_description, key, is_default, domain_type,
            organization_id, allow_assign, role_key, created_at, updated_at
        ) VALUES (
            'Quản trị development',
            'Quản trị đầy đủ chỉ được khởi tạo bởi development workflow',
            'QHPRO_DEVELOPMENT_ADMIN', false, 'SYSTEM',
            0, false, 30, NOW(), NOW()
        ) RETURNING id INTO v_admin_role_id;
    END IF;

    INSERT INTO role_permissions(role_id, permission_id)
    SELECT v_admin_role_id, permission.id
      FROM permissions permission
     WHERE permission.deleted_at IS NULL
       AND permission.key IS NOT NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    INSERT INTO role_profiles(profile_id, role_id, created_at, updated_at)
    SELECT v_admin_profile_id, v_admin_role_id, NOW(), NOW()
     WHERE NOT EXISTS (
         SELECT 1
           FROM role_profiles
          WHERE profile_id = v_admin_profile_id
            AND role_id = v_admin_role_id
            AND deleted_at IS NULL
     );

    SELECT id
      INTO v_admin_auth_id
      FROM auth_method
     WHERE provider = 'ADMIN'
       AND lower(auth_name) = 'admin'
       AND deleted_at IS NULL
     ORDER BY id
     LIMIT 1;

    IF v_admin_auth_id IS NULL THEN
        INSERT INTO auth_method (
            provider, auth_name, password, full_name, email, phone,
            role_key, status, user_id, created_at, updated_at
        ) VALUES (
            'ADMIN', 'admin', crypt('admin123', gen_salt('bf', 10)),
            'BDSPro Development Admin', 'admin@development.invalid',
            '0900000001', 30, 1, v_admin_profile_id, NOW(), NOW()
        ) RETURNING id INTO v_admin_auth_id;
    ELSE
        UPDATE auth_method
           SET password = crypt('admin123', gen_salt('bf', 10)),
               full_name = 'BDSPro Development Admin',
               email = 'admin@development.invalid',
               phone = '0900000001',
               role_key = 30,
               status = 1,
               locked_at = NULL,
               locked_until = NULL,
               lock_reason = NULL,
               user_id = v_admin_profile_id,
               updated_at = NOW()
         WHERE id = v_admin_auth_id;
    END IF;

    UPDATE user_profile
       SET role_id = v_admin_role_id, updated_at = NOW()
     WHERE profile_id = v_admin_profile_id;

    INSERT INTO user_info(profile_id, status, role_id, created_at, updated_at)
    VALUES (v_admin_profile_id, 10, v_admin_role_id, NOW(), NOW())
    ON CONFLICT (profile_id) DO UPDATE
       SET status = EXCLUDED.status,
           role_id = EXCLUDED.role_id,
           updated_at = NOW(),
           deleted_at = NULL;

    INSERT INTO user_status(auth_id, active, verified, updated_at)
    VALUES (v_admin_auth_id, true, true, NOW())
    ON CONFLICT (auth_id) DO UPDATE
       SET active = true,
           verified = true,
           locked_until = NULL,
           updated_at = NOW();

    IF EXISTS (
        SELECT 1 FROM admin_profiles
         WHERE profile_id = v_admin_profile_id AND deleted_at IS NULL
    ) THEN
        UPDATE admin_profiles
           SET auth_id = v_admin_auth_id,
               full_name = 'BDSPro Development Admin',
               email = 'admin@development.invalid',
               phone = '0900000001',
               role_id = v_admin_role_id,
               role_key = 'QHPRO_DEVELOPMENT_ADMIN',
               role_type = 'ROLE_SUPER_ADMIN',
               role_description = 'Development bootstrap administrator',
               status = 10,
               locked_at = NULL,
               updated_at = NOW()
         WHERE profile_id = v_admin_profile_id AND deleted_at IS NULL;
    ELSE
        INSERT INTO admin_profiles (
            auth_id, profile_id, full_name, email, phone, role_id,
            role_key, role_type, role_description, status,
            created_at, updated_at
        ) VALUES (
            v_admin_auth_id, v_admin_profile_id,
            'BDSPro Development Admin', 'admin@development.invalid',
            '0900000001', v_admin_role_id,
            'QHPRO_DEVELOPMENT_ADMIN', 'ROLE_SUPER_ADMIN',
            'Development bootstrap administrator', 10, NOW(), NOW()
        );
    END IF;

    -- End-user profile. The PHONE identity mirrors OTP registration; the
    -- PASSWORD identity lets the same phone enter through the real password
    -- endpoint without requiring an external SMS provider.
    SELECT profile_id
      INTO v_client_profile_id
      FROM user_profile
     WHERE phone = '0900000000' AND deleted_at IS NULL
     ORDER BY profile_id
     LIMIT 1;

    IF v_client_profile_id IS NULL THEN
        INSERT INTO user_profile (
            phone, email, full_name, role_type, role_key, status,
            created_at, updated_at
        ) VALUES (
            '0900000000', 'client@development.invalid',
            'Khách hàng QHPRO Development', 10, 10, 10, NOW(), NOW()
        ) RETURNING profile_id INTO v_client_profile_id;
    ELSE
        UPDATE user_profile
           SET email = 'client@development.invalid',
               full_name = 'Khách hàng QHPRO Development',
               role_type = 10,
               role_key = 10,
               status = 10,
               deleted_at = NULL,
               updated_at = NOW()
         WHERE profile_id = v_client_profile_id;
    END IF;

    SELECT id
      INTO v_client_phone_auth_id
      FROM auth_method
     WHERE provider = 'PHONE'
       AND auth_name = '0900000000'
       AND deleted_at IS NULL
     ORDER BY id
     LIMIT 1;

    IF v_client_phone_auth_id IS NULL THEN
        INSERT INTO auth_method (
            provider, auth_name, full_name, email, phone,
            role_key, status, user_id, created_at, updated_at
        ) VALUES (
            'PHONE', '0900000000', 'Khách hàng QHPRO Development',
            'client@development.invalid', '0900000000',
            10, 1, v_client_profile_id, NOW(), NOW()
        ) RETURNING id INTO v_client_phone_auth_id;
    ELSE
        UPDATE auth_method
           SET full_name = 'Khách hàng QHPRO Development',
               email = 'client@development.invalid',
               phone = '0900000000',
               role_key = 10,
               status = 1,
               locked_at = NULL,
               locked_until = NULL,
               lock_reason = NULL,
               user_id = v_client_profile_id,
               updated_at = NOW()
         WHERE id = v_client_phone_auth_id;
    END IF;

    SELECT id
      INTO v_client_password_auth_id
      FROM auth_method
     WHERE provider = 'PASSWORD'
       AND auth_name = 'client'
       AND deleted_at IS NULL
     ORDER BY id
     LIMIT 1;

    IF v_client_password_auth_id IS NULL THEN
        INSERT INTO auth_method (
            provider, auth_name, password, full_name, email, phone,
            role_key, status, user_id, created_at, updated_at
        ) VALUES (
            'PASSWORD', 'client', crypt('client123', gen_salt('bf', 10)),
            'Khách hàng QHPRO Development', 'client@development.invalid',
            '0900000000', 10, 1, v_client_profile_id, NOW(), NOW()
        ) RETURNING id INTO v_client_password_auth_id;
    ELSE
        UPDATE auth_method
           SET password = crypt('client123', gen_salt('bf', 10)),
               full_name = 'Khách hàng QHPRO Development',
               email = 'client@development.invalid',
               phone = '0900000000',
               role_key = 10,
               status = 1,
               locked_at = NULL,
               locked_until = NULL,
               lock_reason = NULL,
               user_id = v_client_profile_id,
               updated_at = NOW()
         WHERE id = v_client_password_auth_id;
    END IF;

    INSERT INTO user_info(profile_id, status, role_id, created_at, updated_at)
    VALUES (v_client_profile_id, 10, NULL, NOW(), NOW())
    ON CONFLICT (profile_id) DO UPDATE
       SET status = EXCLUDED.status,
           updated_at = NOW(),
           deleted_at = NULL;

    INSERT INTO user_status(auth_id, active, verified, updated_at)
    VALUES
        (v_client_phone_auth_id, true, true, NOW()),
        (v_client_password_auth_id, true, true, NOW())
    ON CONFLICT (auth_id) DO UPDATE
       SET active = true,
           verified = true,
           locked_until = NULL,
           updated_at = NOW();
END
$development_identities$;

COMMIT;

SELECT
    auth.provider,
    auth.auth_name,
    auth.phone,
    auth.user_id AS profile_id,
    CASE WHEN auth.password IS NULL OR auth.password = '' THEN 'identity' ELSE 'password' END AS credential_kind
FROM auth_method auth
WHERE auth.deleted_at IS NULL
  AND (
      (auth.provider = 'ADMIN' AND auth.auth_name = 'admin')
      OR (auth.provider = 'PHONE' AND auth.auth_name = '0900000000')
      OR (auth.provider = 'PASSWORD' AND auth.auth_name = 'client')
  )
ORDER BY auth.provider;
