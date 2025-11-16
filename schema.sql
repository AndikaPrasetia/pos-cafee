-- POS Cafe Complete SQL Schema

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('cashier', 'manager', 'admin')),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for users table
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Create user_sessions table
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_accessed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for user_sessions table
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions(token);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Create categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for categories table
CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_categories_is_active ON categories(is_active);
CREATE INDEX idx_categories_created_at ON categories(created_at);

-- Create menu_items table
CREATE TABLE menu_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    cost DECIMAL(10,2) NOT NULL CHECK (cost >= 0),
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Ensure cost is not greater than price
    CONSTRAINT cost_not_greater_than_price CHECK (cost <= price)
);

-- Create indexes for menu_items table
CREATE INDEX idx_menu_items_category_id ON menu_items(category_id);
CREATE INDEX idx_menu_items_name ON menu_items(name);
CREATE INDEX idx_menu_items_is_available ON menu_items(is_available);
CREATE INDEX idx_menu_items_price ON menu_items(price);
CREATE INDEX idx_menu_items_created_at ON menu_items(created_at);

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

-- Create indexes for inventory table
CREATE INDEX idx_inventory_menu_item_id ON inventory(menu_item_id);
CREATE INDEX idx_inventory_current_stock ON inventory(current_stock);
CREATE INDEX idx_inventory_minimum_stock ON inventory(minimum_stock);
CREATE INDEX idx_inventory_last_updated_at ON inventory(last_updated_at);

-- Create orders table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'pending', 'completed', 'cancelled')) DEFAULT 'draft',
    total_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00 CHECK (total_amount >= 0),
    discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00 CHECK (discount_amount >= 0),
    tax_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00 CHECK (tax_amount >= 0),
    payment_method VARCHAR(20) CHECK (payment_method IN ('cash', 'card', 'qris', 'transfer')),
    payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN ('pending', 'paid', 'failed')) DEFAULT 'pending',
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for orders table
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_completed_at ON orders(completed_at);
CREATE INDEX idx_orders_total_amount ON orders(total_amount);
CREATE INDEX idx_orders_status_completed_at ON orders(status, completed_at);

