-- name: CreateSlackUser :one
INSERT INTO slack_users (
    id,
    name,
    real_name,
    email,
    deleted,
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
) RETURNING *;