BEGIN;
CREATE TABLE IF NOT EXISTS payment_event_inbox (
  event_id VARCHAR(192) PRIMARY KEY,
  event_type VARCHAR(128) NOT NULL,
  schema_version INT NOT NULL,
  order_id BIGINT NOT NULL,
  payload_hash CHAR(64) NOT NULL,
  received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS event_notifications (
  id BIGSERIAL PRIMARY KEY,
  source_event_id VARCHAR(192) NOT NULL UNIQUE REFERENCES payment_event_inbox(event_id),
  order_id BIGINT NOT NULL,
  subject_kind VARCHAR(32) NOT NULL,
  subject_id VARCHAR(128) NOT NULL,
  kind VARCHAR(96) NOT NULL,
  title VARCHAR(255) NOT NULL,
  body TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_event_notifications_subject ON event_notifications(subject_kind, subject_id, created_at DESC);

CREATE TABLE IF NOT EXISTS notification_delivery_intents (
  id BIGSERIAL PRIMARY KEY,
  event_id VARCHAR(192) NOT NULL REFERENCES payment_event_inbox(event_id),
  subject_kind VARCHAR(32) NOT NULL,
  subject_id VARCHAR(128) NOT NULL,
  channel VARCHAR(32) NOT NULL,
  template_code VARCHAR(96) NOT NULL,
  payload TEXT NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','running','retry','sent','failed','unknown')),
  attempt_count INT NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
  available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  locked_by VARCHAR(128),
  claim_version BIGINT NOT NULL DEFAULT 0 CHECK (claim_version >= 0),
  lease_until TIMESTAMPTZ,
  last_error TEXT NOT NULL DEFAULT '',
  sent_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT uq_notification_delivery_event_channel UNIQUE(event_id, channel)
);
CREATE INDEX IF NOT EXISTS idx_notification_delivery_claim ON notification_delivery_intents(status, available_at, lease_until, id);
COMMIT;
