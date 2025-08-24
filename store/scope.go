package store

import (
	"database/sql"
	"fmt"
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
		return 0, fmt.Errorf("failed to save scope %q: %w", scope.Name, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for scope %q: %w", scope.Name, err)
	}
	return id, nil
}

func (s *Store) GetScopes() ([]Scope, error) {
	rows, err := s.Conn.Query("SELECT id, name FROM scopes")
	if err != nil {
		return nil, fmt.Errorf("failed to query scopes: %w", err)
	}
	defer rows.Close()

	var scopes []Scope

	for rows.Next() {
		var scope Scope
		if err := rows.Scan(&scope.ID, &scope.Name); err != nil {
			return nil, fmt.Errorf("failed to scan scope row: %w", err)
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
			return 0, fmt.Errorf("failed to insert new scope %q: %w", name, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("failed to get last insert ID for scope %q: %w", name, err)
		}
		return id, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to query scope by name %q: %w", name, err)
	}
	return id, nil
}
