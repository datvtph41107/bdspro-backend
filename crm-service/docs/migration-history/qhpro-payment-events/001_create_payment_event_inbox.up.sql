BEGIN;
CREATE TABLE payment_event_inbox (
  event_id VARCHAR(192) PRIMARY KEY,
  event_type VARCHAR(128) NOT NULL,
  order_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE payment_crm_activities (
  id BIGSERIAL PRIMARY KEY,
  event_id VARCHAR(192) NOT NULL UNIQUE REFERENCES payment_event_inbox(event_id),
  order_id BIGINT NOT NULL,
  subject_kind VARCHAR(32) NOT NULL,
  subject_id VARCHAR(128) NOT NULL,
  activity_type VARCHAR(96) NOT NULL,
  product_code VARCHAR(128) NOT NULL,
  plan_code VARCHAR(128) NOT NULL,
  amount_minor BIGINT NOT NULL,
  currency VARCHAR(8) NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMIT;
