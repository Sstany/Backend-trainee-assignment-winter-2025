-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_transactions(
		from VARCHAR(255) PRIMARY KEY,
		to VARCHAR(255) NOT NULL,
        amount INTEGER NOT NULL CONSTRAINT positive_amount CHECK (amount > 0)
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_transactions;
-- +goose StatementEnd
