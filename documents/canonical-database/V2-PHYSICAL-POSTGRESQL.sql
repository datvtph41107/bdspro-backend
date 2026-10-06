-- BDSPro Canonical Database — v2 Physical PostgreSQL Reference
--
-- IMPORTANT:
--   1. This is a TARGET reference, not a migration to run wholesale against production.
--   2. Each section belongs to a separate service database.
--   3. Cross-service logical references intentionally have NO PostgreSQL FK.
--   4. Production adoption follows EXPAND → BACKFILL → COMPARE → SWITCH → CONTRACT.

-- =============================================================================
-- USER-SERVICE DATABASE
-- Identity + Organization/IAM authority
-- =============================================================================

CREATE TABLE parties (
    id BIGSERIAL PRIMARY KEY
);

CREATE TABLE persons (
    party_id BIGINT PRIMARY KEY
        REFERENCES parties(id)
);

CREATE TABLE profiles (
    person_id BIGINT PRIMARY KEY
        REFERENCES persons(party_id),
    display_name VARCHAR(255),
    avatar_url TEXT,
    bio TEXT,
    email VARCHAR(320),
    phone VARCHAR(32),
    website_url TEXT
);

CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    person_id BIGINT NOT NULL
        REFERENCES persons(party_id),
    closed_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ux_accounts_one_open_per_person
    ON accounts(person_id)
    WHERE closed_at IS NULL;

CREATE TABLE account_password_credentials (
    account_id BIGINT PRIMARY KEY
        REFERENCES accounts(id),
    login_name VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX ux_account_password_login_name
    ON account_password_credentials(lower(login_name));

CREATE TABLE account_federated_identities (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL
        REFERENCES accounts(id),
    provider VARCHAR(64) NOT NULL,
    provider_subject VARCHAR(512) NOT NULL,
    CONSTRAINT uq_federated_provider_subject
        UNIQUE(provider, provider_subject)
);

CREATE TABLE account_passkeys (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL
        REFERENCES accounts(id),
    credential_id TEXT NOT NULL UNIQUE,
    public_key TEXT NOT NULL,
    sign_count BIGINT NOT NULL DEFAULT 0
        CHECK (sign_count >= 0),
    created_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    CHECK (revoked_at IS NULL OR revoked_at >= created_at)
);

CREATE TABLE account_sessions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL
        REFERENCES accounts(id),
    client_id VARCHAR(128),
    device_key VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    CHECK (last_seen_at IS NULL OR last_seen_at >= created_at),
    CHECK (ended_at IS NULL OR ended_at >= created_at)
);

CREATE INDEX ix_account_sessions_current
    ON account_sessions(account_id, last_seen_at DESC)
    WHERE ended_at IS NULL;

CREATE TABLE account_blocks (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL
        REFERENCES accounts(id),
    blocked_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    CHECK (expires_at IS NULL OR expires_at > blocked_at),
    CHECK (ended_at IS NULL OR ended_at >= blocked_at)
);

CREATE UNIQUE INDEX ux_account_blocks_one_current
    ON account_blocks(account_id)
    WHERE ended_at IS NULL;

CREATE TABLE permissions (
    code VARCHAR(128) PRIMARY KEY
        CHECK (btrim(code) <> '')
);

-- -----------------------------------------------------------------------------
-- Organization
-- -----------------------------------------------------------------------------

CREATE TABLE organizations (
    party_id BIGINT PRIMARY KEY
        REFERENCES parties(id),
    name VARCHAR(255) NOT NULL
        CHECK (btrim(name) <> ''),
    tax_code VARCHAR(64),
    phone VARCHAR(32),
    email VARCHAR(320),
    website TEXT,
    description TEXT,
    archived_at TIMESTAMPTZ
);

-- Tax-code input MUST be canonicalized by the owning write path before storage.
-- Archive does not release uniqueness.
CREATE UNIQUE INDEX ux_organizations_tax_code
    ON organizations(tax_code)
    WHERE tax_code IS NOT NULL AND btrim(tax_code) <> '';

CREATE TABLE organization_roles (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL
        REFERENCES organizations(party_id),
    name VARCHAR(255) NOT NULL
        CHECK (btrim(name) <> '')
);

