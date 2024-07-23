-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN password varchar,
    ADD COLUMN password_confirm varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN password,
    DROP COLUMN password_confirm;
-- +goose StatementEnd
