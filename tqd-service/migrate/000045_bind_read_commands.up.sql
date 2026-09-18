BEGIN;
-- Owner: TQD. Khóa ý định không phụ thuộc meter hay chu kỳ quota.
-- accepted_at và accepted_usage_key được ghi cùng transaction với usage event.
CREATE TABLE quota_read_commands (
    subject_type varchar(32) NOT NULL CHECK (subject_type IN ('profile', 'organization')),
    subject_id varchar(128) NOT NULL,
    operation varchar(128) NOT NULL,
    command_key varchar(128) NOT NULL,
    request_hash varchar(64) CHECK (request_hash ~ '^[0-9a-f]{64}$'),
    accepted_usage_key char(64) REFERENCES quota_usage_events(usage_key),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    accepted_at timestamptz,
    PRIMARY KEY (subject_type, subject_id, operation, command_key),
    CHECK (accepted_usage_key IS NULL OR accepted_at IS NOT NULL)
);
-- Không đoán payload của usage lịch sử. Retry legacy phải bắt đầu ý định mới,
-- thay vì chiếm key cũ rồi truy cập dữ liệu khác mà không tính quota.
INSERT INTO quota_read_commands(subject_type, subject_id, operation, command_key, accepted_usage_key, created_at, accepted_at)
SELECT DISTINCT ON (subject_type, subject_id, operation, command_key)
       subject_type, subject_id, operation, command_key, usage_key, created_at, created_at
FROM quota_usage_events
WHERE operation IN ('workspace.area.analyze', 'workspace.parcel.detail.view')
ORDER BY subject_type, subject_id, operation, command_key, created_at, id;
COMMIT;