CREATE TABLE organization_role_permissions (
    role_id BIGINT NOT NULL
        REFERENCES organization_roles(id),
    permission_code VARCHAR(128) NOT NULL
        REFERENCES permissions(code),
    PRIMARY KEY(role_id, permission_code)
);

CREATE TABLE organization_invitations (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL
        REFERENCES organizations(party_id),

    invitee_person_id BIGINT
        REFERENCES persons(party_id),
    invitee_email VARCHAR(320),
    invitee_phone VARCHAR(32),

    intended_role_id BIGINT NOT NULL
        REFERENCES organization_roles(id),

    invited_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,

    outcome VARCHAR(16),
    resolved_at TIMESTAMPTZ,

    CHECK (
        ((invitee_person_id IS NOT NULL)::int +
         (invitee_email IS NOT NULL)::int +
         (invitee_phone IS NOT NULL)::int) = 1
    ),
    CHECK (expires_at > invited_at),
    CHECK (outcome IS NULL OR outcome IN ('ACCEPTED', 'DECLINED', 'REVOKED')),
    CHECK (
        (outcome IS NULL AND resolved_at IS NULL)
        OR
        (outcome IS NOT NULL AND resolved_at IS NOT NULL AND resolved_at >= invited_at)
    )
);

CREATE INDEX ix_org_invitations_org_unresolved
    ON organization_invitations(organization_id, expires_at)
    WHERE outcome IS NULL;

CREATE TABLE organization_memberships (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL
        REFERENCES organizations(party_id),
    person_id BIGINT NOT NULL
        REFERENCES persons(party_id),
    originating_invitation_id BIGINT UNIQUE
        REFERENCES organization_invitations(id),
    joined_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    CHECK (ended_at IS NULL OR ended_at >= joined_at)
);

CREATE UNIQUE INDEX ux_org_memberships_one_current
    ON organization_memberships(organization_id, person_id)
    WHERE ended_at IS NULL;

CREATE INDEX ix_org_memberships_current_org
    ON organization_memberships(organization_id, id)
    WHERE ended_at IS NULL;

CREATE TABLE organization_membership_blocks (
    id BIGSERIAL PRIMARY KEY,
    membership_id BIGINT NOT NULL
        REFERENCES organization_memberships(id),
    blocked_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    CHECK (ended_at IS NULL OR ended_at >= blocked_at)
);

CREATE UNIQUE INDEX ux_org_membership_blocks_one_current
    ON organization_membership_blocks(membership_id)
    WHERE ended_at IS NULL;

CREATE TABLE organization_role_assignments (
    id BIGSERIAL PRIMARY KEY,
    membership_id BIGINT NOT NULL
        REFERENCES organization_memberships(id),
    role_id BIGINT NOT NULL
        REFERENCES organization_roles(id),
    assigned_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    CHECK (ended_at IS NULL OR ended_at >= assigned_at)
);

CREATE UNIQUE INDEX ux_org_role_assignment_one_current
    ON organization_role_assignments(membership_id)
    WHERE ended_at IS NULL;

CREATE TABLE organization_ownerships (
    organization_id BIGINT PRIMARY KEY
        REFERENCES organizations(party_id),
    membership_id BIGINT NOT NULL UNIQUE
        REFERENCES organization_memberships(id)
);

CREATE TABLE organization_branches (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL
        REFERENCES organizations(party_id),
    name VARCHAR(255) NOT NULL
        CHECK (btrim(name) <> ''),
    manager_membership_id BIGINT
        REFERENCES organization_memberships(id),
    retired_at TIMESTAMPTZ
);

CREATE INDEX ix_organization_branches_current
    ON organization_branches(organization_id, id)
    WHERE retired_at IS NULL;

CREATE TABLE organization_branch_assignments (
    branch_id BIGINT NOT NULL
        REFERENCES organization_branches(id),
    membership_id BIGINT NOT NULL
        REFERENCES organization_memberships(id),
    PRIMARY KEY(branch_id, membership_id)
);

-- Same-Organization/context checks cannot be expressed with an ordinary CHECK.
-- Constraint triggers make the relationship explicit at the DB boundary.

CREATE OR REPLACE FUNCTION assert_org_role_assignment_context()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    membership_org BIGINT;
    role_org BIGINT;
