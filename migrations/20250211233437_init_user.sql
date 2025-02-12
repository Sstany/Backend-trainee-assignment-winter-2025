-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
		username VARCHAR(255) PRIMARY KEY,
		password VARCHAR(255) NOT NULL
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
