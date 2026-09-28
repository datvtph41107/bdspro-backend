-- Product-tab satellites for modules/bdspro (deals-by-product, notes, appointments, property history).
-- Idempotent CREATE IF NOT EXISTS for DBs that already absorbed bdspro/crm/notification dumps.

CREATE TABLE IF NOT EXISTS deal_products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deal_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_deal_products_deal_id ON deal_products (deal_id);
CREATE INDEX IF NOT EXISTS idx_deal_products_product_id ON deal_products (product_id);
CREATE INDEX IF NOT EXISTS idx_deal_products_deleted_at ON deal_products (deleted_at);

CREATE TABLE IF NOT EXISTS product_notes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    product_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    is_pinned BOOLEAN,
    pinned_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_product_notes_product_id ON product_notes (product_id);
CREATE INDEX IF NOT EXISTS idx_product_notes_author_id ON product_notes (author_id);
CREATE INDEX IF NOT EXISTS idx_product_notes_deleted_at ON product_notes (deleted_at);

CREATE TABLE IF NOT EXISTS product_note_mentions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    note_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    start_pos INTEGER,
    length INTEGER
);
CREATE INDEX IF NOT EXISTS idx_product_note_mentions_note_id ON product_note_mentions (note_id);
CREATE INDEX IF NOT EXISTS idx_product_note_mentions_deleted_at ON product_note_mentions (deleted_at);

CREATE TABLE IF NOT EXISTS product_note_files (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    note_id BIGINT NOT NULL,
    file_name VARCHAR(255),
    url TEXT NOT NULL,
    size BIGINT,
    file_type VARCHAR(50)
);
CREATE INDEX IF NOT EXISTS idx_product_note_files_note_id ON product_note_files (note_id);
CREATE INDEX IF NOT EXISTS idx_product_note_files_deleted_at ON product_note_files (deleted_at);

CREATE TABLE IF NOT EXISTS appointments (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    title TEXT,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    place TEXT,
    participant_ids INTEGER[],
    deal_id BIGINT,
    product_ids INTEGER[],
    description TEXT,
    reminder_schedule BIGINT,
    mode BIGINT,
    option_extend BIGINT,
    status INTEGER DEFAULT 10,
    transaction_ids INTEGER[],
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    transaction_id BIGINT,
    transaction_steps INTEGER[]
);
CREATE INDEX IF NOT EXISTS idx_appointments_deleted_at ON appointments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_appointments_deal_id ON appointments (deal_id);

CREATE TABLE IF NOT EXISTS appointment_participants (
    id BIGSERIAL PRIMARY KEY,
    appointment_id BIGINT,
    user_id BIGINT,
    role VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(255) NOT NULL DEFAULT '',
    joined_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    contact_id BIGINT
);
CREATE INDEX IF NOT EXISTS idx_appointment_participants_appointment_id ON appointment_participants (appointment_id);
CREATE INDEX IF NOT EXISTS idx_appointment_participants_deleted_at ON appointment_participants (deleted_at);

CREATE TABLE IF NOT EXISTS property_histories (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    subject_id BIGINT NOT NULL,
    subject_type VARCHAR(50),
    action BIGINT,
    metadata JSONB,
    actor_id BIGINT NOT NULL,
    actor_type INTEGER,
    description TEXT,
    occurred_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_property_histories_subject_id ON property_histories (subject_id);
CREATE INDEX IF NOT EXISTS idx_property_histories_actor_id ON property_histories (actor_id);
CREATE INDEX IF NOT EXISTS idx_property_histories_occurred_at ON property_histories (occurred_at);
CREATE INDEX IF NOT EXISTS idx_property_histories_deleted_at ON property_histories (deleted_at);
