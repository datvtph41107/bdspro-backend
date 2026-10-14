-- Remaining tb_contact columns referenced by modules/crm/entity.ContactEntity.
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS visibility SMALLINT DEFAULT 30;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS status SMALLINT DEFAULT 10;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS lead_id BIGINT;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS contact_tags INTEGER[];
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS priority_pin_at TIMESTAMPTZ;
ALTER TABLE tb_contact ADD COLUMN IF NOT EXISTS last_interaction_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_tb_contact_priority_pin_at ON tb_contact (priority_pin_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tb_contact_last_interaction_at ON tb_contact (last_interaction_at) WHERE deleted_at IS NULL;
