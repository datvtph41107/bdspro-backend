-- QHPRO User Service canonical baseline schema.
-- This directory is the ONLY schema migration authority for User Service.
-- Fresh-clone target: PostgreSQL 17+. Runtime startup performs no DDL.

BEGIN;

-- AdminAccessDomain (internal/domain/access/admin_access.go)
CREATE TABLE IF NOT EXISTS admin_access_control (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    ip_address VARCHAR(45),
    ip_range VARCHAR(50),
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    device_type VARCHAR(50),
    auth_type VARCHAR(20),
    status VARCHAR(20) DEFAULT 'active',
    effective_from TIMESTAMPTZ,
    effective_to TIMESTAMPTZ,
    max_devices INTEGER DEFAULT 3,
    require_2fa BOOLEAN DEFAULT true,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- AdminAccessLogDomain (internal/domain/access/admin_access.go)
CREATE TABLE IF NOT EXISTS admin_access_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    ip_address VARCHAR(45),
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    action VARCHAR(50),
    status VARCHAR(20),
    reason VARCHAR(500),
    user_agent VARCHAR(500),
    created_at TIMESTAMPTZ
);

-- Color (internal/domain/access/color.go)
CREATE TABLE IF NOT EXISTS colors (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(50) NOT NULL,
    hex_code VARCHAR(24) NOT NULL,
    content_color VARCHAR(24) NOT NULL,
    background_color VARCHAR(24) NOT NULL,
    color_key VARCHAR(24) DEFAULT 'default',
    description VARCHAR(255),
    is_active BOOLEAN DEFAULT true
);

-- Permission (internal/domain/access/permission.go)
CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(100) NOT NULL,
    key VARCHAR(50),
    description VARCHAR(500),
    permission_type VARCHAR(20) DEFAULT 'ORGANIZATION',
    module VARCHAR(50),
    parent_id BIGINT,
    priority INTEGER DEFAULT 10,
    required_key VARCHAR(10)
);

-- Role (internal/domain/access/role.go)
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    role_name VARCHAR(100) NOT NULL,
    role_description VARCHAR(255),
    key VARCHAR(50),
    is_default BOOLEAN DEFAULT false,
    domain_type VARCHAR(20) DEFAULT 'ORGANIZATION',
    role_group_id BIGINT,
    color_id BIGINT,
    organization_id BIGINT NOT NULL,
    allow_assign BOOLEAN DEFAULT false,
    role_key INTEGER
);

-- RoleGroup (internal/domain/access/role_group.go)
CREATE TABLE IF NOT EXISTS role_groups (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    group_name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    key INTEGER,
    code VARCHAR(20)
);

-- RoleProfile (internal/domain/access/role_profile.go)
CREATE TABLE IF NOT EXISTS role_profiles (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    profile_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL
);

-- AuthConfig (internal/domain/auth/auth_config.go)
CREATE TABLE IF NOT EXISTS auth_config (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    config_key VARCHAR(20) NOT NULL,
    config_value TEXT,
    is_active BOOLEAN DEFAULT true
);

-- AuthMethod (internal/domain/auth/auth_method.go)
CREATE TABLE IF NOT EXISTS auth_method (
    id BIGSERIAL PRIMARY KEY,
    provider VARCHAR(20) NOT NULL,
    auth_name VARCHAR(255) NOT NULL,
    password TEXT,
    is_sensitive BOOLEAN,
    avatar VARCHAR(255),
    full_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(255),
    role_key INTEGER,
    status SMALLINT DEFAULT 1,
    locked_at TIMESTAMPTZ,
    locked_until TIMESTAMPTZ,
    lock_reason VARCHAR(500),
    locked_by BIGINT,
    private_key VARCHAR(512),
    public_key VARCHAR(512),
    auth_key VARCHAR(512),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id BIGINT
);

