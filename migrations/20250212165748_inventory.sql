-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS inventory (
        id SERIAL PRIMARY KEY,
		username VARCHAR(255),
		item VARCHAR(255)
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS inventory;
-- +goose StatementEnd
