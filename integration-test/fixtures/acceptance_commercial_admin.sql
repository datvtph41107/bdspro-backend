\set ON_ERROR_STOP on

BEGIN;

-- Acceptance-only data. Production IAM remains operator-owned and is never
-- seeded by the canonical User migrations.
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TEMP TABLE qhpro_acceptance_admin_input (password text NOT NULL) ON COMMIT DROP;
INSERT INTO qhpro_acceptance_admin_input(password) VALUES (:'admin_password');

DO $fixture$
DECLARE
    v_profile_id bigint;
    v_role_id bigint;
    v_auth_id bigint;
    v_password text;
BEGIN
    SELECT password INTO v_password FROM qhpro_acceptance_admin_input;
    IF length(v_password) < 12 THEN
        RAISE EXCEPTION 'QHPRO_ACCEPTANCE_ADMIN_PASSWORD must contain at least 12 characters';
    END IF;

    -- Historical/development snapshots may contain rows inserted with explicit IDs,
    -- leaving their BIGSERIAL sequences behind the durable table state. Repair
    -- only the fixture-owned insert paths before requesting new identifiers.
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

    SELECT profile_id
      INTO v_profile_id
      FROM user_profile
     WHERE phone = '0399999999' AND deleted_at IS NULL
     ORDER BY profile_id
     LIMIT 1;

    IF v_profile_id IS NULL THEN
        INSERT INTO user_profile (
            phone, email, full_name, role_type, role_key, status,
            created_at, updated_at
        ) VALUES (
            '0399999999', 'qhpro.acceptance.admin@example.invalid',
            'QHPRO Acceptance Commercial Admin', 20, 20, 10, NOW(), NOW()
        ) RETURNING profile_id INTO v_profile_id;
    ELSE
        UPDATE user_profile
           SET email = 'qhpro.acceptance.admin@example.invalid',
               full_name = 'QHPRO Acceptance Commercial Admin',
               role_type = 20,
               role_key = 20,
               status = 10,
               updated_at = NOW()
         WHERE profile_id = v_profile_id;
    END IF;

    SELECT id
      INTO v_role_id
      FROM roles
     WHERE organization_id = 0
       AND key = 'QHPRO_ACCEPTANCE_COMMERCIAL_ADMIN'
       AND deleted_at IS NULL
     ORDER BY id
     LIMIT 1;

    IF v_role_id IS NULL THEN
        INSERT INTO roles (
            role_name, role_description, key, is_default, domain_type,
            organization_id, allow_assign, role_key, created_at, updated_at
        ) VALUES (
            'QHPRO Acceptance Commercial Admin',
            'Acceptance fixture role for commercial operation verification',
            'QHPRO_ACCEPTANCE_COMMERCIAL_ADMIN', false, 'SYSTEM',
            0, false, 20, NOW(), NOW()
        ) RETURNING id INTO v_role_id;
    END IF;

    INSERT INTO role_permissions(role_id, permission_id)
    SELECT v_role_id, permission.id
      FROM permissions permission
     WHERE permission.deleted_at IS NULL
       AND permission.key IN (
           'CATALOG_PLAN_VIEW',
           'CATALOG_PLAN_MANAGE',
           'CATALOG_PLAN_PUBLISH',
           'COMMERCIAL_SUBSCRIPTION_VIEW',
           'COMMERCIAL_USAGE_VIEW',
           'COMMERCIAL_USAGE_RECONCILE',
           'PAYMENT_ORDER_VIEW',
           'PAYMENT_FULFILLMENT_VIEW',
           'PAYMENT_FULFILLMENT_REDRIVE',
           'USER_ADMIN_VIEW',
           'USER_ADMIN_MANAGE',
           'IAM_ROLE_VIEW',
           'IAM_ROLE_MANAGE',
           'ORGANIZATION_VIEW',
           'ORGANIZATION_MANAGE'
       )
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    INSERT INTO role_profiles(profile_id, role_id, created_at, updated_at)
    SELECT v_profile_id, v_role_id, NOW(), NOW()
     WHERE NOT EXISTS (
         SELECT 1
           FROM role_profiles
          WHERE profile_id = v_profile_id
            AND role_id = v_role_id
            AND deleted_at IS NULL
     );

    SELECT id
      INTO v_auth_id
      FROM auth_method
     WHERE provider = 'ADMIN'
       AND auth_name = 'qhpro.acceptance.admin'
       AND deleted_at IS NULL
     ORDER BY id
     LIMIT 1;

    IF v_auth_id IS NULL THEN
        INSERT INTO auth_method (
            provider, auth_name, password, full_name, email, phone,
            role_key, status, user_id, created_at, updated_at
        ) VALUES (
            'ADMIN', 'qhpro.acceptance.admin', crypt(v_password, gen_salt('bf', 10)),
            'QHPRO Acceptance Commercial Admin',
            'qhpro.acceptance.admin@example.invalid', '0399999999',
            20, 1, v_profile_id, NOW(), NOW()
        ) RETURNING id INTO v_auth_id;
    ELSE
        UPDATE auth_method
           SET password = crypt(v_password, gen_salt('bf', 10)),
               full_name = 'QHPRO Acceptance Commercial Admin',
               email = 'qhpro.acceptance.admin@example.invalid',
               phone = '0399999999',
               role_key = 20,
               status = 1,
               user_id = v_profile_id,
               updated_at = NOW()
         WHERE id = v_auth_id;
    END IF;

    UPDATE user_profile
       SET role_id = v_role_id, updated_at = NOW()
     WHERE profile_id = v_profile_id;

    INSERT INTO user_info(profile_id, status, role_id, created_at, updated_at)
    VALUES (v_profile_id, 10, v_role_id, NOW(), NOW())
    ON CONFLICT (profile_id) DO UPDATE
       SET status = EXCLUDED.status,
           role_id = EXCLUDED.role_id,
           updated_at = NOW(),
           deleted_at = NULL;

    IF EXISTS (SELECT 1 FROM admin_profiles WHERE profile_id = v_profile_id AND deleted_at IS NULL) THEN
        UPDATE admin_profiles
           SET auth_id = v_auth_id,
               full_name = 'QHPRO Acceptance Commercial Admin',
               email = 'qhpro.acceptance.admin@example.invalid',
               phone = '0399999999',
               role_id = v_role_id,
               role_key = 'QHPRO_ACCEPTANCE_COMMERCIAL_ADMIN',
               role_type = 'ROLE_ADMIN',
               status = 10,
               updated_at = NOW()
         WHERE profile_id = v_profile_id AND deleted_at IS NULL;
    ELSE
        INSERT INTO admin_profiles (
            auth_id, profile_id, full_name, email, phone, role_id,
            role_key, role_type, role_description, status,
            created_at, updated_at
        ) VALUES (
            v_auth_id, v_profile_id, 'QHPRO Acceptance Commercial Admin',
            'qhpro.acceptance.admin@example.invalid', '0399999999', v_role_id,
            'QHPRO_ACCEPTANCE_COMMERCIAL_ADMIN', 'ROLE_ADMIN',
            'Acceptance fixture for commercial UI verification', 10, NOW(), NOW()
        );
    END IF;
END
$fixture$;

COMMIT;

SELECT profile.profile_id, auth.id AS auth_id, role.id AS role_id,
       count(permission.permission_id) AS permission_count
  FROM user_profile profile
  JOIN auth_method auth ON auth.user_id = profile.profile_id
                       AND auth.provider = 'ADMIN'
                       AND auth.auth_name = 'qhpro.acceptance.admin'
                       AND auth.deleted_at IS NULL
  JOIN role_profiles assignment ON assignment.profile_id = profile.profile_id
                               AND assignment.deleted_at IS NULL
  JOIN roles role ON role.id = assignment.role_id
                 AND role.key = 'QHPRO_ACCEPTANCE_COMMERCIAL_ADMIN'
                 AND role.deleted_at IS NULL
  LEFT JOIN role_permissions permission ON permission.role_id = role.id
 WHERE profile.phone = '0399999999' AND profile.deleted_at IS NULL
 GROUP BY profile.profile_id, auth.id, role.id;
