CREATE TABLE seed_init(
    seed_init_id SERIAL NOT NULL PRIMARY KEY,
    status VARCHAR(64) NOT NULL DEFAULT 'uninitialized' CHECK (status IN ('uninitialized', 'initialized')),
    use_yn BOOLEAN NOT NULL DEFAULT TRUE,
    create_code VARCHAR(255) NOT NULL,
    create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_code VARCHAR(255),
    update_date TIMESTAMP,
    delete_code VARCHAR(255),
    delete_date TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 0
);