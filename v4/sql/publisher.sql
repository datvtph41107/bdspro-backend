create table profile_transfer (
    profile_id SERIAL PRIMARY KEY,
    full_name TEXT,
    phone TEXT
)
CREATE PUBLICATION profile_pub_transfer FOR TABLE public.profile_transfer;

SELECT trigger_name, event_manipulation, event_object_table
FROM information_schema.triggers
WHERE event_object_table = 'profile_transfer';

SELECT * FROM pg_proc WHERE proname = 'sync_profile_func';
SELECT * FROM pg_class WHERE relname = 'profile_transfer';

ALTER TABLE user_profile ENABLE TRIGGER profile_sync_trigger;

-- drop publication profile_pub_transfer 

-- show wal_level

-- SELECT * FROM pg_publication_tables;
-- SELECT * FROM pg_publication;
-- SELECT * FROM pg_replication_slots;
-- SELECT pg_drop_replication_slot('profile_pub_transfer');

-- SELECT pg_reload_conf();  -- Reload cấu hình mà không cần khởi động lại
