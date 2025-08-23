package store

import (
	"database/sql"
	"time"
)

type Entry struct {
	ID        int64         `db:"id"`
	Note      string        `db:"note"`
	ScopeID   sql.NullInt64 `db:"scope_id"`
	TypeID    sql.NullInt64 `db:"type_id"`
	TicketID  sql.NullInt64 `db:"ticket_id"`
	CreatedAt time.Time     `db:"created_at"`
}

func (s *Store) SaveEntry(entry Entry) error {
	saveEntryStatement := `INSERT INTO entries (id, note, scope_id, type_id, ticket_id, created_at)
	VALUES (?, ?, ?, ?, ?, ?)`
	if entry.ID == 0 {
		entry.ID = time.Now().UnixNano()
	}
	_, err := s.Conn.Exec(saveEntryStatement, entry.ID, entry.Note, entry.ScopeID, entry.TypeID, entry.TicketID, entry.CreatedAt)
	return err
}
