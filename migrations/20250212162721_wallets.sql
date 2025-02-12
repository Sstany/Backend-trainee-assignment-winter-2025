-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallets (
		username VARCHAR(255) PRIMARY KEY,
		balance INTEGER  NOT NULL CONSTRAINT positive_balance CHECK (balance >= 0)
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
