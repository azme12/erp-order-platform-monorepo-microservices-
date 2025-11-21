-- name: CreateOrder :exec
INSERT INTO purchase_orders (id, vendor_id, status, total_amount, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetOrderByID :one
SELECT id, vendor_id, status, total_amount, created_at, updated_at
FROM purchase_orders
WHERE id = $1;

-- name: ListOrders :many
SELECT id, vendor_id, status, total_amount, created_at, updated_at
FROM purchase_orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrder :exec
UPDATE purchase_orders
SET vendor_id = $2,
    status = $3,
    total_amount = $4,
    updated_at = $5
WHERE id = $1;

-- name: UpdateOrderStatus :exec
UPDATE purchase_orders
SET status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

