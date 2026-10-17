-- Deal internal notes (POST /v4/deal/{id}/internal-note).

CREATE TABLE IF NOT EXISTS deal_internal_notes (
    id BIGSERIAL PRIMARY KEY,
    deal_id BIGINT NOT NULL,
    actor_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    action_type BIGINT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_deal_internal_notes_deal_id ON deal_internal_notes (deal_id);
CREATE INDEX IF NOT EXISTS idx_deal_internal_notes_action_type ON deal_internal_notes (action_type);
