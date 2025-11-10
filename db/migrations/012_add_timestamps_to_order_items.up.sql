-- Add created_at and updated_at columns to order_items table
ALTER TABLE order_items ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();
ALTER TABLE order_items ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();

-- Add indexes for performance optimization
CREATE INDEX idx_order_items_created_at ON order_items(created_at);
CREATE INDEX idx_order_items_updated_at ON order_items(updated_at);