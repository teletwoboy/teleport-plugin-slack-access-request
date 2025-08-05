CREATE TABLE teleport_users (
    teleport_user_id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    use_yn BOOLEAN NOT NULL DEFAULT TRUE,
    create_code VARCHAR(255) NOT NULL,
    create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_code VARCHAR(255),
    update_date TIMESTAMP,
    delete_code VARCHAR(255),
    delete_date TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 0
);