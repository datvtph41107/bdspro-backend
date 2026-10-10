-- Extend tb_contact for modules/crm (legacy crm-service / AutoMigrate shape).
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS avatar TEXT;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS email TEXT;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS zalo TEXT;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS company TEXT;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS note TEXT;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS birthday DATE;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS owner_id BIGINT;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS owner_of SMALLINT DEFAULT 30;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS friend_id BIGINT;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS origin_profile_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_tb_contact_owner ON tb_contact (owner_id, owner_of) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tb_contact_origin_profile_id ON tb_contact (origin_profile_id) WHERE deleted_at IS NULL;
