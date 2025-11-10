-- name: GetOrder :one
SELECT id, order_number, user_id, status, total_amount, discount_amount, tax_amount, 
       payment_method, payment_status, completed_at, created_at, updated_at
FROM orders
WHERE id = $1
LIMIT 1;

-- name: GetOrderByNumber :one
SELECT id, order_number, user_id, status, total_amount, discount_amount, tax_amount, 
       payment_method, payment_status, completed_at, created_at, updated_at
FROM orders
WHERE order_number = $1
LIMIT 1;

-- name: ListOrders :many
SELECT id, order_number, user_id, status, total_amount, discount_amount, tax_amount, 
       payment_method, payment_status, completed_at, created_at, updated_at
FROM orders
WHERE ($1 = '' OR status = $1) 
  AND ($2 = '00000000-0000-0000-0000-000000000000'::uuid OR user_id = $2)
  AND ($3 = '0001-01-01'::date OR created_at >= $3) 
  AND ($4 = '0001-01-01'::date OR created_at <= $4)
ORDER BY created_at DESC
LIMIT $5 OFFSET $6;

-- name: CreateOrder :one
INSERT INTO orders (
    order_number, user_id, total_amount, discount_amount, tax_amount
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, order_number, user_id, status, total_amount, discount_amount, tax_amount, 
          payment_method, payment_status, completed_at, created_at, updated_at;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateOrderPayment :exec
UPDATE orders
SET payment_method = $2, payment_status = $3, completed_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: UpdateOrderTotal :exec
UPDATE orders
SET total_amount = $2, discount_amount = $3, tax_amount = $4, updated_at = NOW()
WHERE id = $1;