BEGIN
    SELECT organization_id INTO membership_org
    FROM organization_memberships
    WHERE id = NEW.membership_id AND ended_at IS NULL;

    SELECT organization_id INTO role_org
    FROM organization_roles
    WHERE id = NEW.role_id;

    IF membership_org IS NULL THEN
        RAISE EXCEPTION 'role assignment requires current membership %', NEW.membership_id;
    END IF;

    IF role_org IS NULL OR role_org <> membership_org THEN
        RAISE EXCEPTION 'role % and membership % belong to different organizations',
            NEW.role_id, NEW.membership_id;
    END IF;

    RETURN NEW;
END $$;

CREATE CONSTRAINT TRIGGER ct_org_role_assignment_context
AFTER INSERT OR UPDATE OF membership_id, role_id, ended_at
ON organization_role_assignments
DEFERRABLE INITIALLY IMMEDIATE
FOR EACH ROW
WHEN (NEW.ended_at IS NULL)
EXECUTE FUNCTION assert_org_role_assignment_context();

CREATE OR REPLACE FUNCTION assert_org_ownership_context()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    membership_org BIGINT;
    is_blocked BOOLEAN;
    org_archived TIMESTAMPTZ;
BEGIN
    SELECT organization_id INTO membership_org
    FROM organization_memberships
    WHERE id = NEW.membership_id AND ended_at IS NULL;

    IF membership_org IS NULL OR membership_org <> NEW.organization_id THEN
        RAISE EXCEPTION 'ownership requires current same-organization membership';
    END IF;

    SELECT EXISTS (
        SELECT 1
        FROM organization_membership_blocks
        WHERE membership_id = NEW.membership_id
          AND ended_at IS NULL
    ) INTO is_blocked;

    IF is_blocked THEN
        RAISE EXCEPTION 'blocked membership cannot be current organization owner';
    END IF;

    SELECT archived_at INTO org_archived
    FROM organizations
    WHERE party_id = NEW.organization_id;

    IF org_archived IS NOT NULL THEN
        RAISE EXCEPTION 'archived organization cannot have current ownership';
    END IF;

    RETURN NEW;
END $$;

CREATE CONSTRAINT TRIGGER ct_org_ownership_context
AFTER INSERT OR UPDATE
ON organization_ownerships
DEFERRABLE INITIALLY IMMEDIATE
FOR EACH ROW
EXECUTE FUNCTION assert_org_ownership_context();

CREATE OR REPLACE FUNCTION assert_branch_manager_context()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    membership_org BIGINT;
BEGIN
    IF NEW.manager_membership_id IS NULL THEN
        RETURN NEW;
    END IF;

    SELECT organization_id INTO membership_org
    FROM organization_memberships
    WHERE id = NEW.manager_membership_id
      AND ended_at IS NULL;

    IF membership_org IS NULL OR membership_org <> NEW.organization_id THEN
        RAISE EXCEPTION 'branch manager must be current membership of same organization';
    END IF;

    RETURN NEW;
END $$;

CREATE CONSTRAINT TRIGGER ct_branch_manager_context
AFTER INSERT OR UPDATE OF organization_id, manager_membership_id
ON organization_branches
DEFERRABLE INITIALLY IMMEDIATE
FOR EACH ROW
EXECUTE FUNCTION assert_branch_manager_context();

CREATE OR REPLACE FUNCTION assert_branch_assignment_context()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    branch_org BIGINT;
    membership_org BIGINT;
BEGIN
    SELECT organization_id INTO branch_org
    FROM organization_branches
    WHERE id = NEW.branch_id AND retired_at IS NULL;

    SELECT organization_id INTO membership_org
    FROM organization_memberships
    WHERE id = NEW.membership_id AND ended_at IS NULL;

    IF branch_org IS NULL THEN
        RAISE EXCEPTION 'branch assignment requires current branch';
    END IF;

    IF membership_org IS NULL OR membership_org <> branch_org THEN
        RAISE EXCEPTION 'branch assignment requires current same-organization membership';
    END IF;

    RETURN NEW;
END $$;

CREATE CONSTRAINT TRIGGER ct_branch_assignment_context
AFTER INSERT OR UPDATE
ON organization_branch_assignments
DEFERRABLE INITIALLY IMMEDIATE
FOR EACH ROW
EXECUTE FUNCTION assert_branch_assignment_context();

