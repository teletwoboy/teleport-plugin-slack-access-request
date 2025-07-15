CREATE TABLE access_request_reviews (
    access_request_review_id SERIAL NOT NULL PRIMARY KEY,
    access_request_id BIGINT NOT NULL,
    reviewer_user_id BIGINT NOT NULL,
    reason TEXT,
    decision VARCHAR(64) NOT NULL CHECK (decision IN ('allow', 'deny')),
    use_yn BOOLEAN NOT NULL DEFAULT TRUE,
    create_code VARCHAR(255) NOT NULL,
    create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_code VARCHAR(255) NOT NULL,
    update_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delete_code VARCHAR(255) NOT NULL,
    delete_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 0
);