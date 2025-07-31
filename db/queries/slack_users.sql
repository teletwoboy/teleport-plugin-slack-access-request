-- name: CreateSlackUser :one
INSERT INTO slack_users (
    id,
    name,
    real_name,
    email,
    time_zone,
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

-- name: DeleteSlackUserByName :one
UPDATE slack_users 
SET use_yn = $2,
    delete_code = $3,
    delete_date = $4
WHERE name = $1 AND use_yn = true
RETURNING *;

-- name: ExistsSlackUserByID :one
SELECT EXISTS (
    SELECT 1
    FROM slack_users
    WHERE id = $1 AND use_yn = true
);

-- name: GetSlackUserByID :one
SELECT *
FROM slack_users
WHERE id = $1 AND use_yn = true;

-- name: GetSlackUserByName :one
SELECT *
FROM slack_users
WHERE name = $1 AND use_yn = true;

-- name: GetSlackUserBySlackUserID :one
SELECT *
FROM slack_users
WHERE slack_user_id = $1 AND use_yn = true;