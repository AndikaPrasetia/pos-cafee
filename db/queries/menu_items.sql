-- name: GetMenuItem :one
SELECT id, name, category_id, description, price, cost, is_available, created_at, updated_at
FROM menu_items
WHERE id = $1 AND is_available = true
LIMIT 1;

-- name: ListMenuItems :many
SELECT id, name, category_id, description, price, cost, is_available, created_at, updated_at
FROM menu_items
WHERE is_available = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: ListMenuItemsByCategory :many
SELECT id, name, category_id, description, price, cost, is_available, created_at, updated_at
FROM menu_items
WHERE category_id = $1 AND is_available = true
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: CreateMenuItem :one
INSERT INTO menu_items (
    name, category_id, description, price, cost
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, name, category_id, description, price, cost, is_available, created_at, updated_at;

-- name: UpdateMenuItem :one
UPDATE menu_items
SET name = $2, category_id = $3, description = $4, price = $5, cost = $6, is_available = $7, updated_at = NOW()
WHERE id = $1
RETURNING id, name, category_id, description, price, cost, is_available, created_at, updated_at;

-- name: DeleteMenuItem :exec
UPDATE menu_items
SET is_available = false, updated_at = NOW()
WHERE id = $1;