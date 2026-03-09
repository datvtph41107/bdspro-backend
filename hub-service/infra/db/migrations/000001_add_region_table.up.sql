CREATE TABLE IF NOT EXISTS regions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    code VARCHAR(50),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP
);