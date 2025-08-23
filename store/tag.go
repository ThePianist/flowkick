package store

type Tag struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