-- IMPORTANT invariant not fully representable by declarative DDL:
-- every non-archived Organization must have exactly one usable ownership row.
-- All Organization-governance write paths MUST lock organizations.party_id first,
-- re-read Ownership/Membership/Block state, then transition atomically.


-- =============================================================================
-- ORGANIZATION-SERVICE DATABASE
-- Group + Deal bounded context
-- Person/Organization/Branch/File/BankAccount/Product IDs below are external
-- authoritative references and intentionally have no PostgreSQL FK here.
-- =============================================================================

CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
        CHECK (btrim(name) <> ''),
    description TEXT,
    avatar_url TEXT,
    disbanded_at TIMESTAMPTZ
);

CREATE TABLE group_memberships (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL
        REFERENCES groups(id),
    person_id BIGINT NOT NULL, -- external Person reference
    joined_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    CHECK (ended_at IS NULL OR ended_at >= joined_at)
);

CREATE UNIQUE INDEX ux_group_memberships_one_current
    ON group_memberships(group_id, person_id)
    WHERE ended_at IS NULL;

CREATE TABLE group_leaderships (
    group_id BIGINT PRIMARY KEY
        REFERENCES groups(id),
    membership_id BIGINT NOT NULL UNIQUE
        REFERENCES group_memberships(id)
);

CREATE OR REPLACE FUNCTION assert_group_leadership_context()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    member_group BIGINT;
    group_end TIMESTAMPTZ;
BEGIN
    SELECT group_id INTO member_group
    FROM group_memberships
    WHERE id = NEW.membership_id AND ended_at IS NULL;

    IF member_group IS NULL OR member_group <> NEW.group_id THEN
        RAISE EXCEPTION 'group leader must be current membership of same group';
    END IF;

    SELECT disbanded_at INTO group_end
    FROM groups
    WHERE id = NEW.group_id;

    IF group_end IS NOT NULL THEN
        RAISE EXCEPTION 'disbanded group cannot have current leader';
    END IF;

    RETURN NEW;
END $$;

CREATE CONSTRAINT TRIGGER ct_group_leadership_context
AFTER INSERT OR UPDATE
ON group_leaderships
DEFERRABLE INITIALLY IMMEDIATE
FOR EACH ROW
EXECUTE FUNCTION assert_group_leadership_context();

CREATE TABLE group_documents (
    group_id BIGINT NOT NULL
        REFERENCES groups(id),
    file_id BIGINT NOT NULL, -- external file-service reference
    PRIMARY KEY(group_id, file_id)
);

CREATE TABLE group_chat_links (
    group_id BIGINT PRIMARY KEY
        REFERENCES groups(id),
    conversation_id BIGINT NOT NULL UNIQUE -- external chat reference
);

CREATE TABLE deals (
    id BIGSERIAL PRIMARY KEY,

    organization_id BIGINT, -- external user-service Organization reference
    group_id BIGINT
        REFERENCES groups(id),

    branch_id BIGINT, -- external user-service Branch reference

    name VARCHAR(255) NOT NULL
        CHECK (btrim(name) <> ''),

    kind VARCHAR(16) NOT NULL
        CHECK (kind IN ('BROKERAGE', 'INVESTMENT', 'JOINT')),

    target_profit_amount NUMERIC(20,2),
    currency_code CHAR(3),

    note TEXT,

    from_date TIMESTAMPTZ,
    to_date TIMESTAMPTZ,

    allow_sharing BOOLEAN NOT NULL DEFAULT FALSE,
    participant_can_add_transaction BOOLEAN NOT NULL DEFAULT FALSE,
    only_lead_get_commission BOOLEAN NOT NULL DEFAULT FALSE,
    allow_manual_input BOOLEAN NOT NULL DEFAULT FALSE,

    bank_account_id BIGINT, -- external payment-service BankAccount reference

    outcome VARCHAR(16),
    ended_at TIMESTAMPTZ,
    cancel_reason TEXT,

    CHECK (
        ((organization_id IS NOT NULL)::int +
         (group_id IS NOT NULL)::int) = 1
    ),
    CHECK (branch_id IS NULL OR organization_id IS NOT NULL),
    CHECK (
        (target_profit_amount IS NULL AND currency_code IS NULL)
        OR
        (target_profit_amount IS NOT NULL AND target_profit_amount >= 0 AND currency_code IS NOT NULL)
    ),
    CHECK (from_date IS NULL OR to_date IS NULL OR to_date >= from_date),
    CHECK (outcome IS NULL OR outcome IN ('COMPLETED', 'CANCELED')),
    CHECK (
        (outcome IS NULL AND ended_at IS NULL)
        OR
        (outcome IS NOT NULL AND ended_at IS NOT NULL)
    ),
    CHECK (
        outcome = 'CANCELED'
        OR cancel_reason IS NULL
    )
);

