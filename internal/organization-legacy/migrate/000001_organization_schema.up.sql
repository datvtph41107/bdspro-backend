-- Baseline cho fresh local database. Các table/column khớp registry GORM
-- đang được repository production sử dụng; serving process không chạy
-- file này, golang-migrate/Compose là SQL owner.

CREATE TABLE IF NOT EXISTS colors (
    created_by BIGINT, updated_by BIGINT, id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL, content_color TEXT NOT NULL, background_color TEXT NOT NULL,
    description TEXT, is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS business_domains (
    id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL, description TEXT,
    code TEXT NOT NULL, is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    created_by BIGINT, updated_by BIGINT, deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_business_domains_code ON business_domains(code);

CREATE TABLE IF NOT EXISTS organizations (
    created_by BIGINT, updated_by BIGINT, id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL, tax_code TEXT NOT NULL, business_license_url TEXT,
    address TEXT, detail_address TEXT, phone TEXT, email TEXT, website TEXT,
    logo_url TEXT, status BIGINT, approve_investor BOOLEAN DEFAULT FALSE,
    account_bank_investor BIGINT, founded_at TIMESTAMPTZ, description TEXT,
    owner_id BIGINT
);

CREATE TABLE IF NOT EXISTS organization_business_domains (
    id BIGSERIAL PRIMARY KEY, organization_id BIGINT NOT NULL,
    business_domain_id BIGINT NOT NULL, created_at TIMESTAMPTZ, created_by BIGINT
);

CREATE TABLE IF NOT EXISTS organization_permissions (
    id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL, key TEXT NOT NULL,
    description TEXT, parent_id BIGINT, permission_type BIGINT NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS organization_roles (
    id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL, key TEXT NOT NULL,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    organization_id BIGINT NOT NULL, is_default BOOLEAN DEFAULT FALSE,
    domain_type BIGINT DEFAULT 0, deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE, color_id BIGINT REFERENCES colors(id),
    role_key BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS role_permissions (
    organization_role_model_id BIGINT NOT NULL REFERENCES organization_roles(id),
    organization_permission_model_id BIGINT NOT NULL REFERENCES organization_permissions(id),
    PRIMARY KEY (organization_role_model_id, organization_permission_model_id)
);

CREATE TABLE IF NOT EXISTS organization_members (
    id BIGSERIAL PRIMARY KEY, organization_id BIGINT NOT NULL REFERENCES organizations(id),
    user_id BIGINT NOT NULL, role BIGINT NOT NULL, role_id BIGINT,
    role_key BIGINT, status BIGINT DEFAULT 1, joined_at TIMESTAMPTZ,
    removed_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS organization_branches (
    id BIGSERIAL PRIMARY KEY, organization_id BIGINT NOT NULL, name TEXT,
    address TEXT, phone TEXT, email TEXT, manager_id BIGINT,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, created_by BIGINT,
    deleted_at TIMESTAMPTZ, is_deleted BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE, type TEXT, description TEXT
);

CREATE TABLE IF NOT EXISTS organization_branch_members (
    id BIGSERIAL PRIMARY KEY, organization_branch_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    created_by BIGINT, role_id BIGINT, deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS organization_log_activities (
    id BIGSERIAL PRIMARY KEY, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    created_by BIGINT NOT NULL, updated_by BIGINT NOT NULL,
    organization_id BIGINT NOT NULL, actor_id BIGINT NOT NULL,
    log_type TEXT NOT NULL, log_data TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS groups (
    id BIGSERIAL PRIMARY KEY, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    created_by BIGINT NOT NULL, updated_by BIGINT NOT NULL,
    name TEXT NOT NULL, description TEXT, avatar_url TEXT,
    status TEXT DEFAULT 'active', deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS group_settings (
    id BIGSERIAL PRIMARY KEY, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    created_by BIGINT NOT NULL, updated_by BIGINT NOT NULL,
    group_id BIGINT NOT NULL, config_key TEXT NOT NULL, config_value TEXT NOT NULL,
    deleted_at TIMESTAMPTZ, is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS group_members (
    id BIGSERIAL PRIMARY KEY, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    created_by BIGINT NOT NULL, updated_by BIGINT NOT NULL,
    group_id BIGINT NOT NULL, user_id BIGINT NOT NULL,
    role TEXT NOT NULL, status TEXT NOT NULL, deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE, role_id BIGINT
);

CREATE TABLE IF NOT EXISTS group_log_activities (
    id BIGSERIAL PRIMARY KEY, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    created_by BIGINT NOT NULL, updated_by BIGINT NOT NULL,
    group_id BIGINT NOT NULL, actor_id BIGINT NOT NULL,
    log_type TEXT NOT NULL, log_data TEXT NOT NULL,
    deleted_at TIMESTAMPTZ, is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS group_notifications (
    id BIGSERIAL PRIMARY KEY, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    created_by BIGINT NOT NULL, updated_by BIGINT NOT NULL,
    group_id BIGINT NOT NULL, type TEXT NOT NULL, content TEXT NOT NULL,
    deleted_at TIMESTAMPTZ, is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS group_documents (
    id BIGSERIAL PRIMARY KEY, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    created_by BIGINT NOT NULL, updated_by BIGINT NOT NULL,
    group_id BIGINT NOT NULL, name TEXT NOT NULL, description TEXT,
    file_url TEXT NOT NULL, file_type TEXT NOT NULL, file_size BIGINT NOT NULL,
    deleted_at TIMESTAMPTZ, is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS group_chats (
    id BIGSERIAL PRIMARY KEY, group_id BIGINT NOT NULL, conversation_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, created_by BIGINT,
    deleted_at TIMESTAMPTZ, is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS bank_accounts (
    created_by BIGINT, updated_by BIGINT, id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    owner_id BIGINT NOT NULL, owner_type BIGINT NOT NULL DEFAULT 10,
    bank_name TEXT NOT NULL, bank_id BIGINT, account_number TEXT NOT NULL,
    account_name TEXT NOT NULL, account_type BIGINT,
    account_status BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS deals (
    created_by BIGINT, updated_by BIGINT, id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    owner_id BIGINT, owner_type BIGINT, name TEXT, target_profit DECIMAL DEFAULT 0,
    status BIGINT DEFAULT 10, deal_type BIGINT DEFAULT 10, note TEXT,
    charge_person_id BIGINT, cancel_reason TEXT, is_unilateral BOOLEAN DEFAULT FALSE,
    is_manual BOOLEAN DEFAULT FALSE, bank_account_id BIGINT REFERENCES bank_accounts(id),
    from_date TIMESTAMPTZ, to_date TIMESTAMPTZ, allow_sharing BOOLEAN,
    member_can_add_transaction BOOLEAN, only_owner_get_commission BOOLEAN,
    internal_note TEXT, allow_manual_input BOOLEAN
);

CREATE TABLE IF NOT EXISTS deal_products (
    deal_id BIGINT NOT NULL REFERENCES deals(id), product_id BIGINT NOT NULL,
    PRIMARY KEY (deal_id, product_id)
);

CREATE TABLE IF NOT EXISTS deal_members (
    created_by BIGINT, updated_by BIGINT, id BIGSERIAL,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    deal_id BIGINT NOT NULL REFERENCES deals(id), member_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL, amount_commit DECIMAL DEFAULT 0,
    commission_value DECIMAL DEFAULT 0, commission_type BIGINT DEFAULT 10,
    note TEXT, is_unilateral BOOLEAN DEFAULT FALSE, done_investment BOOLEAN DEFAULT FALSE,
    role_key BIGINT NOT NULL DEFAULT 430, status BIGINT NOT NULL DEFAULT 10,
    message TEXT, invited_at TIMESTAMPTZ, responded_at TIMESTAMPTZ,
    withdrawn_at TIMESTAMPTZ, inviter_id BIGINT, color_id BIGINT REFERENCES colors(id),
    is_owner BOOLEAN DEFAULT FALSE, PRIMARY KEY (id, deal_id, member_id)
);

CREATE TABLE IF NOT EXISTS deal_milestones (
    created_by BIGINT, updated_by BIGINT, id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    deal_id BIGINT NOT NULL REFERENCES deals(id), title TEXT NOT NULL,
    description TEXT, expected_date TIMESTAMPTZ, completed_date TIMESTAMPTZ,
    status BIGINT NOT NULL DEFAULT 10, order_index BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS investments (
    created_by BIGINT, updated_by BIGINT, id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    deal_id BIGINT NOT NULL, member_id BIGINT, invest_type BIGINT NOT NULL DEFAULT 10,
    amount DECIMAL NOT NULL, transfer_time TIMESTAMPTZ, note TEXT NOT NULL,
    transfer_proof_image TEXT, status BIGINT NOT NULL DEFAULT 10,
    changed_at TIMESTAMPTZ, confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    proxy_user_id BIGINT
);

CREATE TABLE IF NOT EXISTS attach_documents (
    created_by BIGINT, updated_by BIGINT, id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
    owner_id BIGINT, doc_owner BIGINT, doc_name TEXT, doc_path TEXT,
    doc_type TEXT, doc_size BIGINT
);

CREATE TABLE IF NOT EXISTS internal_notes (
    id BIGSERIAL PRIMARY KEY, deal_id BIGINT NOT NULL REFERENCES deals(id),
    actor_id BIGINT NOT NULL, content TEXT NOT NULL, action_type BIGINT NOT NULL,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS deal_of_organization (
    id BIGSERIAL PRIMARY KEY, deal_id BIGINT NOT NULL,
    organization_id BIGINT NOT NULL, branch_id BIGINT NOT NULL
);
CREATE TABLE IF NOT EXISTS deal_of_branch (
    id BIGSERIAL PRIMARY KEY, deal_id BIGINT NOT NULL, branch_id BIGINT NOT NULL
);
CREATE TABLE IF NOT EXISTS deal_of_group (
    id BIGSERIAL PRIMARY KEY, deal_id BIGINT NOT NULL, group_id BIGINT NOT NULL
);

-- Compatibility tables còn được investment summary production truy vấn.
CREATE TABLE IF NOT EXISTS cost_types (
    id BIGSERIAL PRIMARY KEY, name TEXT, cost_type BIGINT NOT NULL,
    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS deal_costs (
    id BIGSERIAL PRIMARY KEY, deal_id BIGINT NOT NULL, cost_type_id BIGINT NOT NULL,
    amount DECIMAL NOT NULL DEFAULT 0, created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_org_members_user_org ON organization_members(user_id, organization_id);
CREATE INDEX IF NOT EXISTS idx_org_members_status_org ON organization_members(status, organization_id);
CREATE INDEX IF NOT EXISTS idx_org_roles_organization_id ON organization_roles(organization_id);
CREATE INDEX IF NOT EXISTS idx_org_roles_key ON organization_roles(key);
CREATE INDEX IF NOT EXISTS idx_org_permissions_key ON organization_permissions(key);
CREATE INDEX IF NOT EXISTS idx_org_branches_org_id ON organization_branches(organization_id);
CREATE INDEX IF NOT EXISTS idx_org_branch_members_user_branch ON organization_branch_members(user_id, organization_branch_id);
CREATE INDEX IF NOT EXISTS idx_group_members_group_user ON group_members(group_id, user_id);
CREATE INDEX IF NOT EXISTS idx_deal_members_member_id ON deal_members(member_id);
CREATE INDEX IF NOT EXISTS idx_deal_milestones_deal_id ON deal_milestones(deal_id);
CREATE INDEX IF NOT EXISTS idx_investments_deal_id ON investments(deal_id);
CREATE INDEX IF NOT EXISTS idx_internal_notes_deal_id ON internal_notes(deal_id);
CREATE INDEX IF NOT EXISTS idx_deal_costs_deal_id ON deal_costs(deal_id);
