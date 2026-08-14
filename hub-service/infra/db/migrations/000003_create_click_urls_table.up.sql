CREATE TABLE IF NOT EXISTS click_urls (
    id          BIGSERIAL PRIMARY KEY,
    event_time  BIGINT      NOT NULL,
    domain_id   BIGINT      NOT NULL,
    domain_type INTEGER     NOT NULL,
    profile_id  BIGINT
);

CREATE INDEX IF NOT EXISTS idx_click_urls_event_time
    ON click_urls (event_time);

CREATE INDEX IF NOT EXISTS idx_click_urls_domain_id
    ON click_urls (domain_id);

CREATE INDEX IF NOT EXISTS idx_click_urls_domain_type
    ON click_urls (domain_type);

CREATE INDEX IF NOT EXISTS idx_click_urls_profile_id
    ON click_urls (profile_id);

CREATE INDEX IF NOT EXISTS idx_click_urls_domain_type_id
    ON click_urls (domain_type, domain_id);
