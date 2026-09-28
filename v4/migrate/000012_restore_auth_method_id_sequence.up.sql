BEGIN;

-- Legacy auth_method lost its identity default after the table rename from
-- user_auth. Serving Create() relies on PostgreSQL to allocate id; without a
-- sequence default, inserts fail with SQLSTATE 23502 (null id).

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
          FROM pg_class
         WHERE relkind = 'S'
           AND relname = 'user_auth_id_seq'
    ) THEN
        CREATE SEQUENCE user_auth_id_seq;
    END IF;
END $$;

SELECT setval(
    'user_auth_id_seq',
    GREATEST(
        (SELECT COALESCE(MAX(id), 0) FROM auth_method),
        (SELECT last_value FROM user_auth_id_seq)
    ),
    true
);

-- Do not require OWNED BY: legacy environments may have different table vs
-- sequence owners after the user_auth -> auth_method rename.
ALTER TABLE auth_method
    ALTER COLUMN id SET DEFAULT nextval('user_auth_id_seq'::regclass);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
          FROM pg_constraint
         WHERE conrelid = 'auth_method'::regclass
           AND contype = 'p'
    ) THEN
        ALTER TABLE auth_method
            ADD CONSTRAINT auth_method_pkey PRIMARY KEY (id);
    END IF;
END $$;

COMMIT;
