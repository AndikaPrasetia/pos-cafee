-- Inventory table
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE,
    menu_item_id INTEGER NOT NULL UNIQUE REFERENCES menu_items(id) ON DELETE CASCADE,
    current_stock INTEGER NOT NULL DEFAULT 0 CHECK (current_stock >= 0),
    min_stock INTEGER DEFAULT 10 CHECK (min_stock >= 0),
    unit VARCHAR(20) DEFAULT 'pcs',
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER REFERENCES users(id)
);
