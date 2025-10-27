-- name: AddSale :one
INSERT INTO sales (product_id, qty, sell_price, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ReduceStock :exec
UPDATE stock_levels
SET qty = qty - $2, updated_at = now()
WHERE product_id = $1 AND qty >= $2;
