DROP INDEX IF EXISTS idx_tb_contact_last_interaction_at;
DROP INDEX IF EXISTS idx_tb_contact_priority_pin_at;

ALTER TABLE tb_contact DROP COLUMN IF EXISTS last_interaction_at;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS priority_pin_at;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS contact_tags;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS lead_id;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS status;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS visibility;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS address;
