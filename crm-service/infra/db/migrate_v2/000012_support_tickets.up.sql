CREATE TABLE IF NOT EXISTS support_tickets (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    issue_type INT NOT NULL DEFAULT 70,
    status INT NOT NULL DEFAULT 10,
    priority INT NOT NULL DEFAULT 20,
    product INT NOT NULL DEFAULT 10,
    source INT NOT NULL DEFAULT 10,
    assignee_id BIGINT,
    related_user_id BIGINT,
    related_report_id BIGINT,
    closed_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_support_tickets_code ON support_tickets(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_support_tickets_status ON support_tickets(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_support_tickets_priority ON support_tickets(priority) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_support_tickets_assignee ON support_tickets(assignee_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_support_tickets_related_report ON support_tickets(related_report_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS support_ticket_notes (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL REFERENCES support_tickets(id),
    author_id BIGINT NOT NULL,
    body TEXT NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_support_ticket_notes_ticket ON support_ticket_notes(ticket_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS support_ticket_events (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL REFERENCES support_tickets(id),
    actor_id BIGINT NOT NULL,
    action VARCHAR(64) NOT NULL,
    before_json TEXT,
    after_json TEXT,
    note TEXT,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_support_ticket_events_ticket ON support_ticket_events(ticket_id) WHERE deleted_at IS NULL;
