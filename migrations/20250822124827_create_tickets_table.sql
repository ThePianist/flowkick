-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tickets (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  name TEXT UNIQUE NOT NULL
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tickets;
-- +goose StatementEnd
