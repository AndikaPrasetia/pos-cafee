-- Remove created_at and updated_at columns from order_items table
DROP INDEX IF EXISTS idx_order_items_created_at;
DROP INDEX IF EXISTS idx_order_items_updated_at;

ALTER TABLE order_items DROP COLUMN created_at;
ALTER TABLE order_items DROP COLUMN updated_at;