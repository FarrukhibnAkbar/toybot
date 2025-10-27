-- name: AddPurchase :one
INSERT INTO purchases (product_id, qty, cost_price, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;
