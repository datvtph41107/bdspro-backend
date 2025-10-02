-- Create amenities table
CREATE TABLE IF NOT EXISTS amenities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_amenities_created_at ON amenities(created_at);
CREATE INDEX IF NOT EXISTS idx_amenities_deleted_at ON amenities(deleted_at);
CREATE INDEX IF NOT EXISTS idx_amenities_created_by ON amenities(created_by);
CREATE INDEX IF NOT EXISTS idx_amenities_updated_by ON amenities(updated_by);

-- Create trigger to auto update updated_at
CREATE OR REPLACE FUNCTION update_amenities_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_amenities_updated_at
    BEFORE UPDATE ON amenities
    FOR EACH ROW
    EXECUTE FUNCTION update_amenities_updated_at();

