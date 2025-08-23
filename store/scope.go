package store

import (
	"database/sql"
)

type Scope struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (s *Store) GetScopes() ([]string, error) {
	rows, err := s.Conn.Query("SELECT name FROM scopes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scopes []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		scopes = append(scopes, name)
	}
	return scopes, nil
}

func (s *Store) SaveScope(scopeID int64, name string) error {
	upsertScope := `INSERT INTO scopes (id, name)
	VALUES (?, ?)
	ON CONFLICT(id) DO UPDATE SET name=excluded.name;`
	_, err := s.Conn.Exec(upsertScope, scopeID, name)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetScopeIDByName(name string) (int64, error) {
	var id int64
	err := s.Conn.QueryRow("SELECT id FROM scopes WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		// Insert new scope and return its ID
		res, err := s.Conn.Exec("INSERT INTO scopes (name) VALUES (?)", name)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	return id, err
}
