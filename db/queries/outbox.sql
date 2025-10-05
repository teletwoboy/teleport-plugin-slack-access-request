-- name: CreateOutbox :one
INSERT INTO outbox (
    event_type,
    aggregate_type,
    aggregate_id,
    payload,
    status,
    use_yn,
    create_code,
    create_date,
    version
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING outbox_id;

-- name: ClaimOutboxByOutboxID :one
UPDATE outbox o
SET status=$2,
    attempts = attempts + 1,
    next_try_at = now() + make_interval(secs => $3),
    update_code = $4,
    update_date = $5
WHERE o.outbox_id = $1
RETURNING *;

-- name: ClaimOutboxes :many
WITH picked AS (
    SELECT outbox_id
    FROM outbox o
    WHERE o.status = $1
    ORDER BY create_date
    FOR UPDATE SKIP LOCKED
    LIMIT $2
)
UPDATE outbox o
SET status = $3,
    attempts = attempts + 1,
    next_try_at = now() + make_interval(secs => $4),
    update_code = $5,
    update_date = $6
FROM picked p
WHERE o.outbox_id = p.outbox_id
RETURNING *;

-- name: MarKDeadBatch :many
WITH picked AS (
    SELECT o.outbox_id
    FROM outbox o
    WHERE o.status = ANY($1::text[]) AND o.next_try_at <= now()
    ORDER BY o.create_date
    FOR UPDATE SKIP LOCKED
    LIMIT $2
)
UPDATE outbox o
SET status      = $3,
    attempts    = o.attempts + 1,
    update_code = $4,
    update_date = $5
    FROM picked p
WHERE o.outbox_id = p.outbox_id
RETURNING o.outbox_id, o.attempts, o.last_error;

-- name: MarkStatus :exec
UPDATE outbox
SET status=$3,
    update_code = $4,
    update_date = $5
WHERE outbox_id = $1
  AND status=$2;

-- name: MarkStatusAndNextTry :exec
UPDATE outbox
SET status=$3,
    last_error = $4,
    next_try_at = now() + make_interval(secs => $5),
    update_code = $6,
    update_date = $7
WHERE outbox_id = $1
  AND status=$2;

-- name: Notify :exec
SELECT pg_notify($1, $2);