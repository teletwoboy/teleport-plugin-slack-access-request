-- name: CreateOutbox :exec
INSERT INTO outbox (
    event_type,
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
    $8
);

-- name: ClaimNextOutbox :one
WITH picked AS (
    SELECT outbox_id
    FROM outbox
    WHERE status IN ('pending', 'failed')
        AND (next_try_at IS NULL OR next_try_at <= now())
    ORDER BY create_date
    LIMIT 1
)
UPDATE outbox o
SET status='processing',
    attempts = attempts + 1,
    update_code = 'worker',
    update_date = now()
FROM picked p
WHERE o.outbox_id = p.outbox_id
RETURNING *;

-- name: MarkFailed :exec
UPDATE outbox
SET status='failed',
    last_error = $2,
    next_try_at = now() + interval '4 seconds',
    update_code = 'worker',
    update_date = now()
WHERE outbox_id = $1
  AND status='processing';