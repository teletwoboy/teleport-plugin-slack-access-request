-- name: CreateTeleportUser :one
INSERT INTO teleport_users (
    username,
    use_yn,
    create_code,
    create_date,
    version
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: DeleteTeleportUserUseYnByUsername :one
UPDATE teleport_users 
SET use_yn = false,
    delete_code = $2,
    delete_date = $3
WHERE username = $1 AND use_yn = true
RETURNING *;

-- name: ExistsTeleportUserByUsername :one
SELECT EXISTS (
    SELECT 1
    FROM teleport_users
    WHERE username = $1 AND use_yn = true
);

-- name: GetTeleportUserByTeleportUserID :one
SELECT *
FROM teleport_users
WHERE teleport_user_id = $1 AND use_yn = true;

-- name: GetTeleportUserByUsername :one
SELECT *
FROM teleport_users
WHERE username = $1 AND use_yn = true;