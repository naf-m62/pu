-- migrate create -ext sql -dir migrations migration_name
-- migrate -path migrations -database "postgres://postgres:postgres@127.0.0.1:5432/user?sslmode=disable" up
CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    name          TEXT                     NOT NULL,
    email         TEXT UNIQUE              NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    points        BIGINT,
    password_hash TEXT                     NOT NULL,
    salt          TEXT                     NOT NULL
);