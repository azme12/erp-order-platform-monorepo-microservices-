-- name: CreateOrder :exec
INSERT INTO sales_orders (id, customer_id, status, total_amount, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetOrderByID :one
SELECT id, customer_id, status, total_amount, created_at, updated_at
FROM sales_orders
WHERE id = $1;

-- name: ListOrders :many
SELECT id, customer_id, status, total_amount, created_at, updated_at
FROM sales_orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrder :exec
UPDATE sales_orders
SET customer_id = $2,
    status = $3,
    total_amount = $4,
    updated_at = $5
WHERE id = $1;

-- name: UpdateOrderStatus :exec
UPDATE sales_orders
SET status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

