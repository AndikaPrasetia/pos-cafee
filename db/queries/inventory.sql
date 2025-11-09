-- name: GetInventoryByMenuItem :one
SELECT id, menu_item_id, current_stock, minimum_stock, unit, last_updated_at, last_updated_by
FROM inventory
WHERE menu_item_id = $1
LIMIT 1;

-- name: ListInventory :many
SELECT i.id, i.menu_item_id, mi.name as menu_item_name, i.current_stock, i.minimum_stock, 
       i.unit, i.last_updated_at, u.username as last_updated_by
FROM inventory i
JOIN menu_items mi ON i.menu_item_id = mi.id
LEFT JOIN users u ON i.last_updated_by = u.id
WHERE ($1::boolean IS NULL OR (i.current_stock <= i.minimum_stock) = $1)  -- low_stock_only
ORDER BY i.current_stock
LIMIT $2 OFFSET $3;

-- name: UpdateInventoryStock :exec
UPDATE inventory
SET current_stock = $2, last_updated_at = NOW(), last_updated_by = $3
WHERE menu_item_id = $1;

-- name: CreateInventoryRecord :exec
INSERT INTO inventory (menu_item_id)
VALUES ($1);