package store

import "fmt"

type Type struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (s *Store) GetTypes() ([]string, error) {
	rows, err := s.Conn.Query("SELECT name FROM types ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to query types: %w", err)
	}
	defer rows.Close()
	var types []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan type row: %w", err)
		}
		types = append(types, name)
	}
	return types, nil
}
