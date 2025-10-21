-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users
(
    id        BIGINT PRIMARY KEY,
    username  TEXT NOT NULL,
    name      TEXT NOT NULL,
    photo_url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS friends
(
    user_id   BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    friend_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,

    PRIMARY KEY (user_id, friend_id)
);

CREATE TABLE IF NOT EXISTS friends_requests
(
    user_id_to   BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    user_id_from BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,

    PRIMARY KEY (user_id_to, user_id_from)
);

CREATE INDEX IF NOT EXISTS friends_requests_user_id_from_idx ON friends_requests (user_id_from);

CREATE TABLE IF NOT EXISTS wishlists
(
    id          BIGSERIAL PRIMARY KEY,
    owner_id    BIGINT  NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title       TEXT    NOT NULL,
    description TEXT    NOT NULL,
    is_private  BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS wishlist_access_list
(
    list_id BIGINT NOT NULL REFERENCES wishlists (id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    owner_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,

    PRIMARY KEY (list_id, user_id)
);

CREATE INDEX IF NOT EXISTS wishlist_access_list_owner_id_user_id_idx ON wishlist_access_list (owner_id, user_id);

CREATE TABLE IF NOT EXISTS wishlist_items
(
    id          BIGSERIAL PRIMARY KEY,
    owner_id    BIGINT  NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    wishlist_id BIGINT  NOT NULL REFERENCES wishlists (id) ON DELETE CASCADE,
    title       TEXT    NOT NULL,
    description TEXT    NOT NULL,
    price       TEXT    NOT NULL,
    links       jsonb   NOT NULL,
    reservable  BOOLEAN NOT NULL,
    reserved_by BIGINT REFERENCES users (id) ON DELETE SET NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS wishlist_items;
DROP TABLE IF EXISTS wishlist_access_list;
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS friends_requests;
DROP TABLE IF EXISTS friends;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
