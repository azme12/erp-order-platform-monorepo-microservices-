-- name: CreateItem :exec
INSERT INTO items (id, name, description, sku, unit_price, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetItemByID :one
SELECT id, name, description, sku, unit_price, created_at, updated_at
FROM items
WHERE id = $1;

-- name: GetItemBySKU :one
SELECT id, name, description, sku, unit_price, created_at, updated_at
FROM items
WHERE sku = $1;

-- name: ListItems :many
SELECT id, name, description, sku, unit_price, created_at, updated_at
FROM items
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateItem :exec
UPDATE items
SET name = $2,
    description = $3,
    sku = $4,
    unit_price = $5,
    updated_at = $6
WHERE id = $1;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = $1;

