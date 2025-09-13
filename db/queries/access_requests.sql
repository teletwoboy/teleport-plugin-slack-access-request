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
	state,
    start_date,
    access_duration,
	request_ttl,
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
    $16
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

-- name: GetAccessRequestStateByName :one
SELECT state
FROM access_requests
WHERE name = $1 AND use_yn = true;

-- name: UpdateAccessRequestStateByName :one
UPDATE access_requests
SET state = $1,
    start_date = $2,
    request_ttl = $3,
    update_code = $4,
    update_date = $5
WHERE name = $6 AND use_yn = true
RETURNING *;