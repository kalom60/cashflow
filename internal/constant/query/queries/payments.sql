-- name: CreatePayment :one
INSERT INTO payments (reference, amount, currency, status)
VALUES ($1, $2, $3, $4)
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
