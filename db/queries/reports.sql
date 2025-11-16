-- name: GetDailySalesReportData :one
SELECT
    COALESCE(dss.total_orders, 0) AS total_orders,
    COALESCE(dss.total_sales, '0')::TEXT AS total_sales
FROM daily_sales_summary dss
WHERE dss.sale_date = $1::date;

-- name: GetTopSellingItemsByDateRange :many
SELECT
    mi.name AS menu_item_name,
    SUM(oi.quantity) AS total_quantity_sold,
    SUM(oi.total_price)::TEXT AS total_revenue
FROM order_items oi
JOIN orders o ON oi.order_id = o.id
JOIN menu_items mi ON oi.menu_item_id = mi.id
WHERE o.status = 'completed'
AND o.completed_at IS NOT NULL
AND o.completed_at >= $1::timestamp
AND o.completed_at <= $2::timestamp
GROUP BY mi.id, mi.name
ORDER BY total_quantity_sold DESC
LIMIT $3;

-- name: GetFinancialSummaryByDateRange :one
SELECT
    COUNT(o.id) AS total_orders,
    COALESCE(SUM(o.total_amount), '0')::TEXT AS total_sales,
    COALESCE(SUM(o.discount_amount), '0')::TEXT AS total_discount,
    COALESCE(SUM(o.tax_amount), '0')::TEXT AS total_tax
FROM orders o
WHERE o.status = 'completed'
AND o.completed_at IS NOT NULL
AND o.completed_at >= $1::timestamp
AND o.completed_at <= $2::timestamp;

-- name: GetSalesByCategoryByDateRange :many
SELECT
    c.name AS category_name,
    COUNT(oi.id) AS items_sold,
    SUM(oi.quantity) AS total_quantity,
    SUM(oi.total_price)::TEXT AS total_revenue
FROM order_items oi
JOIN orders o ON oi.order_id = o.id
JOIN menu_items mi ON oi.menu_item_id = mi.id
JOIN categories c ON mi.category_id = c.id
WHERE o.status = 'completed'
AND o.completed_at IS NOT NULL
AND o.completed_at >= $1::timestamp
AND o.completed_at <= $2::timestamp
GROUP BY c.id, c.name
ORDER BY total_revenue DESC;