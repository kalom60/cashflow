-- name: CreateOutboxEvent :one
INSERT INTO outbox_events (payload, status, created_at, updated_at)
VALUES ($1, $2, $3, $3)
RETURNING *;

-- name: GetPendingOutboxEventsForUpdate :many
SELECT *
FROM outbox_events
WHERE status = 'PENDING'
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: UpdateOutboxStatus :execrows
UPDATE outbox_events
SET
    status = $2,
    updated_at = $3
WHERE id = $1;

-- name: DeleteOutboxEvent :execrows
DELETE FROM outbox_events WHERE id = $1;