CREATE INDEX ix_deals_org_current
    ON deals(organization_id, id)
    WHERE organization_id IS NOT NULL AND outcome IS NULL;

CREATE INDEX ix_deals_group_current
    ON deals(group_id, id)
    WHERE group_id IS NOT NULL AND outcome IS NULL;

CREATE INDEX ix_deals_branch
    ON deals(branch_id)
    WHERE branch_id IS NOT NULL;

-- Cross-service invariant:
-- if branch_id is present, user-service must confirm:
-- branch.organization_id = deals.organization_id and branch.retired_at IS NULL.

CREATE TABLE deal_invitations (
    id BIGSERIAL PRIMARY KEY,
    deal_id BIGINT NOT NULL
        REFERENCES deals(id),
    invitee_person_id BIGINT NOT NULL, -- external Person reference

    as_customer BOOLEAN NOT NULL DEFAULT FALSE,
    as_partner BOOLEAN NOT NULL DEFAULT FALSE,
    as_administrator BOOLEAN NOT NULL DEFAULT FALSE,

    message TEXT,
    invited_at TIMESTAMPTZ NOT NULL,

    outcome VARCHAR(16),
    resolved_at TIMESTAMPTZ,

    CHECK (outcome IS NULL OR outcome IN ('ACCEPTED', 'DECLINED', 'REVOKED')),
    CHECK (
        (outcome IS NULL AND resolved_at IS NULL)
        OR
        (outcome IS NOT NULL AND resolved_at IS NOT NULL AND resolved_at >= invited_at)
    )
);

CREATE UNIQUE INDEX ux_deal_invitation_one_unresolved
    ON deal_invitations(deal_id, invitee_person_id)
    WHERE outcome IS NULL;

CREATE TABLE deal_participations (
    id BIGSERIAL PRIMARY KEY,
    deal_id BIGINT NOT NULL
        REFERENCES deals(id),
    person_id BIGINT NOT NULL, -- external Person reference
    originating_invitation_id BIGINT UNIQUE
        REFERENCES deal_invitations(id),
    joined_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    CHECK (ended_at IS NULL OR ended_at >= joined_at)
);

CREATE UNIQUE INDEX ux_deal_participations_one_current
    ON deal_participations(deal_id, person_id)
    WHERE ended_at IS NULL;

CREATE TABLE deal_customer_relationships (
    participation_id BIGINT PRIMARY KEY
        REFERENCES deal_participations(id)
);

CREATE TABLE deal_partner_relationships (
    participation_id BIGINT PRIMARY KEY
        REFERENCES deal_participations(id)
);

CREATE TABLE deal_administrator_assignments (
    participation_id BIGINT PRIMARY KEY
        REFERENCES deal_participations(id)
);

CREATE TABLE deal_leaderships (
    deal_id BIGINT PRIMARY KEY
        REFERENCES deals(id),
    participation_id BIGINT NOT NULL UNIQUE
        REFERENCES deal_participations(id)
);

CREATE OR REPLACE FUNCTION assert_deal_leadership_context()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    participation_deal BIGINT;
    deal_end TIMESTAMPTZ;
BEGIN
    SELECT deal_id INTO participation_deal
    FROM deal_participations
    WHERE id = NEW.participation_id
      AND ended_at IS NULL;

    IF participation_deal IS NULL OR participation_deal <> NEW.deal_id THEN
        RAISE EXCEPTION 'deal lead must be current participation of same deal';
    END IF;

    SELECT ended_at INTO deal_end
    FROM deals
    WHERE id = NEW.deal_id;

    IF deal_end IS NOT NULL THEN
        RAISE EXCEPTION 'terminal deal cannot have current lead authority';
    END IF;

    RETURN NEW;
END $$;

