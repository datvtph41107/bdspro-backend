-- Create province_v2 table
CREATE TABLE IF NOT EXISTS province_v2 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code INTEGER UNIQUE NOT NULL,
    codename VARCHAR(100),
    division_type VARCHAR(100),
    phone_code INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create ward_v2 table
CREATE TABLE IF NOT EXISTS ward_v2 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code INTEGER UNIQUE NOT NULL,
    codename VARCHAR(100),
    division_type VARCHAR(100),
    short_codename VARCHAR(100),
    province_code INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (province_code) REFERENCES province_v2(code) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_province_v2_code ON province_v2(code);
CREATE INDEX idx_province_v2_codename ON province_v2(codename);
CREATE INDEX idx_ward_v2_code ON ward_v2(code);
CREATE INDEX idx_ward_v2_province_code ON ward_v2(province_code);
CREATE INDEX idx_ward_v2_codename ON ward_v2(codename);

-- Create trigger to update updated_at timestamp for province_v2
CREATE OR REPLACE FUNCTION update_province_v2_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER province_v2_updated_at_trigger
BEFORE UPDATE ON province_v2
FOR EACH ROW
EXECUTE FUNCTION update_province_v2_updated_at();

-- Create trigger to update updated_at timestamp for ward_v2
CREATE OR REPLACE FUNCTION update_ward_v2_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ward_v2_updated_at_trigger
BEFORE UPDATE ON ward_v2
FOR EACH ROW
EXECUTE FUNCTION update_ward_v2_updated_at();

