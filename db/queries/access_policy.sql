-- name: CreateAccessPolicy :one
INSERT INTO access_policies (
    user_id,
    input_channel_id,
    input_channel_name,
    title,
    reason,
    start_date,
    end_date,
    effect,
    target_channel_id,
    target_channel_name,
    target_role,
    target_role_name,
    target_slack_id,
    target_real_name,
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

-- name: DeleteAccessPolicyByAccessPolicyID :exec
UPDATE access_policies
SET use_yn = $2,
    delete_code = $3,
    delete_date = $4
WHERE access_policy_id = $1 AND use_yn = true;

-- name: DeleteAccessPolicyByUserID :many
UPDATE access_policies
SET use_yn = $2,
    delete_code = $3,
    delete_date = $4
WHERE user_id = $1 AND use_yn = true
RETURNING *;

-- name: GetAccessPoliciesByAccessPolicyID :one
SELECT *
FROM access_policies
WHERE access_policy_id = $1 AND use_yn = true;

-- name: GetAccessPoliciesByInputChannelID :many
SELECT *
FROM access_policies
WHERE input_channel_id = $1 AND use_yn = true;

-- name: UpdateAccessPolicyMessageTimestamp :exec
UPDATE access_policies
SET message_timestamp = $2,
    update_code = $3,
    update_date = $4
WHERE access_policy_id = $1 AND use_yn = true;