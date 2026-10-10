DROP INDEX IF EXISTS idx_tb_contact_origin_profile_id;
DROP INDEX IF EXISTS idx_tb_contact_owner;

ALTER TABLE tb_contact DROP COLUMN IF EXISTS origin_profile_id;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS friend_id;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS owner_of;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS owner_id;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS birthday;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS note;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS company;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS zalo;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS email;
ALTER TABLE tb_contact DROP COLUMN IF EXISTS avatar;
