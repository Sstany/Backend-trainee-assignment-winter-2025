-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS index_user
ON user_transactions(user_from, user_to);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS index_user
-- +goose StatementEnd
