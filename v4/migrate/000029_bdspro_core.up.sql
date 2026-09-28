-- Core bdspro tables for phase-1 modules/bdspro routes (product / property / deal).
-- Idempotent against a DB that already has legacy bdspro tables (CREATE IF NOT EXISTS
-- + ADD COLUMN IF NOT EXISTS + heal known GORM column-name typos).

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    parent_id BIGINT,
    asset_id BIGINT,
    property_id BIGINT,
    image_id BIGINT,
    property_type_id BIGINT,
    doc_type_id BIGINT,
    doc_type_note TEXT,
    project_id BIGINT,
    visibility SMALLINT NOT NULL DEFAULT 30,
    name VARCHAR(255),
    code VARCHAR(15),
    category_id BIGINT,
    area DOUBLE PRECISION NOT NULL DEFAULT 0,
    position_url TEXT,
    description TEXT,
    last_price_id BIGINT,
    source_type SMALLINT NOT NULL DEFAULT 0,
    source_status SMALLINT NOT NULL DEFAULT 10,
    source_contact_id BIGINT,
    source_contact_note TEXT,
    source_time TIMESTAMPTZ,
    priority SMALLINT NOT NULL DEFAULT 20,
    note TEXT,
    transaction_type SMALLINT NOT NULL DEFAULT 10,
    sale_status SMALLINT NOT NULL DEFAULT 10,
    sale_visibility SMALLINT NOT NULL DEFAULT 30,
    rent_status SMALLINT NOT NULL DEFAULT 10,
    rent_visibility SMALLINT NOT NULL DEFAULT 30,
    owner_id BIGINT NOT NULL DEFAULT 0,
    owner_of SMALLINT NOT NULL DEFAULT 10,
    province_id BIGINT,
    ward_id BIGINT,
    address VARCHAR(255),
    google_map_link VARCHAR(255),
    archived BOOLEAN NOT NULL DEFAULT FALSE
);
ALTER TABLE products ADD COLUMN IF NOT EXISTS property_id BIGINT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS owner_of SMALLINT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS source_status BIGINT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS source_contact_id BIGINT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS source_contact_note TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS source_time TIMESTAMPTZ;
ALTER TABLE products ADD COLUMN IF NOT EXISTS priority SMALLINT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS doc_type_note TEXT;
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products (deleted_at);
CREATE INDEX IF NOT EXISTS idx_products_owner_id ON products (owner_id);
CREATE INDEX IF NOT EXISTS idx_products_property_id ON products (property_id);

CREATE TABLE IF NOT EXISTS product_price (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    product_id BIGINT,
    currency VARCHAR(16),
    sale_price DECIMAL(15,0),
    sale_commission DECIMAL(15,0),
    sale_commission_type INTEGER,
    deposite DECIMAL(15,0),
    rent_price DECIMAL(15,0),
    rent_commission DECIMAL(15,0),
    rent_commission_type INTEGER,
    rent_payment_cycle INTEGER,
    channel_price BOOLEAN NOT NULL DEFAULT FALSE,
    change_note TEXT,
    changed_by BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_product_price_deleted_at ON product_price (deleted_at);
CREATE INDEX IF NOT EXISTS idx_product_price_product_id ON product_price (product_id);

CREATE TABLE IF NOT EXISTS product_user (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    profile_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    is_owner BOOLEAN NOT NULL DEFAULT TRUE,
    role_id BIGINT NOT NULL DEFAULT 0,
    owner_of SMALLINT NOT NULL DEFAULT 10,
    origin_profile_id BIGINT
);
CREATE INDEX IF NOT EXISTS idx_product_user_deleted_at ON product_user (deleted_at);
CREATE INDEX IF NOT EXISTS idx_product_user_profile_id ON product_user (profile_id);
CREATE INDEX IF NOT EXISTS idx_product_user_product_id ON product_user (product_id);
CREATE INDEX IF NOT EXISTS idx_product_user_origin_profile_id ON product_user (origin_profile_id);

CREATE TABLE IF NOT EXISTS property_info (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    title VARCHAR(255),
    identifier VARCHAR(100),
    property_type_id BIGINT,
    project_id BIGINT,
    avatar_id BIGINT
);
CREATE INDEX IF NOT EXISTS idx_property_info_deleted_at ON property_info (deleted_at);

CREATE TABLE IF NOT EXISTS property_location (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    address_detail VARCHAR(500),
    province_id BIGINT,
    district_id BIGINT,
    ward_id BIGINT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    map_url VARCHAR(500)
);
CREATE INDEX IF NOT EXISTS idx_property_location_deleted_at ON property_location (deleted_at);

CREATE TABLE IF NOT EXISTS property_lineage (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    origin_profile_id BIGINT,
    lineage_status SMALLINT NOT NULL DEFAULT 10,
    property_identify_id BIGINT,
    location_id BIGINT,
    property_info_id BIGINT,
    land_info_id BIGINT,
    building_info_id BIGINT,
    statistic_id BIGINT
);
-- Legacy AutoMigrate typo: gorm tag `column: origin_profile_id` (leading space).
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'property_lineage'
          AND column_name = ' origin_profile_id'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'property_lineage'
          AND column_name = 'origin_profile_id'
    ) THEN
        ALTER TABLE property_lineage RENAME COLUMN " origin_profile_id" TO origin_profile_id;
    END IF;
