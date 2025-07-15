CREATE TABLE users (
    user_id SERIAL NOT NULL PRIMARY KEY,
    teleport_user_id BIGINT NOT NULL,
    slack_user_id BIGINT NOT NULL,
    FOREIGN KEY (teleport_user_id) REFERENCES teleport_users(teleport_user_id),
    FOREIGN KEY (slack_user_id) REFERENCES slack_users(slack_user_id),
    use_yn BOOLEAN NOT NULL DEFAULT TRUE,
    create_code VARCHAR(255) NOT NULL,
    create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_code VARCHAR(255) NOT NULL,
    update_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delete_code VARCHAR(255) NOT NULL,
    delete_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 0
);