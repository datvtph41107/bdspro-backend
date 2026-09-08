CREATE TABLE IF NOT EXISTS report_events (
    id BIGSERIAL PRIMARY KEY,
    report_id BIGINT NOT NULL REFERENCES reports(id),
    actor_id BIGINT NOT NULL,
    action VARCHAR(64) NOT NULL,
    before_json TEXT,
    after_json TEXT,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_report_events_report ON report_events(report_id);
CREATE INDEX IF NOT EXISTS idx_report_events_created ON report_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_report_events_action ON report_events(action);
CREATE INDEX IF NOT EXISTS idx_report_events_actor ON report_events(actor_id);
