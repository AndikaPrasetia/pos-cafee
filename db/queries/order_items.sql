-- name: GetOrderItem :one
SELECT id, order_id, menu_item_id, quantity, unit_price, total_price, created_at, updated_at
FROM order_items
WHERE id = $1
LIMIT 1;

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, menu_item_id, quantity, unit_price, total_price, created_at, updated_at
FROM order_items
WHERE order_id = $1
ORDER BY created_at;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id, menu_item_id, quantity, unit_price, total_price
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, order_id, menu_item_id, quantity, unit_price, total_price, created_at, updated_at;

-- name: UpdateOrderItem :one
UPDATE order_items
SET quantity = $2, unit_price = $3, total_price = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, order_id, menu_item_id, quantity, unit_price, total_price, created_at, updated_at;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;

-- name: DeleteOrderItemsByOrderID :exec
DELETE FROM order_items
WHERE order_id = $1;

-- name: GetOrderItemsWithDetails :many
SELECT oi.id, oi.order_id, oi.menu_item_id, mi.name as menu_item_name, oi.quantity, oi.unit_price, oi.total_price, oi.created_at, oi.updated_at
FROM order_items oi
JOIN menu_items mi ON oi.menu_item_id = mi.id
WHERE oi.order_id = $1
ORDER BY oi.created_at;