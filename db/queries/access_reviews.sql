-- name: CreateAccessReview :one
INSERT INTO access_reviews (
    access_request_id,
    reviewer_user_id,
    reason,
    decision,
    use_yn,
    create_code,
    create_date
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) RETURNING *;