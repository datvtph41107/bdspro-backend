-- Migration: Bỏ unique constraint/index trên phone trong bảng tb_contact
-- Date: 2025-01-XX
-- Description: Cho phép nhiều contact có cùng số điện thoại

-- Drop unique index nếu tồn tại (tên index có thể khác nhau tùy database)
DO $$
BEGIN
    -- Thử drop index với các tên có thể có
    IF EXISTS (SELECT 1 FROM pg_indexes WHERE tablename = 'tb_contact' AND indexname LIKE '%phone%owner%') THEN
        EXECUTE (
            SELECT 'DROP INDEX IF EXISTS ' || string_agg(indexname, ', ')
            FROM pg_indexes 
            WHERE tablename = 'tb_contact' 
            AND indexname LIKE '%phone%owner%'
        );
    END IF;
    
    -- Drop unique constraint nếu tồn tại
    IF EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conrelid = 'tb_contact'::regclass 
        AND contype = 'u'
        AND (
            conkey::text LIKE '%phone%' 
            OR pg_get_constraintdef(oid) LIKE '%phone%owner%'
        )
    ) THEN
        EXECUTE (
            SELECT 'ALTER TABLE tb_contact DROP CONSTRAINT IF EXISTS ' || conname
            FROM pg_constraint 
            WHERE conrelid = 'tb_contact'::regclass 
            AND contype = 'u'
            AND (
                conkey::text LIKE '%phone%' 
                OR pg_get_constraintdef(oid) LIKE '%phone%owner%'
            )
            LIMIT 1
        );
    END IF;
END $$;

-- Drop index cụ thể nếu có (thử các tên phổ biến)
DROP INDEX IF EXISTS idx_tb_contact_phone_owner_id_owner_of;
DROP INDEX IF EXISTS uniq_tb_contact_phone_owner_id_owner_of;
DROP INDEX IF EXISTS tb_contact_phone_owner_id_owner_of_idx;
DROP INDEX IF EXISTS tb_contact_phone_owner_id_owner_of_key;
