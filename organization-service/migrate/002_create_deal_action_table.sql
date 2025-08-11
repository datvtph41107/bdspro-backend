-- Migration: Create deal_action table
-- Date: 2025-01-27

CREATE TABLE IF NOT EXISTS deal_action (
    id BIGSERIAL PRIMARY KEY,
    deal_id BIGINT NOT NULL,
    transaction_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    product_ids BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    transaction_step INTEGER NOT NULL,
    transaction_type INTEGER NOT NULL,
    note TEXT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL,
    organization_id BIGINT NOT NULL
);

-- Create indexes
CREATE INDEX idx_deal_action_deal_id ON deal_action(deal_id);
CREATE INDEX idx_deal_action_transaction_id ON deal_action(transaction_id);
CREATE INDEX idx_deal_action_customer_id ON deal_action(customer_id);
CREATE INDEX idx_deal_action_organization_id ON deal_action(organization_id);
CREATE INDEX idx_deal_action_created_at ON deal_action(created_at);

-- Add foreign key constraints if needed
-- ALTER TABLE deal_action ADD CONSTRAINT fk_deal_action_deal_id FOREIGN KEY (deal_id) REFERENCES deals(id);
-- ALTER TABLE deal_action ADD CONSTRAINT fk_deal_action_customer_id FOREIGN KEY (customer_id) REFERENCES customers(id);
