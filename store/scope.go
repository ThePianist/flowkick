package store

import (
	"database/sql"
)

type Scope struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (s *Store) SaveScope(scope Scope) (int64, error) {
	upsertScope := `INSERT INTO scopes (name)
	VALUES (?)
	ON CONFLICT(name) DO UPDATE SET name=excluded.name;`

	result, err := s.Conn.Exec(upsertScope, scope.Name)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (s *Store) GetScopes() ([]Scope, error) {
	rows, err := s.Conn.Query("SELECT id, name FROM scopes")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var scopes []Scope

	for rows.Next() {
		var scope Scope
		if err := rows.Scan(&scope.ID, &scope.Name); err != nil {
			return nil, err
		}
		scopes = append(scopes, scope)
	}

	return scopes, nil
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
