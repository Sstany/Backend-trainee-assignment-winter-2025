-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS inventory_user
ON inventory (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS inventory_user
-- +goose StatementEnd
