-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS entries (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	note TEXT NOT NULL,
	scope_id INTEGER,
	type_id INTEGER,
	ticket_id INTEGER,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(scope_id) REFERENCES scopes(id),
	FOREIGN KEY(type_id) REFERENCES types(id),
	FOREIGN KEY(ticket_id) REFERENCES tickets(id)
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS entries;
-- +goose StatementEnd
