package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	Conn *sql.DB
}

func (s *Store) Init() error {
	var err error
	s.Conn, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		return err
	}

	return nil
}
