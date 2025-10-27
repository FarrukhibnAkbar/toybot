-- name: CreateProduct :one
INSERT INTO products (name, sell_price)
VALUES ($1, $2)
RETURNING *;

-- name: GetProductByName :one
SELECT * FROM products WHERE name = $1;

-- name: ListProducts :many
SELECT p.id, p.name, s.qty, s.cost_price, p.sell_price
FROM products p
LEFT JOIN stock_levels s ON s.product_id = p.id
ORDER BY p.name;

-- name: UpsertStock :exec
INSERT INTO stock_levels (product_id, qty, cost_price)
VALUES ($1, $2, $3)
ON CONFLICT (product_id)
DO UPDATE SET
    qty = stock_levels.qty + EXCLUDED.qty,
    cost_price = EXCLUDED.cost_price,
    updated_at = now();
