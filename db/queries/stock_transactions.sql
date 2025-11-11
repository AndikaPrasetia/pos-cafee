-- name: CreateStockTransaction :one
INSERT INTO stock_transactions (
    menu_item_id, transaction_type, quantity, previous_stock, current_stock,
    reason, reference_type, reference_id, user_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, menu_item_id, transaction_type, quantity, previous_stock, current_stock,
          reason, reference_type, reference_id, user_id, created_at;

-- name: ListStockTransactions :many
SELECT st.id, st.menu_item_id, mi.name as menu_item_name, st.transaction_type,
       st.quantity, st.previous_stock, st.current_stock, st.reason,
       st.reference_type, st.reference_id, st.user_id, u.username as user_name, st.created_at
FROM stock_transactions st
LEFT JOIN menu_items mi ON st.menu_item_id = mi.id
LEFT JOIN users u ON st.user_id = u.id
WHERE ($1 = '00000000-0000-0000-0000-000000000000'::uuid OR st.menu_item_id = $1)
  AND ($2 = '0001-01-01'::date OR st.created_at >= $2)
  AND ($3 = '0001-01-01'::date OR st.created_at <= $3)
ORDER BY st.created_at DESC
LIMIT $4 OFFSET $5;