CREATE CONSTRAINT TRIGGER ct_deal_leadership_context
AFTER INSERT OR UPDATE
ON deal_leaderships
DEFERRABLE INITIALLY IMMEDIATE
FOR EACH ROW
EXECUTE FUNCTION assert_deal_leadership_context();

CREATE OR REPLACE FUNCTION assert_deal_admin_current_participation()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM deal_participations
        WHERE id = NEW.participation_id
          AND ended_at IS NULL
    ) THEN
        RAISE EXCEPTION 'deal admin must be current participation';
    END IF;
    RETURN NEW;
END $$;

CREATE CONSTRAINT TRIGGER ct_deal_admin_current_participation
AFTER INSERT OR UPDATE
ON deal_administrator_assignments
DEFERRABLE INITIALLY IMMEDIATE
FOR EACH ROW
EXECUTE FUNCTION assert_deal_admin_current_participation();

CREATE TABLE deal_participant_commissions (
    participation_id BIGINT PRIMARY KEY
        REFERENCES deal_participations(id),
    kind VARCHAR(16) NOT NULL
        CHECK (kind IN ('PERCENT', 'FIXED_AMOUNT')),
    value NUMERIC(20,6) NOT NULL
        CHECK (value >= 0),
    currency_code CHAR(3),
    note TEXT,
    CHECK (
        (kind = 'PERCENT' AND currency_code IS NULL AND value <= 100)
        OR
        (kind = 'FIXED_AMOUNT' AND currency_code IS NOT NULL)
    )
);

CREATE TABLE deal_investment_commitments (
    participation_id BIGINT PRIMARY KEY
        REFERENCES deal_participations(id),
    amount NUMERIC(20,2) NOT NULL
        CHECK (amount >= 0),
    currency_code CHAR(3) NOT NULL
);

CREATE TABLE deal_investments (
    id BIGSERIAL PRIMARY KEY,
    participation_id BIGINT NOT NULL
        REFERENCES deal_participations(id),

    amount NUMERIC(20,2) NOT NULL
        CHECK (amount > 0),
    currency_code CHAR(3) NOT NULL,

    transferred_at TIMESTAMPTZ,
    transfer_proof_file_id BIGINT, -- external file-service reference
    note TEXT,

    outcome VARCHAR(16),
    resolved_at TIMESTAMPTZ,

    CHECK (outcome IS NULL OR outcome IN ('APPROVED', 'REJECTED')),
    CHECK (
        (outcome IS NULL AND resolved_at IS NULL)
        OR
        (outcome IS NOT NULL AND resolved_at IS NOT NULL)
    )
);

CREATE INDEX ix_deal_investments_participation
    ON deal_investments(participation_id, id);

CREATE TABLE deal_products (
    deal_id BIGINT NOT NULL
        REFERENCES deals(id),
    product_id BIGINT NOT NULL, -- external bdspro-service reference
    PRIMARY KEY(deal_id, product_id)
);

CREATE TABLE deal_milestones (
    id BIGSERIAL PRIMARY KEY,
    deal_id BIGINT NOT NULL
        REFERENCES deals(id),
    title VARCHAR(255) NOT NULL
        CHECK (btrim(title) <> ''),
    description TEXT,
    expected_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    position INTEGER NOT NULL DEFAULT 0
        CHECK (position >= 0)
);

CREATE INDEX ix_deal_milestones_order
    ON deal_milestones(deal_id, position, id);

CREATE TABLE deal_cost_types (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) NOT NULL UNIQUE
        CHECK (btrim(code) <> ''),
    name VARCHAR(255) NOT NULL
        CHECK (btrim(name) <> '')
);

CREATE TABLE deal_costs (
    id BIGSERIAL PRIMARY KEY,
    deal_id BIGINT NOT NULL
        REFERENCES deals(id),
    cost_type_id BIGINT NOT NULL
        REFERENCES deal_cost_types(id),
    amount NUMERIC(20,2) NOT NULL
        CHECK (amount >= 0),
    currency_code CHAR(3) NOT NULL
);

CREATE INDEX ix_deal_costs_deal
    ON deal_costs(deal_id, cost_type_id, id);

