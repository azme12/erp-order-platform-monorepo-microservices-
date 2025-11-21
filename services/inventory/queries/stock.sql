-- name: CreateStock :exec
INSERT INTO stock (id, item_id, quantity, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetStockByItemID :one
SELECT id, item_id, quantity, created_at, updated_at
FROM stock
WHERE item_id = $1;

-- name: GetStockQuantityForUpdate :one
SELECT quantity FROM stock WHERE item_id = $1 FOR UPDATE;

-- name: UpdateStock :exec
UPDATE stock
SET quantity = $2,
    updated_at = $3
WHERE item_id = $1;

-- name: AdjustStock :exec
UPDATE stock
SET quantity = quantity + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE item_id = $1;

