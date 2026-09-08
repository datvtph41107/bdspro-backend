BEGIN;

CREATE TABLE IF NOT EXISTS organizations (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL CHECK (btrim(name) <> ''),
    type VARCHAR(32) NOT NULL DEFAULT 'business'
        CHECK (type IN ('business', 'team', 'other')),
    tax_code VARCHAR(64),
    phone VARCHAR(32),
    email VARCHAR(255),
    address TEXT,
    website VARCHAR(500),
    description TEXT,
    status VARCHAR(32) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'locked', 'archived')),
    verification_status VARCHAR(32) NOT NULL DEFAULT 'unverified'
        CHECK (verification_status IN ('unverified', 'verifying', 'verified')),
    warning_level VARCHAR(32) NOT NULL DEFAULT 'none'
        CHECK (warning_level IN ('none', 'mild', 'medium', 'severe', 'warned')),
    warning_count INTEGER NOT NULL DEFAULT 0 CHECK (warning_count >= 0),
    owner_profile_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_organizations_code_active
    ON organizations (lower(code)) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS ux_organizations_tax_code_active
    ON organizations (lower(tax_code))
    WHERE archived_at IS NULL AND tax_code IS NOT NULL AND btrim(tax_code) <> '';
CREATE INDEX IF NOT EXISTS ix_organizations_owner
    ON organizations (owner_profile_id) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS organization_members (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    profile_id BIGINT NOT NULL,
    role_key VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'active'
        CHECK (status IN ('invited', 'active', 'suspended', 'removed')),
    joined_at TIMESTAMPTZ,
    removed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_organization_members_active
    ON organization_members (organization_id, profile_id)
    WHERE status <> 'removed';
CREATE INDEX IF NOT EXISTS ix_organization_members_profile
    ON organization_members (profile_id, organization_id);

INSERT INTO permissions (name, key, description, permission_type, module, priority)
SELECT seed.name, seed.key, seed.description, 'ORGANIZATION', 'organization', seed.priority
FROM (VALUES
    ('Xem tổ chức', 'ORGANIZATION_VIEW', 'Đọc danh sách và chi tiết tổ chức', 10),
    ('Quản lý tổ chức', 'ORGANIZATION_MANAGE', 'Tạo, sửa, khóa, duyệt và lưu trữ tổ chức', 30)
) AS seed(name, key, description, priority)
WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.key = seed.key AND p.deleted_at IS NULL);

COMMIT;
