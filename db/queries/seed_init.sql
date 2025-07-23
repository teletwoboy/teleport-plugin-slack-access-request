-- name: CreateSeedInit :exec
INSERT INTO seed_init (
    seed_init_id,
    status,
    use_yn,
    create_code,
    create_date,
    version
) VALUES (
    1,
    'uninitialized',
    $1,
    $2,
    $3,
    $4
) ON CONFLICT DO NOTHING;

-- name: GetSeedInitStatus :one
SELECT * FROM seed_init
WHERE seed_init_id = 1;

-- name: UpdateSeedInitStatus :exec
UPDATE seed_init
SET status = 'initialized',
    update_code = $1,
    update_date = $2
WHERE seed_init_id = 1;