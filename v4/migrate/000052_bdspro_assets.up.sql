-- Assets + legal documents for modules/bdspro (property/new asset payload).
-- Idempotent against DBs that already absorbed the legacy bdspro `assets` / `asset_legals` tables.

CREATE TABLE IF NOT EXISTS assets (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    name TEXT NOT NULL DEFAULT '',
    ward_id BIGINT,
    province_id BIGINT,
    product_id BIGINT,
    parent_asset_id BIGINT,
    split_merge_status TEXT NOT NULL DEFAULT '',
    purchase_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    purchase_date TIMESTAMPTZ,
    legal_status BIGINT NOT NULL DEFAULT 10,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    rent_status BIGINT NOT NULL DEFAULT 10,
    address TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    area DOUBLE PRECISION NOT NULL DEFAULT 0,
    property_type_id BIGINT,
    image_id BIGINT,
    owner_id BIGINT,
    owner_of BIGINT NOT NULL DEFAULT 10
);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS product_id BIGINT;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS property_type_id BIGINT;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS image_id BIGINT;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS owner_id BIGINT;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS owner_of BIGINT;
CREATE INDEX IF NOT EXISTS idx_assets_deleted_at ON assets (deleted_at);
CREATE INDEX IF NOT EXISTS idx_assets_owner_id ON assets (owner_id);
CREATE INDEX IF NOT EXISTS idx_assets_product_id ON assets (product_id);

CREATE TABLE IF NOT EXISTS asset_legals (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    asset_id BIGINT NOT NULL,
    document_name TEXT NOT NULL DEFAULT '',
    document_url TEXT NOT NULL DEFAULT '',
    document_type TEXT NOT NULL DEFAULT '',
    issued_date TIMESTAMPTZ,
    expiry_date TIMESTAMPTZ,
    description TEXT NOT NULL DEFAULT '',
    related_split_merge_id BIGINT
);
CREATE INDEX IF NOT EXISTS idx_asset_legals_deleted_at ON asset_legals (deleted_at);
CREATE INDEX IF NOT EXISTS idx_asset_legals_asset_id ON asset_legals (asset_id);
