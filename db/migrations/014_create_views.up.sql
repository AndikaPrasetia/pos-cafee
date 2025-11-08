-- Daily sales summary view
CREATE VIEW daily_sales_summary AS
SELECT
    DATE(created_at) as sale_date,
    COUNT(*) as total_orders,
    SUM(total_amount) as total_sales,
    AVG(total_amount) as average_order_value,
    COUNT(DISTINCT user_id) as active_cashiers
FROM orders
WHERE status = 'completed'
GROUP BY DATE(created_at);

-- Menu performance view
CREATE VIEW menu_performance AS
SELECT
    mi.id,
    mi.name,
    mi.category_id,
    c.name as category_name,
    COUNT(oi.id) as times_ordered,
    SUM(oi.quantity) as total_quantity,
    SUM(oi.subtotal) as total_revenue,
    mi.price,
    mi.cost,
    (mi.price - mi.cost) as profit_margin
FROM menu_items mi
LEFT JOIN categories c ON mi.category_id = c.id
LEFT JOIN order_items oi ON mi.id = oi.menu_item_id
LEFT JOIN orders o ON oi.order_id = o.id AND o.status = 'completed'
GROUP BY mi.id, mi.name, mi.category_id, c.name, mi.price, mi.cost;

-- Inventory alert view
CREATE VIEW inventory_alerts AS
SELECT
    i.id,
    mi.name as menu_item_name,
    i.current_stock,
    i.min_stock,
    i.unit,
    CASE
        WHEN i.current_stock = 0 THEN 'out_of_stock'
        WHEN i.current_stock <= i.min_stock THEN 'low_stock'
        ELSE 'adequate'
    END as stock_status,
    i.last_updated
FROM inventory i
JOIN menu_items mi ON i.menu_item_id = mi.id
WHERE i.current_stock <= i.min_stock OR i.current_stock = 0;

-- Financial summary view
CREATE VIEW financial_summary AS
SELECT
    DATE(o.created_at) as transaction_date,
    COALESCE(SUM(o.total_amount), 0) as daily_income,
    COALESCE(SUM(e.amount), 0) as daily_expenses,
    COALESCE(SUM(o.total_amount), 0) - COALESCE(SUM(e.amount), 0) as daily_profit
FROM orders o
FULL JOIN expenses e ON DATE(o.created_at) = e.expense_date
WHERE o.status = 'completed' OR e.id IS NOT NULL
GROUP BY DATE(o.created_at);
