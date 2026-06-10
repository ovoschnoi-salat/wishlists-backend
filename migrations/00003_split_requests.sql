-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wishlist_items_split_requests
(
    list_id    BIGINT    NOT NULL REFERENCES wishlists (id) ON DELETE CASCADE,
    owner_id   BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    wish_id    BIGINT    NOT NULL REFERENCES wishlist_items (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id    BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,

    PRIMARY KEY (wish_id, user_id)
);

CREATE INDEX wishlist_items_split_requests_list_id_user_id_idx ON wishlist_items_split_requests (list_id, user_id);

CREATE TYPE split_request_privacy as enum ('invisible_to_owner', 'visible_to_owner');

ALTER TABLE wishlists ADD COLUMN split_request_privacy split_request_privacy NOT NULL DEFAULT 'invisible_to_owner';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE wishlists DROP COLUMN split_request_privacy;

DROP TYPE IF EXISTS split_request_privacy;

DROP TABLE IF EXISTS wishlist_items_split_requests;
-- +goose StatementEnd
