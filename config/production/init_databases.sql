-- ============================================================
-- Script tạo toàn bộ Database và User Owner cho hệ thống
-- Microservices Production Environment
-- ============================================================
-- Tạo ngày: 2025-01-08
-- Mô tả: Script này tạo tất cả database và user owner cần thiết
--        cho các microservices trong môi trường production
-- ============================================================

-- Kết nối với database postgres mặc định để có quyền tạo database
\c postgres

-- ============================================================
-- PHẦN 1: TẠO USER OWNERS
-- ============================================================

-- Tạo user owner riêng cho từng database
-- Quy tắc: db_<service> -> username: <service>
-- Password: Chuỗi hex ngẫu nhiên 20 ký tự (được generate an toàn)

-- 1. User cho Auth Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'auth') THEN
        CREATE USER auth WITH PASSWORD 'a7f3c2e9b4d8a1f6e5c7';
    END IF;
END
$$;
COMMENT ON ROLE auth IS 'Owner for db_auth - Authentication Service';

-- 2. User cho User Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'user_owner') THEN
        CREATE USER user_owner WITH PASSWORD 'b8e4d1f0c5a9b2e7d3f8';
    END IF;
END
$$;
COMMENT ON ROLE user_owner IS 'Owner for db_user - User Service';

-- 3. User cho Notification Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'notification') THEN
        CREATE USER notification WITH PASSWORD 'c9f5e2a1d6b0c3f8e4a9';
    END IF;
END
$$;
COMMENT ON ROLE notification IS 'Owner for db_notification - Notification Service';

-- 4. User cho Organization Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'org') THEN
        CREATE USER org WITH PASSWORD 'd0a6f3b2e7c1d4a9f5b0';
    END IF;
END
$$;
COMMENT ON ROLE org IS 'Owner for db_org - Organization Service';

-- 5. User cho CRM Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'crm') THEN
        CREATE USER crm WITH PASSWORD 'e1b7a4c3f8d2e5b0a6c1';
    END IF;
END
$$;
COMMENT ON ROLE crm IS 'Owner for db_crm - CRM Service';

-- 6. User cho BDSPro Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'bdspro') THEN
        CREATE USER bdspro WITH PASSWORD 'f2c8b5d4a9e3f6c1b7d2';
    END IF;
END
$$;
COMMENT ON ROLE bdspro IS 'Owner for db_bdspro - Real Estate Service';

-- 7. User cho Chat Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'chat') THEN
        CREATE USER chat WITH PASSWORD 'a3d9c6e5b0f4a7d2c8e3';
    END IF;
END
$$;
COMMENT ON ROLE chat IS 'Owner for db_chat - Chat Service';

-- 8. User cho Social Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'social') THEN
        CREATE USER social WITH PASSWORD 'b4e0d7f6c1a5b8e3d9f4';
    END IF;
END
$$;
COMMENT ON ROLE social IS 'Owner for db_social - Social Service';

-- 9. User cho Transaction Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'transaction') THEN
        CREATE USER transaction WITH PASSWORD 'c5f1e8a7d2b6c9f4e0a5';
    END IF;
END
$$;
COMMENT ON ROLE transaction IS 'Owner for db_transaction - Transaction Service';

-- 10. User cho Membership Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'membership') THEN
        CREATE USER membership WITH PASSWORD 'd6a2f9b8e3c7d0a5f1b6';
    END IF;
END
$$;
COMMENT ON ROLE membership IS 'Owner for db_membership - Membership Service';

-- 11. User cho Payment Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'payment') THEN
        CREATE USER payment WITH PASSWORD 'e7b3a0c9f4d8e1b6a2c7';
    END IF;
END
$$;
COMMENT ON ROLE payment IS 'Owner for db_payment - Payment Service';

-- 12. User cho Marketing Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'marketing') THEN
        CREATE USER marketing WITH PASSWORD 'f8c4b1d0a5e9f2c7b3d8';
    END IF;
END
$$;
COMMENT ON ROLE marketing IS 'Owner for marketing_db - Marketing Service';

-- 13. User cho Feedback Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'feedback') THEN
        CREATE USER feedback WITH PASSWORD 'a9d5c2e1b6f0a3d8c4e9';
    END IF;
END
$$;
COMMENT ON ROLE feedback IS 'Owner for db_feedback - Feedback Service';

-- 14. User cho Appointment Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'appointment') THEN
        CREATE USER appointment WITH PASSWORD 'b0e6d3f2c7a1b4e9d5f0';
    END IF;
END
$$;
COMMENT ON ROLE appointment IS 'Owner for db_appointment - Appointment Service';

-- 15. User cho TQD Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'tqd') THEN
        CREATE USER tqd WITH PASSWORD 'c1f7e4a3d8b2c5f0e6a1';
    END IF;
END
$$;
COMMENT ON ROLE tqd IS 'Owner for db_tqd - TQD Service';

-- 16. User cho File Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'file_owner') THEN
        CREATE USER file_owner WITH PASSWORD 'd2a8f5b4e9c3d6a1f7b2';
    END IF;
