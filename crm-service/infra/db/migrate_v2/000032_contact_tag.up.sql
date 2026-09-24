-- Contact custom tags + relation to contacts.

CREATE TABLE IF NOT EXISTS tb_contact_tag (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    owner_id BIGINT,
    owner_of INT,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_tb_contact_tag_deleted_at ON tb_contact_tag (deleted_at);
CREATE INDEX IF NOT EXISTS idx_tb_contact_tag_owner ON tb_contact_tag (owner_id, owner_of) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tb_contact_tag_default_active ON tb_contact_tag (is_default, is_active) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS tb_contact_tag_relation (
    contact_id BIGINT NOT NULL REFERENCES tb_contact (id),
    tag_id BIGINT NOT NULL REFERENCES tb_contact_tag (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (contact_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_tb_contact_tag_relation_tag_id ON tb_contact_tag_relation (tag_id);

-- Seed system default tags (idempotent by name + is_default).
INSERT INTO tb_contact_tag (name, is_default, is_active, created_at, updated_at)
SELECT v.name, TRUE, TRUE, NOW(), NOW()
FROM (VALUES
    ('Khách hàng'),
    ('Môi giới'),
    ('Đối tác'),
    ('Chủ nhà')
) AS v(name)
WHERE NOT EXISTS (
    SELECT 1
    FROM tb_contact_tag t
    WHERE t.name = v.name
      AND t.is_default = TRUE
      AND t.deleted_at IS NULL
);
