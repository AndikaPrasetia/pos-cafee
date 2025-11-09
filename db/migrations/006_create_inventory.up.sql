-- Create inventory table
CREATE TABLE inventory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    menu_item_id UUID NOT NULL UNIQUE REFERENCES menu_items(id) ON DELETE CASCADE,
    current_stock INTEGER NOT NULL DEFAULT 0 CHECK (current_stock >= 0),
    minimum_stock INTEGER NOT NULL DEFAULT 0 CHECK (minimum_stock >= 0),
    unit VARCHAR(50) NOT NULL DEFAULT 'pieces',
    last_updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_updated_by UUID REFERENCES users(id)
);

-- Create indexes for performance optimization
CREATE INDEX idx_inventory_menu_item_id ON inventory(menu_item_id);
CREATE INDEX idx_inventory_current_stock ON inventory(current_stock);
CREATE INDEX idx_inventory_minimum_stock ON inventory(minimum_stock);
CREATE INDEX idx_inventory_last_updated_at ON inventory(last_updated_at);