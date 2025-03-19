create table profile_transfer (
    profile_id SERIAL PRIMARY KEY,
    full_name TEXT,
    phone TEXT
)

CREATE SUBSCRIPTION profile_sub_transfer_crm
CONNECTION 'host=14.225.205.222 port=5432 dbname=db_user user=user password=znodVcrXdFtaOLdPaTQr3lmV' 
publication profile_pub_transfer;

-- SELECT * FROM pg_publication_tables;
-- SELECT * FROM pg_subscription;
-- SELECT * FROM pg_publication;

-- ALTER SUBSCRIPTION profile_sub_transfer SET (slot_name = NONE);
-- DROP SUBSCRIPTION profile_sub_transfer;
-- ALTER SUBSCRIPTION profile_sub_transfer DISABLE;

-- SELECT * FROM pg_stat_subscription;

-- ALTER TABLE profile_transfer  REPLICA IDENTITY FULL;
-- SELECT * FROM pg_stat_activity WHERE state <> 'idle';
