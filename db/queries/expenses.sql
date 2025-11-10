-- name: CreateExpense :one
INSERT INTO expenses (
    category, description, amount, date, user_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, category, description, amount, date, user_id, created_at;

-- name: GetExpense :one
SELECT id, category, description, amount, date, user_id, created_at
FROM expenses
WHERE id = $1
LIMIT 1;

-- name: ListExpenses :many
SELECT id, category, description, amount, date, user_id, created_at
FROM expenses
WHERE ($1::date IS NULL OR date >= $1)
  AND ($2::date IS NULL OR date <= $2)
  AND ($3::text IS NULL OR category = $3)
ORDER BY date DESC, created_at DESC
LIMIT $4 OFFSET $5;

-- name: UpdateExpense :one
UPDATE expenses
SET category = $2, description = $3, amount = $4, date = $5
WHERE id = $1
RETURNING id, category, description, amount, date, user_id, created_at;

-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id = $1;