-- Create order_items table
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id UUID NOT NULL REFERENCES menu_items(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    total_price DECIMAL(12,2) NOT NULL CHECK (total_price >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for order_items table
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_menu_item_id ON order_items(menu_item_id);
CREATE INDEX idx_order_items_quantity ON order_items(quantity);
CREATE INDEX idx_order_items_total_price ON order_items(total_price);
CREATE INDEX idx_order_items_created_at ON order_items(created_at);
CREATE INDEX idx_order_items_updated_at ON order_items(updated_at);
CREATE INDEX idx_order_items_order_id_menu_item_id ON order_items(order_id, menu_item_id);

-- Create stock_transactions table
CREATE TABLE stock_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    menu_item_id UUID NOT NULL REFERENCES menu_items(id),
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('in', 'out', 'adjustment')),
    quantity INTEGER NOT NULL,
    previous_stock INTEGER NOT NULL,
    current_stock INTEGER NOT NULL,
    reason VARCHAR(255) NOT NULL,
    reference_type VARCHAR(50),
    reference_id UUID,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for stock_transactions table
CREATE INDEX idx_stock_transactions_menu_item_id ON stock_transactions(menu_item_id);
CREATE INDEX idx_stock_transactions_transaction_type ON stock_transactions(transaction_type);
CREATE INDEX idx_stock_transactions_created_at ON stock_transactions(created_at);
CREATE INDEX idx_stock_transactions_user_id ON stock_transactions(user_id);
CREATE INDEX idx_stock_transactions_quantity ON stock_transactions(quantity);

-- Create expenses table
CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category VARCHAR(100) NOT NULL,
    description TEXT,
    amount DECIMAL(12,2) NOT NULL CHECK (amount >= 0),
    date DATE NOT NULL,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for expenses table
CREATE INDEX idx_expenses_category ON expenses(category);
CREATE INDEX idx_expenses_date ON expenses(date);
CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_amount ON expenses(amount);
CREATE INDEX idx_expenses_created_at ON expenses(created_at);
CREATE INDEX idx_expenses_date_category ON expenses(date, category);

-- Create view for menu items with category information
CREATE VIEW menu_items_with_category AS
SELECT
    mi.id,
    mi.name,
    mi.description,
    mi.price,
    mi.cost,
    mi.is_available,
    mi.created_at,
    mi.updated_at,
    c.id AS category_id,
    c.name AS category_name,
    c.description AS category_description
FROM menu_items mi
JOIN categories c ON mi.category_id = c.id;

-- Create view for inventory with menu item and user information
CREATE VIEW inventory_with_details AS
SELECT
    i.id,
    i.menu_item_id,
    mi.name AS menu_item_name,
    i.current_stock,
    i.minimum_stock,
    i.unit,
    i.last_updated_at,
    u.username AS last_updated_by_username,
    CASE
        WHEN i.current_stock <= i.minimum_stock THEN 'LOW'
        ELSE 'OK'
    END AS stock_status
FROM inventory i
JOIN menu_items mi ON i.menu_item_id = mi.id
LEFT JOIN users u ON i.last_updated_by = u.id;

-- Create view for orders with user information
CREATE VIEW orders_with_user AS
SELECT
    o.id,
    o.order_number,
    o.user_id,
    u.username,
    u.first_name,
    u.last_name,
    o.status,
    o.total_amount,
    o.discount_amount,
    o.tax_amount,
    o.payment_method,
    o.payment_status,
    o.completed_at,
    o.created_at,
    o.updated_at
FROM orders o
JOIN users u ON o.user_id = u.id;

-- Create view for order items with menu item details
CREATE VIEW order_items_with_details AS
SELECT
    oi.id,
    oi.order_id,
    o.order_number,
    oi.menu_item_id,
    mi.name AS menu_item_name,
    oi.quantity,
    oi.unit_price,
    oi.total_price
FROM order_items oi
JOIN orders o ON oi.order_id = o.id
JOIN menu_items mi ON oi.menu_item_id = mi.id;

-- Create view for daily sales summary
CREATE VIEW daily_sales_summary AS
SELECT
    DATE(o.completed_at) AS sale_date,
    COUNT(*) AS total_orders,
    SUM(o.total_amount) AS total_sales,
    SUM(o.discount_amount) AS total_discount,
    SUM(o.tax_amount) AS total_tax
FROM orders o
WHERE o.status = 'completed'
AND o.completed_at IS NOT NULL
GROUP BY DATE(o.completed_at);

-- Create view for monthly sales summary
CREATE VIEW monthly_sales_summary AS
SELECT
    DATE_TRUNC('month', o.completed_at)::date AS sale_month,
    COUNT(*) AS total_orders,
    SUM(o.total_amount) AS total_sales,
    SUM(o.discount_amount) AS total_discount,
    SUM(o.tax_amount) AS total_tax
FROM orders o
WHERE o.status = 'completed'
AND o.completed_at IS NOT NULL
GROUP BY DATE_TRUNC('month', o.completed_at);

-- Create view for top selling items
CREATE VIEW top_selling_items AS
SELECT
    mi.id AS menu_item_id,
    mi.name AS menu_item_name,
    mi.description,
    c.name AS category_name,
    SUM(oi.quantity) AS total_quantity_sold,
    SUM(oi.total_price) AS total_revenue,
    COUNT(DISTINCT oi.order_id) AS times_ordered
FROM order_items oi
JOIN orders o ON oi.order_id = o.id
JOIN menu_items mi ON oi.menu_item_id = mi.id
JOIN categories c ON mi.category_id = c.id
WHERE o.status = 'completed'
AND o.completed_at IS NOT NULL
GROUP BY mi.id, mi.name, mi.description, c.name
ORDER BY total_quantity_sold DESC;