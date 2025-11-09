-- name: GetUser :one
SELECT id, username, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at
FROM users
WHERE id = $1 AND is_active = true
LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at
FROM users
WHERE username = $1 AND is_active = true
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at
FROM users
WHERE email = $1 AND is_active = true
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    username, email, password_hash, role, first_name, last_name
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, username, email, role, first_name, last_name, is_active, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, role = $4, first_name = $5, last_name = $6, updated_at = NOW()
WHERE id = $1
RETURNING id, username, email, role, first_name, last_name, is_active, created_at, updated_at;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserStatus :exec
UPDATE users
SET is_active = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
UPDATE users
SET is_active = false, updated_at = NOW()
WHERE id = $1;