END $$;
ALTER TABLE property_lineage ADD COLUMN IF NOT EXISTS origin_profile_id BIGINT;
ALTER TABLE property_lineage ADD COLUMN IF NOT EXISTS lineage_status SMALLINT;
ALTER TABLE property_lineage ADD COLUMN IF NOT EXISTS property_identify_id BIGINT;
ALTER TABLE property_lineage ADD COLUMN IF NOT EXISTS location_id BIGINT;
ALTER TABLE property_lineage ADD COLUMN IF NOT EXISTS property_info_id BIGINT;
ALTER TABLE property_lineage ADD COLUMN IF NOT EXISTS land_info_id BIGINT;
ALTER TABLE property_lineage ADD COLUMN IF NOT EXISTS building_info_id BIGINT;
ALTER TABLE property_lineage ADD COLUMN IF NOT EXISTS statistic_id BIGINT;
CREATE INDEX IF NOT EXISTS idx_property_lineage_deleted_at ON property_lineage (deleted_at);
CREATE INDEX IF NOT EXISTS idx_property_lineage_origin_profile_id ON property_lineage (origin_profile_id);

CREATE TABLE IF NOT EXISTS property_user (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    owner_origin_id BIGINT NOT NULL,
    property_lineage_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL DEFAULT 0,
    origin_id BIGINT,
    ownered_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
    hidden_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_property_user_deleted_at ON property_user (deleted_at);
CREATE INDEX IF NOT EXISTS idx_property_user_owner_origin_id ON property_user (owner_origin_id);
CREATE INDEX IF NOT EXISTS idx_property_user_property_lineage_id ON property_user (property_lineage_id);

CREATE TABLE IF NOT EXISTS deals (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    owner_id BIGINT NOT NULL DEFAULT 0,
    owner_type SMALLINT NOT NULL DEFAULT 10,
    name VARCHAR(255),
    target_profit DOUBLE PRECISION NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 10,
    deal_type SMALLINT NOT NULL DEFAULT 10,
    note TEXT,
    charge_person_id BIGINT,
    cancel_reason TEXT,
    is_unilateral BOOLEAN NOT NULL DEFAULT FALSE,
    is_manual BOOLEAN NOT NULL DEFAULT FALSE,
    bank_account_id BIGINT,
    from_date TIMESTAMPTZ,
    to_date TIMESTAMPTZ,
    allow_sharing BOOLEAN NOT NULL DEFAULT FALSE,
    member_can_add_transaction BOOLEAN NOT NULL DEFAULT FALSE,
    only_owner_get_commission BOOLEAN NOT NULL DEFAULT FALSE,
    internal_note TEXT,
    allow_manual_input BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT,
    share_visibility SMALLINT NOT NULL DEFAULT 10
);
CREATE INDEX IF NOT EXISTS idx_deals_deleted_at ON deals (deleted_at);
CREATE INDEX IF NOT EXISTS idx_deals_owner_id ON deals (owner_id);

CREATE TABLE IF NOT EXISTS deal_members (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    deal_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    member_type SMALLINT NOT NULL DEFAULT 10,
    status SMALLINT NOT NULL DEFAULT 10,
    role_id BIGINT NOT NULL DEFAULT 0,
    amount_commit DOUBLE PRECISION NOT NULL DEFAULT 0,
    commission_value DOUBLE PRECISION NOT NULL DEFAULT 0,
    commission_type SMALLINT NOT NULL DEFAULT 10,
    note TEXT,
    is_unilateral BOOLEAN NOT NULL DEFAULT FALSE,
    done_investment BOOLEAN NOT NULL DEFAULT FALSE,
    role_key INTEGER NOT NULL DEFAULT 430,
    message TEXT,
    invited_at TIMESTAMPTZ,
    responded_at TIMESTAMPTZ,
    withdrawn_at TIMESTAMPTZ,
    inviter_id BIGINT NOT NULL DEFAULT 0,
    color_id INTEGER,
    is_owner BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_deal_members_deleted_at ON deal_members (deleted_at);
CREATE INDEX IF NOT EXISTS idx_deal_members_deal_id ON deal_members (deal_id);
CREATE INDEX IF NOT EXISTS idx_deal_members_member_id ON deal_members (member_id);
