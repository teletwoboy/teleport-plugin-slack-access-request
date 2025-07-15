CREATE TABLE slack_users (
    slack_user_id SERIAL NOT NULL PRIMARY KEY,
    id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    real_name VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    use_yn BOOLEAN NOT NULL DEFAULT TRUE,
    create_code VARCHAR(255) NOT NULL,
    create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_code VARCHAR(255),
    update_date TIMESTAMP,
    delete_code VARCHAR(255),
    delete_date TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 0
);