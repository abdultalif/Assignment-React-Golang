CREATE TABLE IF NOT EXISTS orders (
    id UUID DEFAULT gen_random_uuid () PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products (id),
    buyer_id VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    total_cents INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_buyer_id ON orders (buyer_id);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);