\set ON_ERROR_STOP on

BEGIN;

-- Acceptance-only end-user credential for exercising the client login and
-- checkout flow. Production migrations never seed identities or passwords.
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TEMP TABLE qhpro_acceptance_user_input (
    username text NOT NULL,
    phone text NOT NULL,
    email text NOT NULL,
    password text NOT NULL
) ON COMMIT DROP;
INSERT INTO qhpro_acceptance_user_input(username, phone, email, password)
VALUES (:'username', :'phone', :'email', :'user_password');

DO $fixture$
DECLARE
    v_profile_id bigint;
    v_auth_id bigint;
    v_username text;
    v_phone text;
    v_email text;
    v_password text;
BEGIN
    SELECT username, phone, email, password
      INTO v_username, v_phone, v_email, v_password
      FROM qhpro_acceptance_user_input;
    IF length(v_password) < 12 THEN
        RAISE EXCEPTION 'QHPRO_ACCEPTANCE_USER_PASSWORD must contain at least 12 characters';
    END IF;
    IF v_username !~ '^qhpro_acceptance_[0-9]{10,20}$'
       OR v_phone !~ '^039[0-9]{7}$'
       OR v_email !~ '^qhpro-[0-9]{10,20}@e[.]invalid$' THEN
        RAISE EXCEPTION 'invalid acceptance identity';
    END IF;

    PERFORM setval(
        pg_get_serial_sequence('user_profile', 'profile_id'),
        GREATEST(COALESCE((SELECT MAX(profile_id) FROM user_profile), 0) + 1, 1),
        false
    );
    PERFORM setval(
        pg_get_serial_sequence('auth_method', 'id'),
        GREATEST(COALESCE((SELECT MAX(id) FROM auth_method), 0) + 1, 1),
        false
    );

    SELECT profile_id
     INTO v_profile_id
      FROM user_profile
     WHERE phone = v_phone AND deleted_at IS NULL
     ORDER BY profile_id
     LIMIT 1;

    IF v_profile_id IS NULL THEN
        INSERT INTO user_profile (
            phone, email, full_name, role_type, role_key, status,
            created_at, updated_at
        ) VALUES (
            v_phone, v_email,
            'QHPRO Acceptance Commercial User', 10, 10, 10, NOW(), NOW()
        ) RETURNING profile_id INTO v_profile_id;
    ELSE
        UPDATE user_profile
           SET email = v_email,
               full_name = 'QHPRO Acceptance Commercial User',
               role_type = 10,
               role_key = 10,
               status = 10,
               updated_at = NOW()
         WHERE profile_id = v_profile_id;
    END IF;

    SELECT id
      INTO v_auth_id
     FROM auth_method
     WHERE provider IN ('PASSWORD', 'ADMIN')
       AND auth_name = v_username
       AND deleted_at IS NULL
     ORDER BY CASE provider WHEN 'PASSWORD' THEN 0 ELSE 1 END, id
     LIMIT 1;

    IF v_auth_id IS NULL THEN
        INSERT INTO auth_method (
            provider, auth_name, password, full_name, email, phone,
            role_key, status, user_id, created_at, updated_at
        ) VALUES (
            'PASSWORD', v_username, crypt(v_password, gen_salt('bf', 10)),
            'QHPRO Acceptance Commercial User',
            v_email, v_phone,
            10, 1, v_profile_id, NOW(), NOW()
        ) RETURNING id INTO v_auth_id;
    ELSE
        UPDATE auth_method
           SET provider = 'PASSWORD',
               password = crypt(v_password, gen_salt('bf', 10)),
               full_name = 'QHPRO Acceptance Commercial User',
               email = v_email,
               phone = v_phone,
               role_key = 10,
               status = 1,
               user_id = v_profile_id,
               updated_at = NOW()
         WHERE id = v_auth_id;
    END IF;

    INSERT INTO user_info(profile_id, status, role_id, created_at, updated_at)
    VALUES (v_profile_id, 10, NULL, NOW(), NOW())
    ON CONFLICT (profile_id) DO UPDATE
       SET status = EXCLUDED.status,
           updated_at = NOW(),
           deleted_at = NULL;

    INSERT INTO user_status(auth_id, active, verified, updated_at)
    VALUES (v_auth_id, true, true, NOW())
    ON CONFLICT (auth_id) DO UPDATE
       SET active = true,
           verified = true,
           locked_until = NULL,
           updated_at = NOW();
END
$fixture$;

COMMIT;

SELECT profile.profile_id, auth.id AS auth_id, profile.role_type, profile.status
  FROM user_profile profile
  JOIN auth_method auth ON auth.user_id = profile.profile_id
                       AND auth.provider = 'PASSWORD'
                       AND auth.auth_name = :'username'
                       AND auth.deleted_at IS NULL
 WHERE profile.phone = :'phone' AND profile.deleted_at IS NULL;