-- DeviceEntity (internal/domain/auth/device.go)
CREATE TABLE IF NOT EXISTS auth_devices (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    device_type VARCHAR(50),
    platform VARCHAR(100),
    os_version VARCHAR(100),
    app_version VARCHAR(50),
    build_number VARCHAR(50),
    manufacturer VARCHAR(100),
    model VARCHAR(100),
    locale VARCHAR(20),
    timezone VARCHAR(50),
    push_token VARCHAR(255),
    last_seen_at TIMESTAMPTZ,
    auth_id BIGINT,
    profile_id BIGINT,
    organization_id BIGINT,
    ip_address VARCHAR(64),
    user_agent VARCHAR(512)
);

-- UserOTPEntity (internal/domain/auth/otp.go)
CREATE TABLE IF NOT EXISTS user_otp (
    auth_id BIGINT PRIMARY KEY,
    otp VARCHAR(10),
    otp_send_time INTEGER DEFAULT 0,
    otp_check_time INTEGER DEFAULT 0,
    otp_date TIMESTAMPTZ,
    expired_time TIMESTAMPTZ,
    activate BOOLEAN DEFAULT false
);

-- UserSessionEntity (internal/domain/auth/session.go)
CREATE TABLE IF NOT EXISTS user_session (
    session_id BIGSERIAL PRIMARY KEY,
    auth_id BIGINT,
    device_id VARCHAR(255),
    platform VARCHAR(50),
    version VARCHAR(50),
    os VARCHAR(50),
    device_name VARCHAR(255),
    created_date TIMESTAMPTZ,
    finished_date TIMESTAMPTZ,
    last_login TIMESTAMPTZ,
    logout_at TIMESTAMPTZ,
    last_req TIMESTAMPTZ,
    ip_request VARCHAR(50),
    total_request BIGINT,
    user_agent TEXT,
    activate BOOLEAN,
    session_key TEXT,
    client_id TEXT
);

