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

-- name: GetUserBySlackUserID :one
SELECT *
FROM users
WHERE slack_user_id = $1 AND use_yn = true;