-- name: CreateOrderItem :exec
INSERT INTO purchase_order_items (id, order_id, item_id, quantity, unit_price, subtotal, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, item_id, quantity, unit_price, subtotal, created_at, updated_at
FROM purchase_order_items
WHERE order_id = $1
ORDER BY created_at ASC;

-- name: DeleteOrderItemsByOrderID :exec
DELETE FROM purchase_order_items
WHERE order_id = $1;

