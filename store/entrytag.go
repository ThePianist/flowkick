package store

type EntryTag struct {
	EntryID int64 `db:"entry_id"`
	TagID   int64 `db:"tag_id"`
}
