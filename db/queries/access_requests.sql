-- name: CreateAccessRequest :one
INSERT INTO access_requests (
	requester_user_id,
	name,
	input_channel_id,
	input_channel_name,
	role,
	reason,
	review_channel_id,
	review_channel_name,
	status,
    expires,
    session_ttl,
    access_duration,
	start_date,
	expiry_date,
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
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17,
    $18
) RETURNING *;

-- name: ExistsAccessRequestByName :one
SELECT EXISTS (
    SELECT 1
    FROM access_requests
    WHERE name = $1 AND use_yn = true
);

-- name: GetAccessRequestByName :one
SELECT *
FROM access_requests
WHERE name = $1 AND use_yn = true;

-- name: UpdateAccessRequestStatusByName :one
UPDATE access_requests
SET status = $1,
    expires = $2,
    session_ttl = $3,
    access_duration = $4,
    start_date = $5,
    expiry_date = $6,
    update_code = $7,
    update_date = $8
WHERE name = $9 AND use_yn = true
RETURNING *;