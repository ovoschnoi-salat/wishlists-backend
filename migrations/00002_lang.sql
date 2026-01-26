-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN language text NOT NULL DEFAULT 'en';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN language;
-- +goose StatementEnd
