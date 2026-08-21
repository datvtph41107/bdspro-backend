BEGIN;
CREATE TABLE payment_event_inbox (
  event_id VARCHAR(192) PRIMARY KEY,
  event_type VARCHAR(128) NOT NULL,
  order_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE payment_notification_intents (
  id BIGSERIAL PRIMARY KEY,
  event_id VARCHAR(192) NOT NULL UNIQUE REFERENCES payment_event_inbox(event_id),
  order_id BIGINT NOT NULL,
  subject_kind VARCHAR(32) NOT NULL,
  subject_id VARCHAR(128) NOT NULL,
  template_code VARCHAR(96) NOT NULL,
  status VARCHAR(32) NOT NULL CHECK (status IN ('pending','sent','failed')),
  payload TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMIT;
