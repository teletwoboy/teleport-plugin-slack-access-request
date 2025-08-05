CREATE TABLE access_reviews (
    access_review_id SERIAL NOT NULL PRIMARY KEY,
    access_request_id INT NOT NULL,
    reviewer_user_id INT NOT NULL,
    reason TEXT,
    decision VARCHAR(64) NOT NULL CHECK (decision IN ('allow', 'deny')),
    use_yn BOOLEAN NOT NULL DEFAULT TRUE,
    create_code VARCHAR(255) NOT NULL,
    create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_code VARCHAR(255),
    update_date TIMESTAMP,
    delete_code VARCHAR(255),
    delete_date TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 0
);