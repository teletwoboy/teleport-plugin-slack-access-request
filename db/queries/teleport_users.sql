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

-- name: GetTeleportUserByUsername :one
SELECT *
FROM teleport_users
WHERE username = $1 AND use_yn = true;