-- UserStatusEntity (internal/domain/auth/status.go)
CREATE TABLE IF NOT EXISTS user_status (
    auth_id BIGINT PRIMARY KEY,
    active BOOLEAN,
    verified BOOLEAN,
    locked_until TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- AuthUser (internal/domain/auth/user_info.go)
CREATE TABLE IF NOT EXISTS user_info (
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    profile_id BIGINT PRIMARY KEY,
    status SMALLINT DEFAULT 10,
    locked_at TIMESTAMPTZ,
    locked_until TIMESTAMPTZ,
    lock_reason VARCHAR(500),
    locked_by BIGINT,
    role_id BIGINT
);

-- UserPINEntity (internal/domain/auth/user_pin.go)
CREATE TABLE IF NOT EXISTS user_pin (
    auth_id BIGINT PRIMARY KEY,
    pin VARCHAR(255),
    pin_check_time INTEGER DEFAULT 0,
    pin_date TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true,
    locked_until TIMESTAMPTZ
);

-- AdminProfile (internal/models/admin_profile.go)
CREATE TABLE IF NOT EXISTS admin_profiles (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    auth_id BIGINT NOT NULL,
    profile_id BIGINT,
    address VARCHAR(255),
    gender INTEGER DEFAULT 0,
    birth TIMESTAMPTZ,
    full_name VARCHAR(255),
    email VARCHAR(50),
    phone VARCHAR(20),
    avatar VARCHAR(255),
    job_title VARCHAR(100),
    work_at TIMESTAMPTZ,
    role_id BIGINT,
    role_key VARCHAR(50),
    role_type VARCHAR(50),
    role_description VARCHAR(255),
    attachments JSONB,
    internal_notes TEXT,
    send_notification BOOLEAN DEFAULT false,
    status INTEGER DEFAULT 10,
    last_login_at TIMESTAMPTZ,
    total_login INTEGER DEFAULT 0,
    verified_at TIMESTAMPTZ,
    locked_at TIMESTAMPTZ
);

-- BlockEntity (internal/models/block.go)
CREATE TABLE IF NOT EXISTS block (
    profile_id BIGINT,
    blocked_id BIGINT,
    PRIMARY KEY (profile_id, blocked_id)
);

-- BookmarkUserEntity (internal/models/bookmark_user.go)
CREATE TABLE IF NOT EXISTS bookmark_user (
    admin_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (admin_id, user_id)
);

-- CertificationEntity (internal/models/certification.go)
CREATE TABLE IF NOT EXISTS certification (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    approve_user_id BIGINT,
    name VARCHAR(255),
    file_name VARCHAR(255),
    file_url TEXT,
    file_type VARCHAR(255),
    verified_status SMALLINT DEFAULT 10,
    issuer VARCHAR(255),
    issue_date DATE
);

-- ContactEntity (internal/models/contact.go)
CREATE TABLE IF NOT EXISTS tb_contact (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    profile_id BIGINT,
    phone TEXT,
    full_name TEXT
);

-- FollowEntity (internal/models/follow.go)
CREATE TABLE IF NOT EXISTS follow (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    following_id BIGINT NOT NULL,
    status INTEGER DEFAULT 1
);

-- FriendEntity (internal/models/friend.go)
CREATE TABLE IF NOT EXISTS friend (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    receiver_id BIGINT NOT NULL,
    status INTEGER DEFAULT 1,
    responded_at TIMESTAMPTZ,
    group_id BIGINT,
    group_receiver_id BIGINT
);

-- GroupEntity (internal/models/group.go)
CREATE TABLE IF NOT EXISTS db_group (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    text_color VARCHAR(12),
    background_color VARCHAR(12)
);

-- KYCEntity (internal/models/kyc.go)
CREATE TABLE IF NOT EXISTS kyc (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    profile_id BIGINT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    identity_card VARCHAR(20) NOT NULL,
    front_image VARCHAR(500) NOT NULL,
    back_image VARCHAR(500) NOT NULL,
    selfie_image VARCHAR(500),
    status INTEGER DEFAULT 10,
    reject_reason VARCHAR(500),
    reviewed_by BIGINT,
    reviewed_at TIMESTAMPTZ,
    id_number VARCHAR(20),
    date_of_birth VARCHAR(20),
    expiry_date VARCHAR(20)
);

-- MainAreaEntity (internal/models/main_area.go)
CREATE TABLE IF NOT EXISTS main_area (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(512),
    number_profile BIGINT NOT NULL DEFAULT 0
);

-- MainAreaProfileEntity (internal/models/main_area_profile.go)
CREATE TABLE IF NOT EXISTS main_area_profile (
    profile_id BIGINT,
    main_area_id BIGINT,
    created_at TIMESTAMPTZ,
    PRIMARY KEY (profile_id, main_area_id)
);

-- PriceTableDomain (internal/models/pack_price_table_domain.go)
CREATE TABLE IF NOT EXISTS pack_price_table (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    price DOUBLE PRECISION NOT NULL DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'VND',
    description TEXT,
    is_active BOOLEAN DEFAULT true
);

-- CatalogProduct (internal/models/product_catalog_v2.go)
CREATE TABLE IF NOT EXISTS catalog_products (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(128) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- CatalogPlan (internal/models/product_catalog_v2.go)
CREATE TABLE IF NOT EXISTS catalog_plans (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    code VARCHAR(128) NOT NULL,
    tier_rank INTEGER,
    status VARCHAR(20) NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- CatalogPlanVersion (internal/models/product_catalog_v2.go)
CREATE TABLE IF NOT EXISTS catalog_plan_versions (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL,
    version VARCHAR(32) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    subject_scope VARCHAR(32) NOT NULL,
    subscription_term_days INTEGER,
    effective_from TIMESTAMPTZ,
    effective_until TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    source VARCHAR(32) NOT NULL,
    terms_checksum CHAR(64) NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    published_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- CatalogPlanEntitlement (internal/models/product_catalog_v2.go)
CREATE TABLE IF NOT EXISTS catalog_plan_entitlements (
    id BIGSERIAL PRIMARY KEY,
    plan_version_id BIGINT NOT NULL,
    code VARCHAR(128) NOT NULL,
    kind VARCHAR(32) NOT NULL,
    feature_code VARCHAR(128),
    meter_code VARCHAR(128),
    amount BIGINT NOT NULL DEFAULT 0,
    unlimited BOOLEAN NOT NULL DEFAULT false,
    period VARCHAR(32) NOT NULL,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL
);

-- CatalogPlanOperationPolicy (internal/models/product_catalog_v2.go)
CREATE TABLE IF NOT EXISTS catalog_plan_operation_policies (
    id BIGSERIAL PRIMARY KEY,
    plan_version_id BIGINT NOT NULL,
    operation_code VARCHAR(128) NOT NULL,
    feature_code VARCHAR(128) NOT NULL,
    meter_code VARCHAR(128),
    units_per_action BIGINT NOT NULL DEFAULT 0,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL
);

-- CatalogPriceItem (internal/models/product_catalog_v2.go)
CREATE TABLE IF NOT EXISTS catalog_price_items (
    id BIGSERIAL PRIMARY KEY,
    plan_version_id BIGINT NOT NULL,
    code VARCHAR(128) NOT NULL,
    kind VARCHAR(20) NOT NULL,
    currency CHAR(3) NOT NULL,
    amount_minor BIGINT NOT NULL,
    billing_unit VARCHAR(128) NOT NULL,
    meter_code VARCHAR(128),
    quantity BIGINT NOT NULL,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL
);

-- CatalogSubscription (internal/models/product_catalog_v2.go)
CREATE TABLE IF NOT EXISTS catalog_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    subscription_key VARCHAR(128) NOT NULL,
    subject_kind VARCHAR(32) NOT NULL,
    subject_id VARCHAR(128) NOT NULL,
    product_id BIGINT NOT NULL,
    plan_version_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    current_period_start TIMESTAMPTZ NOT NULL,
    current_period_end TIMESTAMPTZ NOT NULL,
    access_until TIMESTAMPTZ,
    auto_renew BOOLEAN NOT NULL DEFAULT false,
    order_reference VARCHAR(128),
    pending_plan_version_id BIGINT,
    pending_effective_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- CatalogSubscriptionEvent (internal/models/product_catalog_v2.go)
CREATE TABLE IF NOT EXISTS catalog_subscription_events (
    id BIGSERIAL PRIMARY KEY,
    subscription_id BIGINT NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    from_plan_version_id BIGINT,
    to_plan_version_id BIGINT,
    effective_at TIMESTAMPTZ NOT NULL,
    reason TEXT NOT NULL,
    actor_kind VARCHAR(32) NOT NULL,
    actor_id VARCHAR(128) NOT NULL,
    request_id VARCHAR(128) NOT NULL,
    operation_id VARCHAR(128) NOT NULL,
    metadata JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

-- ProfessionEntity (internal/models/profession.go)
CREATE TABLE IF NOT EXISTS profession (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    approve_user_id BIGINT,
    name VARCHAR(255) NOT NULL,
    issuer VARCHAR(255),
    issue_date DATE,
    verified_status SMALLINT NOT NULL DEFAULT 10,
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- UserProfileEntity (internal/models/profile.go)
CREATE TABLE IF NOT EXISTS user_profile (
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    profile_id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20),
    phone2 VARCHAR(20),
    email VARCHAR(50),
    tax_code VARCHAR(13),
    full_name VARCHAR(255),
    address VARCHAR(255),
    gender INTEGER,
    avatar VARCHAR(255),
    tick_verified BOOLEAN DEFAULT false,
    background_image VARCHAR(255),
    referral_code VARCHAR(255),
    referral_by VARCHAR(255),
    role_type INTEGER DEFAULT 10,
    role_key INTEGER DEFAULT 10,
    role_real_estate INTEGER DEFAULT 10,
    position VARCHAR(255),
    workplace VARCHAR(255),
    role_title VARCHAR(255),
    department_id BIGINT,
    birth TIMESTAMPTZ,
    start_date TIMESTAMPTZ,
    front_identify TEXT,
    back_identify TEXT,
    province_id BIGINT,
    ward_id BIGINT,
    plan_id BIGINT,
    plan_at TIMESTAMPTZ,
    slogan TEXT,
    website_url TEXT,
    zalo_url TEXT,
    facebook_url TEXT,
    instagram_url TEXT,
    twitter_url TEXT,
    linkedin_url TEXT,
    youtube_url TEXT,
    introduction TEXT,
    visibility INTEGER,
    status_online INTEGER,
    profile_visibility INTEGER,
    tab_default INTEGER,
    signature_visible BOOLEAN DEFAULT false,
    visibility_introduce INTEGER,
    visibility_profession INTEGER,
    visibility_main_area INTEGER,
    visibility_friends INTEGER,
    visibility_signature INTEGER,
    view_roles INTEGER[],
    status INTEGER DEFAULT 10,
    last_login_at TIMESTAMPTZ,
    total_login INTEGER,
    warning INTEGER DEFAULT 0,
    role_id BIGINT
);

-- Profile (internal/models/profile.transfer.go)
CREATE TABLE IF NOT EXISTS profile_transfer (
    profile_id BIGINT PRIMARY KEY,
    phone TEXT,
    full_name TEXT
);

-- ProfileDeleted (internal/models/profile_deleted.go)
CREATE TABLE IF NOT EXISTS profile_deleted (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT,
    phone VARCHAR(20),
    phone2 VARCHAR(20),
    email VARCHAR(50),
    tax_code VARCHAR(13),
    full_name VARCHAR(255),
    address VARCHAR(255),
    gender INTEGER,
    avatar VARCHAR(255),
    tick_verified BOOLEAN DEFAULT false,
    background_image VARCHAR(255),
    referral_code VARCHAR(255),
    referral_by VARCHAR(255),
    role_type INTEGER DEFAULT 10,
    role_key INTEGER DEFAULT 10,
    role_real_estate INTEGER DEFAULT 10,
    position VARCHAR(255),
    department_id BIGINT,
    birth TIMESTAMPTZ,
    start_date TIMESTAMPTZ,
    front_identify TEXT,
    back_identify TEXT,
    plan_id BIGINT,
    plan_at TIMESTAMPTZ,
    slogan TEXT,
    website TEXT,
    facebook TEXT,
    instagram TEXT,
    twitter TEXT,
    linkedin TEXT,
    youtube TEXT,
    introduction TEXT,
    visibility INTEGER,
    status_online INTEGER,
    profile_visibility INTEGER,
    tab_default INTEGER,
    status INTEGER DEFAULT 10,
    warning INTEGER DEFAULT 0,
    role_id BIGINT,
    deleted_by BIGINT,
    deleted_reason TEXT,
    deleted_at TIMESTAMPTZ,
    original_created_at TIMESTAMPTZ,
    original_updated_at TIMESTAMPTZ
);

-- ProfileMediaEntity (internal/models/profile_media.go)
CREATE TABLE IF NOT EXISTS profile_media (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255),
    file_name VARCHAR(255),
    file_url TEXT,
    file_type VARCHAR(50),
    status SMALLINT DEFAULT 0
);

-- PurposeUseEntity (internal/models/purpose_use.go)
CREATE TABLE IF NOT EXISTS purpose_use (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    number_profile BIGINT NOT NULL DEFAULT 0
);

-- PurposeUseProfileEntity (internal/models/purpose_use_profile.go)
CREATE TABLE IF NOT EXISTS purpose_use_profile (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    profile_id BIGINT NOT NULL,
    purpose_use_id BIGINT NOT NULL
);

-- TagEntity (internal/models/tag.go)
CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    is_default BOOLEAN DEFAULT false,
    is_active BOOLEAN NOT NULL,
    user_id BIGINT,
    tag_type INTEGER NOT NULL
);

-- TagUserEntity (internal/models/tag_user.go)
CREATE TABLE IF NOT EXISTS tag_user (
    profile_id BIGINT,
    tag_id BIGINT,
    created_at TIMESTAMPTZ,
    PRIMARY KEY (profile_id, tag_id)
);

-- checkoutSnapshotRow (infra/postgres/checkout/snapshot.go)
CREATE TABLE IF NOT EXISTS subscription_checkout_snapshots (
    id BIGSERIAL PRIMARY KEY,
    subject_kind VARCHAR(32) NOT NULL,
    subject_id VARCHAR(128) NOT NULL,
    command_key TEXT NOT NULL,
    request_hash CHAR(64) NOT NULL,
    product_id BIGINT,
    product_code VARCHAR(128) NOT NULL,
    plan_id BIGINT NOT NULL,
    plan_code VARCHAR(128) NOT NULL,
    plan_version_id BIGINT NOT NULL,
    plan_version VARCHAR(128) NOT NULL,
    tier_rank INTEGER NOT NULL CHECK (tier_rank > 0),
    subscription_term_days INTEGER NOT NULL CHECK (subscription_term_days > 0),
    terms_checksum CHAR(64) NOT NULL,
    subject_scope VARCHAR(32) NOT NULL,
    currency CHAR(3) NOT NULL,
    amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS subscription_settlement_receipts (
    effect_key VARCHAR(128) PRIMARY KEY,
    fingerprint CHAR(64) NOT NULL,
    order_id BIGINT NOT NULL,
    subject_kind VARCHAR(32) NOT NULL,
    subject_id VARCHAR(128) NOT NULL,
    product_code VARCHAR(64) NOT NULL,
    plan_code VARCHAR(64) NOT NULL,
    plan_version_id BIGINT NOT NULL,
    plan_version VARCHAR(64) NOT NULL,
    tier_rank INTEGER NOT NULL,
    subscription_term_days INTEGER NOT NULL,
    terms_checksum CHAR(64) NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    action VARCHAR(32) NOT NULL,
    subscription_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (order_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    PRIMARY KEY (role_id, permission_id)
);
CREATE INDEX IF NOT EXISTS ix_role_permissions_permission_id ON role_permissions(permission_id);

CREATE TABLE IF NOT EXISTS group_permissions (
    role_group_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    PRIMARY KEY (role_group_id, permission_id)
);
CREATE INDEX IF NOT EXISTS ix_group_permissions_permission_id ON group_permissions(permission_id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_roles_org_key ON roles(organization_id, key) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS ix_roles_domain_type ON roles(domain_type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS ix_role_profiles_profile_id ON role_profiles(profile_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS ix_role_profiles_role_id ON role_profiles(role_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_products_code ON catalog_products(code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_plans_product_code ON catalog_plans(product_id, code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_plan_versions_plan_version ON catalog_plan_versions(plan_id, version);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_entitlement_code ON catalog_plan_entitlements(plan_version_id, code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_operation_code ON catalog_plan_operation_policies(plan_version_id, operation_code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_price_code ON catalog_price_items(plan_version_id, code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_subscription_key ON catalog_subscriptions(subscription_key);
CREATE INDEX IF NOT EXISTS ix_catalog_subscription_subject ON catalog_subscriptions(subject_kind, subject_id, product_id, status);
CREATE INDEX IF NOT EXISTS ix_catalog_subscription_events_subscription ON catalog_subscription_events(subscription_id, created_at);
CREATE UNIQUE INDEX IF NOT EXISTS ux_checkout_command ON subscription_checkout_snapshots(subject_kind, subject_id, command_key);
CREATE UNIQUE INDEX IF NOT EXISTS ux_role_groups_key ON role_groups(key);
CREATE UNIQUE INDEX IF NOT EXISTS ux_role_groups_code ON role_groups(code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_auth_config_config_key ON auth_config(config_key);
CREATE UNIQUE INDEX IF NOT EXISTS idx_device_unique ON auth_devices(device_id);
CREATE INDEX IF NOT EXISTS ix_admin_profiles_auth_id ON admin_profiles(auth_id);
CREATE INDEX IF NOT EXISTS ix_kyc_profile_id ON kyc(profile_id);
CREATE INDEX IF NOT EXISTS ix_catalog_plans_product_id ON catalog_plans(product_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_plans_code ON catalog_plans(code);
CREATE INDEX IF NOT EXISTS ix_catalog_plan_versions_plan_id ON catalog_plan_versions(plan_id);
CREATE INDEX IF NOT EXISTS ix_catalog_plan_entitlements_plan_version_id ON catalog_plan_entitlements(plan_version_id);
CREATE INDEX IF NOT EXISTS ix_catalog_plan_operation_policies_plan_version_id ON catalog_plan_operation_policies(plan_version_id);
CREATE INDEX IF NOT EXISTS ix_catalog_price_items_plan_version_id ON catalog_price_items(plan_version_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_subscriptions_subscription_key ON catalog_subscriptions(subscription_key);
CREATE INDEX IF NOT EXISTS idx_catalog_subscription_subject ON catalog_subscriptions(subject_kind, subject_id);
CREATE INDEX IF NOT EXISTS ix_catalog_subscriptions_product_id ON catalog_subscriptions(product_id);
CREATE INDEX IF NOT EXISTS ix_catalog_subscriptions_plan_version_id ON catalog_subscriptions(plan_version_id);
CREATE INDEX IF NOT EXISTS ix_catalog_subscription_events_subscription_id ON catalog_subscription_events(subscription_id);
CREATE INDEX IF NOT EXISTS ix_profile_deleted_profile_id ON profile_deleted(profile_id);
CREATE INDEX IF NOT EXISTS ix_purpose_use_profile_profile_id ON purpose_use_profile(profile_id);
CREATE INDEX IF NOT EXISTS ix_purpose_use_profile_purpose_use_id ON purpose_use_profile(purpose_use_id);

-- Permission thương mại/quota là dữ liệu IAM thuộc User Service. Seed tại
-- migration owner để Auth có thể dựng catalog quyền mà không cần scripts rời.
INSERT INTO permissions (name, key, description, permission_type, module, priority)
SELECT seed.name, seed.key, seed.description, seed.permission_type, seed.module, seed.priority
FROM (VALUES
    ('Xem phiên bản catalog', 'CATALOG_PLAN_VIEW', 'Đọc plan version và điều khoản thương mại', 'ORGANIZATION', 'commercial', 10),
    ('Xuất bản phiên bản catalog', 'CATALOG_PLAN_PUBLISH', 'Xuất bản draft plan version', 'ORGANIZATION', 'commercial', 20),
    ('Xem subscription thương mại', 'COMMERCIAL_SUBSCRIPTION_VIEW', 'Đọc subscription và commercial projection của người dùng', 'ORGANIZATION', 'commercial', 10),
    ('Xem usage thương mại', 'COMMERCIAL_USAGE_VIEW', 'Đọc durable/runtime usage projection của TQD', 'ORGANIZATION', 'commercial', 10),
    ('Reconcile usage thương mại', 'COMMERCIAL_USAGE_RECONCILE', 'Yêu cầu TQD đồng bộ runtime projection từ durable usage', 'ORGANIZATION', 'commercial', 20),
    ('Xem payment thương mại', 'PAYMENT_ORDER_VIEW', 'Đọc order, payment attempt và settlement evidence', 'ORGANIZATION', 'commercial', 10),
    ('Xem fulfillment thương mại', 'PAYMENT_FULFILLMENT_VIEW', 'Đọc trạng thái cấp quyền và hàng đợi fulfillment', 'ORGANIZATION', 'commercial', 10),
    ('Redrive fulfillment thương mại', 'PAYMENT_FULFILLMENT_REDRIVE', 'Đưa fulfillment cần review trở lại hàng đợi xử lý', 'ORGANIZATION', 'commercial', 30)
) AS seed(name, key, description, permission_type, module, priority)
WHERE NOT EXISTS (SELECT 1 FROM permissions existing WHERE existing.key = seed.key);

COMMIT;
