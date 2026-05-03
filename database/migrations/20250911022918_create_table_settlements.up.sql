CREATE TABLE IF NOT EXISTS settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    merchant_id VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    gross_cents BIGINT NOT NULL DEFAULT 0,
    fee_cents BIGINT NOT NULL DEFAULT 0,
    net_cents BIGINT NOT NULL DEFAULT 0,
    txn_count INTEGER NOT NULL DEFAULT 0,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    unique_run_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_settlements_merchant_date ON settlements (merchant_id, date);

CREATE INDEX idx_settlements_run_id ON settlements (unique_run_id);