BEGIN;

CREATE TABLE quota_pools (
    subject_type VARCHAR(32) NOT NULL
        CHECK (subject_type IN ('profile', 'organization')),
    subject_id VARCHAR(128) NOT NULL CHECK (subject_id <> ''),
    meter_code VARCHAR(128) NOT NULL CHECK (meter_code <> ''),
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    used BIGINT NOT NULL DEFAULT 0 CHECK (used >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (subject_type, subject_id, meter_code, period_start),
    CONSTRAINT chk_quota_pool_period CHECK (period_end > period_start)
);

-- period_end is deliberately excluded from identity. If two decisions resolve the
-- same pool start with different ends, they must collide and fail closed instead
-- of silently creating two allowance buckets.

COMMIT;