CREATE TABLE deal_internal_notes (
    id BIGSERIAL PRIMARY KEY,
    deal_id BIGINT NOT NULL
        REFERENCES deals(id),
    author_person_id BIGINT, -- external Person reference
    content TEXT NOT NULL
        CHECK (btrim(content) <> ''),
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX ix_deal_internal_notes_deal
    ON deal_internal_notes(deal_id, created_at, id);

CREATE TABLE deal_documents (
    deal_id BIGINT NOT NULL
        REFERENCES deals(id),
    file_id BIGINT NOT NULL, -- external file-service reference
    PRIMARY KEY(deal_id, file_id)
);

CREATE TABLE deal_investment_documents (
    investment_id BIGINT NOT NULL
        REFERENCES deal_investments(id),
    file_id BIGINT NOT NULL, -- external file-service reference
    PRIMARY KEY(investment_id, file_id)
);

-- Required transactional invariants in organization-service:
--
-- Create Deal:
--   create Deal
--   create creator Participation
--   create DealLeadership
--   create Invitations for requested people/capacities
--   attach initial products/financial configuration
--   commit
--
-- Accept Deal Invitation:
--   lock invitation
--   validate caller Person == target
--   validate unresolved
--   validate no current participation
--   create Participation
--   materialize Customer/Partner/Admin terms
--   mark Invitation ACCEPTED at same effective time
--   commit
--
-- Transfer Deal Lead:
--   lock Deal
--   verify expected current Lead
--   verify target Participation current
--   update deal_leaderships
--   commit
--
-- End Participation:
--   lock Deal + Participation
--   if current Lead: reject unless lead transfer is in same transaction
--   remove current Admin assignment
--   set participation.ended_at
--   preserve Customer/Partner/Commission/Investment history
--   commit
--
-- Complete/Cancel Deal:
--   lock Deal
--   validate terminal preconditions
--   set outcome + ended_at (+ cancel_reason if canceled)
--   remove current Lead/Admin authority as policy requires
--   commit


-- =============================================================================
-- PAYMENT-SERVICE DATABASE
-- Existing commerce tables are already close to the target model.
-- This section defines the canonical identity corrections/additions.
-- =============================================================================

CREATE TABLE banks (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(192) NOT NULL,
    logo_url TEXT,
    retired_at TIMESTAMPTZ
);

CREATE TABLE bank_accounts (
    id BIGSERIAL PRIMARY KEY,
    party_id BIGINT NOT NULL, -- external user-service Party reference
    bank_id BIGINT NOT NULL
        REFERENCES banks(id),
    account_number VARCHAR(128) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    closed_at TIMESTAMPTZ,
    CONSTRAINT uq_bank_account_number
        UNIQUE(bank_id, account_number)
);

CREATE TABLE wallets (
    id BIGSERIAL PRIMARY KEY,
    person_id BIGINT NOT NULL UNIQUE, -- external user-service Person reference
    balance_minor BIGINT NOT NULL DEFAULT 0
        CHECK (balance_minor >= 0),
    currency_code VARCHAR(8) NOT NULL
        CHECK (btrim(currency_code) <> ''),
    closed_at TIMESTAMPTZ
);

CREATE TABLE wallet_transactions (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL
        REFERENCES wallets(id),
    transaction_code VARCHAR(192) NOT NULL UNIQUE,
    kind VARCHAR(32) NOT NULL
        CHECK (btrim(kind) <> ''),
    amount_minor BIGINT NOT NULL
        CHECK (amount_minor <> 0),
    status VARCHAR(32) NOT NULL
        CHECK (btrim(status) <> ''),
    external_reference VARCHAR(192),
    created_at TIMESTAMPTZ NOT NULL
);

-- commerce_orders target correction:
--   subject_kind + subject_id  → party_id (external Party reference)
-- Keep the established order/attempt/settlement/fulfillment/outbox separation.

CREATE TABLE commerce_orders (
    id BIGSERIAL PRIMARY KEY,

    party_id BIGINT NOT NULL, -- external Party reference

    product_code VARCHAR(128) NOT NULL
        CHECK (btrim(product_code) <> ''),
    plan_code VARCHAR(128) NOT NULL
        CHECK (btrim(plan_code) <> ''),
    plan_version_id BIGINT NOT NULL
        CHECK (plan_version_id > 0),

    currency_code VARCHAR(8) NOT NULL
        CHECK (btrim(currency_code) <> ''),
    amount_minor BIGINT NOT NULL
        CHECK (amount_minor > 0),

    command_key VARCHAR(128) NOT NULL
        CHECK (btrim(command_key) <> ''),
    command_fingerprint CHAR(64) NOT NULL,

    reference VARCHAR(64) NOT NULL UNIQUE,

    status VARCHAR(32) NOT NULL
        CHECK (status IN ('pending_funds', 'funds_confirmed', 'requires_review')),

    expires_at TIMESTAMPTZ NOT NULL,
    funds_confirmed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT uq_commerce_order_command
        UNIQUE(party_id, product_code, command_key)
);

CREATE TABLE commerce_payment_attempts (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL
        REFERENCES commerce_orders(id),
    command_key VARCHAR(128) NOT NULL,
    command_fingerprint CHAR(64) NOT NULL,
    method VARCHAR(32) NOT NULL
        CHECK (method IN ('card', 'bank_transfer', 'qr')),
    provider VARCHAR(64) NOT NULL,
    provider_reference VARCHAR(192) NOT NULL,
    status VARCHAR(32) NOT NULL
        CHECK (status IN (
            'created', 'pending_action', 'processing', 'authorized',
            'succeeded', 'declined', 'expired', 'canceled', 'failed'
        )),
    next_action_kind VARCHAR(32) NOT NULL
        CHECK (next_action_kind IN ('none', 'redirect', 'display_qr', 'wait')),
    redirect_url TEXT,
    qr_payload TEXT,
    action_expires_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    failure_code VARCHAR(96),
    failure_detail TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE(order_id, command_key),
    UNIQUE(provider, provider_reference)
);

CREATE TABLE commerce_settlements (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT
        REFERENCES commerce_orders(id),
    provider VARCHAR(64) NOT NULL,
    provider_transaction_id VARCHAR(192) NOT NULL,
    reference VARCHAR(128) NOT NULL,
    currency_code VARCHAR(8) NOT NULL,
    amount_minor BIGINT NOT NULL
        CHECK (amount_minor > 0),
    occurred_at TIMESTAMPTZ NOT NULL,
    evidence_hash CHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL
        CHECK (status IN ('funds_confirmed', 'unmatched', 'requires_review', 'late')),
    review_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE(provider, provider_transaction_id)
);

CREATE TABLE commerce_fulfillments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE
        REFERENCES commerce_orders(id),
    status VARCHAR(32) NOT NULL
        CHECK (status IN ('pending', 'running', 'completed', 'retry', 'requires_review')),
    attempt_count INTEGER NOT NULL DEFAULT 0
        CHECK (attempt_count >= 0),
    available_at TIMESTAMPTZ NOT NULL,
    locked_by VARCHAR(128),
    claim_version BIGINT NOT NULL DEFAULT 0
        CHECK (claim_version >= 0),
    lease_until TIMESTAMPTZ,
    last_error TEXT,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE commerce_command_effects (
    id BIGSERIAL PRIMARY KEY,
    effect_type VARCHAR(96) NOT NULL,
    scope_id BIGINT NOT NULL,
    command_key VARCHAR(128) NOT NULL,
    actor_reference VARCHAR(128) NOT NULL,
    reason TEXT,
    outcome VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    UNIQUE(effect_type, scope_id, command_key)
);

CREATE TABLE commerce_outbox_events (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(192) NOT NULL UNIQUE,
    event_type VARCHAR(128) NOT NULL,
    schema_version INTEGER NOT NULL CHECK (schema_version > 0),
    routing_key VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(32) NOT NULL
        CHECK (status IN ('pending', 'running', 'published', 'retry')),
    attempt_count INTEGER NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
    available_at TIMESTAMPTZ NOT NULL,
    locked_by VARCHAR(128),
    claim_version BIGINT NOT NULL DEFAULT 0 CHECK (claim_version >= 0),
    lease_until TIMESTAMPTZ,
    last_error TEXT,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- =============================================================================
-- CROSS-SERVICE ENFORCEMENT RULE
-- =============================================================================
--
-- A logical Ref in v1 becomes one of:
--
-- 1. physical FK
--      only when both facts are in the same authoritative database;
--
-- 2. write-time authoritative lookup + stored immutable external ID
--      when the target belongs to another service;
--
-- 3. local projection with source_version/event cursor
--      when high-volume reads require local data;
--
-- never:
--      runtime SQL join across service databases
--      or duplicated mutable truth with two write owners.
