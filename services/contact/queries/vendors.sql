-- name: CreateVendor :exec
INSERT INTO vendors (id, name, email, phone, address, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetVendorByID :one
SELECT id, name, email, phone, address, created_at, updated_at
FROM vendors
WHERE id = $1;

-- name: GetVendorByEmail :one
SELECT id, name, email, phone, address, created_at, updated_at
FROM vendors
WHERE email = $1;

-- name: ListVendors :many
SELECT id, name, email, phone, address, created_at, updated_at
FROM vendors
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateVendor :exec
UPDATE vendors
SET name = $2,
    email = $3,
    phone = $4,
    address = $5,
    updated_at = $6
WHERE id = $1;

-- name: DeleteVendor :exec
DELETE FROM vendors
WHERE id = $1;

