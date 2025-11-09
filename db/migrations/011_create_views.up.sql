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