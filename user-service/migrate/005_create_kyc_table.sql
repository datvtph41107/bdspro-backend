-- Create KYC table
CREATE TABLE IF NOT EXISTS kyc (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    identity_card VARCHAR(20) NOT NULL,
    front_image VARCHAR(500) NOT NULL,
    back_image VARCHAR(500) NOT NULL,
    selfie_image VARCHAR(500),
    status INT DEFAULT 10, -- 10: Pending, 20: Approved, 30: Rejected
    reject_reason VARCHAR(500),
    reviewed_by BIGINT,
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_kyc_profile_id ON kyc(profile_id);
CREATE INDEX IF NOT EXISTS idx_kyc_status ON kyc(status);
CREATE INDEX IF NOT EXISTS idx_kyc_deleted_at ON kyc(deleted_at);
CREATE INDEX IF NOT EXISTS idx_kyc_created_at ON kyc(created_at);

