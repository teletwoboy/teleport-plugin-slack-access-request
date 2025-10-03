-- name: CreateOutbox :exec
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
);

-- name: ClaimNextOutbox :one
WITH picked AS (
    SELECT outbox_id
    FROM outbox
    WHERE status = ANY($5::text[])
      AND (next_try_at IS NULL OR next_try_at <= now())
    ORDER BY create_date
    FOR UPDATE SKIP LOCKED
    LIMIT $6
)
UPDATE outbox o
SET status=$1,
    attempts = attempts + 1,
    next_try_at = now() + make_interval(secs => $2),
    update_code = $3,
    update_date = $4
FROM picked p
WHERE o.outbox_id = p.outbox_id
RETURNING *;

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