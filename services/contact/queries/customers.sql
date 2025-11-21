-- name: CreateCustomer :exec
INSERT INTO customers (id, name, email, phone, address, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetCustomerByID :one
SELECT id, name, email, phone, address, created_at, updated_at
FROM customers
WHERE id = $1;

-- name: GetCustomerByEmail :one
SELECT id, name, email, phone, address, created_at, updated_at
FROM customers
WHERE email = $1;

-- name: ListCustomers :many
SELECT id, name, email, phone, address, created_at, updated_at
FROM customers
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCustomer :exec
UPDATE customers
SET name = $2,
    email = $3,
    phone = $4,
    address = $5,
    updated_at = $6
WHERE id = $1;

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1;