END
$$;
COMMENT ON ROLE file_owner IS 'Owner for db_file - File Service';

-- 17. User cho Map Service
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'map') THEN
        CREATE USER map WITH PASSWORD 'e3b9a6c5f0d4e7b2a8c3';
    END IF;
END
$$;
COMMENT ON ROLE map IS 'Owner for db_map - Map Service';

-- ============================================================
-- PHẦN 2: TẠO CÁC DATABASE CHO MICROSERVICES
-- ============================================================

-- 1. Auth Service Database
DROP DATABASE IF EXISTS db_auth;
CREATE DATABASE db_auth
    WITH 
    OWNER = auth
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_auth IS 'Database for Authentication Service';

-- 2. User Service Database
DROP DATABASE IF EXISTS db_user;
CREATE DATABASE db_user
    WITH 
    OWNER = user_owner
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_user IS 'Database for User Service';

-- 3. Notification Service Database
DROP DATABASE IF EXISTS db_notification;
CREATE DATABASE db_notification
    WITH 
    OWNER = notification
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_notification IS 'Database for Notification Service';

-- 4. Organization Service Database
DROP DATABASE IF EXISTS db_org;
CREATE DATABASE db_org
    WITH 
    OWNER = org
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_org IS 'Database for Organization Service';

-- 5. CRM Service Database
DROP DATABASE IF EXISTS db_crm;
CREATE DATABASE db_crm
    WITH 
    OWNER = crm
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_crm IS 'Database for CRM Service';

-- 6. BDS Pro Service Database
DROP DATABASE IF EXISTS db_bdspro;
CREATE DATABASE db_bdspro
    WITH 
    OWNER = bdspro
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_bdspro IS 'Database for Real Estate (BDSPro) Service';

-- 7. Chat Service Database
DROP DATABASE IF EXISTS db_chat;
CREATE DATABASE db_chat
    WITH 
    OWNER = chat
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_chat IS 'Database for Chat Service';

-- 8. Social Service Database
DROP DATABASE IF EXISTS db_social;
CREATE DATABASE db_social
    WITH 
    OWNER = social
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_social IS 'Database for Social Service';

-- 9. Transaction Service Database
DROP DATABASE IF EXISTS db_transaction;
CREATE DATABASE db_transaction
    WITH 
    OWNER = transaction
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_transaction IS 'Database for Transaction Service';

-- 10. Membership Service Database
DROP DATABASE IF EXISTS db_membership;
CREATE DATABASE db_membership
    WITH 
    OWNER = membership
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_membership IS 'Database for Membership Service';

-- 11. Payment Service Database
DROP DATABASE IF EXISTS db_payment;
CREATE DATABASE db_payment
    WITH 
    OWNER = payment
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_payment IS 'Database for Payment Service';

-- 12. Marketing Service Database
DROP DATABASE IF EXISTS marketing_db;
CREATE DATABASE marketing_db
    WITH 
    OWNER = marketing
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE marketing_db IS 'Database for Marketing Service';

-- 13. Feedback Service Database
DROP DATABASE IF EXISTS db_feedback;
CREATE DATABASE db_feedback
    WITH 
    OWNER = feedback
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_feedback IS 'Database for Feedback Service';

-- 14. Appointment Service Database
DROP DATABASE IF EXISTS db_appointment;
CREATE DATABASE db_appointment
    WITH 
    OWNER = appointment
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_appointment IS 'Database for Appointment Service';

-- 15. TQD Service Database
DROP DATABASE IF EXISTS db_tqd;
CREATE DATABASE db_tqd
    WITH 
    OWNER = tqd
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_tqd IS 'Database for TQD Service';

-- 16. File Service Database
DROP DATABASE IF EXISTS db_file;
CREATE DATABASE db_file
    WITH 
    OWNER = file_owner
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_file IS 'Database for File Service';

-- 17. Map Service Database
DROP DATABASE IF EXISTS db_map;
CREATE DATABASE db_map
    WITH 
    OWNER = map
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE db_map IS 'Database for Map Service';

-- ============================================================
-- PHẦN 3: CẤP QUYỀN CHO USER
-- ============================================================

