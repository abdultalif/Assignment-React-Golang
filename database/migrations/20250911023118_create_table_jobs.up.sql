CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'QUEUED',
    progress INTEGER DEFAULT 0,
    processed BIGINT DEFAULT 0,
    total BIGINT DEFAULT 0,
    params JSONB,
    result_path VARCHAR(500),
    error_message TEXT,
    unique_run_id VARCHAR(255),
    cancelled BOOLEAN DEFAULT FALSE,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_jobs_status ON jobs (status);

CREATE INDEX idx_jobs_type ON jobs(type);

CREATE INDEX idx_jobs_created_at ON jobs (created_at);