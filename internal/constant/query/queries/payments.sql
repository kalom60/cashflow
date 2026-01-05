-- name: CreatePayment :one
INSERT INTO payments (reference, amount, currency, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $5)
RETURNING *;

-- name: GetPaymentByID :one
SELECT *
FROM payments
WHERE id = $1;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET status = $2
WHERE id = $1
RETURNING *;

-- name: GetPaymentByIDForUpdate :one
SELECT *
FROM payments
WHERE id = $1
FOR UPDATE;
