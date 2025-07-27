-- name: CreateUser :one
INSERT INTO users (
    teleport_user_id,
    slack_user_id,
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
    $6
) RETURNING *;

-- name: DeleteUserByTeleportAndSlackID :one
UPDATE users 
SET use_yn = false,
    delete_code = $3,
    delete_date = $4
WHERE teleport_user_id = $1 AND slack_user_id = $2 AND use_yn = true
RETURNING *;

-- name: GetUserBySlackUserID :one
SELECT *
FROM users
WHERE slack_user_id = $1 AND use_yn = true;