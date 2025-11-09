-- name: GetCategory :one
SELECT id, name, description, is_active, created_at, updated_at
FROM categories
WHERE id = $1 AND is_active = true
LIMIT 1;

-- name: ListCategories :many
SELECT id, name, description, is_active, created_at, updated_at
FROM categories
WHERE is_active = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: CreateCategory :one
INSERT INTO categories (
    name, description
) VALUES (
    $1, $2
)
RETURNING id, name, description, is_active, created_at, updated_at;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, description = $3, is_active = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, is_active, created_at, updated_at;

-- name: DeleteCategory :exec
UPDATE categories
SET is_active = false, updated_at = NOW()
WHERE id = $1;