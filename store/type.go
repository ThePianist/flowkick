package store

type Type struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (s *Store) GetTypes() ([]string, error) {
	rows, err := s.Conn.Query("SELECT name FROM types ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var types []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		types = append(types, name)
	}
	return types, nil
}
