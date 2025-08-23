-- +goose Up
-- +goose StatementBegin
INSERT OR IGNORE INTO types (id, name) VALUES (1, '🏆 Win');
INSERT OR IGNORE INTO types (id, name) VALUES (2, '⛔ Blocker');
INSERT OR IGNORE INTO types (id, name) VALUES (3, '📝 General');
INSERT OR IGNORE INTO types (id, name) VALUES (4, '🔁 Retrospective');
INSERT OR IGNORE INTO types (id, name) VALUES (5, '💡 Idea');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM types WHERE id IN (1, 2, 3, 4, 5);
-- +goose StatementEnd
