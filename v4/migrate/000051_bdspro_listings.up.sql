-- Listing rows attached to products (property/new post payload).
CREATE TABLE IF NOT EXISTS listings (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    product_id BIGINT NOT NULL,
    owner_id BIGINT NOT NULL DEFAULT 0,
    title VARCHAR(500) NOT NULL DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    post_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    transaction_type SMALLINT NOT NULL DEFAULT 10,
    visibility SMALLINT NOT NULL DEFAULT 30,
    expired_at TIMESTAMPTZ,
    num_date INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_listings_deleted_at ON listings (deleted_at);
CREATE INDEX IF NOT EXISTS idx_listings_product_id ON listings (product_id);
CREATE INDEX IF NOT EXISTS idx_listings_owner_id ON listings (owner_id);

ALTER TABLE products ADD COLUMN IF NOT EXISTS post_count INTEGER NOT NULL DEFAULT 0;
