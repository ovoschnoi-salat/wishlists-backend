-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users
(
    id        BIGINT PRIMARY KEY,
    username  TEXT,
    photo_url TEXT
);

CREATE TABLE IF NOT EXISTS friends
(
    user_id   BIGINT NOT NULL REFERENCES users (id),
    friend_id BIGINT NOT NULL REFERENCES users (id),

    PRIMARY KEY (user_id, friend_id)
);

CREATE TABLE IF NOT EXISTS friends_requests
(
    user_id_from BIGINT NOT NULL REFERENCES users (id),
    user_id_to   BIGINT NOT NULL REFERENCES users (id),

    PRIMARY KEY (user_id_from, user_id_to)
);

CREATE INDEX IF NOT EXISTS friends_requests_user_id_to_idx ON friends_requests (user_id_to);

CREATE TABLE IF NOT EXISTS wishlists
(
    id          BIGSERIAL PRIMARY KEY,
    owner_id    BIGINT  NOT NULL REFERENCES users (id),
    title       TEXT    NOT NULL,
    description TEXT    NOT NULL,
    is_private  BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS wishlists_access_list
(
    list_id BIGINT NOT NULL REFERENCES wishlists (id),
    user_id BIGINT NOT NULL REFERENCES users (id),

    PRIMARY KEY (list_id, user_id)
);

CREATE TABLE IF NOT EXISTS wishes
(
    id          BIGSERIAL PRIMARY KEY,
    owner_id    BIGINT  NOT NULL REFERENCES users (id),
    list_id     BIGINT  NOT NULL REFERENCES wishlists (id),
    title       TEXT    NOT NULL,
    description TEXT    NOT NULL,
    price       TEXT    NOT NULL,
    links       jsonb   NOT NULL,
    reservable  BOOLEAN NOT NULL,
    reserved_by BIGINT REFERENCES users (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS wishes;
DROP TABLE IF EXISTS wishlists_access_list;
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS friends_requests;
DROP TABLE IF EXISTS friends;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