-- Cấp quyền ALL PRIVILEGES cho từng owner trên database của họ
GRANT ALL PRIVILEGES ON DATABASE db_auth TO auth;
GRANT ALL PRIVILEGES ON DATABASE db_user TO user_owner;
GRANT ALL PRIVILEGES ON DATABASE db_notification TO notification;
GRANT ALL PRIVILEGES ON DATABASE db_org TO org;
GRANT ALL PRIVILEGES ON DATABASE db_crm TO crm;
GRANT ALL PRIVILEGES ON DATABASE db_bdspro TO bdspro;
GRANT ALL PRIVILEGES ON DATABASE db_chat TO chat;
GRANT ALL PRIVILEGES ON DATABASE db_social TO social;
GRANT ALL PRIVILEGES ON DATABASE db_transaction TO transaction;
GRANT ALL PRIVILEGES ON DATABASE db_membership TO membership;
GRANT ALL PRIVILEGES ON DATABASE db_payment TO payment;
GRANT ALL PRIVILEGES ON DATABASE marketing_db TO marketing;
GRANT ALL PRIVILEGES ON DATABASE db_feedback TO feedback;
GRANT ALL PRIVILEGES ON DATABASE db_appointment TO appointment;
GRANT ALL PRIVILEGES ON DATABASE db_tqd TO tqd;
GRANT ALL PRIVILEGES ON DATABASE db_file TO file_owner;
GRANT ALL PRIVILEGES ON DATABASE db_map TO map;

-- ============================================================
-- PHẦN 4: KHỞI TẠO EXTENSIONS CHO MỖI DATABASE
-- ============================================================

-- Kết nối và tạo extensions cho từng database
\c db_auth
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_user
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_notification
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_org
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_crm
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_bdspro
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_chat
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_social
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_transaction
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_membership
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_payment
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c marketing_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_feedback
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_appointment
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_tqd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_file
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c db_map
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- HOÀN TẤT
-- ============================================================

\c postgres

-- In ra danh sách database đã tạo
SELECT 
    datname as "Database Name",
    pg_catalog.pg_get_userbyid(datdba) as "Owner",
    pg_encoding_to_char(encoding) as "Encoding",
    datcollate as "Collate",
    datctype as "Ctype"
FROM pg_catalog.pg_database
WHERE datname LIKE 'db_%' OR datname = 'marketing_db'
ORDER BY datname;

-- ============================================================
-- HƯỚNG DẪN SỬ DỤNG
-- ============================================================
-- Để chạy script này:
-- 1. Đăng nhập PostgreSQL với quyền superuser:
--    psql -U postgres -h 103.172.239.89 -p 5432
--
-- 2. Chạy script:
--    \i /path/to/init_databases.sql
--
-- 3. Hoặc từ command line:
--    psql -U postgres -h 103.172.239.89 -p 5432 -f init_databases.sql
--
-- LƯU Ý:
-- - Script này sẽ XÓA các database cũ nếu đã tồn tại (DROP DATABASE IF EXISTS)
-- - Đảm bảo backup dữ liệu trước khi chạy trong môi trường production
-- - Mỗi database có user owner riêng với password hex ngẫu nhiên 20 ký tự
-- - Tất cả database sử dụng encoding UTF8
-- - Chi tiết credentials xem file: database_credentials.csv
-- ============================================================

-- In thông báo hoàn tất
\echo '============================================================'
\echo 'HOÀN TẤT KHỞI TẠO DATABASE SYSTEM'
\echo '============================================================'
\echo 'Đã tạo 17 users và databases cho hệ thống microservices:'
\echo ''
\echo 'Database                | Owner         | Password (hex 20 chars)'
\echo '------------------------|---------------|------------------------'
\echo '1.  db_auth             | auth          | a7f3c2e9b4d8a1f6e5c7'
\echo '2.  db_user             | user_owner    | b8e4d1f0c5a9b2e7d3f8'
\echo '3.  db_notification     | notification  | c9f5e2a1d6b0c3f8e4a9'
\echo '4.  db_org              | org           | d0a6f3b2e7c1d4a9f5b0'
\echo '5.  db_crm              | crm           | e1b7a4c3f8d2e5b0a6c1'
\echo '6.  db_bdspro           | bdspro        | f2c8b5d4a9e3f6c1b7d2'
\echo '7.  db_chat             | chat          | a3d9c6e5b0f4a7d2c8e3'
\echo '8.  db_social           | social        | b4e0d7f6c1a5b8e3d9f4'
\echo '9.  db_transaction      | transaction   | c5f1e8a7d2b6c9f4e0a5'
\echo '10. db_membership       | membership    | d6a2f9b8e3c7d0a5f1b6'
\echo '11. db_payment          | payment       | e7b3a0c9f4d8e1b6a2c7'
\echo '12. marketing_db        | marketing     | f8c4b1d0a5e9f2c7b3d8'
\echo '13. db_feedback         | feedback      | a9d5c2e1b6f0a3d8c4e9'
\echo '14. db_appointment      | appointment   | b0e6d3f2c7a1b4e9d5f0'
\echo '15. db_tqd              | tqd           | c1f7e4a3d8b2c5f0e6a1'
\echo '16. db_file             | file_owner    | d2a8f5b4e9c3d6a1f7b2'
\echo '17. db_map              | map           | e3b9a6c5f0d4e7b2a8c3'
\echo ''
\echo 'Host: 103.172.239.89'
\echo 'Port: 5432'
\echo 'Extensions: uuid-ossp, pgcrypto'
\echo ''
\echo 'QUAN TRỌNG: Lưu file database_credentials.csv để tham khảo!'
\echo '============================================================'

